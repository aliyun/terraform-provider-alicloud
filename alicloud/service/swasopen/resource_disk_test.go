// Package swasopen. This file is generated automatically. Please do not modify it manually, thank you!
package swasopen_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	swasopen "github.com/aliyun/terraform-provider-alicloud/alicloud/service/swasopen"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test SimpleApplicationServer Disk. >>> Resource test cases, automatically generated.
// Case testDisk 5761
// lintignore: AT001
func TestAccAliCloudSimpleApplicationServerDisk_basic5761(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_simple_application_server_disk.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSimpleApplicationServerDiskMap5761)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return swasopen.NewSimpleApplicationServerServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSimpleApplicationServerDisk")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsimpleapplicationserver%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSimpleApplicationServerDiskBasicDependence5761)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
			acctest.PreCheckPrePaidResources(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
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
						"instance_id": acctest.CHECKSET,
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
	"create_time": acctest.CHECKSET,
	"region_id":   acctest.CHECKSET,
	"disk_name":   acctest.CHECKSET,
}

func AlicloudSimpleApplicationServerDiskBasicDependence5761(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_simple_application_server_images" "default" {
  platform = "Linux"
}

data "alicloud_simple_application_server_plans" "default" {
  platform = "Linux"
}

resource "alicloud_simple_application_server_instance" "defaultV70JQf" {
  instance_name     = "testwujie"
  status            = "Running"
  plan_id           = data.alicloud_simple_application_server_plans.default.plans.0.id
  image_id          = data.alicloud_simple_application_server_images.default.images.0.id
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
