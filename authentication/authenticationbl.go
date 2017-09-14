package authentication

import (
	//"encoding/json"
	"github.com/gkewl/pulsecheck/apis/actor"
	"github.com/gkewl/pulsecheck/apis/token"
	"github.com/gkewl/pulsecheck/common"
	//	"github.com/gkewl/pulsecheck/config"
	//	"github.com/gkewl/pulsecheck/constant"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
	"github.com/gkewl/pulsecheck/utilities"
	"gopkg.in/guregu/null.v3"
	//"net/http"
	"strings"
	"time"
)

// BizLogic is the interface for all authentication business logic methods
type BizLogic interface {
	LoginUser(common.RequestContext, model.User) (model.TokenInfo, error)
	MachineToken(common.RequestContext, string) (model.TokenInfo, error)
	ValidateToken(common.RequestContext, string) (int64, error)
	//WDAuthToken(common.RequestContext) (model.WDAuthToken, error)
}

// BLAuthentication implements the authentication.BizLogic interface
type BLAuthentication struct {
}

const (
	tokenDuration = 72
	expireOffset  = 3600
)

// LoginUser
func (bl BLAuthentication) LoginUser(reqCtx common.RequestContext, requestUser model.User) (model.TokenInfo, error) {
	authBackend := InitJWTAuthenticationBackend()
	if authBackend.Authenticate(requestUser) {
		//Get user id for loggedin  user from db
		type roleInfo struct {
			Actorid int64       `db:"actorid"`
			Role    null.String `db:"role"`
		}
		userRole, err := DLGetUserRole(reqCtx, requestUser.Username)
		if err != nil {
			if eh.HasNoRowsError(err) {

				//Check again for userRole
				userRole, err = DLGetUserRole(reqCtx, requestUser.Username)
				if err != nil {
					return model.TokenInfo{}, eh.NewError(eh.ErrAuthUserIdNotAvailbleInDB, "")
				}

			} else {
				return model.TokenInfo{}, eh.NewError(eh.ErrAuthUserIdNotAvailbleInDB, "")
			}
		}

		if userRole.Role.Valid == false {
			userRole.Role = null.NewString("Guest", true)
		}

		tk, err := authBackend.GenerateToken(requestUser.Username, userRole.Actorid, "USER", userRole.Role.String)
		if err != nil {
			return model.TokenInfo{}, eh.NewError(eh.ErrAuthGenerateToken, "Generate token error. Username:%s  Role:%s, Error:%s  : ", requestUser.Username, userRole.Role, err.Error())
		}
		reqCtx.SetUserId(userRole.Actorid)
		//update last successful login
		err = actor.BLActor{}.UpdateLastLoginTime(reqCtx, userRole.Actorid)
		if err != nil {
			//ignore this error
		}

		//audit token
		err = AuditToken(reqCtx, tk, userRole.Actorid, requestUser.Username)
		if err != nil {
			//ignore this error
		}

		return tk, nil

	}
	return model.TokenInfo{}, eh.NewError(eh.ErrAuthWrongUseridOrPwd, "Username:%s", requestUser.Username)
}

func AuditToken(reqCtx common.RequestContext, tk model.TokenInfo, actorid int64, username string) error {

	auditTk := model.Token{ActorID: actorid, Token: tk.Token, CreatedBy: model.NameDescription{Id: reqCtx.UserID()},
		ModifiedBy: model.NameDescription{Id: reqCtx.UserID()}}

	if tk.Exp != 0 {
		t, err := utilities.ParseStringtoTime(time.Unix(tk.Exp, 0).Format(time.RFC3339))
		if err == nil {
			auditTk.ExpiresOn = null.TimeFrom(t)
		}

	}

	//store the token in token table
	auditTk, err := token.BLToken{}.Create(reqCtx, auditTk)
	if err != nil {
		return eh.NewError(eh.ErrAuthGenerateToken, "Failed to update token table. Error: "+err.Error())
	}
	return nil
}

//MachineToken
func (bl BLAuthentication) MachineToken(reqCtx common.RequestContext, name string) (model.TokenInfo, error) {
	//Store machine information to DB
	di := actor.BLActor{}
	actor, err := di.Get(reqCtx, name)
	if err != nil {
		if strings.Contains(err.Error(), eh.ErrActorDataNotFound.String()) == true {
			err = eh.AddDetail(err, "Machine name not found: %s", name)
		}
		return model.TokenInfo{}, err
	}

	authBackend := InitJWTAuthenticationBackend()
	tk, err := authBackend.GenerateToken(actor.Name, actor.Id, actor.Type, actor.Role)
	if err != nil {
		return model.TokenInfo{}, eh.NewError(eh.ErrAuthGenerateToken, "Token Error: "+err.Error())
	}

	//audit token
	err = AuditToken(reqCtx, tk, actor.Id, actor.Name)
	if err != nil {
		//ignore this error
	}

	return tk, nil

}

func (bl BLAuthentication) ValidateToken(reqCtx common.RequestContext, token string) (int64, error) {
	return GetValidateToken(token)

}
