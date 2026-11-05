package helper

import (
	"strings"
	"testing"
)

// Test copied from alicloud/common_test.go; test name is unchanged.
func TestUnitCommonBuildClientToken(t *testing.T) {
	token1 := BuildClientToken("CreateInstance")
	token2 := BuildClientToken("CreateInstance")

	if token1 == "" {
		t.Error("buildClientToken returned empty string")
	}

	if len(token1) > 64 {
		t.Error("buildClientToken returned token longer than 64 characters")
	}

	if token1 == token2 {
		t.Error("buildClientToken should return different tokens")
	}

	if !strings.Contains(token1, "TF-CreateInstance") {
		t.Error("buildClientToken should contain action name")
	}
}
