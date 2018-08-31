package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCensDataSource_cen_id(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCensDataSourceCenIdConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cens.cen"),
					resource.TestCheckResourceAttr("data.alicloud_cens.cen", "cens.0.name", "terraformTestAccName"),
					resource.TestCheckResourceAttr("data.alicloud_cens.cen", "cens.0.description", "terraform test"),
				),
			},
		},
	})
}

func TestAccAlicloudCensDataSource_cen_nameRegex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCensDataSourceCenNameRegexConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cens.cen"),
					resource.TestCheckResourceAttr("data.alicloud_cens.cen", "cens.0.name", "terraformTestAccName"),
					resource.TestCheckResourceAttr("data.alicloud_cens.cen", "cens.0.description", "terraform test"),
				),
			},
		},
	})
}

const testAccCheckAlicloudCensDataSourceCenIdConfig = `
resource "alicloud_cen" "cen" {
	name = "terraformTestAccName"
	description = "terraform test"
}

data "alicloud_cens" "cen" {
  cen_ids = ["${alicloud_cen.cen.id}"]
}
`

const testAccCheckAlicloudCensDataSourceCenNameRegexConfig = `
resource "alicloud_cen" "cen" {
	name = "terraformTestAccName"
	description = "terraform test"
}

data "alicloud_cens" "cen" {
	name_regex = "${alicloud_cen.cen.name}"
}
`
