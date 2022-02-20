package require

import (
	"encoding/json"
	"reflect"
	"testing"
)

func JSONEqual(t *testing.T, s1, s2 string) {
	var o1, o2 interface{}

	if err := json.Unmarshal([]byte(s1), &o1); err != nil {
		t.Fatalf("error mashalling string 1: %s", err.Error())
	}

	if err := json.Unmarshal([]byte(s2), &o2); err != nil {
		t.Fatalf("error mashalling string 2: %s", err.Error())
	}

	if !reflect.DeepEqual(o1, o2) {
		t.Fatalf("got json %v, expected %v", o1, o2)
	}
}
