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

func TestAccAliCloudCloudFirewallVpcFirewallCen_basic(t *testing.T) {
	var v map[string]interface{}
	//acctest.CheckoutSupportedRegions(t, true, connectivity.CloudFirewallVpcFirewallCenSupportRegions)
	resourceId := "alicloud_cloud_firewall_vpc_firewall_cen.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudCloudFirewallVpcFirewallCenMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudfirewall.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudFirewallVpcFirewallCen")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scfwCen%d", acctest.DefaultRegionToTest, rand)
	nameUpdate := fmt.Sprintf("tf-testacc%scfwCenup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudCloudFirewallVpcFirewallCenBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_firewall_name": name,
					"cen_id":            "${alicloud_cen_instance_attachment.attach1.instance_id}",
					"vpc_region":        "${data.alicloud_regions.current.ids.0}",
					"status":            "open",
					"local_vpc": []map[string]interface{}{
						{
							"network_instance_id": "${alicloud_cen_instance_attachment.attach1.child_instance_id}",
						},
					},
					"member_uid": "${data.alicloud_account.current.id}",
					"lang":       "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_firewall_name": name,
						"cen_id":            acctest.CHECKSET,
						"vpc_region":        acctest.DefaultRegionToTest,
						"status":            "open",
						"member_uid":        acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_firewall_name": nameUpdate,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_firewall_name": nameUpdate,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "close",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "close",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lang"},
			},
		},
	})
}

var AliCloudCloudFirewallVpcFirewallCenMap = map[string]string{}

func AliCloudCloudFirewallVpcFirewallCenBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_regions" "current" {
  current = true
}
data "alicloud_zones" "zone" {
  available_instance_type = "ecs.sn1ne.large"
  available_resource_creation = "VSwitch"
}

data "alicloud_account" "current" {
}

resource "alicloud_vpc" "foo" {
  vpc_name   = "${var.name}-foo"
  cidr_block = "192.168.0.0/16"
}
resource "alicloud_vpc" "bar" {
  vpc_name   = "${var.name}-bar"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
  vpc_id       = alicloud_vpc.foo.id
  cidr_block   = "192.168.10.0/24"
  zone_id      = data.alicloud_zones.zone.zones.0.id
  vswitch_name = "${var.name}-foo"
}

resource "alicloud_vswitch" "bar" {
  vpc_id       = alicloud_vpc.bar.id
  cidr_block   = "172.16.10.0/24"
  zone_id      = data.alicloud_zones.zone.zones.0.id
  vswitch_name = "${var.name}-bar"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
  description       = var.name
}

resource "alicloud_cen_instance_attachment" "attach1" {
  instance_id = alicloud_cen_instance.default.id
  child_instance_id = alicloud_vpc.foo.id
  child_instance_type = "VPC"
  child_instance_region_id = data.alicloud_regions.current.ids.0
  child_instance_owner_id = data.alicloud_account.current.id
}
resource "alicloud_cen_instance_attachment" "attach2" {
  instance_id = alicloud_cen_instance.default.id
  child_instance_id = alicloud_vpc.bar.id
  child_instance_type = "VPC"
  child_instance_region_id = data.alicloud_regions.current.ids.0
  child_instance_owner_id = data.alicloud_account.current.id
}
`, name)
}
