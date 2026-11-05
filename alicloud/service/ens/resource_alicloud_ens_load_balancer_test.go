package ens_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ens"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Ens LoadBalancer. >>> Resource test cases, automatically generated.
// Case 负载均衡_20240424 6598
func TestAccAliCloudEnsLoadBalancer_basic6598(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_load_balancer.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEnsLoadBalancerMap6598)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ens.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEnsLoadBalancer")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEnsLoadBalancerBasicDependence6598)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":       "PayAsYouGo",
					"ens_region_id":      "cn-chenzhou-telecom_unicom_cmcc",
					"load_balancer_spec": "elb.s1.small",
					"vswitch_id":         "${alicloud_ens_vswitch.switch.id}",
					"network_id":         "${alicloud_ens_network.network.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":       "PayAsYouGo",
						"ens_region_id":      "cn-chenzhou-telecom_unicom_cmcc",
						"load_balancer_spec": "elb.s1.small",
						"vswitch_id":         acctest.CHECKSET,
						"network_id":         acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name,
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

func TestAccAliCloudEnsLoadBalancer_basic6598_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_load_balancer.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEnsLoadBalancerMap6598)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ens.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEnsLoadBalancer")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEnsLoadBalancerBasicDependence6598)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"vswitch_id":         acctest.CHECKSET,
						"network_id":         acctest.CHECKSET,
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

var AliCloudEnsLoadBalancerMap6598 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

func AliCloudEnsLoadBalancerBasicDependence6598(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ens_network" "network" {
  network_name  = "LoadBalancerNetworkNetworkName_autotest"
  description   = var.name
  cidr_block    = "192.168.0.0/16"
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_vswitch" "switch" {
  description   = "LoadBalancerVSwitchDescription_autotest"
  cidr_block    = "192.168.2.0/24"
  vswitch_name  = var.name
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
  network_id    = alicloud_ens_network.network.id
}
`, name)
}

// Case 负载均衡_添加后端服务器_20240429 6626
func TestAccAliCloudEnsLoadBalancer_basic6626(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_load_balancer.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEnsLoadBalancerMap6626)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ens.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEnsLoadBalancer")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEnsLoadBalancerBasicDependence6626)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":       "PayAsYouGo",
					"ens_region_id":      "cn-chenzhou-telecom_unicom_cmcc",
					"load_balancer_spec": "elb.s1.small",
					"vswitch_id":         "${alicloud_ens_vswitch.switch.id}",
					"network_id":         "${alicloud_ens_network.network.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":       "PayAsYouGo",
						"ens_region_id":      "cn-chenzhou-telecom_unicom_cmcc",
						"load_balancer_spec": "elb.s1.small",
						"vswitch_id":         acctest.CHECKSET,
						"network_id":         acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${alicloud_ens_instance.default.0.id}",
							"type":      "ens",
						},
						{
							"server_id": "${alicloud_ens_instance.default.1.id}",
							"type":      "ens",
							"weight":    "50",
							"ip":        "${alicloud_ens_instance.default.1.private_ip_address}",
							"port":      "0",
						},
						{
							"server_id": "${alicloud_ens_instance.default.2.id}",
							"type":      "ens",
							"weight":    "100",
							"ip":        "${alicloud_ens_instance.default.2.private_ip_address}",
							"port":      "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${alicloud_ens_instance.default.1.id}",
							"type":      "ens",
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
				Config: testAccConfig(map[string]interface{}{
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${alicloud_ens_instance.default.0.id}",
							"type":      "ens",
							"weight":    "20",
							"ip":        "${alicloud_ens_instance.default.0.private_ip_address}",
							"port":      "0",
						},
						{
							"server_id": "${alicloud_ens_instance.default.1.id}",
							"type":      "ens",
							"weight":    "50",
							"ip":        "${alicloud_ens_instance.default.1.private_ip_address}",
							"port":      "0",
						},
						{
							"server_id": "${alicloud_ens_instance.default.2.id}",
							"type":      "ens",
							"weight":    "100",
							"ip":        "${alicloud_ens_instance.default.2.private_ip_address}",
							"port":      "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "3",
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

func TestAccAliCloudEnsLoadBalancer_basic6626_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_load_balancer.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEnsLoadBalancerMap6626)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ens.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEnsLoadBalancer")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEnsLoadBalancerBasicDependence6626)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_name": name,
					"payment_type":       "PayAsYouGo",
					"ens_region_id":      "cn-chenzhou-telecom_unicom_cmcc",
					"load_balancer_spec": "elb.s1.small",
					"vswitch_id":         "${alicloud_ens_vswitch.switch.id}",
					"network_id":         "${alicloud_ens_network.network.id}",
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${alicloud_ens_instance.default.0.id}",
							"type":      "ens",
							"weight":    "20",
							"ip":        "${alicloud_ens_instance.default.0.private_ip_address}",
							"port":      "0",
						},
						{
							"server_id": "${alicloud_ens_instance.default.1.id}",
							"type":      "ens",
							"weight":    "50",
							"ip":        "${alicloud_ens_instance.default.1.private_ip_address}",
							"port":      "0",
						},
						{
							"server_id": "${alicloud_ens_instance.default.2.id}",
							"type":      "ens",
							"weight":    "100",
							"ip":        "${alicloud_ens_instance.default.2.private_ip_address}",
							"port":      "0",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_name": name,
						"payment_type":       "PayAsYouGo",
						"ens_region_id":      "cn-chenzhou-telecom_unicom_cmcc",
						"load_balancer_spec": "elb.s1.small",
						"vswitch_id":         acctest.CHECKSET,
						"network_id":         acctest.CHECKSET,
						"backend_servers.#":  "3",
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

var AliCloudEnsLoadBalancerMap6626 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

func AliCloudEnsLoadBalancerBasicDependence6626(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
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
  vswitch_name  = var.name
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
  network_id    = alicloud_ens_network.network.id
}

resource "alicloud_ens_instance" "default" {
  count = 3
  system_disk {
    size     = "20"
    category = "cloud_efficiency"
  }
  scheduling_strategy      = "Concentrate"
  schedule_area_level      = "Region"
  image_id                 = "centos_6_08_64_20G_alibase_20171208"
  payment_type             = "Subscription"
  instance_type            = "ens.sn1.stiny"
  password                 = "12345678abcABC"
  status                   = "Running"
  amount                   = "1"
  vswitch_id               = alicloud_ens_vswitch.switch.id
  internet_charge_type     = "95BandwidthByMonth"
  instance_name            = var.name
  auto_use_coupon          = "true"
  instance_charge_strategy = "PriceHighPriority"
  ens_region_id            = "cn-chenzhou-telecom_unicom_cmcc"
  period_unit              = "Month"
}
`, name)
}

// Test Ens LoadBalancer. <<< Resource test cases, automatically generated.
