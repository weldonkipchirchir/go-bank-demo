package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides all functions to execute db queries and transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResults, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}

// SQLStore provides all functions to execute db queries and transactions
type SQLStore struct {
	*Queries               //  holds the methods to execute queries
	connPool *pgxpool.Pool // Needed to create a new db transaction
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
