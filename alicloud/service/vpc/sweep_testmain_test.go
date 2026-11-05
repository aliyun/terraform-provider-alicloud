package vpc_test

import (
	"flag"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestMain drives the VPC test package. The alicloud_vpc multi-region sweep
// orchestration (runVpcSweep / testSweepVpcs) moved here together with the VPC
// resource tests; the root alicloud TestMain used to host this special path but
// the symbols are no longer visible there once the tests left the root package.
//
// The SDK matches -sweep-run by substring, so alicloud_vpc would also select
// alicloud_vpc_* sweepers with region-wide deletion logic. Invoke exactly the
// registered VPC sweeper callback for VPC cleanup, mirroring the prior root
// behaviour, now scoped to the VPC test binary.
func TestMain(m *testing.M) {
	flag.Parse()
	if flag.Lookup("sweep").Value.String() != "" && flag.Lookup("sweep-run").Value.String() == "alicloud_vpc" {
		if err := runVpcSweep(flag.Lookup("sweep").Value.String(), flag.Lookup("sweep-run").Value.String(), testSweepVpcs); err != nil {
			log.Print(err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	resource.TestMain(m)
}
