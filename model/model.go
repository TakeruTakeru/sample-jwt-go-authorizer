package model

type User struct {
	Id    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Pass  string `json:"pass"`
}

type JWT struct {
	Token string `json:"token"`
}
