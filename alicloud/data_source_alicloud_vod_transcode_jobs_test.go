package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudVodTranscodeJobs_basic0(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.VODSupportRegions)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svodtranscodejobs%d", defaultRegionToTest, rand)
	videoId := os.Getenv("ALICLOUD_VOD_VIDEO_ID")
	if videoId == "" {
		t.Skip("Skipping the test case without ALICLOUD_VOD_VIDEO_ID")
	}
	testAccConfig := resourceTestAccConfigFunc("data.alicloud_vod_transcode_jobs.default", name, func(name string) string {
		return fmt.Sprintf(`
variable "video_id" {
  default = "%s"
}
`, videoId)
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"video_id": "${var.video_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_vod_transcode_jobs.default", "ids.#"),
				),
			},
		},
	})
}
