package privatezone_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudPvtzServiceDataSource(t *testing.T) {
	resourceId := "data.alicloud_pvtz_service.default"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
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
						"id":     acctest.CHECKSET,
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
