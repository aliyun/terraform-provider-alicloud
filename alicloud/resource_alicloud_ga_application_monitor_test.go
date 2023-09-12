package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ga ApplicationMonitor. >>> Resource test cases, automatically generated.
// Case 4415
func TestAccAliCloudGaApplicationMonitor_basic4415(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_application_monitor.default"
	ra := resourceAttrInit(resourceId, AlicloudGaApplicationMonitorMap4415)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaApplicationMonitor")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgaapplicationmonitor%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaApplicationMonitorBasicDependence4415)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"address":        "www.baidu.com",
					"task_name":      "aaaa",
					"accelerator_id": "ga-bp1l8bs1z8gw0u6y1ag5g",
					"listener_id":    "lsr-bp1jovsyx377pboiel4p2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address":        "www.baidu.com",
						"task_name":      "aaaa",
						"accelerator_id": "ga-bp1l8bs1z8gw0u6y1ag5g",
						"listener_id":    "lsr-bp1jovsyx377pboiel4p2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address": "www.baidu.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address": "www.baidu.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"task_name": "aaaa",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"task_name": "aaaa",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"listener_id": "lsr-bp1jovsyx377pboiel4p2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"listener_id": "lsr-bp1jovsyx377pboiel4p2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"task_name": "bbbb",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"task_name": "bbbb",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"task_name": "aaaa",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"task_name": "aaaa",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address":        "www.baidu.com",
					"task_name":      "aaaa",
					"accelerator_id": "ga-bp1l8bs1z8gw0u6y1ag5g",
					"listener_id":    "lsr-bp1jovsyx377pboiel4p2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address":        "www.baidu.com",
						"task_name":      "aaaa",
						"accelerator_id": "ga-bp1l8bs1z8gw0u6y1ag5g",
						"listener_id":    "lsr-bp1jovsyx377pboiel4p2",
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

var AlicloudGaApplicationMonitorMap4415 = map[string]string{
	"status": CHECKSET,
}

func AlicloudGaApplicationMonitorBasicDependence4415(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 4415  twin
func TestAccAliCloudGaApplicationMonitor_basic4415_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_application_monitor.default"
	ra := resourceAttrInit(resourceId, AlicloudGaApplicationMonitorMap4415)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaApplicationMonitor")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgaapplicationmonitor%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaApplicationMonitorBasicDependence4415)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"address":        "www.baidu.com",
					"task_name":      "aaaa",
					"accelerator_id": "ga-bp1l8bs1z8gw0u6y1ag5g",
					"listener_id":    "lsr-bp1jovsyx377pboiel4p2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address":        "www.baidu.com",
						"task_name":      "aaaa",
						"accelerator_id": "ga-bp1l8bs1z8gw0u6y1ag5g",
						"listener_id":    "lsr-bp1jovsyx377pboiel4p2",
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

// Test Ga ApplicationMonitor. <<< Resource test cases, automatically generated.
