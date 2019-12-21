package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TakeruTakeru/auth-sample/dao"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%+v\n", r)
}

func mysqlHandler(w http.ResponseWriter, r *http.Request) {
	var users []dao.User
	userTable := &dao.User{}
	err := userTable.GetAllUsers(&users)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "%v", users)
}

func singupHandler(w http.ResponseWriter, r *http.Request) {
	var user dao.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Printf("invalid req")
		return
	}
	if user.Name == "" {
		fmt.Fprintln(w, "nameがないよ")
		return
	}

	if user.Email == "" {
		fmt.Fprintln(w, "Emailがないよ")
		return
	}
	if user.Pass == "" {
		fmt.Fprintln(w, "passがないよ")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Pass), 10)
	if err != nil {
		fmt.Printf("failed to generate pass")
		return
	}
	user.Pass = string(hash)
	dao.Insert(&user)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user dao.User
	json.NewDecoder(r.Body).Decode(&user)
	token, err := createToken(&user)
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Fprintln(w, token)
}

func main() {
	http.HandleFunc("/", handler) // ハンドラを登録してウェブページを表示させる
	http.HandleFunc("/mysql/test", mysqlHandler)
	http.HandleFunc("/signup", singupHandler)
	http.HandleFunc("/login", loginHandler)
	http.ListenAndServe(":8080", nil)
}

func createToken(user *dao.User) (string, error) {
	secret := "hogehoge123"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"iss":   "__init__", // JWT の発行者が入る(文字列(__init__)は任意)
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
