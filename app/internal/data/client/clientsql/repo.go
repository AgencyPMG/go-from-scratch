package clientsql

import (
	"context"
	"fmt"

	"github.com/AgencyPMG/go-from-scratch/app/internal/data"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/client"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/storage/sqlrepo"
)

var _ client.Repo = &Repo{} //Ensure *Repo is a client.Repo.

const Table = "clients"

const SelectFrom = `SELECT
	c.id,
	c.name,
	c.created_at,
	c.updated_at
	FROM ` + Table + ` AS c`

type Repo struct {
	db *sqlrepo.Repo
}

func New(repo *sqlrepo.Repo) *Repo {
	return &Repo{
		db: repo,
	}
}

func (r *Repo) Get(ctx context.Context, id data.Id) (*client.Client, error) {
	row := r.db.QueryRowContext(
		ctx,
		fmt.Sprintf("%s WHERE id = ?", SelectFrom),
		id,
	)
	return scan(row)
}

func (r *Repo) List(ctx context.Context) ([]*client.Client, error) {
	query := orderQuery(SelectFrom)
	clients, err := r.list(ctx, query)

	return clients, err
}

func (r *Repo) ListByIds(ctx context.Context, ids []data.Id) ([]*client.Client, error) {
	if len(ids) == 0 {
		return []*client.Client{}, nil
	}

	placeholders, args := sqlrepo.IdsPlaceholdersArgs(ids)

	query := fmt.Sprintf("%s WHERE id IN (%s)", SelectFrom, placeholders)
	query = orderQuery(query)

	return r.list(ctx, query, args...)
}

func orderQuery(query string) string {
	return fmt.Sprintf("%s ORDER BY name ASC", query)
}

func (r *Repo) list(ctx context.Context, query string, args ...interface{}) ([]*client.Client, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*client.Client{}
	for rows.Next() {
		client, err := scan(rows)
		if err != nil {
			return nil, err
		}

		result = append(result, client)
	}

	return result, rows.Err()
}

func scan(s sqlrepo.Scanner) (*client.Client, error) {
	client := &client.Client{}

	createdAt, updatedAt := sqlrepo.UTCTime{}, sqlrepo.UTCTime{}

	err := s.Scan(
		&client.Id,
		&client.Name,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	client.CreatedAt, client.UpdatedAt = createdAt.Time, updatedAt.Time

	return client, nil
}

func (r *Repo) Add(ctx context.Context, c client.Client) error {
	work := func(qec sqlrepo.QueryExecerContext) error {
		_, err := qec.ExecContext(
			ctx,
			fmt.Sprintf("INSERT INTO %s (id, name, created_at, updated_at) VALUES (?, ?, ?, ?)", Table),
			c.Id,
			c.Name,
			sqlrepo.UTCTime{c.CreatedAt},
			sqlrepo.UTCTime{c.UpdatedAt},
		)
		return err
	}
	return r.db.TxWorkContext(ctx, work)
}

func (r *Repo) Set(ctx context.Context, c client.Client) error {
	work := func(qec sqlrepo.QueryExecerContext) error {
		_, err := qec.ExecContext(
			ctx,
			fmt.Sprintf("UPDATE %s SET name = ?, created_at = ?, updated_at = ? WHERE id = ?", Table),
			c.Name,
			sqlrepo.UTCTime{c.CreatedAt},
			sqlrepo.UTCTime{c.UpdatedAt},
			c.Id,
		)
		return err
	}
	return r.db.TxWorkContext(ctx, work)
}

func (r *Repo) Remove(ctx context.Context, id data.Id) error {
	work := func(qec sqlrepo.QueryExecerContext) error {
		_, err := qec.ExecContext(
			ctx,
			fmt.Sprintf("DELETE FROM %s WHERE id = ?", Table),
			id,
		)
		return err
	}
	return r.db.TxWorkContext(ctx, work)
}
