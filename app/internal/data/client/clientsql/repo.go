package clientsql

import (
	"context"
	"fmt"

	"github.com/AgencyPMG/go-from-scratch/app/internal/data"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/client"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/storage/sqlrepo"
)

var _ client.Repo = &Repo{} //Ensure *Repo is a client.Repo.

//Table is table we are querying for our client entities.
const Table = "clients"

//SelectFrom is the query to get all of our columns without any filtering, ordering, etc.
const SelectFrom = `SELECT
	c.id,
	c.name,
	c.created_at,
	c.updated_at
	FROM ` + Table + ` AS c`

//Repo is a client.Repo implementation that uses a SQL database as storage.
type Repo struct {
	db *sqlrepo.Repo
}

//New returns a new Repo that uses repo as its helper for talking with the database.
func New(repo *sqlrepo.Repo) *Repo {
	return &Repo{
		db: repo,
	}
}

//Get is the client.QueryRepo implementation.
func (r *Repo) Get(ctx context.Context, id data.Id) (*client.Client, error) {
	row := r.db.QueryRowContext(
		ctx,
		fmt.Sprintf("%s WHERE id = ?", SelectFrom),
		id,
	)
	return scan(row)
}

//List is the client.QueryRepo implementation.
func (r *Repo) List(ctx context.Context) ([]*client.Client, error) {
	query := orderQuery(SelectFrom)
	clients, err := r.list(ctx, query)

	return clients, err
}

//List is the client.QueryRepo implementation.
func (r *Repo) ListByIds(ctx context.Context, ids []data.Id) ([]*client.Client, error) {
	if len(ids) == 0 {
		return []*client.Client{}, nil
	}

	placeholders, args := sqlrepo.IdsPlaceholdersArgs(ids)

	query := fmt.Sprintf("%s WHERE id IN (%s)", SelectFrom, placeholders)
	query = orderQuery(query)

	return r.list(ctx, query, args...)
}

//orderQuery is a helper function to take an unordered select query and add
//ordering to it that the application expects.
func orderQuery(query string) string {
	return fmt.Sprintf("%s ORDER BY name ASC", query)
}

//list is a helper method to query for a list of Clients independent of the
//actual query.
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

//scan is a helper function to scan a new Client from a sql row or rows.
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

//Add is the client.Repo implementation.
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

//Set is the client.Repo implementation.
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

//Remove is the client.Repo implementation.
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
