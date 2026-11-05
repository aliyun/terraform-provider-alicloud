package cloudfirewall_test

import (
	"fmt"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	connectivity "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	cloudfirewall "github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudfirewall"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	resource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test CloudFirewall VpcFirewallControlPolicyOrder. >>> Resource test cases, automatically generated.
// Case test1 12717
func TestAccAliCloudCloudFirewallVpcFirewallControlPolicyOrder_basic12717(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_vpc_firewall_control_policy_order.test.0"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCloudFirewallVpcFirewallControlPolicyOrderMap12717)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallVpcFirewallControlPolicyOrder")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudfirewall%d", rand)
	acctest.ResourceTestAccConfig(resourceId, name, alicloudCloudFirewallVpcFirewallControlPolicyOrderBasicDependence12717)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: alicloudCloudFirewallVpcFirewallControlPolicyOrderBasicDependence12717(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order":           acctest.CHECKSET,
						"vpc_firewall_id": acctest.CHECKSET,
						"lang":            "zh",
						"acl_uuid":        acctest.CHECKSET,
					}),
				),
			},
			{
				Config: alicloudCloudFirewallVpcFirewallControlPolicyOrderUpdateDependence12717(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order":           "3",
						"vpc_firewall_id": acctest.CHECKSET,
						"lang":            "en",
						"acl_uuid":        acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

var AlicloudCloudFirewallVpcFirewallControlPolicyOrderMap12717 = map[string]string{}

func alicloudCloudFirewallVpcFirewallControlPolicyOrderBasicDependence12717(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_account" "test" {
}

resource "alicloud_cen_instance" "test" {
  cen_instance_name = var.name
}

resource "alicloud_cloud_firewall_vpc_firewall_control_policy" "test" {
  count            = 3
  order            = "-1"
  destination      = "127.0.0.2/32"
  application_name = "ANY"
  description      = "example_value"
  source_type      = "net"
  dest_port        = "80/88"
  acl_action       = "accept"
  lang             = "zh"
  destination_type = "net"
  source           = "127.0.0.1/32"
  dest_port_type   = "port"
  proto            = "TCP"
  release          = true
  vpc_firewall_id  = alicloud_cen_instance.test.id
  lifecycle {
    ignore_changes = [order]
  }
}

resource "alicloud_cloud_firewall_vpc_firewall_control_policy_order" "test" {
  count           = 3
  order           = "${count.index + 1}"
  vpc_firewall_id = alicloud_cen_instance.test.id
  lang            = "zh"
  acl_uuid        = split(":", alicloud_cloud_firewall_vpc_firewall_control_policy.test[count.index].id)[1]
}`, name)
}

func alicloudCloudFirewallVpcFirewallControlPolicyOrderUpdateDependence12717(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_account" "test" {
}

resource "alicloud_cen_instance" "test" {
  cen_instance_name = var.name
}

resource "alicloud_cloud_firewall_vpc_firewall_control_policy" "test" {
  count            = 3
  order            = "-1"
  destination      = "127.0.0.2/32"
  application_name = "ANY"
  description      = "example_value"
  source_type      = "net"
  dest_port        = "80/88"
  acl_action       = "accept"
  lang             = "zh"
  destination_type = "net"
  source           = "127.0.0.1/32"
  dest_port_type   = "port"
  proto            = "TCP"
  release          = true
  vpc_firewall_id  = alicloud_cen_instance.test.id
  lifecycle {
    ignore_changes = [order]
  }
}

resource "alicloud_cloud_firewall_vpc_firewall_control_policy_order" "test" {
  count           = 3
  order           = "${3 - count.index}"
  vpc_firewall_id = alicloud_cen_instance.test.id
  lang            = "en"
  acl_uuid        = split(":", alicloud_cloud_firewall_vpc_firewall_control_policy.test[count.index].id)[1]
}`, name)
}

// Test CloudFirewall VpcFirewallControlPolicyOrder. <<< Resource test cases, automatically generated.
