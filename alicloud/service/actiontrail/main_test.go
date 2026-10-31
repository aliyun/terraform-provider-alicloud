// Package actiontrail_test.
package actiontrail_test

import (
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/recording"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMain(m *testing.M) {
	if recording.Enabled() {
		os.Exit(recording.RunMain(m))
	}
	resource.TestMain(m)
}
