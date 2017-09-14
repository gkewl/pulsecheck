package authentication_test

import (
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gkewl/pulsecheck/config"
	"github.com/gkewl/pulsecheck/constant"
	"github.com/gkewl/pulsecheck/protocol"

	auth "github.com/gkewl/pulsecheck/authentication"
)

var (
	authBackendInstance *auth.JWTAuthenticationBackend
)

func init() {
	authBackendInstance = auth.InitJWTAuthenticationBackend()
}

var _ = Describe("Authorization/authentication tests", func() {
	var recorder *httptest.ResponseRecorder
	var request *http.Request

	BeforeEach(func() {
		Expect(authBackendInstance).NotTo(BeNil())
		recorder = httptest.NewRecorder()
		request = httptest.NewRequest("GET", "/", nil)
	})

	var allTests = func(caller protocol.TestCaller) {
		It("Can get request info without a token", func() {
			Expect(auth.GetCurrentUserId(request)).To(BeZero())
			Expect(auth.GetCurrentUsername(request)).To(BeEmpty())
			Expect(auth.GetCurrentUserRole(request)).To(BeEmpty())
		})

		It("Can get authenticated requests", func() {
			token, err := authBackendInstance.GenerateToken("bob", 12, "fake", "Operator")
			Expect(err).To(BeNil())
			request.Header.Add("Authorization", "Bearer "+token.Token)

			req := auth.AnnotateRequestWithAuthorizedUser(recorder, request, constant.Guest)
			Expect(recorder.Body.Len()).To(BeZero()) // no error response written

			Expect(auth.GetCurrentUserId(req)).To(Equal(int64(12)))
			Expect(auth.GetCurrentUsername(req)).To(Equal("bob"))
			Expect(auth.GetCurrentUserRole(req)).To(Equal("Operator"))
		})

	}
	allTests(&httpCaller)
})
