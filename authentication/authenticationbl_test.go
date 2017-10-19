package authentication_test

import (
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gkewl/pulsecheck/apis/authuser"
	"github.com/gkewl/pulsecheck/authentication"

	"github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
	"github.com/gkewl/pulsecheck/protocol"
)

var (
	reader io.Reader //Ignore this for now
)

var tokenInfo model.TokenInfo

// sampleUser constructs a User
func sampleuser() model.AuthenticateUser {
	return model.AuthenticateUser{
		Email:    "rgunari@gmail.com",
		Password: "test123",
	}
}

func expectError(expected *errorhandler.NamedError, got error) {
	Expect(got).ToNot(BeNil())
	Expect(got.Error()).To(ContainSubstring(expected.String()))
}

var _ = Describe("Authentication Business Logic Tests", func() {
	var reqCtx *protocol.TestRequestContext
	var logic = authentication.BLAuthentication{}
	var empToken = "testtokenthatfails"
	BeforeEach(func() {
		reqCtx = &protocol.TestRequestContext{Userid: 1, Username: "admin"}
		reqCtx.AppContext()
		authuser.TestingBizLogic = &authuser.MockBLAuthUser{}
	})
	AfterEach(func() {
		reqCtx.Complete(false)
	})

	It("Login Authentication", func() {

		output, err := logic.LoginUser(reqCtx, sampleuser())
		Expect(err).To(BeNil())
		Expect(len(output.Token)).To(BeNumerically(">", 0))
		Expect(output.Exp).To(BeNumerically(">", 0))
	})

	PIt("Login Authentication Failure", func() {
		user := sampleuser()
		user.Password = ""
		output, err := logic.LoginUser(reqCtx, user)
		Expect(err).NotTo(BeNil())
		Expect(len(output.Token)).To(Equal(0))
	})

	It("ValidateToken", func() {
		output, err := logic.LoginUser(reqCtx, sampleuser())
		Expect(err).To(BeNil())

		output1, err := logic.ValidateToken(reqCtx, output.Token)
		Expect(err).To(BeNil())
		Expect(output1).To(BeNumerically(">", int64(0)))
	})
	PIt("ValidateToken failure", func() {
		output, err := logic.ValidateToken(reqCtx, empToken)
		Expect(err).NotTo(BeNil())
		Expect(output).To(Equal(int64(0)))
	})

})
