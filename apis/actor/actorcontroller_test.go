package actor_test

import (
	"fmt"
	"github.com/gkewl/pulsecheck/apis/actor"
	"github.com/gkewl/pulsecheck/common"
	"github.com/gkewl/pulsecheck/dbhandler"
	"gopkg.in/guregu/null.v3"
	"reflect"
	"testing"
)

type MockActor struct {
}

var ctx common.AppContext

func init() {
	var err error
	ctx = common.AppContext{}
	ctx.Db, err = dbhandler.CreateConnection()
	if err != nil {
		//handle error
	}

}

func (m MockActor) GetAll(ctx *common.AppContext) ([]actor.Actor, error) {

	actorName := "TEST"
	adusername := "test"
	email := "test@gmail.com"
	actorType := "USER"

	var mock []actor.Actor

	mock = append(mock, actor.Actor{0, actorName, adusername, email, actorType, null.NewString("", false), null.NewString("", false), null.NewInt(0, false), 1, ""})

	//	if actorName == "TEST"{
	//	     mock = actor.Actor{ 0 , actorName,adusername, email , actorType, "" ,  "" , 0 , 1, "" }
	//	     } else {
	//	     mock = actor.Actor{ 0 , actorName,adusername, email , actorType, "" ,  "" , 0 , 1, "" }
	//	     }
	return mock, nil
}

func (m MockActor) Get(name string, ctx *common.AppContext) (actor.Actor, error) {

	actorName := "TEST"
	adusername := "test"
	email := "test@gmail.com"
	actorType := "USER"

	var mock actor.Actor
	if name == "TEST" {
		mock = actor.Actor{0, actorName, adusername, email, actorType, null.NewString("", false), null.NewString("", false), null.NewInt(0, false), 1, ""}
	} else {
		mock = actor.Actor{0, actorName, adusername, email, actorType, null.NewString("", false), null.NewString("", false), null.NewInt(0, false), 1, ""}
	}
	return mock, nil
}

func (m MockActor) Save(id int64, actor actor.Actor, ctx *common.AppContext) (int64, error) {

	return 1, nil
}

func (di MockActor) Delete(ID int64, ctx *common.AppContext) (int64, error) {
	return 1, nil
}

func (di MockActor) Create(actor actor.Actor, ctx *common.AppContext) (actor.Actor, error) {

	actor.Id = 100
	return actor, nil
}

func (di MockActor) Search(ctx *common.AppContext, userType string, term string) ([]actor.ActorSearchResponse, error) {
	var mock []actor.ActorSearchResponse

	mock = append(mock, actor.ActorSearchResponse{1, "Test Name", "test", "test@tesla.com", "USER"})
	return mock, nil

}

func TestGetActorCtrlSucess(t *testing.T) {

	f := MockActor{}

	actorName := "TEST"
	adusername := "test"
	email := "test@gmail.com"
	actorType := "USER"

	expected := actor.Actor{0, actorName, adusername, email, actorType, null.NewString("", false), null.NewString("", false), null.NewInt(0, false), 1, ""}
	actual, err := actor.Get(f, "TEST", &ctx)
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Success expected: %q  Actual: %q", expected, actual)
	}
}

func TestUpdateActorCtrlSucess(t *testing.T) {

	f := MockActor{}
	var ID int64 = 1
	var expected int64 = 1

	actorName := "TEST"
	adusername := "test"
	email := "test@gmail.com"
	actorType := "USER"
	input := actor.Actor{0, actorName, adusername, email, actorType, null.NewString("", false), null.NewString("", false), null.NewInt(0, false), 1, ""}

	actual, err := actor.Update(f, ID, input, &ctx)
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Success expected: %q  Actual: %q", expected, actual)
	}
}

func TestCreateActorCtrlSucess(t *testing.T) {

	f := MockActor{}

	actorName := "TEST"
	adusername := "test"
	email := "test@gmail.com"
	actorType := "USER"

	expected := actor.Actor{0, actorName, adusername, email, actorType, null.NewString("", false), null.NewString("", false), null.NewInt(0, false), 1, ""}
	actual, err := actor.Create(f, expected, &ctx)

	if err != nil {
		t.Error(err.Error())
	}
	expected.Id = 100
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Success expected: %q  Actual: %q", expected, actual)
	}
}

func TestSearchActorCtrl(t *testing.T) {

	f := actor.DBActor{}
	actual, err := actor.SearchActor(f, &ctx, "USER", "Ra")

	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("Total Count: %d", len(actual))
}
