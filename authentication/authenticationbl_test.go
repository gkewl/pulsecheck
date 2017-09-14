package authentication_test

import (
	"github.com/gkewl/pulsecheck/apis/actor"
	"github.com/gkewl/pulsecheck/apis/token"
	"github.com/gkewl/pulsecheck/authentication"
	"github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
	"github.com/gkewl/pulsecheck/protocol"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
)

var (
	reader io.Reader //Ignore this for now
)

var tokenInfo model.TokenInfo

// sampleUser constructs a User
func sampleuser() model.User {
	return model.User{
		Username: "test",
		Password: "test1",
	}
}

func expectError(expected *errorhandler.NamedError, got error) {
	Expect(got).ToNot(BeNil())
	Expect(got.Error()).To(ContainSubstring(expected.String()))
}

var _ = Describe("Authentication Business Logic Tests", func() {
	var reqCtx *protocol.TestRequestContext
	var logic = authentication.BLAuthentication{}
	var logicactor = actor.BLActor{}
	var name = "test"
	var empToken = "testtokenthatfails"
	BeforeEach(func() {
		reqCtx = &protocol.TestRequestContext{Userid: 1, Username: "admin"}
	})
	AfterEach(func() {
		reqCtx.Complete(false)
	})

	It("Login Authentication", func() {
		check, err2 := logicactor.Get(reqCtx, name)
		if err2 == nil && check.Id > 0 {
			//Delete actor if exist in DB.
			result, err2 := logicactor.Delete(reqCtx, check.Id)
			Expect(err2).To(BeNil())
			Expect(result).To(Equal("ok"))
		}
		output, err := logic.LoginUser(reqCtx, sampleuser())
		Expect(err).To(BeNil())
		Expect(len(output.Token)).To(BeNumerically(">", 0))
		Expect(output.Exp).To(BeNumerically(">", 0))
		check, err2 = logicactor.Get(reqCtx, name)
		Expect(err2).To(BeNil())
		Expect(check.Name).To(Equal(name))
		Expect(check.LastLoginTime.Valid).To(BeTrue())

		tk, err := token.BLToken{}.GetByToken(reqCtx, output.Token)
		Expect(err).To(BeNil())
		Expect(tk.AdUsername).To(Equal(name))

	})

	It("Login Authentication Failure", func() {
		user := sampleuser()
		user.Password = ""
		output, err := logic.LoginUser(reqCtx, user)
		Expect(err).NotTo(BeNil())
		Expect(len(output.Token)).To(Equal(0))
	})

	It("MachineToken Authentication", func() {
		output, err := logic.MachineToken(reqCtx, glTestStore.MachineActor.Name)
		Expect(err).To(BeNil())
		Expect(len(output.Token)).To(BeNumerically(">", 0))
		Expect(output.Exp).To(Equal(int64(0)))

		tk, err := token.BLToken{}.GetByToken(reqCtx, output.Token)
		Expect(err).To(BeNil())
		Expect(tk.AdUsername).To(Equal(glTestStore.MachineActor.Name))

	})

	It("MachineToken Fail Authentication", func() {
		output, err := logic.MachineToken(reqCtx, "")
		Expect(err).NotTo(BeNil())
		Expect(len(output.Token)).To(Equal(0))
	})

	It("ValidateToken", func() {
		output, err := logic.ValidateToken(reqCtx, tokenInfo.Token)
		Expect(err).To(BeNil())
		Expect(output).To(BeNumerically(">", int64(0)))
	})
	It("ValidateToken failure", func() {
		output, err := logic.ValidateToken(reqCtx, empToken)
		Expect(err).NotTo(BeNil())
		Expect(output).To(Equal(int64(0)))
	})

})
