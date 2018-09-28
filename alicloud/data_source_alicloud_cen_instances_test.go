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
					resource.TestCheckNoResourceAttr("data.alicloud_cen_instances.tf-testAccCen", "instances.0.child_instance_ids"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_instances.tf-testAccCen", "instances.0.bandwidth_package_ids"),
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
					resource.TestCheckNoResourceAttr("data.alicloud_cen_instances.tf-testAccCen", "instances.0.child_instance_ids"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_instances.tf-testAccCen", "instances.0.bandwidth_package_ids"),
				),
			},
		},
	})
}

func TestAccAlicloudCenInstancesDataSource_multi_cen_ids(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCenInstancesDataSourceMultiCenIdsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cen_instances.tf-testAccCen"),
					resource.TestCheckResourceAttrSet("data.alicloud_cen_instances.tf-testAccCen", "instances.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_cen_instances.tf-testAccCen", "instances.1.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_cen_instances.tf-testAccCen", "instances.2.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_cen_instances.tf-testAccCen", "instances.3.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_cen_instances.tf-testAccCen", "instances.4.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_cen_instances.tf-testAccCen", "instances.5.id"),
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

const testAccCheckAlicloudCenInstancesDataSourceMultiCenIdsConfig = `
resource "alicloud_cen_instance" "tf-testAccCen1" {
	name = "tf-testAccCenConfig1"
	description = "tf-testAccCenConfigDescription"
}
resource "alicloud_cen_instance" "tf-testAccCen2" {
	name = "tf-testAccCenConfig2"
	description = "tf-testAccCenConfigDescription"
}
resource "alicloud_cen_instance" "tf-testAccCen3" {
	name = "tf-testAccCenConfig3"
	description = "tf-testAccCenConfigDescription"
}
resource "alicloud_cen_instance" "tf-testAccCen4" {
	name = "tf-testAccCenConfig4"
	description = "tf-testAccCenConfigDescription"
}
resource "alicloud_cen_instance" "tf-testAccCen5" {
	name = "tf-testAccCenConfig5"
	description = "tf-testAccCenConfigDescription"
}
resource "alicloud_cen_instance" "tf-testAccCen6" {
	name = "tf-testAccCenConfig6"
	description = "tf-testAccCenConfigDescription"
}

data "alicloud_cen_instances" "tf-testAccCen" {
	ids = [
           "${alicloud_cen_instance.tf-testAccCen1.id}",
           "${alicloud_cen_instance.tf-testAccCen2.id}",
           "${alicloud_cen_instance.tf-testAccCen3.id}",
           "${alicloud_cen_instance.tf-testAccCen4.id}",
           "${alicloud_cen_instance.tf-testAccCen5.id}",
           "${alicloud_cen_instance.tf-testAccCen6.id}",
           ]
}
`
