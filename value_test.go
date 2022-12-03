package vert

import (
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"testing"
	"time"
)

const isoStringLength = len("yyyy-MM-ddThh:mm:ss")

func TestValueOfTime(t *testing.T) {
	now := time.Now()

	jsVal := ValueOf(now)

	// cut of the timezone information, because seconds precision is sufficient for this test
	jsValIsoString := jsVal.Call("toISOString").String()
	jsValIsoString = jsValIsoString[:isoStringLength]

	then.AssertThat(t, jsValIsoString, is.EqualTo(now.UTC().Format(time.RFC3339)[:isoStringLength]))
}
