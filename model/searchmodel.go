package model

import ()

type SearchResponse struct {
	Result []NameDescription `db:"result" json:"result"`
}
