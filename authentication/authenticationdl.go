package authentication

import (
	"github.com/jmoiron/sqlx"

	"github.com/gkewl/pulsecheck/common"
	//"github.com/gkewl/pulsecheck/constant"
	"github.com/gkewl/pulsecheck/model"
)

const (
	getQuery = `select a.id  actorid, role
									  from Actor a
									   `
)

// DLGet retreives the specified UserRole
func DLGetUserRole(reqCtx common.RequestContext, name string) (userRole model.RoleInfo, err error) {
	return dlGetUserRoleTx(reqCtx.Tx(), name)
}

func dlGetUserRoleTx(tx *sqlx.Tx, name string) (userRole model.RoleInfo, err error) {
	err = tx.Get(&userRole, getQuery+` where a.name=? `, name)
	return
}

func DLGetUserRoleDB(db *sqlx.DB, name string) (userRole model.RoleInfo, err error) {
	tx := db.MustBegin()
	userRole, err = dlGetUserRoleTx(tx, name)
	tx.Rollback()
	return
}
