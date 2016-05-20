package stateful

import (
	"regexp"
	"testing"
	"time"
)

func Test_valueTypeOf(t *testing.T) {
	type expectation struct {
		value     interface{}
		valueType ValueType
	}

	expectations := []expectation{
		{value: float64(0), valueType: TFloat64},
		{value: int64(0), valueType: TInt64},
		{value: "Kapacitor Rulz", valueType: TString},
		{value: true, valueType: TBool},
		{value: regexp.MustCompile("\\d"), valueType: TRegex},
		{value: time.Duration(5), valueType: TDuration},
		{value: time.Time{}, valueType: TTime},
		{value: t, valueType: InvalidType},
	}

	for _, expect := range expectations {
		result := valueTypeOf(expect.value)

		if result != expect.valueType {
			t.Errorf("Got unexpected result for valueTypeOf(%T):\ngot: %s\nexpected: %s", expect.value, result, expect.valueType)
		}

	}
}
