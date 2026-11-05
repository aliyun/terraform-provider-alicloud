package expressconnect_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/expressconnect"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccCheckRouterInterfaceConnectionExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No interface ID is set")
		}

		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		vpcService := expressconnect.NewExpressConnectService(client)

		response, err := vpcService.DescribeRouterInterfaceConnection(rs.Primary.ID, client.RegionId)
		if err != nil {
			return fmt.Errorf("Error finding interface %s: %#v", rs.Primary.ID, err)
		}
		if response.Status != string(expressconnect.Active) {
			return fmt.Errorf("Error connection router interface id %s is not Active.", response.RouterInterfaceId)
		}

		return nil
	}
}

func testAccCheckRouterInterfaceConnectionDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_router_interface_connection" {
			continue
		}

		// Try to find the interface
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		vpcService := expressconnect.NewExpressConnectService(client)

		ri, err := vpcService.DescribeRouterInterfaceConnection(rs.Primary.ID, client.RegionId)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return errs.WrapError(err)
		}

		if ri.Status == string(expressconnect.Active) {
			return errs.WrapError(errs.Error("Interface connection %s still exists.", rs.Primary.ID))
		}
	}
	return nil
}

func TestAccAliCloudVPCRouterInterfaceConnectionBasic(t *testing.T) {
	resourceId := "alicloud_router_interface_connection.foo"
	ra := acctest.ResourceAttrInit(resourceId, testAccRouterInterfaceConnectionCheckMap)
	rand := tfacctest.RandInt()
	testAccCheck := ra.ResourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
		},

		// module name

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      testAccCheckRouterInterfaceConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouterInterfaceConnectionConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouterInterfaceConnectionExists(resourceId),
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccRouterInterfaceConnectionConfigBasic(rand int) string {
	return fmt.Sprintf(
		`
provider "alicloud" {
  region = "${var.region}"
}
variable "region" {
  default = "cn-hangzhou"
}
variable "name" {
  default = "tf-testAccAlicloudRIConnection_basic%d"
}
data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING"
}
resource "alicloud_router_interface" "initiate" {
  opposite_region = "${var.region}"
  router_type = "VRouter"
  router_id = "${data.alicloud_vpcs.default.vpcs.0.router_id}"
  role = "InitiatingSide"
  specification = "Large.2"
  name = "${var.name}"
	description = "${var.name}"
	instance_charge_type = "PostPaid"
}
resource "alicloud_router_interface" "opposite" {
  provider = "alicloud"
  opposite_region = "${var.region}"
  router_type = "VRouter"
  router_id = "${data.alicloud_vpcs.default.vpcs.1.router_id}"
  role = "AcceptingSide"
  specification = "Large.1"
  name = "${var.name}-opposite"
  description = "${var.name}-opposite"
}
resource "alicloud_router_interface_connection" "foo" {
  interface_id = "${alicloud_router_interface.initiate.id}"
  opposite_interface_id = "${alicloud_router_interface.opposite.id}"
  depends_on = ["alicloud_router_interface_connection.bar"]
}
resource "alicloud_router_interface_connection" "bar" {
  provider = "alicloud"
  interface_id = "${alicloud_router_interface.opposite.id}"
  opposite_interface_id = "${alicloud_router_interface.initiate.id}"
}
`, rand)
}

var testAccRouterInterfaceConnectionCheckMap = map[string]string{
	"interface_id":                acctest.CHECKSET,
	"opposite_interface_id":       acctest.CHECKSET,
	"opposite_router_type":        "VRouter",
	"opposite_router_id":          acctest.CHECKSET,
	"opposite_interface_owner_id": acctest.CHECKSET,
}
