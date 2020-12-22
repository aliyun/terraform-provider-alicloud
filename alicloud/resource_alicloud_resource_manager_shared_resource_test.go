package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudResourceManagerSharedResource_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_shared_resource.default"
	ra := resourceAttrInit(resourceId, AlicloudResourceManagerSharedResourceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourcesharingService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerSharedResource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccResourceManagerSharedResource%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudResourceManagerSharedResourceBasicDependence)
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
					"resource_share_id": "${alicloud_resource_manager_resource_share.default.id}",
					"resource_id":       "${alicloud_vswitch.default.id}",
					"resource_type":     "VSwitch",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_share_id": CHECKSET,
						"resource_id":       CHECKSET,
						"resource_type":     "VSwitch",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudResourceManagerSharedResourceMap = map[string]string{
	"status": "Associated",
}

func AlicloudResourceManagerSharedResourceBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
resource "alicloud_vpc" "default" {
  name = var.name
  cidr_block = "192.168.0.0/16"
}
resource "alicloud_vswitch" "default" {
  availability_zone = data.alicloud_zones.default.ids.0
  cidr_block = "192.168.0.0/16"
  vpc_id = alicloud_vpc.default.id
}
resource "alicloud_resource_manager_resource_share" "default" {
	resource_share_name = var.name
}
`, name)
}
