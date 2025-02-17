package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudPhone Image. >>> Resource test cases, automatically generated.
// Case chuyuan_createImage_prod 9933
func TestAccAliCloudCloudPhoneImage_basic9933(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_phone_image.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudPhoneImageMap9933)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudPhoneServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudPhoneImage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudphone%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudPhoneImageBasicDependence9933)
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
					"image_name":  name,
					"instance_id": "${alicloud_cloud_phone_instance.default04hhXk.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name":  name,
						"instance_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"image_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"image_name": name + "_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id"},
			},
		},
	})
}

var AlicloudCloudPhoneImageMap9933 = map[string]string{
	"status": CHECKSET,
}

func AlicloudCloudPhoneImageBasicDependence9933(name string) string {
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

resource "alicloud_cloud_phone_instance" "default04hhXk" {
  android_instance_group_id = alicloud_cloud_phone_instance_group.defaultYHMlTO.id
  android_instance_name     = "CreateInstanceName"
}


`, name)
}

// Test CloudPhone Image. <<< Resource test cases, automatically generated.
