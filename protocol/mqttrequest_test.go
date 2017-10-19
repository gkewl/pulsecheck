package protocol_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/dbhandler"
	"github.com/gkewl/pulsecheck/protocol"
)

var _ = Describe("MQTTRequest", func() {
	var ctx = common.AppContext{}
	var rc protocol.MQTTRequestContext
	var err error
	var meta = `"_meta_":{"xid":"xid","token":"token","reply_topic":"reply_to","action":"action"}`

	ctx.Db, _ = dbhandler.CreateConnection()

	AfterEach(func() {
		_ = rc.Tx().Rollback()
	})

	It("makes parameters available", func() {
		payload := "{" + meta + `, "name":"foo", "id": 42}`
		rc, err = protocol.NewMQTTRequestContext(&ctx, "/topic", payload)
		Expect(err).To(BeNil())
		Expect(rc.Value("name", "")).To(Equal("foo"))
		Expect(rc.Value("notthere", "bar")).To(Equal("bar"))
		Expect(rc.IntValue("id", 0)).To(Equal(int64(42)))
		Expect(rc.IntValue("notthere", 24)).To(Equal(int64(24)))
		Expect(rc.IntValue("name", 24)).To(Equal(int64(24)))
	})

	It("scans posted content", func() {
		payload := "{" + meta + `, "model":{"name":"foo", "vals":[1,2,3]}}`
		rc, err = protocol.NewMQTTRequestContext(&ctx, "/topic", payload)
		Expect(err).To(BeNil())
		dest := TestStruct{}
		err = rc.Scan("model", &dest)
		Expect(err).To(BeNil())
		Expect(dest.Name).To(Equal("foo"))
		Expect(dest.Vals).To(Equal([]int{1, 2, 3}))
	})

	It("captures meta fields", func() {
		payload := "{" + meta + "}"
		rc, err = protocol.NewMQTTRequestContext(&ctx, "/topic", payload)
		Expect(err).To(BeNil())
		Expect(rc.Xid()).To(Equal("xid"))
		Expect(rc.Action()).To(Equal("action"))
		Expect(rc.ReplyTopic()).To(Equal("reply_to"))
	})
})
