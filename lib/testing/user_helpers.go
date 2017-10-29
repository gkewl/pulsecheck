package testing

import (
	//. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	null "gopkg.in/guregu/null.v3"

	"github.com/gkewl/pulsecheck/apis/user"
	"github.com/gkewl/pulsecheck/model"
)

// sampleUser constructs a user
func (t *T) SampleUser() (usr model.RegisterUser) {
	return model.RegisterUser{
		User: model.User{
			Email:      "foo@email.com",
			FirstName:  "foo",
			MiddleName: null.StringFrom("foo"),
			LastName:   "foo",
			CompanyID:  1,
			IsActive:   true,
			CreatedBy:  "admin",
			ModifiedBy: "admin",
		},
		Password: "fooPwd",
	}
}

// User expects a user and an error and verifies the error is nil
// and returns the user
func (t *T) User(usr model.User, e error) model.User {
	Expect(e).To(BeNil(), callers())
	return usr
}

// GetUser fetches a user using an ID and verifies no error
func (t *T) GetUser(ID int) model.User {
	Expect(t.ReqCtx).ToNot(BeNil(), callers())
	usr, e := user.BLUser{}.Get(t.ReqCtx, ID)
	Expect(e).To(BeNil(), callers())
	return usr
}

// ReGetUser expects a user and an error and verifies the error is nil
// and re-gets the user and returns it
func (t *T) ReGetUser(usr model.User, e error) model.User {
	Expect(e).To(BeNil(), callers())
	Expect(t.ReqCtx).ToNot(BeNil(), callers())
	usr, e = user.BLUser{}.Get(t.ReqCtx, usr.ID)
	Expect(e).To(BeNil(), callers())
	return usr
}

// UpdateUser expects a user struct and updates it, verifying no error
func (t *T) UpdateUser(usr model.User) model.User {
	Expect(t.ReqCtx).ToNot(BeNil(), callers())
	usr, e := user.BLUser{}.Update(t.ReqCtx, usr.ID, usr)
	Expect(e).To(BeNil(), callers())
	return usr
}

// Users expects a user slice and an error and verifies the error is nil
// and returns the users
func (t *T) Users(usrs []model.User, e error) []model.User {
	Expect(e).To(BeNil(), callers())
	return usrs
}

// UserErr expects a user and an error and verifies the error is
// not nil and returns the error
func (t *T) UserErr(usr model.User, e error) error {
	Expect(e).ToNot(BeNil(), callers())
	return e
}
