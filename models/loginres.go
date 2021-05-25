package models

type Login struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRes struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}
