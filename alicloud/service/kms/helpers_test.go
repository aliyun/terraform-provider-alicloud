package kms_test

import (
	"os"
	"strings"
	"testing"
)

// testAccPreCheckKMSForKeyIdImport skips import tests unless ALICLOUD_KMS_KEY_ID is set.
func testAccPreCheckKMSForKeyIdImport(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("ALICLOUD_KMS_KEY_ID")); v == "" {
		t.Skipf("Skipping tests without KEY_ID set.")
		t.Skipped()
	}
}
