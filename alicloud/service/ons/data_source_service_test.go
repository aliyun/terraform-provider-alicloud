package ons_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudOnsServiceDataSource exercises the Optional, non-ForceNew
// `enable` argument that the split-layout coverage checker requires for the
// alicloud_ons_service data source: step 1 opens the service (enable=On,
// status=Opened) and step 2 turns it back off (enable=Off, status=""), so
// `enable` appears in both the must-set and the modify sets. `status` is
// Computed; on the Off step the read sets it to "" and the assertion is
// updated accordingly, because the shared check map accumulates across steps.
func TestAccAliCloudOnsServiceDataSource(t *testing.T) {
	resourceId := "data.alicloud_ons_service.current"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudOnsServiceDataSourceOn,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     CHECKSET,
						"status": "Opened",
						"enable": "On",
					}),
				),
			},
			{
				Config: testAccCheckAlicloudOnsServiceDataSourceOff,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     CHECKSET,
						"enable": "Off",
						"status": "",
					}),
				),
			},
		},
	})
}

const testAccCheckAlicloudOnsServiceDataSourceOn = `
data "alicloud_ons_service" "current" {
	enable = "On"
}
`

const testAccCheckAlicloudOnsServiceDataSourceOff = `
data "alicloud_ons_service" "current" {
	enable = "Off"
}
`
