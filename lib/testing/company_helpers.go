package testing

import (
	//. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gkewl/pulsecheck/apis/company"
	"github.com/gkewl/pulsecheck/model"
	"github.com/gkewl/pulsecheck/utilities"
)

// sampleCompany constructs a company
func (t *T) SampleCompany() (comp model.Company) {
	return model.Company{
		Name:       "UT_comp_" + utilities.GenerateRandomString(8),
		IsActive:   true,
		CreatedBy:  "admin",
		ModifiedBy: "admin",
	}
}

// Company expects a company and an error and verifies the error is nil
// and returns the company
func (t *T) Company(comp model.Company, e error) model.Company {
	Expect(e).To(BeNil(), callers())
	return comp
}

// GetCompany fetches a company using an ID and verifies no error
func (t *T) GetCompany(ID int) model.Company {
	Expect(t.ReqCtx).ToNot(BeNil(), callers())
	comp, e := company.BLCompany{}.Get(t.ReqCtx, ID)
	Expect(e).To(BeNil(), callers())
	return comp
}

// ReGetCompany expects a company and an error and verifies the error is nil
// and re-gets the company and returns it
func (t *T) ReGetCompany(comp model.Company, e error) model.Company {
	Expect(e).To(BeNil(), callers())
	Expect(t.ReqCtx).ToNot(BeNil(), callers())
	comp, e = company.BLCompany{}.Get(t.ReqCtx, comp.ID)
	Expect(e).To(BeNil(), callers())
	return comp
}

// UpdateCompany expects a company struct and updates it, verifying no error
func (t *T) UpdateCompany(comp model.Company) model.Company {
	Expect(t.ReqCtx).ToNot(BeNil(), callers())
	comp, e := company.BLCompany{}.Update(t.ReqCtx, comp.ID, comp)
	Expect(e).To(BeNil(), callers())
	return comp
}

// Companys expects a company slice and an error and verifies the error is nil
// and returns the companys
func (t *T) Companys(comps []model.Company, e error) []model.Company {
	Expect(e).To(BeNil(), callers())
	return comps
}

// CompanyErr expects a company and an error and verifies the error is
// not nil and returns the error
func (t *T) CompanyErr(comp model.Company, e error) error {
	Expect(e).ToNot(BeNil(), callers())
	return e
}
