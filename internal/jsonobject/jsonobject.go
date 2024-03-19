package jsonobject

import (
	"database/sql"
	"time"
)

//easyjson:json
type User struct {
	ID           int    `json:"-" db:"ID"`
	Login        string `json:"login" db:"login"`
	Password     string `json:"password" db:"-"`
	HashPassword []byte `json:"-" db:"password"`
}

//easyjson:json
type Orders []Order

//easyjson:json
type Order struct {
	Number       string         `json:"number" db:"number"`
	Status       string         `json:"status" db:"OrderStatus"`
	AccrualDB    sql.NullString `json:"-" db:"accrual"`
	Accrual      string         `json:"accrual,omitempty" db:"-"`
	UploadDateDB time.Time      `json:"-" db:"uploadDate"`
	UploadDate   string         `json:"uploaded_at" db:"-"`
}

type Balance struct {
	AccrualDB      sql.NullFloat64 `json:"-" db:"ACCRUAL"`
	AccrualCurrent float64         `json:"current" db:"-"`
	WithdrawnDB    sql.NullFloat64 `json:"-" db:"WITHDRAWN"`
	Withdrawn      float64         `json:"withdrawn" db:"-"`
}
