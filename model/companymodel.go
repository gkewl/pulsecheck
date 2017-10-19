package model

import (
	"time"
)

// Company represents the model for companys
type Company struct {
	ID         int       `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	IsActive   bool      `json:"isactive" db:"isactive"`
	CreatedBy  string    `json:"createdby" db:"createdby"`
	Created    time.Time `json:"created" db:"created"`
	ModifiedBy string    `json:"modifiedby" db:"modifiedby"`
	Modified   time.Time `json:"modified" db:"modified"`
}
