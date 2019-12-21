package drivers

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlClient struct {
	user     string
	pass     string
	ip       string
	database string
}

func (mc *MysqlClient) Connect() (db *sql.DB) {
	endpoint := fmt.Sprintf("%s:%s@tcp(%s)/%s", mc.user, mc.pass, mc.ip, mc.database)
	db, err := sql.Open("mysql", endpoint)
	if err != nil {
		panic(err.Error())
	}
	return
}

func NewMysqlClient(user string, pass string, ip string, database string) (mc *MysqlClient) {
	mc = &MysqlClient{
		user:     user,
		pass:     pass,
		ip:       ip,
		database: database,
	}
	return
}
