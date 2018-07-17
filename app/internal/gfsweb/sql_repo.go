package gfsweb

import (
	"database/sql"
	"io"

	"github.com/AgencyPMG/go-from-scratch/app/internal/data/storage/sqlrepo"
	"github.com/gogolfing/config"

	_ "github.com/lib/pq"
)

const (
	ConfigKeyDatabaseURL = "databases.app.url"

	SQLDatasourceName = "postgres"
)

var sqlrepoDialect = sqlrepo.PostgresqlDialect{}

//CreateSQLRepo returns a new sqlrepo.Repo and an io.Closer that will
//close the connection to the database.
//
//It connects to a Postgresql database at the url specified in config at key ConfigKeyDatabaseURL.
func CreateSQLRepo(config *config.Config) (*sqlrepo.Repo, io.Closer, error) {
	db, err := sql.Open(
		SQLDatasourceName,
		config.GetString(ConfigKeyDatabaseURL),
	)
	if err != nil {
		return nil, nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, nil, err
	}

	sqlRepo := sqlrepo.New(db, sqlrepoDialect)

	return sqlRepo, db, nil
}
