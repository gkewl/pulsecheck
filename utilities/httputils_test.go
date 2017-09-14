package utilities_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gkewl/pulsecheck/utilities"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "utilities Test Suite")
}

var _ = Describe("Getbytes test", func() {

	BeforeEach(func() {

	})

	AfterEach(func() {

	})

	It("convert int to bytes", func() {
		_, err := utilities.GetBytes(100)
		Expect(err).To(BeNil())
	})

	It("convert string to bytes", func() {
		_, err := utilities.GetBytes("utilities Test Suite")
		Expect(err).To(BeNil())
	})

	It("convert struct to bytes", func() {
		type TestData struct {
			ID  int
			Msg string
		}
		_, err := utilities.GetBytes(TestData{ID: 1, Msg: "Hello"})
		Expect(err).To(BeNil())
	})

	It("convert float to bytes", func() {
		_, err := utilities.GetBytes(100.96)
		Expect(err).To(BeNil())
	})

	It("convert string array to bytes", func() {
		data := []string{"X", "Y", "z"}
		_, err := utilities.GetBytes(data)
		Expect(err).To(BeNil())
	})

	It("convert int array to bytes", func() {
		data := []int{1, 2, 3}
		_, err := utilities.GetBytes(data)
		Expect(err).To(BeNil())
	})
})
