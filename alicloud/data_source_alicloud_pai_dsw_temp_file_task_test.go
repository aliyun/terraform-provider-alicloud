package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudPaiDswTempFileTaskDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-dstft%d", rand)
	// Build an env-based client via sharedClientForRegion so the PAI DSW
	// Instance fixture can be created BEFORE resource.Test runs. PreCheck
	// runs before the provider Configure, so testAccProvider.Meta() would
	// be nil inside PreCheck; sharedClientForRegion reads env vars directly.
	testAccPreCheck(t)
	rawClient, err := sharedClientForRegion(defaultRegionToTest)
	if err != nil {
		t.Fatalf("failed to build shared client for region %s: %s", defaultRegionToTest, err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	instanceId := createPaiDswTempFileTaskInstanceForTest(t, client, name)
	t.Cleanup(func() { deletePaiDswTempFileTaskInstanceForTest(client, instanceId) })

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAlicloudPaiDswTempFileTaskConfig(name, instanceId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_pai_dsw_temp_file_task.default", "instance_id", instanceId),
					resource.TestCheckResourceAttrSet("data.alicloud_pai_dsw_temp_file_task.default", "temp_file_task_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_pai_dsw_temp_file_task.default", "user_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_pai_dsw_temp_file_task.default", "owner_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_pai_dsw_temp_file_task.default", "create_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_pai_dsw_temp_file_task.default", "region_id"),
				),
			},
		},
	})
}

func testAccDataSourceAlicloudPaiDswTempFileTaskConfig(name string, instanceId string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_pai_dsw_temp_file_task" "default" {
  instance_id = "%s"
}

data "alicloud_pai_dsw_temp_file_task" "default" {
  temp_file_task_id = "${alicloud_pai_dsw_temp_file_task.default.id}"
}
`, name, instanceId)
}
