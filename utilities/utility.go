package utilities

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/kylelemons/godebug/pretty"

	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/constant"
)

// NullInt64 -
type NullInt64 sql.NullInt64

var (
	// DefaultFormatter is the default set of overrides for stringification.
	DefaultFormatter = map[reflect.Type]interface{}{
		reflect.TypeOf(time.Time{}):          fmt.Sprint,
		reflect.TypeOf(net.IP{}):             fmt.Sprint,
		reflect.TypeOf((*error)(nil)).Elem(): fmt.Sprint,
	}

	// CompareConfig is the default configuration used for Compare.
	CompareConfig = &pretty.Config{
		Compact:           true,
		Diffable:          true,
		IncludeUnexported: true,
		Formatter:         DefaultFormatter,
	}

	// DefaultConfig is the default configuration used for all other top-level functions.
	DefaultConfig = &pretty.Config{
		Formatter: DefaultFormatter,
	}
)

// PrintSlice -
func PrintSlice(slice interface{}) {

	pretty.Print(slice)
}

// CompareSlice -
func CompareSlice(source interface{}, dest interface{}) string {
	return pretty.Compare(source, dest)
}

// MarshalJSON -
func (i *NullInt64) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return nil, nil
	}
	return []byte(strconv.FormatInt(i.Int64, 10)), nil
}

// UnmarshalJSON -
func (i *NullInt64) UnmarshalJSON(b []byte) error {
	s := string(b)
	if v, err := strconv.ParseInt(s, 10, 64); err != nil {
		i.Int64 = int64(v)
		return nil
	}
	if s == "null" {
		return nil
	}
	return errors.New("Invalid NullINt64: " + s)
}

// ParseInt64 -
func ParseInt64(v string) int64 {
	s, _ := strconv.ParseInt(v, 10, 64)
	return s
}

// ParseInt -
func ParseInt(v string) int {
	i, _ := strconv.Atoi(v)
	return i
}

// ParseInt64ToString -
func ParseInt64ToString(v int64) string {
	s := strconv.FormatInt(v, 10)
	return s
}

// ParseFloatToString -
func ParseFloatToString(v float64) string {
	s := strconv.FormatFloat(v, 'E', -1, 64)
	return s
}

// ParseIntToString -
func ParseIntToString(v int) string {
	s := strconv.Itoa(v)
	return s
}

// StructuredResponseTest -
type StructuredResponseTest struct {
	Xid        string                `json:"xid"`
	StatusCode int                   `json:"statuscode"`
	Response   types.JSONText        `json:"response,omitempty"`
	Error      *common.ErrorResponse `json:"error,omitempty"`
}

// ScanResponseObject -
func ScanResponseObject(resp io.ReadCloser, dest interface{}) (err error) {
	var response StructuredResponseTest
	if err = json.NewDecoder(resp).Decode(&response); err != nil {
		return
	}
	err = json.Unmarshal(response.Response, &dest)
	return
}

// GenerateRandomString -
func GenerateRandomString(s int) string {
	b := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b)
}

// all the formats from the time module
var formats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.UnixDate,
	time.RubyDate,
	time.ANSIC,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
}

// tries to parse time using a couple of formats before giving up
func tryParseTime(value string) (time.Time, bool) {
	var t time.Time
	var err error
	for _, layout := range formats {
		t, err = time.Parse(layout, value)
		if err == nil {
			return t, true
		}
	}
	return t, false
}

// ParseTime converts a string to a time but allows for a range of different formats
func ParseTime(v string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		var ok bool
		t, ok = tryParseTime(v)
		if ok {
			err = nil
		}
	}
	return t, err
}

// MaxDate returns the maximum date in UTC
func MaxDate() time.Time {
	maxTime, _ := ParseStringtoTime("9999-12-31T00:00:00Z")
	return maxTime // since the time string is constant, this should never return an error
}

// MinDate returns the maximum date in UTC supported by MySQL timestamp
func MinDate() time.Time {
	minTime, _ := ParseStringtoTime("1970-01-01T00:00:01Z")
	return minTime // since the time string is constant, this should never return an error
}

// ParseStringtoTime converts a RFC3339 date strictly
func ParseStringtoTime(v string) (time.Time, error) {
	return time.Parse(time.RFC3339, v)
}

const dateFormat = "2006-01-02" // yyyy-mm-dd

// ParseStringToDate -
func ParseStringToDate(val string) (time.Time, error) {
	return time.Parse(dateFormat, val)
}

// ParseDateToString -
func ParseDateToString(t time.Time) string {
	return t.Format(dateFormat)
}

// ParseStringtoFloat64 -
func ParseStringtoFloat64(v string) (float64, error) {
	f, err := strconv.ParseFloat(v, 64)
	return f, err
}

// Now returns a UTC time truncated to the second for db compatibility
func Now() time.Time {
	return time.Now().Truncate(1 * time.Second).UTC()
}

// DurationTruncated returns floating point seconds with N digits of precision
// between two times
func DurationTruncated(start, end time.Time, precision float64) float64 {
	dur := float64(end.Sub(start)) / float64(1*time.Second)
	return float64(int64(dur*math.Pow(10, precision)+0.5)) / float64(math.Pow(10, precision))
}

