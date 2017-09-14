package actor

import (
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/model"
)

// MockBL will store inbound parameters for verification and returns
// the provided error
type MockBL struct {
	Name string
	ID   int64
	Term string
	Type string
	Err  error
}

// Create stubs create method
func (mb *MockBL) Create(rc common.RequestContext, s model.Actor) (model.Actor, error) {
	return s, mb.Err
}

// Get stubs get method
func (mb *MockBL) Get(rc common.RequestContext, name string) (model.Actor, error) {
	mb.Name = name
	return model.Actor{}, mb.Err
}

// GetAll stubs getting all steps method
func (mb *MockBL) GetAll(rc common.RequestContext, usertype string) ([]model.Actor, error) {
	mb.Type = usertype
	return []model.Actor{}, mb.Err
}

// Update stubs update method
func (mb *MockBL) Update(rc common.RequestContext, ID int64, s model.Actor) (model.Actor, error) {
	mb.ID = ID
	return s, mb.Err
}

// UpdateLastLoginTime stubs UpdateLastLoginTime method
func (mb *MockBL) UpdateLastLoginTime(reqCtx common.RequestContext, ID int64) error {
	mb.ID = ID
	return mb.Err
}

// Delete stubs delete method
func (mb *MockBL) Delete(rc common.RequestContext, ID int64) (string, error) {
	mb.ID = ID
	return "ok", mb.Err
}

// Search stubs search method
func (mb *MockBL) Search(rc common.RequestContext, usertype string, term string) ([]model.ActorSearchResponse, error) {
	mb.Term = term
	mb.Type = usertype
	return []model.ActorSearchResponse{}, mb.Err
}

//GetDeprecated is needed for the interface.
func (mb *MockBL) GetDeprecated(name string, ctx *common.AppContext) (model.Actor, error) {
	mb.Name = name
	return model.Actor{}, mb.Err
}
