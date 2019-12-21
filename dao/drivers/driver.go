package drivers

import "database/sql"

type Connector interface {
	Connect() *sql.DB
}

type MysqlConnector struct{}

func (mc *MysqlConnector) Connect() *sql.DB {
	return NewClient(MYSQL, "dev", "Password123>", "127.0.0.1:3306", "authsample")
}