// AllErrorsNil returns true if a list of errors is all nils
func AllErrorsNil(errs ...error) bool {
	for _, e := range errs {
		if e != nil {
			return false
		}
	}
	return true
}

// Blank returns true if trimmed string is empty
func Blank(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// ContainString -
func ContainString(s []string, e string) bool {
	for _, a := range s {
		if strings.EqualFold(a, e) {
			return true
		}
	}
	return false
}

// ContainsInt -
func ContainsInt(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// ContainsCI case insensitive string contains
func ContainsCI(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// Remove from slice safely
func Remove(a []int64, i int64) []int64 {
	newSlice := []int64{}

	for _, val := range a {
		if val != i {
			newSlice = append(newSlice, val)
		}
	}
	return newSlice
}

// RemoveLabelPrefix removes 1J/5J/6J from a label
func RemoveLabelPrefix(label string) string {
	if len(label) > 2 {
		prefix := strings.ToUpper(label[:2])
		if prefix == "1J" || prefix == "6J" || prefix == "5J" {
			return label[2:len(label)]
		}
	}
	return label
}

// ValidateSearchParams validates params
func ValidateSearchParams(params map[string]string) (map[string]string, error) {
	var err error
	var strError string
	for _, value := range params {
		if len(strings.Split(value, " ")) > 1 {
			strError = strError + "Search value has multiple values. Values: " + value
		}
	}

	fromdate, fok := params["fromdate"]
	todate, tok := params["todate"]

	if !fok || len(fromdate) == 0 {
		strError += " FromDate cannot be null" + constant.Newline
	} else {
		if !tok || len(todate) == 0 { // If todate is null, make it to current date
			params["todate"] = fmt.Sprintf("%s", Now())
		} else {

			_, err := time.Parse(time.RFC3339, fromdate)
			if err != nil {
				strError += " Invalid format for from date" + constant.Newline
			}
			_, err = time.Parse(time.RFC3339, todate)
			if err != nil {
				strError += " Invalid format for to date" + constant.Newline
			}

		}
	}

	if len(params) == 0 {
		strError += "search parameters are not available"
	}
	if len(strError) > 0 {
		err = errors.New(strError)
	}

	return params, err
}

// ConfirmValuesInSlice checks for a variable number of values in
// a slice of structs. Usage: ConfirmValuesInSlice(mydata, "ID", crd1.ID, crd2.ID)
// Returns nil if successful, string with missing values otherwise
// TO BE USED IN TESTING ONLY
func ConfirmValuesInSlice(data interface{}, fieldName string, values ...interface{}) interface{} {
	if len(values) == 0 {
		return true
	}
	if reflect.TypeOf(data).Kind() != reflect.Slice {
		panic("must send slice")
	}
	dSlice := reflect.ValueOf(data)
	found := map[string]bool{}
	for i := 0; i < dSlice.Len(); i++ {
		r := dSlice.Index(i)
		f := reflect.Indirect(r).FieldByName(fieldName)
		if !f.IsValid() {
			panic("field " + fieldName + " not found in struct")
		}
		found[fmt.Sprintf("%v", f.Interface())] = true
	}
	missing := []string{}
	for _, v := range values {
		val := fmt.Sprintf("%v", v)
		if _, ok := found[val]; !ok {
			missing = append(missing, val)
		}
	}
	if len(missing) == 0 {
		return nil
	}
	return "missing values: " + strings.Join(missing, ", ")
}

// ExtractValuesFromSlice assumes data is a slice of structs and returns
// an array of values of the field from those structs
// Returns empty slice if there are typing problems or missing field
func ExtractValuesFromSlice(data interface{}, fieldName string) []interface{} {
	if reflect.TypeOf(data).Kind() != reflect.Slice {
		return []interface{}{}
	}
	dSlice := reflect.ValueOf(data)
	results := []interface{}{}
	for i := 0; i < dSlice.Len(); i++ {
		r := dSlice.Index(i)
		f := reflect.Indirect(r).FieldByName(fieldName)
		if !f.IsValid() {
			return []interface{}{}
		}
		results = append(results, f.Interface())
	}
	return results
}

// ChunkBy splits a string s into chunks of length n
func ChunkBy(s string, n int) (chunks []string) {
	i := 0
	for i+n < len(s) {
		chunks = append(chunks, s[i:i+n])
		i += n
	}
	if i < len(s) {
		chunks = append(chunks, s[i:len(s)])
	}
	return
}

// ToCsvString takes a slice of ints and returns its CSV string representation
// if slice if empty or nil an empty string is returned
func ToCsvString(numbers []int64) string {
	strs := []string{}
	if numbers == nil || len(numbers) <= 0 {
		return ""
	}

	for _, n := range numbers {
		strs = append(strs, fmt.Sprintf("%d", n))
	}

	return strings.Join(strs, ",")
}

// FirstError returns first non-nil error from variable list of errors
func FirstError(errs ...error) error {
	for _, e := range errs {
		if e != nil {
			return e
		}
	}
	return nil
}

// ConvertInchesToMillimeters -
func ConvertInchesToMillimeters(v float64) float64 {
	mmPerInch := 25.4
	return v * mmPerInch
}

// ConvertCentimetersToMillimeters -
func ConvertCentimetersToMillimeters(v float64) float64 {
	return v * 10
}

// GenerateRandomBytes -
func GenerateRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil
	}

	return b
}
