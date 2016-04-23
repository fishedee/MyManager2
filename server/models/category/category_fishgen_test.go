package category

import (
	. "github.com/fishedee/web"
	"testing"
)

type testFishGenStruct struct{}

func TestCategory(t *testing.T) {
	RunBeegoValidateTest(t, &testFishGenStruct{})
}
