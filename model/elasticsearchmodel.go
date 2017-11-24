package model

import (
	"github.com/gkewl/pulsecheck/constant"
)

//	"fmt"
//	"time"

//	"gopkg.in/guregu/null.v3"

// ElasticSearchResult -
type ElasticSearchResult struct {
	EmployeeID  int
	Firstname   string
	Middlename  string
	Lastname    string
	DateOfBirth string
	CompanyName string
	Addresss    string
	City        string
	State       string
	ZipCode     string
	Source      string
	Excltype    string
}

// IsMatched -
func (e *ElasticSearchResult) IsMatched(emp Employee) bool {
	if e.Firstname == emp.Firstname &&
		e.Lastname == emp.Lastname &&
		e.DateOfBirth == emp.Dateofbirth {
		return true
	}

	return false
}

// OIGSearch -
type OIGSearch struct {
	UniqueUser  string `json:"uniqueuser"`
	Firstname   string `json:"firstname"`
	Middlename  string `json:"middlename"`
	Lastname    string `json:"lastname"`
	DateOfBirth string `json:"dateofbirth"`
	CompanyName string `json:"companyname"`
	Addresss    string `json:"adddress"`
	City        string `json:"city"`
	State       string `json:"state"`
	ZipCode     string `json:"zipcode"`
	SourceDate  string `json:"sourcedate"`
	Created     string `json:"created"`
	// excldate    time.Time `json:"excldate"`
	// reindate    time.Time `json:"reindate"`
	// waivdate    time.Time `json:"waivdate"`
	Excltype string `json:"excltype"`
}

type OIG struct {
	UniqueUser  string `json:"uniqueuser"`
	Firstname   string `json:"firstname"`
	Middlename  string `json:"middlename"`
	Lastname    string `json:"lastname"`
	DateOfBirth string `json:"dateofbirth"`
	CompanyName string `json:"companyname"`
	Addresss    string `json:"adddress"`
	City        string `json:"city"`
	State       string `json:"state"`
	ZipCode     string `json:"zipcode"`
	SourceDate  string `json:"sourcedate"`
	Created     string `json:"created"`
	// excldate time.Time  `json:"excldate"`
	// reindate time.Time  `json:"reindate"`
	// waivdate time.Time  `json:"waivdate"`
	Excltype string `json:"excltype"`
}

// ToResult -
func (s *OIGSearch) ToResult() ElasticSearchResult {
	return ElasticSearchResult{
		Source:      constant.Source_OIG,
		Firstname:   s.Firstname,
		Middlename:  s.Middlename,
		Lastname:    s.Lastname,
		DateOfBirth: s.DateOfBirth,
		CompanyName: s.CompanyName,
		Addresss:    s.Addresss,
		City:        s.City,
		State:       s.State,
		ZipCode:     s.ZipCode,
		Excltype:    s.Excltype,
	}
}
