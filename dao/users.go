package dao

import (
	"database/sql"
	"fmt"

	"github.com/TakeruTakeru/auth-sample/dao/drivers"
)

const (
	SELECT_ALL                   = "SELECT * FROM users"
	SELECT_EMAIL_PASS_WITH_EMAIL = "SELECT email, password FROM users WHERE email='%s'"
	INSERT                       = "INSERT INTO users(username, email, password) VALUES(?,?,?)"
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

func (u *User) executeQuery(q string, result *[]User, setter func(result *sql.Rows, user *User) error) (err error) {
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
		err = setter(results, &user)
		// err = results.Scan(&user.Id, &user.Name, &user.Email, &user.Pass)
		if err != nil {
			return
		}
		buf = append(buf, user)
	}
	*result = buf
	return
}

func (u *User) GetEmailAndPassByEmail(result *[]User) (err error) {
	return u.executeQuery(fmt.Sprintf(SELECT_EMAIL_PASS_WITH_EMAIL, u.Email), result, func(results *sql.Rows, user *User) error {
		return results.Scan(&user.Email, &user.Pass)
	})
}

func (u *User) GetAllUsers(result *[]User) (err error) {
	return u.executeQuery(SELECT_ALL, result, func(results *sql.Rows, user *User) error {
		return results.Scan(&user.Id, &user.Name, &user.Email, &user.Pass)
	})
}
