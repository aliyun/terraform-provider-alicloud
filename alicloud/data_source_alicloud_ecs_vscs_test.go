package alicloud

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccAliCloudEcsVscsDataSource creates a Vsc and queries the
// alicloud_ecs_vscs data source filtered by instance_id, asserting the
// created Vsc is listed. Requires TF_ECS_VSC_INSTANCE_ID.
func TestAccAliCloudEcsVscsDataSource(t *testing.T) {
	if v := os.Getenv("TF_ECS_VSC_INSTANCE_ID"); v == "" {
		t.Skip("TF_ECS_VSC_INSTANCE_ID is not set, skip TestAccAliCloudEcsVscsDataSource")
	}
	instanceId := os.Getenv("TF_ECS_VSC_INSTANCE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEcsVscDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsVscsDataSourceConfig(instanceId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_ecs_vscs.default", "vscs.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_ecs_vscs.default", "vscs.0.vsc_id"),
					resource.TestCheckResourceAttr("data.alicloud_ecs_vscs.default", "vscs.0.instance_id", instanceId),
				),
			},
		},
	})
}

func testAccEcsVscsDataSourceConfig(instanceId string) string {
	return `
resource "alicloud_ecs_vsc" "default" {
  instance_id = "` + instanceId + `"
  vsc_name    = "tf-testacc-vsc-ds"
  description = "data source test"
  vsc_type    = "Primary"
}

data "alicloud_ecs_vscs" "default" {
  instance_id = alicloud_ecs_vsc.default.instance_id
  ids         = [alicloud_ecs_vsc.default.vsc_id]
}
`
}
