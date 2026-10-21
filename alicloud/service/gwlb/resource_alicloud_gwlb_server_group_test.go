package gwlb_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/gwlb"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Gwlb ServerGroup. >>> Resource test cases, automatically generated.
// Case ServerGroup Test_instance 8419
func TestAccAliCloudGwlbServerGroup_basic8419(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gwlb_server_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGwlbServerGroupMap8419)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gwlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeServerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccgwlb%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudGwlbServerGroupBasicDependence8419)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": "5TCH",
					"health_check_config": []map[string]interface{}{
						{
							"health_check_protocol":        "HTTP",
							"health_check_connect_port":    "80",
							"health_check_connect_timeout": "5",
							"health_check_domain":          "www.domain.com",
							"health_check_enabled":         "true",
							"health_check_http_code": []string{
								"http_2xx", "http_3xx", "http_4xx"},
							"health_check_interval": "10",
							"health_check_path":     "/health-check",
							"healthy_threshold":     "2",
							"unhealthy_threshold":   "2",
						},
					},
					"protocol":          "GENEVE",
					"server_group_type": "Instance",
					"connection_drain_config": []map[string]interface{}{
						{
							"connection_drain_enabled": "true",
							"connection_drain_timeout": "1",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"vpc_id":            "${alicloud_vpc.defaultEaxcvb.id}",
					"dry_run":           "false",
					"servers": []map[string]interface{}{
						{
							"server_id":   "${alicloud_instance.defaultH6McvC.id}",
							"server_type": "Ecs",
						},
					},
					"server_group_name":    name,
					"server_failover_mode": "Rebalance",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler":            "5TCH",
						"protocol":             "GENEVE",
						"server_group_type":    "Instance",
						"resource_group_id":    acctest.CHECKSET,
						"vpc_id":               acctest.CHECKSET,
						"dry_run":              "false",
						"servers.#":            "1",
						"server_group_name":    name,
						"server_failover_mode": "Rebalance",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": "3TCH",
					"health_check_config": []map[string]interface{}{
						{
							"health_check_protocol":        "HTTP",
							"health_check_connect_port":    "81",
							"health_check_connect_timeout": "6",
							"health_check_domain":          "www.domain-update.com",
							"health_check_enabled":         "false",
							"health_check_http_code": []string{
								"http_5xx"},
							"health_check_interval": "11",
							"health_check_path":     "/health-check-update",
							"healthy_threshold":     "3",
							"unhealthy_threshold":   "3",
						},
					},
					"connection_drain_config": []map[string]interface{}{
						{
							"connection_drain_enabled": "false",
							"connection_drain_timeout": "2",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"servers": []map[string]interface{}{
						{
							"server_type": "Eni",
							"server_id":   "${alicloud_instance.defaultH6McvC.network_interface_id}",
						},
					},
					"server_group_name":    name + "_update",
					"server_failover_mode": "NoRebalance",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler":            "3TCH",
						"resource_group_id":    acctest.CHECKSET,
						"servers.#":            "1",
						"server_group_name":    name + "_update",
						"server_failover_mode": "NoRebalance",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudGwlbServerGroupMap8419 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

func AlicloudGwlbServerGroupBasicDependence8419(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-wulanchabu"
}


data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  instance_type_family = "ecs.sn1ne"
}

data "alicloud_images" "default" {
  name_regex  = "^aliyun_2_1903_x64_20G_alibase"
  most_recent = true
  owners      = "system"
}

resource "alicloud_vpc" "defaultEaxcvb" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = "tf-gwlb-vpc"
}

resource "alicloud_vswitch" "defaultc3uVID" {
  vpc_id       = alicloud_vpc.defaultEaxcvb.id
  zone_id      = data.alicloud_zones.default.zones.0.id
  cidr_block   = "10.0.0.0/24"
  vswitch_name = "tf-test-vsw1"
}

resource "alicloud_security_group" "default7NNxRl" {
  description         = "sg"
  security_group_name = "sg_name"
  vpc_id              = alicloud_vpc.defaultEaxcvb.id
  security_group_type = "normal"
}

resource "alicloud_instance" "defaultH6McvC" {
	vswitch_id = alicloud_vswitch.defaultc3uVID.id
	image_id = data.alicloud_images.default.images.0.id
	
	instance_type = data.alicloud_instance_types.default.instance_types.0.id
	system_disk_category = "cloud_efficiency"
	
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 5
	instance_name = format("%%s4", var.name)
    description   = "tf-test-ecs"
    security_groups = [alicloud_security_group.default7NNxRl.id]

  	availability_zone          = alicloud_vswitch.defaultc3uVID.zone_id
  	instance_charge_type       = "PostPaid"
}


`, name)
}

// Case ServerGroup Test_IP 8500
func TestAccAliCloudGwlbServerGroup_basic8500(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gwlb_server_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGwlbServerGroupMap8500)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gwlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeServerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccgwlb%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudGwlbServerGroupBasicDependence8500)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": "5TCH",
					"health_check_config": []map[string]interface{}{
						{
							"health_check_protocol":        "HTTP",
							"health_check_connect_port":    "80",
							"health_check_connect_timeout": "5",
							"health_check_domain":          "www.domain.com",
							"health_check_enabled":         "true",
							"health_check_http_code": []string{
								"http_2xx", "http_3xx", "http_4xx"},
							"health_check_interval": "10",
							"health_check_path":     "/health-check",
							"healthy_threshold":     "2",
							"unhealthy_threshold":   "2",
						},
					},
					"protocol":          "GENEVE",
					"server_group_type": "Ip",
					"connection_drain_config": []map[string]interface{}{
						{
							"connection_drain_enabled": "true",
							"connection_drain_timeout": "1",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"vpc_id":            "${alicloud_vpc.defaultEaxcvb.id}",
					"dry_run":           "false",
					"servers": []map[string]interface{}{
						{
							"server_id":   "10.0.0.1",
							"server_ip":   "10.0.0.1",
							"server_type": "Ip",
						},
						{
							"server_id":   "10.0.0.2",
							"server_ip":   "10.0.0.2",
							"server_type": "Ip",
						},
						{
							"server_id":   "10.0.0.3",
							"server_ip":   "10.0.0.3",
							"server_type": "Ip",
						},
					},
					"server_group_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler":         "5TCH",
						"protocol":          "GENEVE",
						"server_group_type": "Ip",
						"resource_group_id": acctest.CHECKSET,
						"vpc_id":            acctest.CHECKSET,
						"dry_run":           "false",
						"servers.#":         "3",
						"server_group_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": "3TCH",
					"health_check_config": []map[string]interface{}{
						{
							"health_check_protocol":        "HTTP",
							"health_check_connect_port":    "81",
							"health_check_connect_timeout": "6",
							"health_check_domain":          "www.domain-update.com",
							"health_check_enabled":         "false",
							"health_check_http_code": []string{
								"http_5xx"},
							"health_check_interval": "11",
							"health_check_path":     "/health-check-update",
							"healthy_threshold":     "3",
							"unhealthy_threshold":   "3",
						},
					},
					"connection_drain_config": []map[string]interface{}{
						{
							"connection_drain_enabled": "false",
							"connection_drain_timeout": "2",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"servers": []map[string]interface{}{
						{
							"server_id":   "10.0.0.4",
							"server_ip":   "10.0.0.4",
							"server_type": "Ip",
						},
					},
					"server_group_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler":         "3TCH",
						"resource_group_id": acctest.CHECKSET,
						"servers.#":         "1",
						"server_group_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": "5TCH",
					"health_check_config": []map[string]interface{}{
						{
							"health_check_protocol":        "TCP",
							"health_check_connect_port":    "80",
							"health_check_connect_timeout": "5",
							"health_check_enabled":         "true",
							"health_check_http_code":       []string{},
							"health_check_interval":        "10",
							"healthy_threshold":            "2",
							"unhealthy_threshold":          "2",
						},
					},
					"connection_drain_config": []map[string]interface{}{
						{
							"connection_drain_enabled": "true",
							"connection_drain_timeout": "5",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"servers":           acctest.REMOVEKEY,
					"server_group_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler":         "5TCH",
						"resource_group_id": acctest.CHECKSET,
						"servers.#":         "0",
						"server_group_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudGwlbServerGroupMap8500 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

func AlicloudGwlbServerGroupBasicDependence8500(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-wulanchabu"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultEaxcvb" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = "tf-gwlb-vpc"
}


`, name)
}

