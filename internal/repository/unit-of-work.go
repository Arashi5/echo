package repository

import (
	"context"
	"fmt"
	"github.com/arashi5/echo/internal/repository/echo"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Instance struct {
	pool *pgxpool.Pool
}

func NewRepositoryUnitOfWork(pool *pgxpool.Pool) *Instance {
	return &Instance{pool: pool}
}

func (i *Instance) New(ctx context.Context) (ofw UnitOfWork, err error) {
	tx, err := i.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &unitOfWork{
		tx: tx,
	}, nil
}

type unitOfWork struct {
	db          *pgxpool.Pool
	tx          pgx.Tx
	isCommitted bool

	echo echo.Repository
}

type DB interface {
	Transaction(ctx context.Context, fn func() error) error
	Rollback(ctx context.Context) error
	Commit(ctx context.Context) error
}

// Commit фиксирует изменения транзакции.
func (u *unitOfWork) Commit(ctx context.Context) error {
	u.isCommitted = true
	return u.tx.Commit(ctx)
}

// Rollback откатывает изменения транзакции.
func (u *unitOfWork) Rollback(ctx context.Context) error {
	if !u.isCommitted {
		return u.tx.Rollback(ctx)
	} else {
		u.isCommitted = false
		return nil
	}
}

func (u *unitOfWork) Transaction(ctx context.Context, fn func() error) error {
	if err := fn(); err != nil {
		if rollbackErr := u.tx.Rollback(ctx); rollbackErr != nil {
			return fmt.Errorf("rollback tx: %w", rollbackErr)
		}
		return fmt.Errorf("rollback: %w", err)
	}
	return u.tx.Commit(ctx)
}

type UnitOfWork interface {
	PoolRepository
	DB
}

type PoolRepository interface {
	GetEchoRepository() echo.Repository
}

func (u *unitOfWork) GetEchoRepository() echo.Repository {
	return u.echo
}
