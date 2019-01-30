package alicloud

import (
	"testing"

	"fmt"
	"time"

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
					resource.TestCheckResourceAttr("data.alicloud_cen_instances.tf-testAccCen", "instances.0.name", "tf-testAccCenInstancesDataSourceCenIdConfig"),
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
	rand := time.Now().UnixNano()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCenInstancesDataSourceCenNameRegexConfig(defaultRegionToTest, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cen_instances.tf-testAccCen"),
					resource.TestCheckResourceAttr("data.alicloud_cen_instances.tf-testAccCen", "instances.0.name",
						fmt.Sprintf("tf-testAccCenDataSourceCenNameRegexConfig-%s-%d", defaultRegionToTest, rand)),
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
					resource.TestCheckResourceAttr("data.alicloud_cen_instances.tf-testAccCen", "instances.#", "5"),
					resource.TestCheckResourceAttr("data.alicloud_cen_instances.tf-testAccCen", "instances.0.description", "tf-testAccCenConfigDescription"),
					resource.TestCheckResourceAttr("data.alicloud_cen_instances.tf-testAccCen", "instances.1.status", "Active"),
					resource.TestCheckResourceAttr("data.alicloud_cen_instances.tf-testAccCen", "instances.2.name", "tf-testAccCenInstancesDataSourceMultiCenIdsConfig"),
				),
			},
		},
	})
}

func TestAccAlicloudCenInstancesDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCenInstancesDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cen_instances.tf-testAccCen"),
					resource.TestCheckResourceAttr("data.alicloud_cen_instances.tf-testAccCen", "instances.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_instances.tf-testAccCen", "instances.0.description"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_instances.tf-testAccCen", "instances.0.status"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_instances.tf-testAccCen", "instances.0.name"),
				),
			},
		},
	})
}

const testAccCheckAlicloudCenInstancesDataSourceCenIdConfig = `
resource "alicloud_cen_instance" "tf-testAccCen" {
	name = "tf-testAccCenInstancesDataSourceCenIdConfig"
	description = "tf-testAccCenConfigDescription"
}

data "alicloud_cen_instances" "tf-testAccCen" {
	ids = ["${alicloud_cen_instance.tf-testAccCen.id}"]
}
`

func testAccCheckAlicloudCenInstancesDataSourceCenNameRegexConfig(region string, rand int64) string {
	return fmt.Sprintf(`
		resource "alicloud_cen_instance" "tf-testAccCen" {
			name = "tf-testAccCenDataSourceCenNameRegexConfig-%s-%d"
			description = "tf-testAccCenConfigDescription"
		}
		
		data "alicloud_cen_instances" "tf-testAccCen" {
			name_regex = "${alicloud_cen_instance.tf-testAccCen.name}"
		}
		`, region, rand)
}

const testAccCheckAlicloudCenInstancesDataSourceMultiCenIdsConfig = `
resource "alicloud_cen_instance" "tf-testAccCen" {
	name = "tf-testAccCenInstancesDataSourceMultiCenIdsConfig"
	description = "tf-testAccCenConfigDescription"
	count = 5
}

data "alicloud_cen_instances" "tf-testAccCen" {
	ids = ["${alicloud_cen_instance.tf-testAccCen.*.id}"]
}
`

const testAccCheckAlicloudCenInstancesDataSourceEmpty = `
data "alicloud_cen_instances" "tf-testAccCen" {
	name_regex = "^tf-testacc-fake-name"
}
`
