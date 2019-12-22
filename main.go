package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"strings"

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
	var res []dao.User
	json.NewDecoder(r.Body).Decode(&user)
	noneHashedPass := user.Pass
	err := user.GetEmailAndPassByEmail(&res)
	if err != nil {
		fmt.Fprintf(w, "1: %s", err.Error())
		return
	}
	resultNum := len(res)
	if resultNum > 1 {
		fmt.Fprintf(w, "Problem occured.")
		return
	}
	if resultNum == 0 {
		fmt.Fprintf(w, "No users matched.")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(res[0].Pass), []byte(noneHashedPass))
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	token, err := createToken(&user)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprint(w, token)
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	bearerToken := strings.Split(authHeader, " ")
	fmt.Println(bearerToken)
	if len(bearerToken) == 2 {
		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("エラーです")
			}
			return []byte("hogehoge123"), nil
		})
		if err != nil {
			fmt.Fprintf(w, "%+v", errors.WithStack(err))
		}

		if token.Valid {
			fmt.Fprintf(w, "ok")
		} else {
			fmt.Fprintf(w, err.Error())
			return
		}
	} else {
		fmt.Fprintf(w, "Invalid request")
		return
	}
}

func main() {
	http.HandleFunc("/", handler) // ハンドラを登録してウェブページを表示させる
	http.HandleFunc("/mysql/test", mysqlHandler)
	http.HandleFunc("/signup", singupHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/verify", verifyHandler)
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
