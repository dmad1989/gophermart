package jsonobject

//go:generate easyjson -all jsonobject.go
type User struct {
	ID           int    `json:"-" db:"ID"`
	Login        string `json:"login" db:"login"`
	Password     string `json:"password"`
	HashPassword []byte `json:"-" db:"password"`
}
