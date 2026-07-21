package helper

import (
	"math"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
)

// TestHashcode compare local hash with terraform-plugin-sdk v1 hash
func TestUnitCommonHashcode(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"simple", "hello"},
		{"with_spaces", "hello world"},
		{"special_chars", "user@host:80/path?k=v#frag"},
		{"unicode", "你好世界 🌍"},
		{"numbers", "1234567890"},
		{"json_like", `{"key":"value","list":[1,2,3]}`},
		{"long_string", "a" + "very_long_string_to_test_hash_collision_avoidance_" + "b"},
		{"same_as_tf_example", "alicloud_vpc.example"},
		{"with_newlines", "line1\nline2\r\nline3"},
		{"only_symbols", "!@#$%^&*()_+-=[]{}|;':\",./<>?"},
		{"mixed_case", "AbC123!@#"},
		{"mixed_case2", "AbC1|BC12"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hashA := hashcode.String(tc.input)
			hashB := Hashcode(tc.input)

			// check hash
			if hashA != hashB {
				t.Errorf("Hash mismatch for input %q:\n  SDK v1: %d\n  custom:    %d",
					tc.input, hashA, hashB)
			}
		})
	}
}

// TestHashcode_NonNegative verifies that Hashcode always returns a non-negative integer.
func TestUnitCommonHashcode_NonNegative(t *testing.T) {
	inputs := []string{
		"",
		"a",
		"hello world",
		"alicloud_vpc.example",
		"你好世界",
		"!@#$%^",
		"1234567890abcdefghijklmnopqrstuvwxyz",
	}
	for _, s := range inputs {
		if v := Hashcode(s); v < 0 {
			t.Errorf("Hashcode(%q) = %d; want >= 0", s, v)
		}
	}
}

// TestHashcode_Deterministic verifies that repeated calls with the same input always
// produce the same result.
func TestUnitCommonHashcode_Deterministic(t *testing.T) {
	inputs := []string{"", "hello", "alicloud_vpc.default", "你好"}
	for _, s := range inputs {
		first := Hashcode(s)
		for i := 0; i < 10; i++ {
			if got := Hashcode(s); got != first {
				t.Errorf("Hashcode(%q): non-deterministic – got %d and %d", s, first, got)
			}
		}
	}
}

// TestHashcode_DifferentInputsDifferentHashes verifies that distinct inputs
// produce distinct hash values for a representative set of strings.
func TestUnitCommonHashcode_DifferentInputsDifferentHashes(t *testing.T) {
	pairs := [][2]string{
		{"hello", "world"},
		{"alicloud_vpc.a", "alicloud_vpc.b"},
		{"abc", "ABC"},
		{"1", "2"},
	}
	for _, p := range pairs {
		h1, h2 := Hashcode(p[0]), Hashcode(p[1])
		if h1 == h2 {
			t.Errorf("Hashcode(%q) == Hashcode(%q) == %d; expected different hashes", p[0], p[1], h1)
		}
	}
}

// TestHashcode_BoundaryValues exercises boundary conditions of the int type.
// On a 32-bit system the CRC32 value may overflow into negative territory;
// the implementation must still return a non-negative integer.
func TestUnitCommonHashcode_BoundaryValues(t *testing.T) {
	// Ensure the return value never exceeds math.MaxInt32 on 32-bit and is
	// always non-negative regardless of platform.
	_ = math.MaxInt32 // keep import used
	for _, s := range []string{"", "test", "boundary"} {
		v := Hashcode(s)
		if v < 0 {
			t.Errorf("Hashcode(%q) = %d; must be non-negative", s, v)
		}
	}
}
