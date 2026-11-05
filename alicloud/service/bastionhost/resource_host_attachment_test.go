package bastionhost_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/bastionhost"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudBastionhostHostAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_bastionhost_host_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudBastionhostHostAttachmentMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return bastionhost.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeBastionhostHostAttachment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sbastionattachment%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudBastionhostHostAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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

var AliCloudBastionhostHostAttachmentMap0 = map[string]string{
	"host_id":       acctest.CHECKSET,
	"instance_id":   acctest.CHECKSET,
	"host_group_id": acctest.CHECKSET,
}

func AliCloudBastionhostHostAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(` 
	variable "name" {
  		default = "%s"
	}

	data "alicloud_bastionhost_instances" "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_security_group" "default" {
  		count  = length(data.alicloud_bastionhost_instances.default.ids) > 0 ? 0 : 1
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_bastionhost_instance" "default" {
  		count              = length(data.alicloud_bastionhost_instances.default.ids) > 0 ? 0 : 1
  		description        = var.name
  		license_code       = "bhah_ent_50_asset"
  		plan_code          = "cloudbastion"
  		storage            = "5"
  		bandwidth          = "5"
  		period             = "1"
  		vswitch_id         = data.alicloud_vswitches.default.ids.0
  		security_group_ids = [alicloud_security_group.default.0.id]
	}

	resource "alicloud_bastionhost_host_group" "default" {
  		instance_id     = local.instance_id
  		host_group_name = var.name
	}

	resource "alicloud_bastionhost_host" "default" {
  		instance_id          = alicloud_bastionhost_host_group.default.instance_id
  		host_name            = var.name
  		active_address_type  = "Private"
  		host_private_address = "172.16.0.10"
  		os_type              = "Linux"
  		source               = "Local"
	}

	locals {
  		instance_id = length(data.alicloud_bastionhost_instances.default.ids) > 0 ? data.alicloud_bastionhost_instances.default.ids.0 : alicloud_bastionhost_instance.default.0.id
	}
`, name)
}
