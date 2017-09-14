package token_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gopkg.in/guregu/null.v3"

	"github.com/gkewl/pulsecheck/apis/token"

	"github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
	"github.com/gkewl/pulsecheck/protocol"
	"github.com/gkewl/pulsecheck/utilities"
)

// sampleToken constructs a token
func sampleToken() (tk model.Token) {
	return model.Token{
		ActorID:    glTestStore.Actor.Id,
		AdUsername: glTestStore.Actor.Name,
		Token:      "foo123456",
		ExpiresOn:  null.TimeFrom(time.Now()),
		Blocked:    false,
		Modified:   time.Now(),
	}
}

func expectError(expected *errorhandler.NamedError, got error) {
	Expect(got).ToNot(BeNil())
	Expect(got.Error()).To(ContainSubstring(expected.String()))
}

var _ = Describe("Token Business Logic Tests", func() {
	var reqCtx *protocol.TestRequestContext
	var logic = token.BLToken{}

	BeforeEach(func() {
		reqCtx = &protocol.TestRequestContext{Userid: glTestStore.Actor.Id}
	})
	AfterEach(func() {
		reqCtx.Complete(false)
	})

	It("creates a good token", func() {
		tk, err := logic.Create(reqCtx, sampleToken())
		Expect(err).To(BeNil())
		check, err := logic.Get(reqCtx, tk.ID)
		Expect(err).To(BeNil())
		Expect(check.AdUsername).To(Equal(tk.AdUsername))
		Expect(check.Token).To(Equal(tk.Token))
		Expect(check.Token).ToNot(Equal("foo123456"))

		check1, err := logic.GetByToken(reqCtx, "foo123456")
		Expect(err).To(BeNil())
		Expect(check1.AdUsername).To(Equal(tk.AdUsername))

	})

	It("gets all tokens", func() {
		_, err := logic.Create(reqCtx, sampleToken())
		Expect(err).To(BeNil())
		result, err := logic.GetAll(reqCtx, 10)
		Expect(err).To(BeNil())
		Expect(len(result)).To(BeNumerically(">", 0))
	})

	It("updates an existing token", func() {
		By("creating a token")
		tk, err := logic.Create(reqCtx, sampleToken())
		tk.Blocked = !tk.Blocked

		By("updating the token")
		_, err = logic.Update(reqCtx, tk.ID, tk)
		Expect(err).To(BeNil())

		By("retrieving the token")
		check, err := logic.Get(reqCtx, tk.ID)
		Expect(err).To(BeNil())

		Expect(check.Blocked).To(Equal(tk.Blocked))
	})

	It("does not update a token that isn't there", func() {
		_, err := logic.Update(reqCtx, 0, sampleToken())
		Expect(err).ToNot(BeNil())
	})

	It("deletes a token", func() {
		tk, err := logic.Create(reqCtx, sampleToken())
		Expect(err).To(BeNil())
		result, err := logic.Delete(reqCtx, tk.ID)
		Expect(err).To(BeNil())
		Expect(result).To(Equal("ok"))
		_, err = logic.Get(reqCtx, tk.ID)
		Expect(err).ToNot(BeNil())
	})

	It("does not delete a token that isn't there", func() {
		_, err := logic.Delete(reqCtx, 0)
		Expect(err).ToNot(BeNil())
	})

	It("searches for tokens", func() {
		baseName := utilities.GenerateRandomString(10)
		for i := 0; i < 5; i++ {
			tk := sampleToken()
			tk.AdUsername = baseName + utilities.GenerateRandomString(10)
			_, err := logic.Create(reqCtx, tk)
			Expect(err).To(BeNil())
		}
		tks, err := logic.Search(reqCtx, baseName, 4)
		Expect(err).To(BeNil())
		Expect(len(tks)).To(Equal(4))
	})

})
