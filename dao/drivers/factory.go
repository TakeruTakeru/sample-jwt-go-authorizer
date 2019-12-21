package drivers

import "database/sql"

func NewClient(driver DRIVER, id string, pass string, ip string, database string) (client *sql.DB){
	switch driver {
	case MYSQL:
		client = NewMysqlClient(id, pass, ip, database).Connect()
	default:
		panic("Not implemented yet.")
	}

	return
}