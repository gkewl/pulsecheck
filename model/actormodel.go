package model

import (
	"gopkg.in/guregu/null.v3"
)

type Actor struct {
	BasicModel
	Name        string      `db:"name" json:"name"`
	Description null.String `db:"description" json:"description"`
	//ADUsername    null.String             `db:"adusername" json:"adusername"`
	EmailAddress  null.String             `db:"email" json:"email"`
	Type          string                  `db:"type" json:"actortype" `
	IPAddress     null.String             `db:"ipconfig" json:"ipaddress"`
	MACAddress    null.String             `db:"macaddress" json:"macaddress"`
	Manager       NullableNameDescription `db:"manager" json:"manager"`
	Role          string                  `db:"role" json:"role"`
	LastLoginTime null.Time               `db:"lastlogintime" json:"lastlogintime"`
}
type Actors []Actor

type ActorSearchResponse struct {
	Id   int64  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	//ADUsername   null.String `db:"adusername" json:"adusername"`
	EmailAddress null.String `db:"email" json:"email"`
	Description  null.String `db:"description" json:"description"`
}

type CustomActor struct {
	CreatorId           int64       `db:"creatorid" json:"creatorid"`
	CreatorName         string      `db:"creatorname" json:"creatorname"`
	CreatorDescription  null.String `db:"creatordescription" json:"creatordescription"`
	ModifierId          int64       `db:"modifierid" json:"modifierid"`
	ModifierName        string      `db:"modifiername" json:"modifiername"`
	ModifierDescription null.String `db:"modifierdescription" json:"modifierdescription"`
}