// Case ServerGroup Test_IP_依赖资源 8564
func TestAccAliCloudGwlbServerGroup_basic8564(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gwlb_server_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGwlbServerGroupMap8564)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gwlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeServerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccgwlb%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudGwlbServerGroupBasicDependence8564)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": "5TCH",
					"health_check_config": []map[string]interface{}{
						{
							"health_check_protocol":        "HTTP",
							"health_check_connect_port":    "80",
							"health_check_connect_timeout": "5",
							"health_check_domain":          "www.domain.com",
							"health_check_enabled":         "true",
							"health_check_http_code": []string{
								"http_2xx", "http_3xx", "http_4xx"},
							"health_check_interval": "10",
							"health_check_path":     "/health-check",
							"healthy_threshold":     "2",
							"unhealthy_threshold":   "2",
						},
					},
					"protocol":          "GENEVE",
					"server_group_type": "Ip",
					"connection_drain_config": []map[string]interface{}{
						{
							"connection_drain_enabled": "true",
							"connection_drain_timeout": "1",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"vpc_id":            "${alicloud_vpc.defaultEaxcvb.id}",
					"dry_run":           "false",
					"servers": []map[string]interface{}{
						{
							"server_id":   "10.0.0.1",
							"server_ip":   "10.0.0.1",
							"server_type": "Ip",
						},
						{
							"server_id":   "10.0.0.2",
							"server_ip":   "10.0.0.2",
							"server_type": "Ip",
						},
						{
							"server_id":   "10.0.0.3",
							"server_ip":   "10.0.0.3",
							"server_type": "Ip",
						},
					},
					"server_group_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler":         "5TCH",
						"protocol":          "GENEVE",
						"server_group_type": "Ip",
						"resource_group_id": acctest.CHECKSET,
						"vpc_id":            acctest.CHECKSET,
						"dry_run":           "false",
						"servers.#":         "3",
						"server_group_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": "3TCH",
					"health_check_config": []map[string]interface{}{
						{
							"health_check_protocol":        "HTTP",
							"health_check_connect_port":    "81",
							"health_check_connect_timeout": "6",
							"health_check_domain":          "www.domain-update.com",
							"health_check_enabled":         "false",
							"health_check_http_code": []string{
								"http_5xx"},
							"health_check_interval": "11",
							"health_check_path":     "/health-check-update",
							"healthy_threshold":     "3",
							"unhealthy_threshold":   "3",
						},
					},
					"connection_drain_config": []map[string]interface{}{
						{
							"connection_drain_enabled": "false",
							"connection_drain_timeout": "2",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"servers": []map[string]interface{}{
						{
							"server_id":   "10.0.0.4",
							"server_ip":   "10.0.0.4",
							"server_type": "Ip",
						},
					},
					"server_group_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler":         "3TCH",
						"resource_group_id": acctest.CHECKSET,
						"servers.#":         "1",
						"server_group_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": "5TCH",
					"health_check_config": []map[string]interface{}{
						{
							"health_check_protocol":        "TCP",
							"health_check_connect_port":    "80",
							"health_check_connect_timeout": "5",
							"health_check_enabled":         "true",
							"health_check_http_code":       []string{},
							"health_check_interval":        "10",
							"healthy_threshold":            "2",
							"unhealthy_threshold":          "2",
						},
					},
					"connection_drain_config": []map[string]interface{}{
						{
							"connection_drain_enabled": "true",
							"connection_drain_timeout": "5",
						},
					},
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"servers":           acctest.REMOVEKEY,
					"server_group_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler":         "5TCH",
						"resource_group_id": acctest.CHECKSET,
						"servers.#":         "0",
						"server_group_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": acctest.REMOVEKEY,
						"tags.For":     acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudGwlbServerGroupMap8564 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

func AlicloudGwlbServerGroupBasicDependence8564(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-wulanchabu"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultEaxcvb" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = "tf-gwlb-vpc"
}


`, name)
}

// Case ServerGroup Test_draining 8570
func TestAccAliCloudGwlbServerGroupDraining8570(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gwlb_server_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGwlbServerGroupMap8570)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gwlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeServerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccgwlb%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudGwlbServerGroupBasicDependence8570)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// create two servers with connection drain enabled (1h timeout keeps the
				// Draining status stable across all the following steps)
				Config: testAccConfig(map[string]interface{}{
					"protocol":          "GENEVE",
					"server_group_type": "Ip",
					"connection_drain_config": []map[string]interface{}{
						{
							"connection_drain_enabled": "true",
							"connection_drain_timeout": "3600",
						},
					},
					"vpc_id": "${alicloud_vpc.defaultEaxcvb.id}",
					"servers": []map[string]interface{}{
						{
							"server_id":   "10.0.0.1",
							"server_ip":   "10.0.0.1",
							"server_type": "Ip",
						},
						{
							"server_id":   "10.0.0.2",
							"server_ip":   "10.0.0.2",
							"server_type": "Ip",
						},
					},
					"server_group_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"protocol":           "GENEVE",
						"server_group_type":  "Ip",
						"vpc_id":             acctest.CHECKSET,
						"servers.#":          "2",
						"draining_servers.#": "0",
						"server_group_name":  name,
					}),
				),
			},
			{
				// remove 10.0.0.1: with drain enabled it enters Draining and is no longer
				// kept in servers (only Available servers are persisted); it is exposed
				// through the draining_servers field instead
				Config: testAccConfig(map[string]interface{}{
					"connection_drain_config": []map[string]interface{}{
						{
							"connection_drain_enabled": "true",
							"connection_drain_timeout": "3600",
						},
					},
					"vpc_id": "${alicloud_vpc.defaultEaxcvb.id}",
					"servers": []map[string]interface{}{
						{
							"server_id":   "10.0.0.2",
							"server_ip":   "10.0.0.2",
							"server_type": "Ip",
						},
					},
					"server_group_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#":          "1",
						"draining_servers.#": "1",
					}),
				),
			},
			{
				// remove 10.0.0.2 while 10.0.0.1 is still Draining: the provider must skip the
				// Draining server in RemoveServersFromServerGroup, otherwise the API rejects the
				// request with an IncorrectStatus error. Both servers are now out of state.
				// acctest.REMOVEKEY drops the servers key from the generated config; an empty map
				// (acctest.CLEARMAP) or list renders invalid HCL for a block-typed attribute.
				Config: testAccConfig(map[string]interface{}{
					"connection_drain_config": []map[string]interface{}{
						{
							"connection_drain_enabled": "true",
							"connection_drain_timeout": "3600",
						},
					},
					"vpc_id":            "${alicloud_vpc.defaultEaxcvb.id}",
					"servers":           acctest.REMOVEKEY,
					"server_group_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#":          "0",
						"draining_servers.#": "2",
					}),
				),
			},
			{
				// re-add the Draining server 10.0.0.1: it must be rescued via AddServers
				Config: testAccConfig(map[string]interface{}{
					"connection_drain_config": []map[string]interface{}{
						{
							"connection_drain_enabled": "true",
							"connection_drain_timeout": "3600",
						},
					},
					"vpc_id": "${alicloud_vpc.defaultEaxcvb.id}",
					"servers": []map[string]interface{}{
						{
							"server_id":   "10.0.0.1",
							"server_ip":   "10.0.0.1",
							"server_type": "Ip",
						},
					},
					"server_group_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#":          "1",
						"draining_servers.#": "1",
					}),
				),
			},
			{
				// re-add 10.0.0.2 as well: both servers back to Available
				Config: testAccConfig(map[string]interface{}{
					"connection_drain_config": []map[string]interface{}{
						{
							"connection_drain_enabled": "true",
							"connection_drain_timeout": "3600",
						},
					},
					"vpc_id": "${alicloud_vpc.defaultEaxcvb.id}",
					"servers": []map[string]interface{}{
						{
							"server_id":   "10.0.0.1",
							"server_ip":   "10.0.0.1",
							"server_type": "Ip",
						},
						{
							"server_id":   "10.0.0.2",
							"server_ip":   "10.0.0.2",
							"server_type": "Ip",
						},
					},
					"server_group_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#":          "2",
						"draining_servers.#": "0",
					}),
				),
			},
			{
				// disable connection drain and clear the servers so the destroy is clean;
				// servers dropped from the config via acctest.REMOVEKEY (empty map/list renders invalid HCL)
				Config: testAccConfig(map[string]interface{}{
					"connection_drain_config": []map[string]interface{}{
						{
							"connection_drain_enabled": "false",
							"connection_drain_timeout": "3600",
						},
					},
					"vpc_id":            "${alicloud_vpc.defaultEaxcvb.id}",
					"servers":           acctest.REMOVEKEY,
					"server_group_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#":          "0",
						"draining_servers.#": acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

var AlicloudGwlbServerGroupMap8570 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

func AlicloudGwlbServerGroupBasicDependence8570(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultEaxcvb" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = "tf-gwlb-vpc"
}


`, name)
}

// Case ServerGroup Test_health_check_zero_value 8801
// Regression for the health_check_config zero-value defect: when only
// health_check_enabled is configured, the jsonpath lookup returned 0 for the
// unconfigured integer sub-attributes and 0 != "" is true, so the Provider sent
// HealthCheckConnectPort/Timeout/Interval/Healthy/UnhealthyThreshold = 0 in the
// Create/Update request. The guards below assert the server-side defaults
// (80/5/10/2/2 per docs) are used instead of 0, and that an explicit
// health_check_enabled=false is still sent.
func TestAccAliCloudGwlbServerGroup_healthCheckZeroValue8801(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gwlb_server_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGwlbServerGroupMap8801)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gwlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeServerGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccgwlb%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudGwlbServerGroupBasicDependence8801)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// AC1: create with only health_check_enabled=false; the unconfigured
				// integer sub-attributes must NOT be sent as 0 and the server applies its
				// documented defaults (port 80 / timeout 5 / interval 10 / 2 / 2).
				Config: testAccConfig(map[string]interface{}{
					"protocol":          "GENEVE",
					"server_group_type": "Ip",
					"vpc_id":            "${alicloud_vpc.defaultEaxcvb.id}",
					"dry_run":           "false",
					"health_check_config": []map[string]interface{}{
						{
							"health_check_enabled": "false",
						},
					},
					"servers": []map[string]interface{}{
						{
							"server_id":   "10.0.0.1",
							"server_ip":   "10.0.0.1",
							"server_type": "Ip",
						},
					},
					"server_group_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_config.0.health_check_enabled":         "false",
						"health_check_config.0.health_check_connect_port":    "80",
						"health_check_config.0.health_check_connect_timeout": "5",
						"health_check_config.0.health_check_interval":        "10",
						"health_check_config.0.healthy_threshold":            "2",
						"health_check_config.0.unhealthy_threshold":          "2",
						"health_check_config.0.health_check_http_code.#":     "0",
					}),
				),
			},
			{
				// AC3: full health_check_config regression (HTTP protocol, port, timeout,
				// domain, http codes, interval, path, thresholds, enabled).
				Config: testAccConfig(map[string]interface{}{
					"health_check_config": []map[string]interface{}{
						{
							"health_check_protocol":        "HTTP",
							"health_check_connect_port":    "80",
							"health_check_connect_timeout": "5",
							"health_check_domain":          "www.domain.com",
							"health_check_enabled":         "true",
							"health_check_http_code":       []string{"http_2xx", "http_3xx", "http_4xx"},
							"health_check_interval":        "10",
							"health_check_path":            "/health-check",
							"healthy_threshold":            "2",
							"unhealthy_threshold":          "2",
						},
					},
					"servers": []map[string]interface{}{
						{
							"server_id":   "10.0.0.2",
							"server_ip":   "10.0.0.2",
							"server_type": "Ip",
						},
					},
					"server_group_name": name + "_full",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_config.0.health_check_enabled":         "true",
						"health_check_config.0.health_check_protocol":        "HTTP",
						"health_check_config.0.health_check_connect_port":    "80",
						"health_check_config.0.health_check_connect_timeout": "5",
						"health_check_config.0.health_check_interval":        "10",
						"health_check_config.0.healthy_threshold":            "2",
						"health_check_config.0.unhealthy_threshold":          "2",
						"health_check_config.0.health_check_http_code.#":     "3",
						"health_check_config.0.health_check_path":            "/health-check",
						"health_check_config.0.health_check_domain":          "www.domain.com",
						"server_group_name":                                  name + "_full",
					}),
				),
			},
			{
				// Clear health_check_http_code while enabled and on TCP (the API clears the
				// code list only when the check is active); this prevents a permanent diff
				// when the next step drops the http_code attribute while disabled.
				Config: testAccConfig(map[string]interface{}{
					"health_check_config": []map[string]interface{}{
						{
							"health_check_protocol":        "TCP",
							"health_check_connect_port":    "80",
							"health_check_connect_timeout": "5",
							"health_check_enabled":         "true",
							"health_check_http_code":       []string{},
							"health_check_interval":        "10",
							"healthy_threshold":            "2",
							"unhealthy_threshold":          "2",
						},
					},
					"servers": []map[string]interface{}{
						{
							"server_id":   "10.0.0.2",
							"server_ip":   "10.0.0.2",
							"server_type": "Ip",
						},
					},
					"server_group_name": name + "_full",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_config.0.health_check_protocol":    "TCP",
						"health_check_config.0.health_check_enabled":     "true",
						"health_check_config.0.health_check_http_code.#": "0",
						// TCP clears the HTTP-only domain/path attributes server-side;
						// drop them from the accumulated expected set so the step does
						// not keep asserting the prior HTTP values.
						"health_check_config.0.health_check_domain": acctest.REMOVEKEY,
						"health_check_config.0.health_check_path":   acctest.REMOVEKEY,
					}),
				),
			},
			{
				// AC2: switch back to only health_check_enabled=false. The integer
				// sub-attributes are now omitted from the config; the Provider must not
				// send 0, so the server keeps the previously applied non-zero values
				// instead of being reset to 0.
				Config: testAccConfig(map[string]interface{}{
					"health_check_config": []map[string]interface{}{
						{
							"health_check_enabled": "false",
						},
					},
					"servers": []map[string]interface{}{
						{
							"server_id":   "10.0.0.2",
							"server_ip":   "10.0.0.2",
							"server_type": "Ip",
						},
					},
					"server_group_name": name + "_full",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_config.0.health_check_enabled":         "false",
						"health_check_config.0.health_check_connect_port":    "80",
						"health_check_config.0.health_check_connect_timeout": "5",
						"health_check_config.0.health_check_interval":        "10",
						"health_check_config.0.healthy_threshold":            "2",
						"health_check_config.0.unhealthy_threshold":          "2",
						"health_check_config.0.health_check_http_code.#":     "0",
						// domain/path were already dropped from the accumulated
						// expected set by the prior TCP step; do not re-reference them
						// here, or REMOVEKEY re-adds the literal sentinel value.
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dry_run"},
			},
		},
	})
}

var AlicloudGwlbServerGroupMap8801 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

func AlicloudGwlbServerGroupBasicDependence8801(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-wulanchabu"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultEaxcvb" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = "tf-gwlb-vpc"
}


`, name)
}

// Test Gwlb ServerGroup. <<< Resource test cases, automatically generated.
