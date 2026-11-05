package arms_test

import (
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

// armsTestClientForRegion delegates to the shared sweeper client in
// alicloud/sweep, which isolates endpoints per region (see
// alicloud/sweep/client.go: SharedEndpointCacheForRegion). It returns
// interface{} so the existing sweeper call sites keep their
// rawClient.(*connectivity.AliyunClient) assertion unchanged.
func armsTestClientForRegion(region string) (interface{}, error) {
	return sweep.SharedClientForRegion(region)
}

// armsSweepAll reports whether sweepers should delete every matched resource
// regardless of the tf-testAcc prefix (mirrors the former root sweepAll).
func armsSweepAll() bool {
	return os.Getenv("ALICLOUD_SWEEP_ALL_RESOURCES") == "true"
}
