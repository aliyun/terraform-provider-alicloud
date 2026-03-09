// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test SimpleApplicationServer Disk. >>> Resource test cases, automatically generated.
// Case testDisk 5761
func TestAccAliCloudSimpleApplicationServerDisk_basic5761(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_simple_application_server_disk.default"
	ra := resourceAttrInit(resourceId, AlicloudSimpleApplicationServerDiskMap5761)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SimpleApplicationServerServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSimpleApplicationServerDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsimpleapplicationserver%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSimpleApplicationServerDiskBasicDependence5761)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_size":   "20",
					"instance_id": "${alicloud_simple_application_server_instance.defaultV70JQf.id}",
					"remark":      "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_size":   "20",
						"instance_id": CHECKSET,
						"remark":      "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "testwujie",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "testwujie",
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

var AlicloudSimpleApplicationServerDiskMap5761 = map[string]string{
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
	"disk_name":   CHECKSET,
}

func AlicloudSimpleApplicationServerDiskBasicDependence5761(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_simple_application_server_instance" "defaultV70JQf" {
  instance_name     = "testwujie"
  status            = "Running"
  plan_id           = "swas.s1.c2m2s50b3"
  image_id          = "21e9617bd4754f77a090d2fbc94916a4"
  period            = "1"
  data_disk_size    = "0"
  password          = "@3612568Wj"
  payment_type      = "Subscription"
  auto_renew        = true
  auto_renew_period = "1"
}


`, name)
}

// Test SimpleApplicationServer Disk. <<< Resource test cases, automatically generated.
