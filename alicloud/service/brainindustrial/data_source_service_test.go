package brainindustrial_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudBrainIndustrialServiceDataSource(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.BrainIndustrialSupportRegions)
	resourceId := "data.alicloud_brain_industrial_service.current"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudBrainIndustrialServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     acctest.CHECKSET,
						"status": "Opened",
						"enable": "On",
					}),
				),
			},
			{
				Config: testAccCheckAliCloudBrainIndustrialServiceDataSourceDisable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     "BrainIndustrialServiceHasNotBeenOpened",
						"enable": "Off",
					}),
				),
			},
		},
	})
}

const testAccCheckAliCloudBrainIndustrialServiceDataSource = `
data "alicloud_brain_industrial_service" "current" {
	enable = "On"
}
`

const testAccCheckAliCloudBrainIndustrialServiceDataSourceDisable = `
data "alicloud_brain_industrial_service" "current" {
	enable = "Off"
}
`
