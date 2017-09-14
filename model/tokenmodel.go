package model

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

// Token represents the model for tokens
type Token struct {
	ID      int64 `json:"id" db:"id"`
	ActorID int64 `json:"actorid" db:"actorid"`
	//AdUsername string          `json:"adusername" db:"adusername"`
	Token      string          `json:"token" db:"token"`
	ExpiresOn  null.Time       `json:"expireson" db:"expireson"`
	Blocked    bool            `json:"blocked" db:"blocked"`
	Created    time.Time       `json:"created" db:"created"`
	CreatedBy  NameDescription `json:"createdby" db:"createdby"`
	Modified   time.Time       `json:"modified" db:"modified"`
	ModifiedBy NameDescription `json:"modifiedby" db:"modifiedby"`
	RowVersion int64           `json:"rowversion" db:"rowversion"`
}
