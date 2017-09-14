package model

import (
	"gopkg.in/guregu/null.v3"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Xid      string
}

type Machine struct {
	Actorid     string `DB:"actorid" json:"actorid"`
	MachineName string `DB:"name" json:"machinename"`
	ActorType   string `DB:"type" json:"actortype"`
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
	Actorid int64       `db:"actorid"`
	Role    null.String `db:"role"`
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
