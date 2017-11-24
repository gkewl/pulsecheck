package elasticsearch_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gkewl/pulsecheck/errorhandler"
	tst "github.com/gkewl/pulsecheck/lib/testing"
	"github.com/gkewl/pulsecheck/protocol"
)

func expectError(expected *errorhandler.NamedError, got error) {
	Expect(got).ToNot(BeNil())
	Expect(got.Error()).To(ContainSubstring(expected.String()))
}

var _ = Describe("Company Business Logic Tests", func() {
	var reqCtx *protocol.TestRequestContext
	//var logic = elasticsearch.BLElasticSearch{}
	var t tst.T

	BeforeEach(func() {
		reqCtx = &protocol.TestRequestContext{Userid: 1}
		t = tst.T{ReqCtx: reqCtx}
	})
	AfterEach(func() {
		reqCtx.Complete(false)
	})

	// It("creates a good company", func() {
	// 	comp := t.Company(logic.Create(reqCtx, t.SampleCompany()))
	// 	check := t.GetCompany(comp.ID)
	// 	Expect(check.Name).To(Equal(comp.Name))
	// })

})
