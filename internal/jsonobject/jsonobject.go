package jsonobject

//go:generate easyjson -all jsonobject.go
type User struct {
	Login        string `json:"login"`
	Password     string `json:"password"`
	HashPassword string `json:"-"`
}
