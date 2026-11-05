package ims_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestMain wires the package into the sweep infrastructure so the orphan
// sweepers registered in the package's test files can be invoked and so the
// shared -sweep flag is registered exactly once for the ims test binary.
func TestMain(m *testing.M) {
	resource.TestMain(m)
}
