package resourcemanager_test

import (
	"os"
	"strings"
	"testing"
)

// testAccPreCheckWithResourceManagerAccount mirrors the root-package helper of
// the same name (alicloud/provider_test.go). It is reproduced here because the
// migrated external test package (resourcemanager_test) cannot reach the root
// package's unexported test helpers, and the resource manager account tests
// rely on the RESOURCE_MANAGER_* credential override semantics.
func testAccPreCheckWithResourceManagerAccount(t *testing.T) {
	if os.Getenv("RESOURCE_MANAGER_ALICLOUD_ACCESS_KEY") == "" || os.Getenv("RESOURCE_MANAGER_ALICLOUD_SECRET_KEY") == "" {
		t.Skip("Skipping the test case with RESOURCE_MANAGER_ALICLOUD_ACCESS_KEY or RESOURCE_MANAGER_ALICLOUD_SECRET_KEY is empty")
		t.Skipped()
	} else {
		os.Setenv("ALICLOUD_ACCESS_KEY", os.Getenv("RESOURCE_MANAGER_ALICLOUD_ACCESS_KEY"))
		os.Setenv("ALICLOUD_SECRET_KEY", os.Getenv("RESOURCE_MANAGER_ALICLOUD_SECRET_KEY"))
	}
}

// testAccPreCheckWithResourceManagerHandshakesSetting mirrors the root-package
// helper of the same name (alicloud/provider_test.go). It skips the handshake
// tests when INVITED_ALICLOUD_ACCOUNT_ID is not configured.
func testAccPreCheckWithResourceManagerHandshakesSetting(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("INVITED_ALICLOUD_ACCOUNT_ID")); v == "" {
		t.Skipf("Skipping the test case with there is no \"INVITED_ALICLOUD_ACCOUNT_ID\" setting")
		t.Skipped()
	}
}
