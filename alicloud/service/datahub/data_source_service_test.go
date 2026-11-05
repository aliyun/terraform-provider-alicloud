package datahub_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudDatahubServiceDataSource(t *testing.T) {
	resourceId := "data.alicloud_datahub_service.current"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDatahubServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     acctest.CHECKSET,
						"status": "Opened",
						"enable": "On",
					}),
				),
			},
			{
				Config: testAccCheckAlicloudDatahubServiceDataSourceOff,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     "DatahubServiceHasNotBeenOpened",
						"status": "",
						"enable": "Off",
					}),
				),
			},
		},
	})
}

const testAccCheckAlicloudDatahubServiceDataSource = `
data "alicloud_datahub_service" "current" {
	enable = "On"
}
`

const testAccCheckAlicloudDatahubServiceDataSourceOff = `
data "alicloud_datahub_service" "current" {
	enable = "Off"
}
`
