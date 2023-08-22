package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCloudFirewallVpcFirewallCen_basic(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.CloudFirewallVpcFirewallCenSupportRegions)
	resourceId := "alicloud_cloud_firewall_vpc_firewall_cen.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudFirewallVpcFirewallCenMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallVpcFirewallCen")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scfwCen%d", defaultRegionToTest, rand)
	nameUpdate := fmt.Sprintf("tf-testacc%scfwCenup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudFirewallVpcFirewallCenBasicDependence)
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
					"vpc_firewall_name": name,
					"cen_id":            "${data.alicloud_cen_instances.cen_instances_ds.instances.0.id}",
					"vpc_region":        defaultRegionToTest,
					"status":            "open",
					"local_vpc": []map[string]interface{}{
						{
							"network_instance_id": "${data.alicloud_vpcs.vpcs_ds.vpcs.0.id}",
						},
					},
					"member_uid": "${data.alicloud_account.current.id}",
					"lang":       "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_firewall_name": name,
						"cen_id":            CHECKSET,
						"vpc_region":        defaultRegionToTest,
						"status":            "open",
						"member_uid":        CHECKSET,
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

	data "alicloud_account" "current" {
	}

	data "alicloud_cen_instances" "cen_instances_ds" {
  		name_regex = "^cfw-test-no-deleting"
	}

	data "alicloud_vpcs" "vpcs_ds" {
  		name_regex = "^cfw-test-no-deleting"
	}
`, name)
}
