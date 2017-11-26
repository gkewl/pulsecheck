package model

import (
	"fmt"
	"time"

	"github.com/gkewl/pulsecheck/constant"
	"github.com/gkewl/pulsecheck/lib/exclusion"
	"github.com/gkewl/pulsecheck/logger"
	"gopkg.in/guregu/null.v3"
)

// EmployeeStatus represents the model for employeestatuss
type EmployeeStatus struct {
	ID             int64          `json:"id" db:"id"`
	EmployeeID     int64          `json:"employeeid" db:"employeeid"`
	Consider       bool           `json:"consider" db:"consider"`
	OIG            bool           `json:"-" db:"oig"`
	OIGLastSearch  null.Time      `json:"-" db:"oiglastsearch"`
	OIGReference   null.String    `json:"-" db:"oigreference"`
	Sam            bool           `json:"-" db:"sam"`
	SamLastSearch  null.Time      `json:"-" db:"samlastsearch"`
	SamReference   null.String    `json:"-" db:"samreference"`
	Ofac           bool           `json:"-" db:"ofac"`
	OfacLastSearch null.Time      `json:"-" db:"ofaclastsearch"`
	OfacReference  null.String    `json:"-" db:"ofacreference"`
	IsActive       bool           `json:"isactive" db:"isactive"`
	CreatedBy      string         `json:"createdby" db:"createdby"`
	Created        time.Time      `json:"created" db:"created"`
	ModifiedBy     string         `json:"modifiedby" db:"modifiedby"`
	Modified       time.Time      `json:"modified" db:"modified"`
	Sources        []SourceDetail `json:"sources"`
}

// SourceInfo -
type SourceDetail struct {
	Source     string    `json:"source"`
	Consider   bool      `json:"consider"`
	LastSearch null.Time `json:"lastsearch"`
	Reference  string    `json:"reference"`
	Details    string    `json:"details"`
}

// ToSourceDetail -
func (es *EmployeeStatus) ToSourceDetail(sources []string) []SourceDetail {
	si := []SourceDetail{}

	for _, src := range sources {

		s := SourceDetail{Source: src}
		switch src {

		case constant.Source_OIG:
			s.Consider = es.OIG
			s.LastSearch = es.OIGLastSearch
			s.Reference = es.OIGReference.String
		case constant.Source_OFAC:
			s.Consider = es.Ofac
			s.LastSearch = es.OfacLastSearch
			s.Reference = es.OfacReference.String
		case constant.Source_SAM:
			s.Consider = es.Sam
			s.LastSearch = es.SamLastSearch
			s.Reference = es.SamReference.String
		default:
			logger.LogError(fmt.Sprintf("Source type %s not found for getting source details", src), "")
			continue
		}
		s.Details = exclusion.GetDetails(s.Reference)
		si = append(si, s)
	}
	return si
}
