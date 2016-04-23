package user

import (
	. "github.com/fishedee/web"
	"testing"
)

type testFishGenStruct struct{}

func TestUser(t *testing.T) {
	RunBeegoValidateTest(t, &testFishGenStruct{})
}
