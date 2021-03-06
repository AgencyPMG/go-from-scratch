package usersql

import (
	"context"
	"fmt"

	"github.com/AgencyPMG/go-from-scratch/app/internal/data"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/storage/sqlrepo"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/user"
)

var _ user.Repo = &Repo{} //Ensure *Repo is a user.Repo.

const (
	//Table is the table we query from to get our user entities.
	Table = "users"

	//TableUserClients is the table for our user, client relationships.
	TableUserClients = "user_clients"
)

//SelectFrom is our select query without filtering, ordering, etc.
const SelectFrom = `SELECT
	u.id,
	u.email,
	u.created_at,
	u.updated_at,
	u.enabled
	FROM ` + Table + ` AS u`

//SelectFromUserClients is our select query for user, client relationships.
const SelectFromUserClients = `SELECT user_id, client_id FROM ` + TableUserClients

//Repo is a user.Repo implementation that uses a SQL database as its storage.
type Repo struct {
	db *sqlrepo.Repo
}

//New returns a new Repo that uses repo to talk to the database.
func New(repo *sqlrepo.Repo) *Repo {
	return &Repo{
		db: repo,
	}
}

//Get is the user.QueryRepo implementation.
func (r *Repo) Get(ctx context.Context, id data.Id) (*user.User, error) {
	return r.get(
		ctx,
		fmt.Sprintf("%s WHERE id = ?", SelectFrom),
		id,
	)
}

//GetEmail is the user.QueryRepo implementation.
func (r *Repo) GetEmail(ctx context.Context, email string) (*user.User, error) {
	return r.get(
		ctx,
		fmt.Sprintf("%s WHERE email = ?", SelectFrom),
		email,
	)
}

//get is a helper method to get a user from an independent get query.
func (r *Repo) get(ctx context.Context, query string, args ...interface{}) (*user.User, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	u, err := scan(row)
	if err != nil {
		return nil, err
	}

	err = r.populateClientIds(ctx, []*user.User{u})

	return u, err
}

//List is the user.QueryRepo implementation.
func (r *Repo) List(ctx context.Context) ([]*user.User, error) {
	rows, err := r.db.QueryContext(
		ctx,
		fmt.Sprintf("%s ORDER BY email ASC", SelectFrom),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*user.User{}
	for rows.Next() {
		u, err := scan(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, u)
	}

	if err := r.populateClientIds(ctx, result); err != nil {
		return nil, err
	}

	return result, rows.Err()
}

//scan is a helper function to scan a single User entity from a sql row or rows.
func scan(s sqlrepo.Scanner) (*user.User, error) {
	user := &user.User{}

	createdAt, updatedAt := sqlrepo.UTCTime{}, sqlrepo.UTCTime{}

	if err := s.Scan(
		&user.Id,
		&user.Email,
		&createdAt,
		&updatedAt,
		&user.Enabled,
	); err != nil {
		return nil, err
	}

	user.CreatedAt, user.UpdatedAt = createdAt.Time, updatedAt.Time

	return user, nil
}

//populateCilentIds is a helper method to populate all client ids on all users
//in users.
//
//It should be called with all available users to avoid the N+1 problem.
func (r *Repo) populateClientIds(ctx context.Context, users []*user.User) error {
	if len(users) == 0 {
		return nil
	}

	usersById := make(map[data.Id]*user.User, len(users))
	ids := make([]data.Id, len(users))
	for i, u := range users {
		usersById[u.Id] = u
		ids[i] = u.Id
	}

	placeholders, args := sqlrepo.IdsPlaceholdersArgs(ids)
	rows, err := r.db.QueryContext(
		ctx,
		fmt.Sprintf("%s WHERE user_id IN (%s) ORDER BY user_id, client_id", SelectFromUserClients, placeholders),
		args...,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		userId, clientId, err := scanUserIdClientId(rows)
		if err != nil {
			return err
		}

		user := usersById[userId]
		user.ClientIds = append(user.ClientIds, clientId)
	}

	return rows.Err()
}

func scanUserIdClientId(s sqlrepo.Scanner) (data.Id, data.Id, error) {
	userId, clientId := data.EmptyId(), data.EmptyId()
	err := s.Scan(&userId, &clientId)
	return userId, clientId, err
}

//Add is the user.Repo implementation.
func (r *Repo) Add(ctx context.Context, user user.User) error {
	work := func(qec sqlrepo.QueryExecerContext) error {
		_, err := qec.ExecContext(
			ctx,
			fmt.Sprintf(
				`INSERT INTO %s (
					id,
					email,
					created_at,
					updated_at,
					enabled
				) VALUES (?, ?, ?, ?, ?)`,
				Table,
			),
			user.Id,
			user.Email,
			sqlrepo.UTCTime{user.CreatedAt},
			sqlrepo.UTCTime{user.UpdatedAt},
			user.Enabled,
		)
		if err != nil {
			return err
		}
		return saveClientIds(qec, ctx, user)
	}

	return r.db.TxWorkContext(ctx, work)
}

//Set is the user.Repo implementation.
func (r *Repo) Set(ctx context.Context, user user.User) error {
	work := func(qec sqlrepo.QueryExecerContext) error {
		_, err := qec.ExecContext(
			ctx,
			fmt.Sprintf(
				`UPDATE %s SET
					email = ?,
					created_at = ?,
					updated_at = ?,
					enabled = ?
					WHERE id = ?`,
				Table,
			),
			user.Email,
			sqlrepo.UTCTime{user.CreatedAt},
			sqlrepo.UTCTime{user.UpdatedAt},
			user.Enabled,
			user.Id,
		)
		if err != nil {
			return err
		}
		return saveClientIds(qec, ctx, user)
	}

	return r.db.TxWorkContext(ctx, work)
}

//saveClientIds is a helper function to do a full update of client ids for a
//single user.
func saveClientIds(qec sqlrepo.QueryExecerContext, ctx context.Context, user user.User) error {
	if err := removeClientIds(qec, ctx, user); err != nil {
		return err
	}
	return setClientIds(qec, ctx, user)
}

//removeClientIds removes all client ids for user.
func removeClientIds(qec sqlrepo.QueryExecerContext, ctx context.Context, user user.User) error {
	_, err := qec.ExecContext(
		ctx,
		fmt.Sprintf("DELETE FROM %s WHERE user_id = ?", TableUserClients),
		user.Id,
	)
	return err
}

//setClientIds adds all client ids for user.
func setClientIds(qec sqlrepo.QueryExecerContext, ctx context.Context, user user.User) error {
	if len(user.ClientIds) == 0 {
		return nil
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (user_id, client_id) VALUES %s",
		TableUserClients,
		sqlrepo.List("(?,?)", len(user.ClientIds)),
	)

	args := make([]interface{}, 0, 2*len(user.ClientIds))
	for _, clientId := range user.ClientIds {
		args = append(args, user.Id, clientId)
	}

	_, err := qec.ExecContext(ctx, query, args...)
	return err
}
