package model

import (
	"gopkg.in/guregu/null.v3"
)

type AuthenticateUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Xid      string 
}

type TokenInfo struct {
	Token string `json:"token"`
	Exp   int64  `json:"exp"`
}

type UserDetail struct {
	UserID int64 `json:"UserID"`
	//ADUsername     null.String `json:"ADUsername"`
	VisualUsername null.String `json:"VisualUsername"`
	FirstName      null.String `json:"FirstName"`
	LastName       null.String `json:"LastName"`
	DisplayName    null.String `json:"DisplayName"`
	Department     string      `json:"Department"`
	CreateDate     string      `json:"CreateDate"`
}

type RoleInfo struct {
	UserID int       `db:"userid"`
	Email  null.String `db:"email"`
}

type UserCompany struct {
	UserID    int    `db:"userid" json:"userid"`
	CompanyID int    `db:"companyid" json:"companyid"`
	Role      string `db:"role" json:"role"`
}

type UserDetailResponse struct {
	UserDetails []UserDetail `json:"Items"`
}
type TokenAuthentication struct {
	Token string `json:"token"`
}

type Settings struct {
	JWTExpirationDelta int
}
