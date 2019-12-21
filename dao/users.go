package dao

import (
	"database/sql"
	"fmt"

	"github.com/TakeruTakeru/auth-sample/dao/drivers"
)

const (
	SELECT_ALL = "SELECT * FROM users"
	INSERT     = "INSERT INTO users(username, email, password) VALUES(?,?,?)"
)

type User struct {
	Id    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Pass  string `json:"pass"`
	DaoBase
}

func (u *User) String() string {
	return fmt.Sprintf("user[%d, %s, %s, %s]\n", u.Id, u.Name, u.Email, u.Pass)
}

func (u *User) getconnection() (db *sql.DB) {
	var conn drivers.Connector
	conn = &drivers.MysqlConnector{}
	return conn.Connect()
}

func (u *User) InsertSelf() (res bool, err error) {
	client := u.getconnection()
	stat, err := client.Prepare(INSERT)
	if err != nil {
		return
	}
	_, err = stat.Exec(u.Name, u.Email, u.Pass)
	if err != nil {
		return
	}
	res = true
	return
}

func (u *User) executeQuery(q string, result *[]User) (err error) {
	var user User
	client := u.getconnection()
	results, err := client.Query(q)
	buf := make([]User, 0, 10)
	if err != nil {
		return
	}
	defer client.Close()
	defer results.Close()

	for results.Next() {
		err = results.Scan(&user.Id, &user.Name, &user.Email, &user.Pass)
		if err != nil {
			return
		}
		fmt.Println(len(buf), cap(buf))
		buf = append(buf, user)
	}
	*result = buf
	return
}

func (u *User) GetAllUsers(result *[]User) (err error) {
	return u.executeQuery(SELECT_ALL, result)
}
