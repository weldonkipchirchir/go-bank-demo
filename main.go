package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/weldonkipchirchir/simple_bank/api"
	db "github.com/weldonkipchirchir/simple_bank/db/sqlc"
	"github.com/weldonkipchirchir/simple_bank/grpc_api"
	"github.com/weldonkipchirchir/simple_bank/mail"
	"github.com/weldonkipchirchir/simple_bank/pb"
	"github.com/weldonkipchirchir/simple_bank/util"
	"github.com/weldonkipchirchir/simple_bank/worker"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	//connect to db
	connPool, err := pgxpool.New(context.Background(), config.DBSource)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect to db")
	}
	store := db.NewStore(connPool)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	//create new task distributor
	taskDistributor := worker.NewRedisDistributor(redisOpt)

	//run servers and workers in separate go routines
	waitGroup, ctx := errgroup.WithContext(ctx)

	runTaskProcessor(ctx, waitGroup, config, redisOpt, store) //go runTaskProcessor(config, redisOpt, store)
	// runGinServer(config, store)
	runGatewayServer(ctx, waitGroup, config, store, taskDistributor) //go runGatewayServer(config, store, taskDistributor)
	runGrpcServer(ctx, waitGroup, config, store, taskDistributor)    //runGrpcServer(config, store, taskDistributor)

	//wait from all process
	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("error from waitgroup")
	}
}

func runTaskProcessor(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config util.Config,
	redisOpt asynq.RedisClientOpt,
	store db.Store,
) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcess(redisOpt, store, mailer)

	log.Info().Msg("start task processor")
	//taskProcessor has already started a go routine
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown task processor")

		taskProcessor.Shutdown()
		log.Info().Msg("task processor is stopped")

		return nil
	})
}

func runGrpcServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config util.Config,
	store db.Store,
	taskDistributor worker.TaskDistributor) {
	server, err := grpc_api.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
	}

	grpcLogger := grpc.UnaryInterceptor(grpc_api.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)

	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err)
	}

	//start server in a go routine
	waitGroup.Go(func() error {
		log.Info().Msgf("start gRPC server at %s", listener.Addr().String())

		if err := grpcServer.Serve(listener); err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			log.Error().Err(err)
			return err
		}
		return nil
	})

	//listen to interrupt signals
	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown gRPC server")
		grpcServer.GracefulStop()
		log.Info().Msg("gRPC server is stopped")
		return nil
	})
}
func runGatewayServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config util.Config,
	store db.Store,
	taskDistributor worker.TaskDistributor,
) {
	server, err := grpc_api.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create handler server:")
	}

	mux := http.NewServeMux() //receive http request
	mux.Handle("/", grpcMux)  //route to grpc

	httpServer := &http.Server{
		Handler: grpc_api.HttpLogger(mux),
		Addr:    config.ServerAddress,
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("start HTTP gateway server at %s", httpServer.Addr)
		err = httpServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Error().Err(err).Msg("HTTP gateway server failed to serve")
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown HTTP gateway server")

		err := httpServer.Shutdown(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("failed to shutdown HTTP gateway server")
			return err
		}

		log.Info().Msg("HTTP gateway server is stopped")
		return nil
	})

}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start server:")
	}
}
