package postgres

import (
	// imports the driver, don't remove this comment, golint requires.
	_ "github.com/lib/pq"
	"github.com/l306287405/go-zero/core/stores/sqlx"
)

const postgresDriverName = "postgres"

// New returns a postgres connection.
func New(datasource string, opts ...sqlx.SqlOption) sqlx.SqlConn {
	return sqlx.NewSqlConn(postgresDriverName, datasource, opts...)
}
