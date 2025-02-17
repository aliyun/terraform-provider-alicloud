package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudPhone Instance. >>> Resource test cases, automatically generated.
// Case chuyuan_CreateInstance_prod 9932
func TestAccAliCloudCloudPhoneInstance_basic9932(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_phone_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudPhoneInstanceMap9932)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudPhoneServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudPhoneInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudphone%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudPhoneInstanceBasicDependence9932)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"android_instance_group_id": "${alicloud_cloud_phone_instance_group.defaultYHMlTO.id}",
					"android_instance_name":     "CreateInstanceName",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"android_instance_group_id": CHECKSET,
						"android_instance_name":     "CreateInstanceName",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"android_instance_name": "AndroidInstanceName",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"android_instance_name": "AndroidInstanceName",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"android_instance_name": "NewAndroidInstanceName",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"android_instance_name": "NewAndroidInstanceName",
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

var AlicloudCloudPhoneInstanceMap9932 = map[string]string{}

func AlicloudCloudPhoneInstanceBasicDependence9932(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_cloud_phone_policy" "defaultjZ1gi0" {
}

resource "alicloud_cloud_phone_instance_group" "defaultYHMlTO" {
  instance_group_spec = "acp.basic.small"
  policy_group_id     = alicloud_cloud_phone_policy.defaultjZ1gi0.id
  instance_group_name = "AutoCreateGroupName"
  period              = "1"
  number_of_instances = "1"
  charge_type         = "PostPaid"
  image_id            = "imgc-075cllfeuazh03tg9"
  period_unit         = "Hour"
  auto_renew          = false
  amount              = "1"
  auto_pay            = false
  gpu_acceleration    = false
}


`, name)
}

// Test CloudPhone Instance. <<< Resource test cases, automatically generated.
