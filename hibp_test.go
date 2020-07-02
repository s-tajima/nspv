package nspv

import (
	"reflect"
	"testing"
)

func TestHibpClient(t *testing.T) {
	cases := []struct {
		password string
		result   bool
	}{
		{"password", true},
		{"D24932AF-E30F-46B0-8570-DF69B82B0BF8", false},
	}

	for _, tt := range cases {
		hc := newHibpClient()
		result, _ := hc.pwnedCount(tt.password)

		if result > 0 != tt.result {
			t.Errorf("for %s ... got: %t, want: %t", tt.password, result > 0, tt.result)
		}
	}
}

func TestParseRange(t *testing.T) {
	cases := []struct {
		body   string
		hashes map[string]int
	}{
		{"000000:1", map[string]int{"000000": 1}},
	}

	for _, tt := range cases {
		hc := newHibpClient()
		result, _ := hc.parseRange(tt.body)

		if !reflect.DeepEqual(result, tt.hashes) {
			t.Errorf("for %s ... got: %v, want: %v", tt.body, result, tt.hashes)
		}
	}
}
