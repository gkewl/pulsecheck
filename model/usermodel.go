package model

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

// User represents the model for users
type User struct {
	ID          int         `json:"id" db:"id"`
	Email       string      `json:"email" db:"email"`
	FirstName   string      `json:"firstname" db:"firstname"`
	MiddleName  null.String `json:"middlename" db:"middlename"`
	LastName    string      `json:"lastname" db:"lastname"`
	CompanyID   int         `json:"companyid" db:"companyid"`
	CompanyName string      `json:"companyname" db:"companyname"`
	IsActive    bool        `json:"isactive" db:"isactive"`
	CreatedBy   string      `json:"createdby" db:"createdby"`
	Created     time.Time   `json:"created" db:"created"`
	ModifiedBy  string      `json:"modifiedby" db:"modifiedby"`
	Modified    time.Time   `json:"modified" db:"modified"`
}

type RegisterUser struct {
	User
	Password string `json:"password"`
}
