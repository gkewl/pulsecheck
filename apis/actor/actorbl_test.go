package actor_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/guregu/null.v3"

	"github.com/gkewl/pulsecheck/apis/actor"
	"github.com/gkewl/pulsecheck/constant"
	"github.com/gkewl/pulsecheck/errorhandler"
	"github.com/gkewl/pulsecheck/model"
	"github.com/gkewl/pulsecheck/protocol"
	"github.com/gkewl/pulsecheck/utilities"
)

// newActor constructs an Actor with a random name for write operations
func newActor() (actor model.Actor) {
	actorName := "UT" + utilities.GenerateRandomString(10)
	actorDesc := actorName + " " + "Description"
	userType := constant.ActorType_Equipment
	role := model.NullableNameDescription{Name: null.NewString("Operator", true)}
	actor = model.Actor{Name: actorName, Description: null.NewString(actorDesc, true), Type: userType, Role: role}

	return
}

func expectError(expected *errorhandler.NamedError, got error) {
	Expect(got).ToNot(BeNil())
	Expect(got.Error()).To(ContainSubstring(expected.String()))
}

var _ = Describe("Actor Business Logic Tests", func() {
	var reqCtx *protocol.TestRequestContext
	var logic = actor.BLActor{}

	BeforeEach(func() {
		reqCtx = &protocol.TestRequestContext{Userid: 1}
	})
	AfterEach(func() {
		reqCtx.Complete(false)
	})

	It("creates an actor", func() {
		actor, err := logic.Create(reqCtx, newActor())
		Expect(err).To(BeNil())
		check, err := logic.Get(reqCtx, actor.Name)
		Expect(err).To(BeNil())
		Expect(check.Name).To(Equal(actor.Name))
		Expect(check.Role.Id.Int64).To(BeNumerically(">", 0))
	})

	It("creates an actor for machine and verify role ", func() {
		a := newActor()
		a.Role.Name = null.NewString("Designer", true)
		actor, err := logic.Create(reqCtx, a)
		Expect(err).To(BeNil())
		check, err := logic.Get(reqCtx, actor.Name)
		Expect(err).To(BeNil())
		Expect(check.Name).To(Equal(actor.Name))
		Expect(check.Role.Name.String).To(Equal("Operator"))

	})

	It("Does not create actor with wrong role name", func() {
		a := newActor()
		a.Type = "USER"
		a.Role = model.NullableNameDescription{Name: null.NewString("OperatorX", true)}
		_, err := logic.Create(reqCtx, a)
		Expect(err).ToNot(BeNil())

	})

	It("does not create a duplicate actor", func() {
		s := newActor()
		_, err := logic.Create(reqCtx, s)
		Expect(err).To(BeNil())
		_, err = logic.Create(reqCtx, s)
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(ContainSubstring(errorhandler.ErrActorInsert.String()))
	})

	It("updates an existing actor", func() {
		By("creating an actor")
		actor, err := logic.Create(reqCtx, newActor())
		Expect(err).To(BeNil())
		actor.Description = null.NewString(actor.Description.String+"modified", true)
		actor.Type = "SYSTEM"

		By("updating the actor")
		_, err = logic.Update(reqCtx, actor.Id, actor)
		Expect(err).To(BeNil())

		By("retrieving the actor")
		check, err := logic.Get(reqCtx, actor.Name)
		Expect(err).To(BeNil())
		Expect(check.Description.String).To(Equal(actor.Description.String))
		Expect(check.Type).To(Equal(actor.Type))

	})

	It("update last login for actor", func() {
		By("creating an actor")
		actor, err := logic.Create(reqCtx, newActor())
		Expect(err).To(BeNil())

		By("updating last login for actor")
		err = logic.UpdateLastLoginTime(reqCtx, actor.Id)
		Expect(err).To(BeNil())

		By("retrieving the actor")
		check, err := logic.Get(reqCtx, actor.Name)
		Expect(err).To(BeNil())
		Expect(check.LastLoginTime.Valid).To(BeTrue())

	})

	It("does not update an actor that isn't there", func() {
		_, err := logic.Update(reqCtx, 0, newActor())
		Expect(err).ToNot(BeNil())
	})

	It("deletes an actor", func() {
		actor, err := logic.Create(reqCtx, newActor())
		Expect(err).To(BeNil())
		result, err := logic.Delete(reqCtx, actor.Id)
		Expect(err).To(BeNil())
		Expect(result).To(Equal("ok"))
		check, err := logic.Get(reqCtx, actor.Name)
		Expect(err).To(BeNil())
		Expect(check.IsActive).To(Equal(false))
	})

	It("does not delete an actor that isn't there", func() {
		_, err := logic.Delete(reqCtx, 0)
		Expect(err).ToNot(BeNil())
	})

	It("searches for actor", func() {
		baseName := utilities.GenerateRandomString(10)
		for i := 0; i < 5; i++ {
			s := newActor()
			s.Name = baseName + utilities.GenerateRandomString(2)
			s.Description = null.NewString("search", true)
			_, err := logic.Create(reqCtx, s)
			Expect(err).To(BeNil())
		}
		actors, err := logic.Search(reqCtx, "EQUIPMENT", "search")
		Expect(err).To(BeNil())
		Expect(len(actors)).To(Equal(5))
	})

	It("does not update if no change", func() {
		actor, err := logic.Create(reqCtx, newActor())
		Expect(err).To(BeNil())
		rows, err := logic.Update(reqCtx, actor.Id, actor)
		Expect(err).To(BeNil())
		Expect(rows.Modified).To(Equal(actor.Modified))
	})

	It("rejects a bad type", func() {
		actor, err := logic.Create(reqCtx, newActor())
		Expect(err).To(BeNil())
		actor.Type = "MEH"
		_, err = logic.Update(reqCtx, actor.Id, actor)
		Expect(err).ToNot(BeNil())
		check, err := logic.Get(reqCtx, actor.Name)
		Expect(check.Type).To(Equal("EQUIPMENT"))
	})
})
