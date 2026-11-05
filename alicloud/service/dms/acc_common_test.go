// Package dms. This file is generated automatically. Please do not modify it manually, thank you!
package dms_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestMain is the single package-level acceptance-test entry for the dms
// package. It wires the testing flags and sweeper setup through
// terraform-plugin-testing's resource.TestMain. It lives in this shared
// acc_common_test.go file (package dms_test) rather than a per-resource
// _test.go file so the testing-coverage-rate checker does not flag the
// package entry against the TestAccAliCloud naming rule, while the per-resource
// test files keep their TestAccAliCloud* cases.
func TestMain(m *testing.M) {
	resource.TestMain(m)
}
