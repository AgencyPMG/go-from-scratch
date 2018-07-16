package sqlrepo

import "fmt"

var _ Dialect = PostgresqlDialect{} //Ensure PostgresqlDialect{} is a Dialect.

//PostgresqlDialect is a Dialect that understands the Postgresql DBMS.
type PostgresqlDialect struct{}

//Placeholder is the Dialect implementation.
func (p PostgresqlDialect) Placeholder(index int) string {
	return fmt.Sprintf("$%d", index+1)
}
