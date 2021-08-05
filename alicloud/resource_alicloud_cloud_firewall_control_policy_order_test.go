package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCloudFirewallControlPolicyOrder_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_firewall_control_policy_order.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudFirewallControlPolicyOrderMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudfwService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudFirewallControlPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scloudfirewallcontrolpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudFirewallControlPolicyOrderBasicDependence0)
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
					"acl_uuid":  "${alicloud_cloud_firewall_control_policy.default.acl_uuid}",
					"direction": "${alicloud_cloud_firewall_control_policy.default.direction}",
					"order":     "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order":  "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"order":     "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"order":  "2",
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

var AlicloudCloudFirewallControlPolicyOrderMap0 = map[string]string{
	"acl_uuid":  CHECKSET,
	"direction": CHECKSET,
	"order":     CHECKSET,
}

func AlicloudCloudFirewallControlPolicyOrderBasicDependence0(name string) string {
	return fmt.Sprintf(` 

resource "alicloud_cloud_firewall_control_policy" "default" {
	application_name =  "ANY"
	acl_action       =  "accept"
	description      =  "%s"
	destination_type =  "net"
	destination      =  "100.1.1.0/24"
	direction        =  "out"
	proto            =  "ANY"
	source           =  "1.2.3.0/24"
	source_type      =  "net"
}
`, name)
}
