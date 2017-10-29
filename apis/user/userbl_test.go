package user_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gkewl/pulsecheck/apis/authuser"
	"github.com/gkewl/pulsecheck/apis/user"

	"github.com/gkewl/pulsecheck/errorhandler"
	tst "github.com/gkewl/pulsecheck/lib/testing"
	"github.com/gkewl/pulsecheck/protocol"
	"github.com/gkewl/pulsecheck/utilities"
)

func expectError(expected *errorhandler.NamedError, got error) {
	Expect(got).ToNot(BeNil())
	Expect(got.Error()).To(ContainSubstring(expected.String()))
}

var _ = Describe("User Business Logic Tests", func() {
	var reqCtx *protocol.TestRequestContext
	var logic = user.BLUser{}
	var t tst.T

	BeforeEach(func() {
		reqCtx = &protocol.TestRequestContext{Userid: 1}
		t = tst.T{ReqCtx: reqCtx}
		reqCtx.Companyid = 1
	})
	AfterEach(func() {
		reqCtx.Complete(false)
	})

	It("creates a good user", func() {
		usr := t.User(logic.Create(reqCtx, t.SampleUser()))
		check := t.GetUser(usr.ID)
		Expect(check.Email).To(Equal(usr.Email))

		a, err := authuser.BLAuthUser{}.Get(reqCtx, usr.ID)
		Expect(err).To(BeNil())
		Expect(a.Password).ToNot(BeNil())

	})

	It("gets all users", func() {
		t.User(logic.Create(reqCtx, t.SampleUser()))
		result := t.Users(logic.GetAll(reqCtx, 10))
		Expect(len(result)).To(BeNumerically(">", 0))
	})

	It("updates an existing user", func() {
		By("creating a user")
		usr := t.User(logic.Create(reqCtx, t.SampleUser()))
		usr.LastName = usr.LastName + "Test"

		By("updating the user")
		t.UpdateUser(usr)

		By("retrieving the user")
		check := t.GetUser(usr.ID)
		Expect(check.LastName).To(Equal(usr.LastName))
	})

	It("does not update a user that isn't there", func() {
		usr := t.User(logic.Create(reqCtx, t.SampleUser()))

		_ = t.UserErr(logic.Update(reqCtx, 0, usr))
	})

	It("deletes a user", func() {
		usr := t.User(logic.Create(reqCtx, t.SampleUser()))
		result, err := logic.Delete(reqCtx, usr.ID)
		Expect(err).To(BeNil())
		Expect(result).To(Equal("Ok"))
		_ = t.UserErr(logic.Get(reqCtx, usr.ID))
	})

	It("does not delete a user that isn't there", func() {
		_, err := logic.Delete(reqCtx, 0)
		Expect(err).ToNot(BeNil())
	})

	PIt("searches for users", func() {
		baseName := utilities.GenerateRandomString(10)
		for i := 0; i < 5; i++ {
			usr := t.SampleUser()
			usr.Email = baseName + utilities.GenerateRandomString(10)
			t.User(logic.Create(reqCtx, usr))
		}
		usrs := t.Users(logic.Search(reqCtx, baseName, 4))
		Expect(len(usrs)).To(Equal(4))
	})

})
