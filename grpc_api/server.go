package grpc_api

import (
	"fmt"

	db "github.com/weldonkipchirchir/simple_bank/db/sqlc"
	"github.com/weldonkipchirchir/simple_bank/pb"
	"github.com/weldonkipchirchir/simple_bank/token"
	"github.com/weldonkipchirchir/simple_bank/util"
)

// server serves HTTP requests for the banking service
type Server struct {
	pb.UnimplementedSimpleBankServer //enables the create and login to accept calls before it is implementedt blocking
	config                           util.Config
	store                            db.Store
	tokenMaker                       token.Maker
}

// NewServer creates a new grpc server
func NewServer(config util.Config, store db.Store) (*Server, error) {
	//use either NewJWTMaker or NewPasetoMaker - comment one to use the other
	// tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{store: store, tokenMaker: tokenMaker, config: config}

	return server, nil
}
