package model

import ()

type Ping struct {
	Status    string `json:"status"`
	DB        string `json:"db"`
	Version   string `json:"version"`
	BuildTime string `json:"buildtime"`
	Githash   string `json:"githash"`
}
