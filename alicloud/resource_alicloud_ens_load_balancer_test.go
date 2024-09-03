package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

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

// Test Ens LoadBalancer. >>> Resource test cases, automatically generated.
// Case 负载均衡_添加后端服务器_20240429 6626
func TestAccAliCloudEnsLoadBalancer_basic6626(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_load_balancer.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsLoadBalancerMap6626)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsLoadBalancer")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sensloadbalancer%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEnsLoadBalancerBasicDependence6626)
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
				Config: testAccConfig(map[string]interface{}{
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${alicloud_ens_instance.defaultfGH5i7.id}",
							"type":      "ens",
							"weight":    "100",
							"ip":        "${alicloud_ens_instance.defaultfGH5i7.private_ip_address}",
							"port":      "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "1",
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

var AlicloudEnsLoadBalancerMap6626 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudEnsLoadBalancerBasicDependence6626(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "ens_region_id" {
  default = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_network" "network" {
  network_name  = var.name
  description   = "LoadBalancerNetworkDescription_autotest"
  cidr_block    = "192.168.0.0/16"
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_vswitch" "switch" {
  description   = "LoadBalancerVSwitchDescription_autotest"
  cidr_block    = "192.168.2.0/24"
  vswitch_name  = format("%%s1", var.name)
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
  network_id    = alicloud_ens_network.network.id
}

resource "alicloud_ens_instance" "defaultfGH5i7" {
  system_disk {
    size     = "20"
    category = "cloud_efficiency"
  }
  scheduling_strategy        = "Concentrate"
  schedule_area_level        = "Region"
  image_id                   = "centos_6_08_64_20G_alibase_20171208"
  payment_type               = "Subscription"
  instance_type              = "ens.sn1.stiny"
  password                   = "12345678abcABC"
  status                     = "Running"
  amount                     = "1"
  vswitch_id                 = alicloud_ens_vswitch.switch.id
  internet_charge_type       = "95BandwidthByMonth"
  instance_name              = format("%%s2", var.name)
  auto_use_coupon            = "true"
  instance_charge_strategy   = "PriceHighPriority"
  ens_region_id              = var.ens_region_id
  period_unit                = "Month"
}


`, name)
}

// Test Ens LoadBalancer. <<< Resource test cases, automatically generated.
