package dbhandler_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gkewl/pulsecheck/config"
	"github.com/gkewl/pulsecheck/dbhandler"
	eh "github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
	"github.com/gkewl/pulsecheck/protocol"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Protocol Test Suite")
}

var _ = BeforeSuite(func() {
	config.LoadConfigurations()

})

var _ = Describe("Collation Tests", func() {

	It("Fails when any column has a collation of utf8_general_ci", func() {
		var count int
		db, err := dbhandler.CreateConnection()
		Expect(err).To(BeNil())

		badCollationRows := db.QueryRow(`
			SELECT COUNT(*) FROM information_schema.columns WHERE TABLE_SCHEMA='pulsecheck' AND COLLATION_NAME = 'utf8_general_ci'
		`)

		badCollationRows.Scan(&count)
		if count > 0 {
			By(`
				There are columns in dev0 with the incorrect collation. This can happen when columns are added manually without using goose.
				Run the following query in dev0 to find the offending columns:
				SELECT * FROM information_schema.columns WHERE TABLE_SCHEMA='github.com/gkewl/pulsecheck' AND COLLATION_NAME = 'utf8_general_ci'
				All columns should have collation set to "utf8_unicode_ci". Please use goose for any database modifications.
			`)
		}
		Expect(count).To(Equal(0))
	})
})
