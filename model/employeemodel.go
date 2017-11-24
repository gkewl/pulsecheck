package model

import (
	"fmt"
	"time"

	"gopkg.in/guregu/null.v3"
)

// Employee represents the model for employees
type Employee struct {
	ID           int64       `json:"id" db:"id"`
	CompanyID    int         `json:"companyid" db:"companyid"`
	Firstname    string      `json:"firstname" db:"firstname"`
	Middlename   null.String `json:"middlename" db:"middlename"`
	Lastname     string      `json:"lastname" db:"lastname"`
	Dateofbirth  string      `json:"dateofbirth" `
	DateofbirthT time.Time   `json:"-" db:"dateofbirtht"`
	Type         int         `json:"type" db:"type"`
	IsActive     bool        `json:"isactive" db:"isactive"`
	CreatedBy    string      `json:"createdby" db:"createdby"`
	Created      time.Time   `json:"created" db:"created"`
	ModifiedBy   string      `json:"modifiedby" db:"modifiedby"`
	Modified     time.Time   `json:"modified" db:"modified"`
	Consider     bool        `json:"consider" db:"consider"`
}

// ToString converts information to string
func (e *Employee) ToString() string {
	return fmt.Sprintf("FirstName %s MiddleName %s LastName %s Dateofbirth %v CompanyId %d",
		e.Firstname, e.Middlename.String, e.Lastname, e.Dateofbirth, e.CompanyID)
}

// ToOIG -
func (e *Employee) ToOIG() OIGSearch {
	return OIGSearch{
		UniqueUser:  e.Firstname + "_" + e.Lastname + "_" + e.Dateofbirth,
		Firstname:   e.Firstname,
		Middlename:  e.Middlename.String,
		Lastname:    e.Lastname,
		DateOfBirth: e.Dateofbirth,
	}
}

// type Date struct {
// 	time.Time `json:",string"`
// }

// const DateFormat = "2006-01-02" // yyyy-mm-dd

// func (d Date) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(d.Format(DateFormat))
// }

// func (t *Date) UnmarshalJSON(data []byte) error {

// 	// var s string
// 	// if err := json.Unmarshal(data, &s); err != nil {
// 	// 	return fmt.Errorf("birthdate should be a string, got %s", data)
// 	// }
// 	// t, err := time.Parse(DateFormat, s)
// 	// if err != nil {
// 	// 	return fmt.Errorf("invalid date: %v", err)
// 	// }
// 	// d.Time = t

// 	// return nil

// 	var err error
// 	var v interface{}
// 	if err = json.Unmarshal(data, &v); err != nil {
// 		return err
// 	}
// 	switch x := v.(type) {
// 	case string:
// 		//err = t.Time.UnmarshalJSON(data)
// 		t1, err := time.Parse(DateFormat, v.(string))
// 		if err != nil {
// 			return fmt.Errorf("invalid date: %v", err)
// 		}
// 		t.Time = t1

// 	case map[string]interface{}:
// 		ti, tiOK := x["Time"].(string)
// 		if !tiOK {
// 			return fmt.Errorf(`json: unmarshalling object into Go value of type null.Time requires key "Time" to be of type string and key "Valid" to be of type bool; found %T and %T, respectively`, x["Time"], x["Valid"])
// 		}
// 		//err = t.Time.UnmarshalText([]byte(ti))
// 		t1, err := time.Parse(DateFormat, ti)
// 		if err != nil {
// 			return fmt.Errorf("invalid date: %v", err)
// 		}
// 		t.Time = t1

// 		return err
// 	case nil:
// 		return nil
// 	default:
// 		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type null.Time", reflect.TypeOf(v).Name())
// 	}
// 	return err

// }
