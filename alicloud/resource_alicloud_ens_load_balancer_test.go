package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ens LoadBalancer. >>> Resource test cases, automatically generated.
// Case 5071
func TestAccAliCloudEnsLoadBalancer_basic5071(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_load_balancer.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsLoadBalancerMap5071)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensloadbalancer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsLoadBalancerBasicDependence5071)
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
					"payment_type":       "PayAsYouGo",
					"ens_region_id":      "cn-chenzhou-telecom_unicom_cmcc",
					"load_balancer_spec": "elb.s1.small",
					"vswitch_id":         "${alicloud_ens_vswitch.switch.id}",
					"network_id":         "${alicloud_ens_network.network.id}",
					"load_balancer_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":       "PayAsYouGo",
						"ens_region_id":      "cn-chenzhou-telecom_unicom_cmcc",
						"load_balancer_spec": "elb.s1.small",
						"vswitch_id":         CHECKSET,
						"network_id":         CHECKSET,
						"load_balancer_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name + "_update",
					"payment_type":       "PayAsYouGo",
					"ens_region_id":      "cn-chenzhou-telecom_unicom_cmcc",
					"load_balancer_spec": "elb.s1.small",
					"vswitch_id":         "${alicloud_ens_vswitch.switch.id}",
					"network_id":         "${alicloud_ens_network.network.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name + "_update",
						"payment_type":       "PayAsYouGo",
						"ens_region_id":      "cn-chenzhou-telecom_unicom_cmcc",
						"load_balancer_spec": "elb.s1.small",
						"vswitch_id":         CHECKSET,
						"network_id":         CHECKSET,
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

var AlicloudEnsLoadBalancerMap5071 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudEnsLoadBalancerBasicDependence5071(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ens_network" "network" {
  network_name = var.name

  description   = "LoadBalancerNetworkDescription_autotest"
  cidr_block    = "192.168.2.0/24"
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_vswitch" "switch" {
  description  = "LoadBalancerVSwitchDescription_autotest"
  cidr_block   = "192.168.2.0/24"
  vswitch_name = var.name

  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
  network_id    = alicloud_ens_network.network.id
}


`, name)
}

// Case 5071  twin
func TestAccAliCloudEnsLoadBalancer_basic5071_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_load_balancer.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsLoadBalancerMap5071)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensloadbalancer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsLoadBalancerBasicDependence5071)
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
					"load_balancer_name": name,
					"payment_type":       "PayAsYouGo",
					"ens_region_id":      "cn-chenzhou-telecom_unicom_cmcc",
					"load_balancer_spec": "elb.s1.small",
					"vswitch_id":         "${alicloud_ens_vswitch.switch.id}",
					"network_id":         "${alicloud_ens_network.network.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name,
						"payment_type":       "PayAsYouGo",
						"ens_region_id":      "cn-chenzhou-telecom_unicom_cmcc",
						"load_balancer_spec": "elb.s1.small",
						"vswitch_id":         CHECKSET,
						"network_id":         CHECKSET,
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

// Test Ens LoadBalancer. <<< Resource test cases, automatically generated.
