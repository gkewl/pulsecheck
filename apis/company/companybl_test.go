package company_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gkewl/pulsecheck/apis/company"

	"github.com/gkewl/pulsecheck/errorhandler"
	tst "github.com/gkewl/pulsecheck/lib/testing"
	"github.com/gkewl/pulsecheck/model"
	"github.com/gkewl/pulsecheck/protocol"
	"github.com/gkewl/pulsecheck/utilities"
)

// sampleCompany constructs a company
func sampleCompany() (comp model.Company) {
	return model.Company{
		Name:       "company_" + utilities.GenerateRandomString(8),
		IsActive:   true,
		CreatedBy:  "admin",
		ModifiedBy: "admin",
	}
}

func expectError(expected *errorhandler.NamedError, got error) {
	Expect(got).ToNot(BeNil())
	Expect(got.Error()).To(ContainSubstring(expected.String()))
}

var _ = Describe("Company Business Logic Tests", func() {
	var reqCtx *protocol.TestRequestContext
	var logic = company.BLCompany{}
	var t tst.T

	BeforeEach(func() {
		reqCtx = &protocol.TestRequestContext{Userid: 1}
		t = tst.T{ReqCtx: reqCtx}
	})
	AfterEach(func() {
		reqCtx.Complete(false)
	})

	It("creates a good company", func() {
		comp := t.Company(logic.Create(reqCtx, sampleCompany()))
		check := t.GetCompany(comp.ID)
		Expect(check.Name).To(Equal(comp.Name))
	})

	It("gets all companys", func() {
		t.Company(logic.Create(reqCtx, sampleCompany()))
		result := t.Companys(logic.GetAll(reqCtx, 10))
		Expect(len(result)).To(BeNumerically(">", 0))
	})

	It("updates an existing company", func() {
		By("creating a company")
		comp := t.Company(logic.Create(reqCtx, sampleCompany()))
		comp.Name = comp.Name + utilities.GenerateRandomString(8)

		By("updating the company")
		t.UpdateCompany(comp)

		By("retrieving the company")
		check := t.GetCompany(comp.ID)
		Expect(check.Name).To(Equal(comp.Name))
	})

	It("does not update a company that isn't there", func() {
		_ = t.CompanyErr(logic.Update(reqCtx, 0, sampleCompany()))
	})

	It("deletes a company", func() {
		comp := t.Company(logic.Create(reqCtx, sampleCompany()))
		result, err := logic.Delete(reqCtx, comp.ID)
		Expect(err).To(BeNil())
		Expect(result).To(Equal("Ok"))
		_ = t.CompanyErr(logic.Get(reqCtx, comp.ID))
	})

	It("does not delete a company that isn't there", func() {
		_, err := logic.Delete(reqCtx, 0)
		Expect(err).ToNot(BeNil())
	})

	// It("searches for companys", func() {
	// 	baseName := utilities.GenerateRandomString(10)
	// 	for i := 0; i < 5; i++ {
	// 		comp := sampleCompany()
	// 		comp. = baseName + utilities.GenerateRandomString(10)
	// 		t.Company(logic.Create(reqCtx, comp))
	// 	}
	// 	comps := t.Companys(logic.Search(reqCtx, baseName, 4))
	// 	Expect(len(comps)).To(Equal(4))
	// })

})
