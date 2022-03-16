package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudSimpleApplicationServerCustomImage_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_simple_application_server_custom_image.default"
	ra := resourceAttrInit(resourceId, AlicloudSimpleApplicationServerCustomImageMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SwasOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSimpleApplicationServerCustomImage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-swascustomimage%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSimpleApplicationServerCustomImageBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.SimpleApplicationServerNotSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_image_name":  "${var.name}",
					"system_snapshot_id": "${alicloud_simple_application_server_snapshot.default.id}",
					"instance_id":        "${data.alicloud_simple_application_server_disks.default.disks.0.instance_id}",
					"description":        "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"custom_image_name": name,
						"description":       name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Share",
				}),
				Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{})),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "UnShare",
				}),
				Check: resource.ComposeTestCheckFunc(testAccCheck(map[string]string{})),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "system_snapshot_id", "status"},
			},
		},
	})
}

var AlicloudSimpleApplicationServerCustomImageMap0 = map[string]string{}

func AlicloudSimpleApplicationServerCustomImageBasicDependence0(name string) string {
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

resource "alicloud_simple_application_server_instance" "default" {
  payment_type   = "Subscription"
  plan_id        = data.alicloud_simple_application_server_plans.default.plans.0.id
  instance_name  = var.name
  image_id       = data.alicloud_simple_application_server_images.default.images.0.id
  period         = 1
  data_disk_size = 100
}

data "alicloud_simple_application_server_disks" "default" {
  disk_type = "System"
  instance_id = alicloud_simple_application_server_instance.default.id
}

resource "alicloud_simple_application_server_snapshot" "default" {
  disk_id       = data.alicloud_simple_application_server_disks.default.ids.0
  snapshot_name = var.name
}
`, name)
}
