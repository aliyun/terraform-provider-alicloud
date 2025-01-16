package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudAckOneMembershipAttachment_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ack_one_membership_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudAckOneMembershipAttachmentMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AckOneServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAckOneMembershipAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testAccAckOneMembershipAttachment-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAckOneMembershipAttachmentBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":     "${alicloud_ack_one_cluster.default.id}",
					"sub_cluster_id": "${alicloud_cs_managed_kubernetes.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(
						map[string]string{
							"cluster_id":     CHECKSET,
							"sub_cluster_id": CHECKSET,
						},
					),
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

var AliCloudAckOneMembershipAttachmentMap = map[string]string{
	"cluster_id":     CHECKSET,
	"sub_cluster_id": CHECKSET,
}

func AliCloudAckOneMembershipAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultVpc" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "defaultyVSwitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  cidr_block   = "172.16.2.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_ack_one_cluster" "default" {
  network {
    vpc_id    = alicloud_vpc.defaultVpc.id
    vswitches = ["${alicloud_vswitch.defaultyVSwitch.id}"]
  }
}

# leave it to empty would create a new one
variable "vpc_id" {
  description = "Existing vpc id used to create several vswitches and other resources."
  default     = ""
}

variable "vpc_cidr" {
  description = "The cidr block used to launch a new vpc when 'vpc_id' is not specified."
  default     = "10.0.0.0/8"
}

# leave it to empty then terraform will create several vswitches
variable "vswitch_ids" {
  description = "List of existing vswitch id."
  type        = list(string)
  default     = []
}

variable "vswitch_cidrs" {
  description = "List of cidr blocks used to create several new vswitches when 'vswitch_ids' is not specified."
  type        = list(string)
  default     = ["10.1.0.0/16", "10.2.0.0/16"]
}

# options: between 24-28
variable "node_cidr_mask" {
  description = "The node cidr block to specific how many pods can run on single node."
  default     = 24
}

# options: ipvs|iptables
variable "proxy_mode" {
  description = "Proxy mode is option of kube-proxy."
  default     = "ipvs"
}

variable "service_cidr" {
  description = "The kubernetes service cidr block. It cannot be equals to vpc's or vswitch's or pod's and cannot be in them."
  default     = "192.168.0.0/16"
}

variable "terway_vswitch_ids" {
  description = "List of existing vswitch ids for terway."
  type        = list(string)
  default     = []
}

variable "terway_vswitch_cidrs" {
  description = "List of cidr blocks used to create several new vswitches when 'terway_vswitch_cidrs' is not specified."
  type        = list(string)
  default     = ["10.4.0.0/16", "10.5.0.0/16"]
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {}

# If there is not specifying vpc_id, the module will launch a new vpc
resource "alicloud_vpc" "vpc" {
  count      = var.vpc_id == "" ? 1 : 0
  cidr_block = var.vpc_cidr
}

# According to the vswitch cidr blocks to launch several vswitches
resource "alicloud_vswitch" "vswitches" {
  count      = length(var.vswitch_ids) > 0 ? 0 : length(var.vswitch_cidrs)
  vpc_id     = var.vpc_id == "" ? join("", alicloud_vpc.vpc.*.id) : var.vpc_id
  cidr_block = element(var.vswitch_cidrs, count.index)
  zone_id    = data.alicloud_enhanced_nat_available_zones.enhanced.zones[count.index].zone_id
}

# According to the vswitch cidr blocks to launch several vswitches
resource "alicloud_vswitch" "terway_vswitches" {
  count      = length(var.terway_vswitch_ids) > 0 ? 0 : length(var.terway_vswitch_cidrs)
  vpc_id     = var.vpc_id == "" ? join("", alicloud_vpc.vpc.*.id) : var.vpc_id
  cidr_block = element(var.terway_vswitch_cidrs, count.index)
  zone_id    = data.alicloud_enhanced_nat_available_zones.enhanced.zones[count.index].zone_id
}

resource "alicloud_cs_managed_kubernetes" "default" {
  cluster_spec = "ack.pro.small"
  # version can not be defined in variables.tf.
  # version            = "1.26.3-aliyun.1"
  vswitch_ids     = length(var.vswitch_ids) > 0 ? split(",", join(",", var.vswitch_ids)) : length(var.vswitch_cidrs) < 1 ? [] : split(",", join(",", alicloud_vswitch.vswitches.*.id))
  pod_vswitch_ids = length(var.terway_vswitch_ids) > 0 ? split(",", join(",", var.terway_vswitch_ids)) : length(var.terway_vswitch_cidrs) < 1 ? [] : split(",", join(",", alicloud_vswitch.terway_vswitches.*.id))
  new_nat_gateway = true
  node_cidr_mask  = var.node_cidr_mask
  proxy_mode      = var.proxy_mode
  service_cidr    = var.service_cidr

  is_enterprise_security_group = true

  addons {
    name = "terway-eniip"
  }
}
`, name)
}
