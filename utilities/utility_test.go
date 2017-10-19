package utilities_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gkewl/pulsecheck/utilities"
)

type local struct {
	ID   int64
	Name string
}

var _ = Describe("utility tests", func() {
	It("generates random strings", func() {
		probes := map[string]bool{}
		for i := 0; i < 1000; i++ {
			s := utilities.GenerateRandomString(10)
			if _, exists := probes[s]; exists {
				Expect(true).To(BeFalse(), "Unexpectedly found random string")
			}
			probes[s] = true
		}
	})

	It("does case insensitive compare", func() {
		Expect(utilities.ContainsCI("MYTEST", "yTes")).To(BeTrue())
		Expect(utilities.ContainsCI("MYTEST", "foo")).To(BeFalse())
	})

	It("tests blank", func() {
		Expect(utilities.Blank("")).To(BeTrue())
		Expect(utilities.Blank("   ")).To(BeTrue())
		Expect(utilities.Blank("foo")).To(BeFalse())
	})

	It("confirms values in slice", func() {
		localData := []local{{ID: 10, Name: "foo"}, {ID: 20, Name: "bar"}}
		Expect(utilities.ConfirmValuesInSlice(localData, "ID", 10, 20)).To(BeNil())
		Expect(utilities.ConfirmValuesInSlice(localData, "ID", 10, 20, 30)).ToNot(BeNil())
		Expect(utilities.ConfirmValuesInSlice(localData, "Name", "foo", "bar")).To(BeNil())
		Expect(utilities.ConfirmValuesInSlice(localData, "Name", "foo", "bar", "foobar")).ToNot(BeNil())
	})

	Describe("ChunkBy", func() {
		It("splits a string 's' into chunks of length 'n', with a final smaller chunk if necessary", func() {
			original := "abcdefghijklmnopqrstuvwxyz"
			chunkedBy2 := []string{"ab", "cd", "ef", "gh", "ij", "kl", "mn", "op", "qr", "st", "uv", "wx", "yz"}
			chunkedBy7 := []string{"abcdefg", "hijklmn", "opqrstu", "vwxyz"}
			chunkedBy13 := []string{"abcdefghijklm", "nopqrstuvwxyz"}

			Expect(utilities.ChunkBy(original, 2)).To(Equal(chunkedBy2))
			Expect(utilities.ChunkBy(original, 7)).To(Equal(chunkedBy7))
			Expect(utilities.ChunkBy(original, 13)).To(Equal(chunkedBy13))
		})
	})

	Describe("first error", func() {
		It("returns nil", func() {
			Expect(utilities.FirstError()).To(BeNil())
			Expect(utilities.FirstError(nil)).To(BeNil())
			Expect(utilities.FirstError(nil, nil)).To(BeNil())
		})

		It("returns first", func() {
			e1 := fmt.Errorf("foo")
			e2 := fmt.Errorf("bar")
			Expect(utilities.FirstError(e1).Error()).To(Equal("foo"))
			Expect(utilities.FirstError(nil, e1).Error()).To(Equal("foo"))
			Expect(utilities.FirstError(nil, e1, nil).Error()).To(Equal("foo"))
			Expect(utilities.FirstError(nil, e1, e2).Error()).To(Equal("foo"))
		})
	})
})
