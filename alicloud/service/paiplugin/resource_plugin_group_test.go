package paiplugin_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/paiplugin"
	tfacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccAliCloudPaiPluginGroup covers the Create -> Read -> Delete lifecycle
// of alicloud_pai_plugin_group and the PluginGroupId import path. A group has
// no in-place Update, so any argument change is exercised as a replacement.
//
// This test drives the real muxed provider through acctest.ProtoV5ProviderFactories,
// which (per alicloud/acctest) works on either the SDK v2 or the plugin-testing
// test harness. We use the SDK v2 helper/resource harness on purpose, not the
// plugin-testing one: a non-test file in package alicloud currently imports
// terraform-plugin-sdk/v2/helper/resource, so any test binary that also pulls
// terraform-plugin-testing/helper/resource links both packages and dies at init
// with "flag redefined: sweep" before any test runs. The SDK v2 harness links
// that package only once and avoids the collision; the framework resource is
// still exercised end-to-end through the muxed provider server.
//
// NOTE: the remote ACC + ci-check-quick pre-PR check for this function is
// pending; the create input shape (source=0, a text body) follows the
// CreateGroup 2022-01-12 contract and is verified by the recovery run, not
// by a local Go build.
func TestAccAliCloudPaiPluginGroup(t *testing.T) {
	r := tfacctest.RandIntRange(100000, 999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckAlicloudPaiPluginGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudPaiPluginGroupConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_plugin_group.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_plugin_group.default", "plugin_group_id"),
					resource.TestCheckResourceAttr("alicloud_pai_plugin_group.default", "name", fmt.Sprintf("tf-testacc-pai-plugin-group-%d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_plugin_group.default", "source", "0"),
					resource.TestCheckResourceAttr("alicloud_pai_plugin_group.default", "remark", fmt.Sprintf("tf-testacc-%d", r)),
					resource.TestCheckResourceAttrSet("alicloud_pai_plugin_group.default", "phone_number"),
					resource.TestCheckResourceAttrSet("alicloud_pai_plugin_group.default", "created_time"),
				),
			},
			{
				ResourceName:            "alicloud_pai_plugin_group.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status", "updated_time", "amount", "campaign_id", "inference_job_id", "history"},
			},
		},
	})
}

func testAccAliCloudPaiPluginGroupConfig(r int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testacc-pai-plugin-group-%[1]d"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_pai_plugin_group" "default" {
  name   = var.name
  source = 0
  text   = "13900000000"
  remark = "tf-testacc-%[1]d"
}
`, r)
}

// testAccCheckAlicloudPaiPluginGroupDestroy verifies every group created by the
// test is gone after the destroy step. It reuses paiplugin.FindPluginGroup so
// the not-found contract (genuine HTTP 404 only, not routing/auth-layer 404s)
// matches the resource Read/Delete exactly; a real not-found yields nil data,
// while a still-existing group or any non-404 error fails the destroy check.
func testAccCheckAlicloudPaiPluginGroupDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_pai_plugin_group" {
			continue
		}
		data, err := paiplugin.FindPluginGroup(client, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("checking destroy of pai plugin group %s: %w", rs.Primary.ID, err)
		}
		if data != nil {
			return fmt.Errorf("pai plugin group %s still exists after destroy", rs.Primary.ID)
		}
	}
	return nil
}
