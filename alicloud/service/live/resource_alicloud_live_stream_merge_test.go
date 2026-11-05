// Package live_test.
package live_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/live"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAliCloudLiveStreamMerge_basic covers create, an in-place update that
// re-issues the idempotent AddLiveStreamMerge, and import.
//
// The resource targets an existing live play domain registered to the test
// account; it cannot create one itself. The domain is taken from the
// ALICLOUD_LIVE_DOMAIN env var and defaults to the shared "demo.alicloud.com"
// used by the live domain test suite.
func TestAccAliCloudLiveStreamMerge_basic(t *testing.T) {
	domain := os.Getenv("ALICLOUD_LIVE_DOMAIN")
	if domain == "" {
		domain = "demo.alicloud.com"
	}
	// Fixed start/end window in the near future (UTC, ISO8601) shared by every
	// step so that ImportStateVerify matches the value stored at create time.
	startTime := time.Now().Add(5 * time.Minute).UTC().Format("2006-01-02T15:04:05Z")
	endTime := time.Now().Add(1 * time.Hour).Add(5 * time.Minute).UTC().Format("2006-01-02T15:04:05Z")

	resourceId := "alicloud_live_stream_merge.default"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shanghai"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckLiveStreamMergeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudLiveStreamMergeBasic(domain, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "domain_name", domain),
					resource.TestCheckResourceAttr(resourceId, "app_name", "liveApp"),
					resource.TestCheckResourceAttr(resourceId, "stream_name", "liveStream"),
					resource.TestCheckResourceAttr(resourceId, "in_app_name_1", "liveApp1"),
					resource.TestCheckResourceAttr(resourceId, "in_stream_name_1", "liveStream1"),
					resource.TestCheckResourceAttr(resourceId, "in_app_name_2", "liveApp2"),
					resource.TestCheckResourceAttr(resourceId, "in_stream_name_2", "liveStream2"),
					resource.TestCheckResourceAttr(resourceId, "protocol", "rtmp"),
					resource.TestCheckResourceAttr(resourceId, "live_merger", "off"),
					resource.TestCheckResourceAttr(resourceId, "switch_mode", "manually"),
					resource.TestCheckResourceAttrSet(resourceId, "id"),
				),
			},
			{
				Config: testAccAliCloudLiveStreamMergeUpdate(domain, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "switch_mode", "automatically"),
					resource.TestCheckResourceAttr(resourceId, "merge_parameters", "ali_low_frame_rate_threshold=10&ali_max_no_frame_timeout=5&block_all_jitter=0"),
					resource.TestCheckResourceAttr(resourceId, "live_merger", "on"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// select_app_name/select_stream_name are write-only (Describe omits
				// them); app_using/stream_using/extra_in_app_streams reflect
				// volatile runtime state that can change between apply and import.
				ImportStateVerifyIgnore: []string{"select_app_name", "select_stream_name", "app_using", "stream_using", "extra_in_app_streams"},
			},
		},
	})
}

func testAccCheckLiveStreamMergeDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_live_stream_merge" {
			continue
		}
		parts := strings.Split(rs.Primary.ID, ":")
		if len(parts) != 3 {
			return fmt.Errorf("unexpected resource id %q: expected DomainName:AppName:StreamName", rs.Primary.ID)
		}
		domainName, appName, streamName := parts[0], parts[1], parts[2]
		_, err := live.NewService(client).DescribeLiveStreamMerge(domainName, appName, streamName, "")
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing live stream merge %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("live stream merge %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccAliCloudLiveStreamMergeBasic(domain, startTime, endTime string) string {
	return fmt.Sprintf(`
variable "domain" {
  default = "%[1]s"
}

provider "alicloud" {
  region = "cn-shanghai"
}

resource "alicloud_live_stream_merge" "default" {
  domain_name      = var.domain
  app_name         = "liveApp"
  stream_name      = "liveStream"
  in_app_name_1    = "liveApp1"
  in_stream_name_1 = "liveStream1"
  in_app_name_2    = "liveApp2"
  in_stream_name_2 = "liveStream2"
  start_time       = "%[2]s"
  end_time         = "%[3]s"
  protocol         = "rtmp"
  live_merger      = "off"
  switch_mode      = "manually"
}
`, domain, startTime, endTime)
}

func testAccAliCloudLiveStreamMergeUpdate(domain, startTime, endTime string) string {
	return fmt.Sprintf(`
variable "domain" {
  default = "%[1]s"
}

provider "alicloud" {
  region = "cn-shanghai"
}

resource "alicloud_live_stream_merge" "default" {
  domain_name      = var.domain
  app_name         = "liveApp"
  stream_name      = "liveStream"
  in_app_name_1    = "liveApp1"
  in_stream_name_1 = "liveStream1"
  in_app_name_2    = "liveApp2"
  in_stream_name_2 = "liveStream2"
  start_time       = "%[2]s"
  end_time         = "%[3]s"
  protocol         = "rtmp"
  live_merger      = "on"
  merge_parameters = "ali_low_frame_rate_threshold=10&ali_max_no_frame_timeout=5&block_all_jitter=0"
  switch_mode      = "automatically"
}
`, domain, startTime, endTime)
}
