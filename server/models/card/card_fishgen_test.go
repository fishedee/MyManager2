package card

import (
	. "github.com/fishedee/web"
	"testing"
)

type testFishGenStruct struct{}

func TestCard(t *testing.T) {
	RunBeegoValidateTest(t, &testFishGenStruct{})
}
