package model

import (
	"time"
)

type AuthUser struct {
	ID         int64     `db:"id" json:"-"`
	UserID     int       `db:"userid" json:"userid"`
	Password   string    `db:"password"`
	IsActive   bool      `db:"isactive"`
	CreatedBy  string    `json:"createdby" db:"createdby"`
	Created    time.Time `json:"created" db:"created"`
	ModifiedBy string    `json:"modifiedby" db:"modifiedby"`
	Modified   time.Time `json:"modified" db:"modified"`
}
