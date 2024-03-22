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

//easyjson:json
type Balance struct {
	AccrualDB      sql.NullFloat64 `json:"-" db:"accrual"`
	AccrualCurrent float64         `json:"current" db:"-"`
	WithdrawnDB    sql.NullFloat64 `json:"-" db:"withdrawn"`
	Withdrawn      float64         `json:"withdrawn" db:"-"`
}

//easyjson:json
type Withdrawls []Withdraw

//easyjson:json
type Withdraw struct {
	Order           string    `json:"order" db:"ordernum"`
	OrderNum        int       `json:"-" db:"-"`
	Sum             float64   `json:"sum" db:"pointsSum"`
	ProcessedDateDB time.Time `json:"-" db:"processedDate"`
	ProcessedDate   string    `json:"processed_at,omitempty" db:"-"`
}

//easyjson:json
type AccrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

//easyjson:json
type OrdersCalc []OrderCalc

//easyjson:json
type OrderCalc struct {
	Number     int            `db:"number"`
	CalcStatus sql.NullString `db:"calcstatus"`
	Accrual    float64        `db:"accrual"`
}
