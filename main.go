package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/weldonkipchirchir/simple_bank/api"
	db "github.com/weldonkipchirchir/simple_bank/db/sqlc"
	"github.com/weldonkipchirchir/simple_bank/grpc_api"
	"github.com/weldonkipchirchir/simple_bank/pb"
	"github.com/weldonkipchirchir/simple_bank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Failed to load config")
	}
	//connect to db
	connection, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to db")
	}
	store := db.NewStore(connection)

	// runGinServer(config, store)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := grpc_api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {

		panic(err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())

	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}
func runGatewayServer(config util.Config, store db.Store) {
	server, err := grpc_api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot create handler server:", err)
	}

	mux := http.NewServeMux() //receive http request
	mux.Handle("/", grpcMux)  //route to grpc

	listener, err := net.Listen("tcp", config.ServerAddress)
	if err != nil {

		panic(err)
	}

	log.Printf("start HTTP gateway server at %s", listener.Addr().String())

	err = http.Serve(listener, mux)

	if err != nil {
		log.Fatal("cannot create server:", err)
	}

}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
