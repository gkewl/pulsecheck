package errorhandler_test

import (
	"errors"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	eh "github.com/gkewl/pulsecheck/errorhandler"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ErrorHandler Test Suite")
}

var _ = Describe("Basic tests", func() {
	It("creates a named error", func() {
		// private so we can't call it, just confirm existing one
		Expect(eh.TestingNamedError.String()).To(Equal("0 - testing"))
		Expect(eh.TestingNamedError.HTTPStatus()).To(Equal(500))
	})

	It("wraps with details", func() {
		err := eh.NewError(eh.TestingNamedError, "details %s", "foo")
		Expect(err.Error()).To(Equal("0 - testing"))
		Expect(err.DetailStack(";")).To(Equal("0: details foo"))
		Expect(err.LocationStack(";")).To(ContainSubstring("errorhandler_test.go"))
	})

	It("adds details", func() {
		err := eh.NewError(eh.TestingNamedError, "details %s", "foo")
		newErr := eh.AddDetail(err, "further %s", "bar").(eh.Error)
		Expect(newErr.DetailStack(";")).To(ContainSubstring("further bar"))
		Expect(newErr.LocationStack(";")).To(MatchRegexp("_test.go.*_test.go"))
	})

	It("wraps an Error", func() {
		err := eh.NewError(eh.TestingNamedError, "inner %s", "foo")
		newErr := eh.WrapError(eh.TestingNamedError, err, "primary %s", "bar")
		Expect(newErr.DetailStack(";")).To(Equal("0: primary bar;0: inner foo"))
	})

	It("wraps an error", func() {
		err := errors.New("foo")
		newErr := eh.WrapError(eh.TestingNamedError, err, "primary %s", "bar")
		Expect(newErr.DetailStack(";")).To(Equal("0: primary bar;foo"))
		Expect(newErr.LocationStack(";")).To(MatchRegexp("_test.go.*;unknown"))
	})

	It("detects sql no rows", func() {
		err := eh.ErrDBNoRows
		Expect(eh.HasNoRowsError(err)).To(BeTrue())
		Expect(eh.NotNoRowsError(err)).To(BeFalse())
		newErr := eh.NewError(eh.TestingNamedError, "DB Error:"+err.Error())
		Expect(eh.HasNoRowsError(newErr)).To(BeTrue())
		Expect(eh.NotNoRowsError(newErr)).To(BeFalse())
	})

	It("handles not found errors", func() {
		err := eh.NewErrorNotFound(eh.TestingNamedError, nil, "details %s", "foo")
		Expect(err.DetailStack(";")).To(ContainSubstring("details foo"))
		err = eh.NewErrorNotFound(eh.TestingNamedError, eh.ErrDBNoRows, "details %s", "foo")
		Expect(err.DetailStack(";")).To(ContainSubstring("details foo"))
		err = eh.NewErrorNotFound(eh.TestingNamedError, errors.New("bad sql problem"), "details %s", "foo")
		Expect(err.DetailStack(";")).To(ContainSubstring("details foo DB Error: bad sql problem"))
	})

	It("reports error text", func() {
		err := eh.NewError(eh.TestingNamedError, "inner %s", "foo")
		newErr := eh.WrapError(eh.TestingNamedError, err, "primary %s", "bar")
		Expect(eh.ContainsErrorText(newErr, "bar")).To(BeTrue())
		Expect(eh.ContainsErrorText(newErr, "nope")).To(BeFalse())
		Expect(eh.ContainsErrorText(nil, "test")).To(BeFalse())
	})
})
