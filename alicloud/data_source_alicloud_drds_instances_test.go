package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDRDSInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDRDSInstancesDataSourceConfig,
				Check:  resource.ComposeTestCheckFunc(
				//resource.TestCheckResourceAttr("data.alicloud_drds_instance.instance", "type", "1"),
				//resource.TestCheckResourceAttr("data.alicloud_drds_instance.instance", "description", "tf-testAccCheckAlicloudDRDSInstancesDataSourceConfig"),
				),
			},
		},
	})
}

const testAccCheckAlicloudDRDSInstancesDataSourceConfig = `

	data "alicloud_drds_instances" "drds_instances_ds" {
  		name_regex = "drds-\\d+"
  		region_id     = "cn-hangzhou"
	}
`
