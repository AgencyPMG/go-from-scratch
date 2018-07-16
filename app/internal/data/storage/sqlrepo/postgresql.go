package sqlrepo

import "fmt"

var _ Dialect = PostgresDialect{} //Ensure PostgresDialect{} is a Dialect.

//PostgresDialect is a Dialect that understands the Postgresql DBMS.
type PostgresDialect struct{}

//Placeholder is the Dialect implementation.
func (p PostgresDialect) Placeholder(index int) string {
	return fmt.Sprintf("$%d", index+1)
}
