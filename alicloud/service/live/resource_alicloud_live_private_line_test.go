package live_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/live"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAliCloudLivePrivateLine_basic0 covers the reuse="no" lifecycle: the
// API creates a GA acceleration channel and binds it to a live stream, then
// terraform destroy unbinds it. The domain must already be a live domain in
// the test account (ALICLOUD_LIVE_DOMAIN, default demo.aliyundoc.com), and the
// account must have GA PostPay open with the SLR permission required by
// CreateLivePrivateLine.
func TestAccAliCloudLivePrivateLine_basic0(t *testing.T) {
	resourceId := "alicloud_live_private_line.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	domain := os.Getenv("ALICLOUD_LIVE_DOMAIN")
	if domain == "" {
		domain = "demo.aliyundoc.com"
	}
	appName := fmt.Sprintf("tf-testacc-app%d", rand)
	streamName := fmt.Sprintf("tf-testacc-stream%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckLivePrivateLineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudLivePrivateLineBasic0(domain, appName, streamName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "domain_name", domain),
					resource.TestCheckResourceAttr(resourceId, "app_name", appName),
					resource.TestCheckResourceAttr(resourceId, "stream_name", streamName),
					resource.TestCheckResourceAttr(resourceId, "video_center", "cn-shanghai"),
					resource.TestCheckResourceAttr(resourceId, "acceleration_type", "play"),
					resource.TestCheckResourceAttr(resourceId, "acceleration_area", "ap-southeast-1"),
					resource.TestCheckResourceAttr(resourceId, "reuse", "no"),
					resource.TestCheckResourceAttr(resourceId, "max_bandwidth", "20"),
					resource.TestCheckResourceAttrSet(resourceId, "instance_id"),
					resource.TestCheckResourceAttrSet(resourceId, "status"),
					resource.TestCheckResourceAttrSet(resourceId, "binding_status"),
					resource.TestCheckResourceAttrSet(resourceId, "ip"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// reuse and max_bandwidth are write-only: the Describe entry never
				// echoes them back, so they cannot be reconstructed after import.
				ImportStateVerifyIgnore: []string{"reuse", "max_bandwidth"},
			},
		},
	})
}

func testAccCheckLivePrivateLineDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_live_private_line" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		_, err := live.NewService(client).DescribeLivePrivateLine(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing live private line %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("live private line %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccAliCloudLivePrivateLineBasic0(domain, appName, streamName string) string {
	return fmt.Sprintf(`
variable "domain_name" {
  default = "%s"
}

variable "app_name" {
  default = "%s"
}

variable "stream_name" {
  default = "%s"
}

resource "alicloud_live_private_line" "default" {
  domain_name       = var.domain_name
  app_name          = var.app_name
  stream_name       = var.stream_name
  video_center      = "cn-shanghai"
  acceleration_type = "play"
  acceleration_area = "ap-southeast-1"
  reuse             = "no"
  max_bandwidth     = 20
}
`, domain, appName, streamName)
}
