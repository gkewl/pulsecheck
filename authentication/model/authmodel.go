package authmodel

import (

)

type User struct {
	Actorid     string `json:"actorid"` 
	Username string `json:"username"`
	Password string `json:"password"`
	Type string `json:"type"`
	}