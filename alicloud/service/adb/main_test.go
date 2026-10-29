package adb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestMain lets the resource test harness register the -sweep flag for an
// external _test package, mirroring the actiontrail migration pattern.
func TestMain(m *testing.M) { resource.TestMain(m) }
