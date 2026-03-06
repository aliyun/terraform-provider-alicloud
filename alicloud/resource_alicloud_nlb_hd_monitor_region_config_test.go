// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Nlb HdMonitorRegionConfig. >>> Resource test cases, automatically generated.
// Case hdMonitor 9477
func TestAccAliCloudNlbHdMonitorRegionConfig_basic9477(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_hd_monitor_region_config.default"
	ra := resourceAttrInit(resourceId, AlicloudNlbHdMonitorRegionConfigMap9477)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNlbHdMonitorRegionConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnlb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNlbHdMonitorRegionConfigBasicDependence9477)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"metric_store": "test",
					"log_project":  "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metric_store": "test",
						"log_project":  "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"metric_store": "test2",
					"log_project":  "test2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metric_store": "test2",
						"log_project":  "test2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudNlbHdMonitorRegionConfigMap9477 = map[string]string{
	"region_id": CHECKSET,
}

func AlicloudNlbHdMonitorRegionConfigBasicDependence9477(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Nlb HdMonitorRegionConfig. <<< Resource test cases, automatically generated.
