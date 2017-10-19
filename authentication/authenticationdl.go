package authentication

import (
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/model"
	//"github.com/gkewl/pulsecheck/constant"
)

// dlGetUserInfoForToken -
func dlGetUserInfoForToken(reqCtx common.RequestContext, email string) (info model.UserCompany, err error) {
	info = model.UserCompany{}
	err = reqCtx.Tx().Get(&info, `select id as userid,companyid, role from user  where email=? `, email)
	return
}
