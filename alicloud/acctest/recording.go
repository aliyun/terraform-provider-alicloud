package acctest

import (
	"fmt"
	"os"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/recording"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

// Initialize before test functions generate names. PreCheck attaches finalization
// to the selected test, after Terraform teardown and CheckDestroy have finished.
// This package is test infrastructure; production Provider startup is unchanged.
func init() {
	if err := recording.StartProcess(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "recording initialization failed:", err)
		os.Exit(1)
	}
}

// RecordingRandInt is deterministic only in record/replay mode.
func RecordingRandInt() int { return recording.RandInt(sdkacctest.RandInt) }
