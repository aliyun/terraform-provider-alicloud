// Package privatezone_test is the external test package for the privatezone
// service. TestMain delegates to resource.TestMain so the sweepers registered
// by init() in this package's *_test.go files are distributed to the test
// binary; without TestMain the plugin-testing harness skips sweeper dispatch
// for external service test packages, so leaked tf-test* endpoints are never
// cleaned up.
package privatezone_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
