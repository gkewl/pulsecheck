package testing

import (
	"fmt"
	"regexp"
	"runtime"

	. "github.com/onsi/gomega"

	"github.com/gkewl/pulsecheck/common"
)

func init() {

}

// T -
type T struct {
	ReqCtx common.RequestContext
}

func callers() string {
	_, f1, ln1, _ := runtime.Caller(2) // first func outside of helper
	_, f2, ln2, _ := runtime.Caller(3)
	_, f3, ln3, _ := runtime.Caller(4)
	re := regexp.MustCompile(".*/")
	return fmt.Sprintf("*CALL STACK: %s:%d, %s:%d, %s:%d",
		re.ReplaceAllString(f1, ""), ln1,
		re.ReplaceAllString(f2, ""), ln2,
		re.ReplaceAllString(f3, ""), ln3)
}

// ExpectError -
func ExpectError(expectedErrorText string, got error) {
	Expect(got).ToNot(BeNil())
	Expect(got.Error()).To(ContainSubstring(expectedErrorText))
}
