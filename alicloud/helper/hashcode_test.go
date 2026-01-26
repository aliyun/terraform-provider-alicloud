package helper

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
)

// TestHashcode compare local hash with terraform-plugin-sdk v1 hash
func TestHashcode(t *testing.T) {
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
