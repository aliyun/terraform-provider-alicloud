package hbr_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudHbrServiceDataSource(t *testing.T) {
	resourceId := "data.alicloud_hbr_service.current"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudHbrServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     acctest.CHECKSET,
						"enable": "On",
						"status": "Opened",
					}),
				),
			},
			{
				Config: testAccCheckAlicloudHbrServiceDataSourceOff,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     acctest.CHECKSET,
						"enable": "Off",
						"status": "",
					}),
				),
			},
		},
	})
}

const testAccCheckAlicloudHbrServiceDataSource = `
data "alicloud_hbr_service" "current" {
	enable = "On"
}
`

const testAccCheckAlicloudHbrServiceDataSourceOff = `
data "alicloud_hbr_service" "current" {
	enable = "Off"
}
`
