package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Eais ClientInstanceAttachment. >>> Resource test cases, automatically generated.
// Case 4152
func TestAccAliCloudEaisClientInstanceAttachment_basic4152(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eais_client_instance_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudEaisClientInstanceAttachmentMap4152)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EaisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEaisClientInstanceAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%seaisclientinstanceattachment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEaisClientInstanceAttachmentBasicDependence4152)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VPCGatewayEndpointSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":        "${alicloud_eais_instance.default.id}",
					"client_instance_id": "${alicloud_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":        CHECKSET,
						"client_instance_id": CHECKSET,
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

var AlicloudEaisClientInstanceAttachmentMap4152 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudEaisClientInstanceAttachmentBasicDependence4152(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	vpc_name = var.name
}

resource "alicloud_vswitch" "default" {
	vpc_id = alicloud_vpc.default.id
	cidr_block = "172.16.0.0/21"
	availability_zone = "cn-hangzhou-j"
	name = var.name
}

resource "alicloud_security_group" "default" {
	name        = var.name
	description = "tf test"
	vpc_id      = alicloud_vpc.default.id
}

resource "alicloud_eais_instance" "default" {
	instance_type     = "eais.ei-a6.2xlarge"
	instance_name     = var.name
	security_group_id = alicloud_security_group.default.id
	vswitch_id        = alicloud_vswitch.default.id
}

data "alicloud_instance_types" "default" {
  availability_zone = "cn-hangzhou-j"
  system_disk_category = "cloud_efficiency"
  cpu_core_count = 4
  minimum_eni_ipv6_address_quantity = 1
}

data "alicloud_images" "default" {
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners = "system"
}

resource "alicloud_instance" "default" {
  image_id                   = "${data.alicloud_images.default.images.0.id}"
  instance_type              = "${data.alicloud_instance_types.default.instance_types.0.id}"
  instance_name              = "${var.name}"
  security_groups            = "${alicloud_security_group.default.*.id}"
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = "${alicloud_vswitch.default.id}"
}

`, name)
}

// Test Eais ClientInstanceAttachment. <<< Resource test cases, automatically generated.
