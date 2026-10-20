package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test Gwlb ServerGroup. >>> Resource test cases, automatically generated.
// Case ServerGroup Test_instance 8419
func TestAccAliCloudGwlbServerGroup_basic8419(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gwlb_server_group.default"
	ra := resourceAttrInit(resourceId, AlicloudGwlbServerGroupMap8419)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GwlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGwlbServerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccgwlb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGwlbServerGroupBasicDependence8419)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
						"resource_group_id":    CHECKSET,
						"vpc_id":               CHECKSET,
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
						"resource_group_id":    CHECKSET,
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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
	"status":      CHECKSET,
	"create_time": CHECKSET,
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
	ra := resourceAttrInit(resourceId, AlicloudGwlbServerGroupMap8500)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GwlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGwlbServerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccgwlb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGwlbServerGroupBasicDependence8500)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
						"resource_group_id": CHECKSET,
						"vpc_id":            CHECKSET,
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
						"resource_group_id": CHECKSET,
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
					"servers":           REMOVEKEY,
					"server_group_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler":         "5TCH",
						"resource_group_id": CHECKSET,
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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
	"status":      CHECKSET,
	"create_time": CHECKSET,
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
	ra := resourceAttrInit(resourceId, AlicloudGwlbServerGroupMap8564)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GwlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGwlbServerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccgwlb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGwlbServerGroupBasicDependence8564)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
						"resource_group_id": CHECKSET,
						"vpc_id":            CHECKSET,
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
						"resource_group_id": CHECKSET,
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
					"servers":           REMOVEKEY,
					"server_group_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler":         "5TCH",
						"resource_group_id": CHECKSET,
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
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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
	"status":      CHECKSET,
	"create_time": CHECKSET,
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
	ra := resourceAttrInit(resourceId, AlicloudGwlbServerGroupMap8570)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GwlbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGwlbServerGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccgwlb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGwlbServerGroupBasicDependence8570)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"vpc_id":             CHECKSET,
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
				// REMOVEKEY drops the servers key from the generated config; an empty map
				// (CLEARMAP) or list renders invalid HCL for a block-typed attribute.
				Config: testAccConfig(map[string]interface{}{
					"connection_drain_config": []map[string]interface{}{
						{
							"connection_drain_enabled": "true",
							"connection_drain_timeout": "3600",
						},
					},
					"vpc_id":            "${alicloud_vpc.defaultEaxcvb.id}",
					"servers":           REMOVEKEY,
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
				// servers dropped from the config via REMOVEKEY (empty map/list renders invalid HCL)
				Config: testAccConfig(map[string]interface{}{
					"connection_drain_config": []map[string]interface{}{
						{
							"connection_drain_enabled": "false",
							"connection_drain_timeout": "3600",
						},
					},
					"vpc_id":            "${alicloud_vpc.defaultEaxcvb.id}",
					"servers":           REMOVEKEY,
					"server_group_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#":          "0",
						"draining_servers.#": CHECKSET,
					}),
				),
			},
		},
	})
}

var AlicloudGwlbServerGroupMap8570 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
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

// Test Gwlb ServerGroup. <<< Resource test cases, automatically generated.
