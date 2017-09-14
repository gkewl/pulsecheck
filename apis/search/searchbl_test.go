package search_test

import (
	"testing"

	"github.com/gkewl/pulsecheck/apis/search"
	"github.com/gkewl/pulsecheck/model"
)

var token model.TokenInfo

func TestSearch_Actor_Bl(t *testing.T) {

	di := search.DBSearch{"sspade", 1}
	actual, err := di.Search("actor", "raj", &ctx)

	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(actual) == 0 {
		t.Errorf("expected actor data and returned empty data")
	}
}
