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
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "terraform-example"
}

variable "key_name" {
  default = "%s"
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {
}

data "alicloud_instance_types" "cloud_efficiency" {
  availability_zone    = data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id
  cpu_core_count       = 4
  memory_size          = 8
  kubernetes_node_role = "Worker"
  system_disk_category = "cloud_efficiency"
}

resource "alicloud_vpc" "default" {
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  cidr_block = "10.4.0.0/24"
  vpc_id     = alicloud_vpc.default.id
  zone_id    = data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id
}

resource "alicloud_cs_managed_kubernetes" "default" {
  cluster_spec         = "ack.pro.small"
  vswitch_ids          = [alicloud_vswitch.default.id]
  new_nat_gateway      = true
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true

  is_enterprise_security_group = true
}

resource "alicloud_key_pair" "default" {
  key_pair_name = var.key_name
}

resource "alicloud_cs_kubernetes_node_pool" "default" {
  node_pool_name       = var.name
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.cloud_efficiency.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alicloud_key_pair.default.key_pair_name
  desired_size         = 1
}

resource "alicloud_ack_one_cluster" "default" {
  depends_on = [alicloud_cs_managed_kubernetes.default]
  network {
    vpc_id    = alicloud_vpc.default.id
    vswitches = [alicloud_vswitch.default.id]
  }
  argocd_enabled = false
}
`, name)
}
