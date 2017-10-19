package authuser_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gkewl/pulsecheck/apis/authuser"

	"github.com/gkewl/pulsecheck/errorhandler"
	tst "github.com/gkewl/pulsecheck/lib/testing"
	"github.com/gkewl/pulsecheck/model"
	"github.com/gkewl/pulsecheck/protocol"
)

// sampleAuthuser constructs a authuser
func sampleAuthuser() (au model.AuthUser) {
	return model.AuthUser{
		UserID:   42,
		Password: "MyPassword",
	}
}

func expectError(expected *errorhandler.NamedError, got error) {
	Expect(got).ToNot(BeNil())
	Expect(got.Error()).To(ContainSubstring(expected.String()))
}

var _ = Describe("Authuser Business Logic Tests", func() {
	var reqCtx *protocol.TestRequestContext
	var logic = authuser.BLAuthUser{}
	var t tst.T

	BeforeEach(func() {
		reqCtx = &protocol.TestRequestContext{Userid: 1}
		t = tst.T{ReqCtx: reqCtx}
	})
	AfterEach(func() {
		reqCtx.Complete(false)
	})

	It("creates a good authuser", func() {
		au, err := logic.CreateOrUpdate(reqCtx, sampleAuthuser())
		Expect(err).To(BeNil())
		check, err := logic.Get(reqCtx, au.UserID)
		Expect(err).To(BeNil())
		Expect(check.UserID).To(Equal(au.UserID))

	})

	It("Authenticate user", func() {
		input := sampleAuthuser()
		usr, err := logic.CreateOrUpdate(reqCtx, input)
		Expect(err).To(BeNil())

		check, err := logic.Authenticate(reqCtx, usr.UserID, input.Password)
		Expect(check).To(BeTrue())
	})

	It("updates an existing authuser", func() {
		By("creating a authuser")
		input := sampleAuthuser()
		au, err := logic.CreateOrUpdate(reqCtx, input)
		Expect(err).To(BeNil())

		By("updating the authuser")
		au.Password = "NewFoo"
		au, err = logic.CreateOrUpdate(reqCtx, au)
		Expect(err).To(BeNil())

		input.Password = "NewFoo"
		check, err := logic.Authenticate(reqCtx, au.UserID, input.Password)
		Expect(check).To(BeTrue())

	})

	It("deletes a authuser", func() {
		au, err := logic.CreateOrUpdate(reqCtx, sampleAuthuser())
		Expect(err).To(BeNil())

		result, err := logic.Delete(reqCtx, au.UserID)
		Expect(err).To(BeNil())
		Expect(result).To(Equal("Ok"))

		check, err := logic.Get(reqCtx, au.UserID)
		Expect(err).To(BeNil())
		Expect(check.IsActive).To(BeFalse())
	})

	It("does not delete a authuser that isn't there", func() {
		_, err := logic.Delete(reqCtx, 0)
		Expect(err).ToNot(BeNil())
	})

})
