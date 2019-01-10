package alicloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccAlicloudLogtailToMachineGroupDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudLogTailConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudLogtailToMachineGroupDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_logtail_to_machine_group.example"),
					resource.TestCheckResourceAttr("data.alicloud_logtail_to_machine_group.example", "project", "tf-logproject1"),
					resource.TestCheckResourceAttr("data.alicloud_logtail_to_machine_group.example", "logtail_config.0", "evan-terraform-config"),
					resource.TestCheckResourceAttr("data.alicloud_logtail_to_machine_group.example", "machine_group.0", "evan-machine-group"),
				),
			},
		},
	})
}

const testAccCheckAlicloudLogtailToMachineGroupDataSource = `
data "alicloud_logtail_to_machine_group" "example" {
   project = "tf-logproject1"
   output_file = "~/newdata/map.json"
}
`
