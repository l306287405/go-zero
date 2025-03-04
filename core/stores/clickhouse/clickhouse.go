package clickhouse

import (
	// imports the driver, don't remove this comment, golint requires.
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/l306287405/go-zero/core/stores/sqlx"
)

const clickHouseDriverName = "clickhouse"

// New returns a clickhouse connection.
func New(datasource string, opts ...sqlx.SqlOption) sqlx.SqlConn {
	return sqlx.NewSqlConn(clickHouseDriverName, datasource, opts...)
}
