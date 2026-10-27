package acctest

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/recording"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// RecordingTestMain preserves the acceptance framework's sweeper handling and
// finalizes a process-scoped recording when explicitly enabled. A recording run
// must select exactly one anchored test name. Product TestMain functions may call
// this helper after migration; normal live testing is unchanged.
func RecordingTestMain(m *testing.M) {
	flag.Parse()
	if recording.Enabled() && flag.Lookup("sweep").Value.String() != "" {
		fmt.Fprintln(os.Stderr, "recording cannot be combined with resource sweeping")
		os.Exit(1)
	}
	resource.TestMain(recording.WrapMain(m))
}

// RecordingRandInt is deterministic only in record/replay mode.
func RecordingRandInt() int { return recording.RandInt(sdkacctest.RandInt) }
