package model

//	"fmt"
//	"time"

//	"gopkg.in/guregu/null.v3"

// ElasticSearchResult -
type ElasticSearchResult struct {
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
}

// OIGSearch -
type OIGSearch struct {
	Firstname   string `json:"firstname"`
	Middlename  string `json:"middilename"`
	Lastname    string `json:"lastname"`
	DateOfBirth string `json:"dateofbirth"`
	CompanyName string `json:"companyname"`
	Addresss    string `json:"adddress"`
	City        string `json:"city"`
	State       string `json:"state"`
	ZipCode     string `json:"zipcode"`
}
