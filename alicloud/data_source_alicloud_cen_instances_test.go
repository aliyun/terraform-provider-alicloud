package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCenInstancesDataSource_cen_id(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCenInstancesDataSourceCenIdConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cen_instances.tf-testAccCen"),
					resource.TestCheckResourceAttr("data.alicloud_cen_instances.tf-testAccCen", "instances.0.name", "tf-testAccCenConfig"),
					resource.TestCheckResourceAttr("data.alicloud_cen_instances.tf-testAccCen", "instances.0.description", "tf-testAccCenConfigDescription"),
					resource.TestCheckResourceAttr("data.alicloud_cen_instances.tf-testAccCen", "instances.0.status", "Active"),
					resource.TestCheckResourceAttrSet("data.alicloud_cen_instances.tf-testAccCen", "instances.0.id"),
				),
			},
		},
	})
}

func TestAccAlicloudCenInstancesDataSource_name_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCenInstancesDataSourceCenNameRegexConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cen_instances.tf-testAccCen"),
					resource.TestCheckResourceAttr("data.alicloud_cen_instances.tf-testAccCen", "instances.0.name", "tf-testAccCenConfig"),
					resource.TestCheckResourceAttr("data.alicloud_cen_instances.tf-testAccCen", "instances.0.description", "tf-testAccCenConfigDescription"),
					resource.TestCheckResourceAttr("data.alicloud_cen_instances.tf-testAccCen", "instances.0.status", "Active"),
					resource.TestCheckResourceAttrSet("data.alicloud_cen_instances.tf-testAccCen", "instances.0.id"),
				),
			},
		},
	})
}

const testAccCheckAlicloudCenInstancesDataSourceCenIdConfig = `
resource "alicloud_cen_instance" "tf-testAccCen" {
	name = "tf-testAccCenConfig"
	description = "tf-testAccCenConfigDescription"
}

data "alicloud_cen_instances" "tf-testAccCen" {
	ids = ["${alicloud_cen_instance.tf-testAccCen.id}"]
}
`

const testAccCheckAlicloudCenInstancesDataSourceCenNameRegexConfig = `
resource "alicloud_cen_instance" "tf-testAccCen" {
	name = "tf-testAccCenConfig"
	description = "tf-testAccCenConfigDescription"
}

data "alicloud_cen_instances" "tf-testAccCen" {
	name_regex = "${alicloud_cen_instance.tf-testAccCen.name}"
}
`