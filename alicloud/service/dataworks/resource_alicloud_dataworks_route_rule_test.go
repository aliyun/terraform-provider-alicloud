package dataworks_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dataworks"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// The route rule lives on a DataWorks network resource inside a DataWorks
// resource group, so the ACC fixture reuses the same cn-shenzhen default VPC /
// VSwitch / dw_resource_group scaffolding as alicloud_data_works_network and
// creates the route rule against the freshly created network. The resource is
// Framework-native, so it uses the ProtoV5 provider factories (which serve both
// the Framework resources and the SDKv2 dependency resources through one provider).
const routeRuleResourceId = "alicloud_dataworks_route_rule.default"

func TestAccAliCloudDataworksRouteRule_basic(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_dwrtrule%d", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-shenzhen"})
			acctest.PreCheck(t)
		},
		IDRefreshName:            routeRuleResourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckDataworksRouteRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataworksRouteRuleBasicConfig(name, "192.168.10.0/24", "VPC"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(routeRuleResourceId, "id"),
					resource.TestCheckResourceAttrSet(routeRuleResourceId, "route_rule_id"),
					resource.TestCheckResourceAttrSet(routeRuleResourceId, "resource_group_id"),
					resource.TestCheckResourceAttrSet(routeRuleResourceId, "resource_id"),
					resource.TestCheckResourceAttr(routeRuleResourceId, "destination_cidr", "192.168.10.0/24"),
					resource.TestCheckResourceAttr(routeRuleResourceId, "destination_type", "VPC"),
					resource.TestCheckResourceAttrSet(routeRuleResourceId, "destination_addressing_type"),
				),
			},
			{
				// destination_cidr and destination_type are the only fields
				// UpdateRouteRule accepts, so this step updates them in place
				// without recreating the rule.
				Config: testAccDataworksRouteRuleBasicConfig(name, "192.168.20.0/24", "IDC"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(routeRuleResourceId, "destination_cidr", "192.168.20.0/24"),
					resource.TestCheckResourceAttr(routeRuleResourceId, "destination_type", "IDC"),
					resource.TestCheckResourceAttrSet(routeRuleResourceId, "route_rule_id"),
				),
			},
			{
				ResourceName:      routeRuleResourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDataworksRouteRuleDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	service := dataworks.NewService(client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_dataworks_route_rule" {
			continue
		}
		_, err := service.DescribeRouteRule(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("route rule %s still exists", rs.Primary.ID)
		}
		if errs.IsProviderNotFound(err) {
			continue
		}
		return fmt.Errorf("describing route rule %s failed: %s", rs.Primary.ID, err)
	}
	return nil
}

// testAccDataworksRouteRuleBasicConfig builds the route rule with its DataWorks
// dependencies: a dw_resource_group (the resource_group_id), a network in the
// default-NODELETING-2 VPC (the resource_id), and the route rule itself. The
// dependency block mirrors alicloud_data_works_network so the two resources share
// the same region scaffolding and resource group.
func testAccDataworksRouteRuleBasicConfig(name, destinationCidr, destinationType string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

data "alicloud_vpcs" "default2" {
  name_regex = "^default-NODELETING-2$"
}

data "alicloud_vswitches" "default2" {
  vpc_id = data.alicloud_vpcs.default2.ids.0
}

resource "alicloud_data_works_dw_resource_group" "default" {
  payment_duration_unit = "Month"
  payment_type          = "PostPaid"
  specification         = "500"
  default_vswitch_id    = data.alicloud_vswitches.default.ids.0
  remark                = "OpenAPI测试用资源组"
  resource_group_name   = var.name
  default_vpc_id        = data.alicloud_vpcs.default.ids.0
  auto_renew            = false
}

resource "alicloud_data_works_network" "default" {
  dw_resource_group_id = alicloud_data_works_dw_resource_group.default.id
  vpc_id               = data.alicloud_vpcs.default2.ids.0
  vswitch_id           = data.alicloud_vswitches.default2.ids.0
}

resource "alicloud_dataworks_route_rule" "default" {
  resource_group_id = alicloud_data_works_dw_resource_group.default.id
  resource_id       = alicloud_data_works_network.default.id
  destination_cidr  = "%s"
  destination_type  = "%s"
}
`, name, destinationCidr, destinationType)
}
