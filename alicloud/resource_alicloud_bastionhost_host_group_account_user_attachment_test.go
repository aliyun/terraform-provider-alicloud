package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudBastionhostHostGroupAccountUserAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bastionhost_host_group_account_user_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudBastionhostHostGroupAccountUserAttachmentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &YundunBastionhostService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBastionhostHostGroupAccountUserAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sbastionhosthostgroupaccount%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBastionhostHostGroupAccountUserAttachmentBasicDependence0)
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
					"user_id":            "${alicloud_bastionhost_user.default.user_id}",
					"host_group_id":      "${alicloud_bastionhost_host_group.default.host_group_id}",
					"instance_id":        "${alicloud_bastionhost_host_account.default[0].instance_id}",
					"host_account_names": []string{"${alicloud_bastionhost_host_account.default[0].host_account_name}", "${alicloud_bastionhost_host_account.default[1].host_account_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_id":              CHECKSET,
						"host_group_id":        CHECKSET,
						"instance_id":          CHECKSET,
						"host_account_names.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_account_names": []string{"${alicloud_bastionhost_host_account.default[0].host_account_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_account_names.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"host_account_names": []string{"${alicloud_bastionhost_host_account.default[0].host_account_name}", "${alicloud_bastionhost_host_account.default[1].host_account_name}", "${alicloud_bastionhost_host_account.default[2].host_account_name}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"host_account_names.#": "3",
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

var AlicloudBastionhostHostGroupAccountUserAttachmentMap0 = map[string]string{
	"instance_id": CHECKSET,
}

func AlicloudBastionhostHostGroupAccountUserAttachmentBasicDependence0(name string) string {
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
resource "alicloud_bastionhost_host_account" "default" {
 count = 3
 instance_id       = alicloud_bastionhost_host.default.instance_id
 host_account_name = "${var.name}-${count.index}"
 host_id           = alicloud_bastionhost_host.default.host_id
 protocol_name     = "SSH"
 password          = "YourPassword12345"
}
resource "alicloud_bastionhost_host_group" "default" {
 instance_id     = data.alicloud_bastionhost_instances.default.ids.0
 host_group_name = var.name
}
resource "alicloud_bastionhost_user" "default" {
  instance_id    = data.alicloud_bastionhost_instances.default.ids.0
  mobile         = "13312345678"
  mobile_country_code = "CN"
  password       = "YourPassword-123"
  source         = "Local"
  user_name      = var.name
}
`, name)
}
