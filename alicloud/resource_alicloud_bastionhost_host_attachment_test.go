package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudBastionhostHostAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bastionhost_host_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudBastionhostHostAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBastionhostHostAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sbastionhosthostattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBastionhostHostAttachmentBasicDependence0)
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
					"instance_id":   "${alicloud_bastionhost_host_group.default.instance_id}",
					"host_group_id": "${alicloud_bastionhost_host_group.default.host_group_id}",
					"host_id":       "${alicloud_bastionhost_host.default.host_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
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

var AlicloudBastionhostHostAttachmentMap0 = map[string]string{
	"host_id":       CHECKSET,
	"instance_id":   CHECKSET,
	"host_group_id": CHECKSET,
}

func AlicloudBastionhostHostAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_bastionhost_instances" "default" {}

resource "alicloud_bastionhost_host" "default" {
 instance_id          = data.alicloud_bastionhost_instances.default.ids.0
 host_name            = var.name
 active_address_type  = "Private"
 host_private_address = "172.16.0.10"
 os_type              = "Linux"
 source               = "Local"
}
resource "alicloud_bastionhost_host_group" "default" {
 instance_id          = data.alicloud_bastionhost_instances.default.ids.0
 host_group_name = var.name
}
`, name)
}
