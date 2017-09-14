package token

import (
	"crypto/md5"
	"fmt"
	"github.com/gkewl/pulsecheck/common"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
	"io"
)

// BizLogic is the interface for all token business logic methods
type BizLogic interface {
	Create(common.RequestContext, model.Token) (model.Token, error)
	Get(common.RequestContext, int64) (model.Token, error)
	GetByToken(common.RequestContext, string) (model.Token, error)
	GetAll(common.RequestContext, int64) ([]model.Token, error)
	Update(common.RequestContext, int64, model.Token) (model.Token, error)
	Delete(common.RequestContext, int64) (string, error)
	Search(common.RequestContext, string, int64) ([]model.Token, error)
}

// BLToken implements the token.BizLogic interface
type BLToken struct {
}

//getMd5Token generates md5 hash for gievn token
func (bl BLToken) getMd5Token(token string) string {
	h := md5.New()
	io.WriteString(h, token)
	return fmt.Sprintf("%x", h.Sum(nil))

}

// Create will insert a new token into the db
func (bl BLToken) Create(reqCtx common.RequestContext, tk model.Token) (model.Token, error) {

	tk.Token = bl.getMd5Token(tk.Token)
	tk, err := dlCreate(reqCtx, tk)
	if err != nil {
		return model.Token{}, eh.NewError(eh.ErrTokenInsert, "DB Error: "+err.Error())
	}
	return tk, err
}

// Get returns a single token by primary key
func (bl BLToken) Get(reqCtx common.RequestContext, id int64) (tk model.Token, err error) {
	tk, err = dlGet(reqCtx, id)

	if err != nil || tk.ID == 0 {
		return model.Token{}, eh.NewErrorNotFound(eh.ErrTokenDataNotFound, err, `Token not found: id %d`, id)
	}
	return
}

// Get returns a single token by primary key
func (bl BLToken) GetByToken(reqCtx common.RequestContext, tokenString string) (tk model.Token, err error) {
	hash := bl.getMd5Token(tokenString)
	tk, err = dlGetByToken(reqCtx, hash)

	if err != nil || tk.ID == 0 {
		return model.Token{}, eh.NewErrorNotFound(eh.ErrTokenDataNotFound, err, `Token not found: tken string %s`, tokenString)
	}
	return
}

// GetAll will return all tokens
func (bl BLToken) GetAll(reqCtx common.RequestContext, limit int64) (tks []model.Token, err error) {
	tks, err = dlGetAll(reqCtx, limit)
	if err != nil {
		return []model.Token{}, eh.NewError(eh.ErrTokenDataNotFound, "DB Error: "+err.Error())
	}
	return
}

// Update updates a single token
func (bl BLToken) Update(reqCtx common.RequestContext, id int64, tk model.Token) (model.Token, error) {
	// todo: add validation here

	result, err := dlUpdate(reqCtx, id, tk)
	if err != nil {
		return model.Token{}, eh.NewError(eh.ErrTokenUpdate, "DB Error: "+err.Error())
	}
	return result, err
}

// Delete marks a single token inactive
func (bl BLToken) Delete(reqCtx common.RequestContext, id int64) (string, error) {
	// todo: add validation here
	err := dlDelete(reqCtx, id)

	if err != nil {
		return "failed", eh.NewError(eh.ErrTokenDelete, "DB Error: "+err.Error())
	}
	return "ok", err
}

// Search finds tokens matching the term
func (bl BLToken) Search(reqCtx common.RequestContext, term string, limit int64) (tks []model.Token, err error) {
	tks, err = dlSearch(reqCtx, term, limit)
	if err != nil {
		return []model.Token{}, eh.NewError(eh.ErrTokenDataNotFound, "DB Error: "+err.Error())
	}
	return
}
