package echo

import (
	"context"
	"github.com/pkg/errors"

	"github.com/jackc/pgx/v5"
)

type Echo struct {
	ID       uint32
	Title    string
	Reminder string
}

type Repository interface {
	FindAll(ctx context.Context) ([]Echo, error)
	CreateOrUpdate(ctx context.Context, e Echo) (id uint32, err error)
}

func New(tx pgx.Tx) Repository {
	return &repository{tx: tx}
}

type repository struct {
	tx pgx.Tx
}

func (r *repository) FindAll(ctx context.Context) ([]Echo, error) {
	rows, err := r.tx.Query(ctx, `select id, title, reminder from echo`)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query echos")
	}
	defer rows.Close()
	result := make([]Echo, 0, rows.CommandTag().RowsAffected())
	for rows.Next() {
		var e Echo
		if err := rows.Scan(&e.ID, &e.Title, &e.Reminder); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		result = append(result, e)
	}
	return result, nil
}

func (r *repository) CreateOrUpdate(ctx context.Context, e Echo) (uint32, error) {
	if err := r.tx.QueryRow(ctx,
		`INSERT INTO echo(id, title, reminder) VALUES ($1, $2, $3) 
         		ON CONFLICT ON CONSTRAINT unq_unq_id_title_echo 
				DO UPDATE SET reminder = excluded.reminder
				RETURNING id`,
		e.ID, e.Title, e.Reminder).Scan(&e.ID); err != nil {
		return 0, errors.Wrap(err, "failed to insert row")
	}
	return e.ID, nil
}
