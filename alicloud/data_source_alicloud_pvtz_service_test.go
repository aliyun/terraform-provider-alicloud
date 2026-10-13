package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAliCloudPvtzServiceDataSource(t *testing.T) {
	resourceId := "data.alicloud_pvtz_service.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudPvtzServiceDataSourceNil,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     "PvtzServiceHasNotBeenOpened",
						"status": "",
					}),
				),
			},
			{
				Config: testAccCheckAliCloudPvtzServiceDataSourceWithOff,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     "PvtzServiceHasNotBeenOpened",
						"status": "",
					}),
				),
			},
			{
				Config: testAccCheckAliCloudPvtzServiceDataSourceWithOn,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     CHECKSET,
						"status": "Opened",
					}),
				),
			},
		},
	})
}

const testAccCheckAliCloudPvtzServiceDataSourceNil = `
	data "alicloud_pvtz_service" "default" {
	}
`

const testAccCheckAliCloudPvtzServiceDataSourceWithOff = `
	data "alicloud_pvtz_service" "default" {
  		enable = "Off"
	}
`

const testAccCheckAliCloudPvtzServiceDataSourceWithOn = `
	data "alicloud_pvtz_service" "default" {
  		enable = "On"
	}
`
