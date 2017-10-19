package authentication

import (
	"fmt"

	"github.com/gkewl/pulsecheck/logger"
	//"encoding/json"

	"github.com/gkewl/pulsecheck/apis/authuser"
	"github.com/gkewl/pulsecheck/common"
	//	"github.com/gkewl/pulsecheck/config"
	//	"github.com/gkewl/pulsecheck/constant"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
	//"net/http"
)

// BizLogic is the interface for all authentication business logic methods
type BizLogic interface {
	LoginUser(common.RequestContext, model.AuthenticateUser) (model.TokenInfo, error)
	ValidateToken(common.RequestContext, string) (int64, error)
}

// BLAuthentication implements the authentication.BizLogic interface
type BLAuthentication struct {
}

const (
	tokenDuration = 24
	expireOffset  = 3600
)

// LoginUser -
func (bl BLAuthentication) LoginUser(reqCtx common.RequestContext, requestUser model.AuthenticateUser) (model.TokenInfo, error) {
	authBackend := InitJWTAuthenticationBackend()
	if bl.Authenticate(reqCtx, requestUser) {

		info, err := dlGetUserInfoForToken(reqCtx, requestUser.Email)
		if err != nil {
			return model.TokenInfo{}, eh.NewError(eh.ErrAuthUserNameNotAvailable, "Username:%s  Not available in pulse system, Error:%s  : ", requestUser.Email, err.Error())
		}

		tk, err := authBackend.GenerateToken(requestUser.Email, info, "USER")
		if err != nil {
			return model.TokenInfo{}, eh.NewError(eh.ErrAuthGenerateToken, "Generate token error. Username:%s   Error:%s  : ", requestUser.Email, err.Error())
		}

		return tk, nil

	}
	return model.TokenInfo{}, eh.NewError(eh.ErrAuthWrongUseridOrPwd, "Username:%s", requestUser.Email)
}

func (bl BLAuthentication) Authenticate(reqCtx common.RequestContext, usr model.AuthenticateUser) bool {
	info, err := dlGetUserInfoForToken(reqCtx, usr.Email)
	if err != nil {
		logger.LogError(fmt.Sprintf("Username:%s  Not available in pulse system, Error:%s  : ", usr.Email, err.Error()), reqCtx.Xid())
		return false
	}

	result, err := authuser.GetInterface().Authenticate(reqCtx, info.UserID, usr.Password)
	if err != nil {
		logger.LogError(fmt.Sprintf("Username:%s  Failed to authenticate in pulse system, Error:%s  : ", usr.Email, err.Error()), reqCtx.Xid())
		return false
	}
	return result
}

// AuditToken -
// func AuditToken(reqCtx common.RequestContext, tk model.TokenInfo, actorid int64, username string) error {

// 	auditTk := model.Token{ActorID: actorid, Token: tk.Token, CreatedBy: model.NameDescription{Id: reqCtx.UserID()},
// 		ModifiedBy: model.NameDescription{Id: reqCtx.UserID()}}

// 	if tk.Exp != 0 {
// 		t, err := utilities.ParseStringtoTime(time.Unix(tk.Exp, 0).Format(time.RFC3339))
// 		if err == nil {
// 			auditTk.ExpiresOn = null.TimeFrom(t)
// 		}

// 	}

// 	return nil
// }

// MachineToken -
// func (bl BLAuthentication) MachineToken(reqCtx common.RequestContext, name string) (model.TokenInfo, error) {
// 	//Store machine information to DB
// 	di := actor.BLActor{}
// 	actor, err := di.Get(reqCtx, name)
// 	if err != nil {
// 		if strings.Contains(err.Error(), eh.ErrActorDataNotFound.String()) == true {
// 			err = eh.AddDetail(err, "Machine name not found: %s", name)
// 		}
// 		return model.TokenInfo{}, err
// 	}

// 	authBackend := InitJWTAuthenticationBackend()
// 	tk, err := authBackend.GenerateToken(actor.Name, actor.Id, actor.Type, actor.Role)
// 	if err != nil {
// 		return model.TokenInfo{}, eh.NewError(eh.ErrAuthGenerateToken, "Token Error: "+err.Error())
// 	}

// 	//audit token
// 	err = AuditToken(reqCtx, tk, actor.Id, actor.Name)
// 	if err != nil {
// 		//ignore this error
// 	}

// 	return tk, nil

// }

// ValidateToken -
func (bl BLAuthentication) ValidateToken(reqCtx common.RequestContext, token string) (int64, error) {
	return GetValidateToken(token)

}
