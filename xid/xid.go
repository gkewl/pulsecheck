package xid

import (
	"github.com/rs/xid"
)

// UniqueIdGenerator -
func UniqueIdGenerator() string {
	ugen := xid.New()

	return ugen.String()
}
