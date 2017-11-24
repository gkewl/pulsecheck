package model

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

// EmployeeStatus represents the model for employeestatuss
type EmployeeStatus struct {
	ID             int64     `json:"id" db:"id"`
	EmployeeID     int64     `json:"employeeid" db:"employeeid"`
	Consider       bool      `json:"consider" db:"consider"`
	OIG            bool      `json:"oig" db:"oig"`
	OIGLastSearch  null.Time `json:"oiglastsearch" db:"oiglastsearch"`
	Sam            bool      `json:"sam" db:"sam"`
	SamLastSearch  null.Time `json:"samlastsearch" db:"samlastsearch"`
	Ofac           bool      `json:"ofac" db:"ofac"`
	OfacLastSearch null.Time `json:"ofaclastsearch" db:"ofaclastsearch"`
	IsActive       bool      `json:"isactive" db:"isactive"`
	CreatedBy      string    `json:"createdby" db:"createdby"`
	Created        time.Time `json:"created" db:"created"`
	ModifiedBy     string    `json:"modifiedby" db:"modifiedby"`
	Modified       time.Time `json:"modified" db:"modified"`
}
