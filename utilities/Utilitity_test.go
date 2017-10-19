package utilities_test

import (
	"fmt"
	"github.com/gkewl/pulsecheck/utilities"
	"testing"
)

func TestGenerateRandomString(t *testing.T) {

	fmt.Println(utilities.GenerateRandomString(5))
}
