package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
	name := fmt.Sprintf("tf-testacc%sgwlbservergroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGwlbServerGroupBasicDependence8419)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-wulanchabu"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler":         "5TCH",
					"protocol":          "GENEVE",
					"server_group_type": "Instance",
					"vpc_id":            "${alicloud_vpc.defaultEaxcvb.id}",
					"dry_run":           "false",
					"server_group_name": name,
					"servers": []map[string]interface{}{
						{
							"server_id":   "${alicloud_instance.default5DqP8f.id}",
							"server_type": "Ecs",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler":         "5TCH",
						"protocol":          "GENEVE",
						"server_group_type": "Instance",
						"vpc_id":            CHECKSET,
						"dry_run":           "false",
						"server_group_name": name,
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
							"health_check_domain":          "www.domain.com",
							"health_check_enabled":         "true",
							"health_check_http_code": []string{
								"http_2xx", "http_4xx", "http_3xx"},
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
					"vpc_id":            "${alicloud_vpc.defaultEaxcvb.id}",
					"dry_run":           "false",
					"server_group_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler":         "5TCH",
						"protocol":          "GENEVE",
						"server_group_type": "Instance",
						"vpc_id":            CHECKSET,
						"dry_run":           "false",
						"server_group_name": name,
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
							"health_check_domain":          "www.domain.com",
							"health_check_enabled":         "true",
							"health_check_http_code": []string{
								"http_2xx", "http_4xx", "http_3xx"},
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
					"vpc_id":  "${alicloud_vpc.defaultEaxcvb.id}",
					"dry_run": "false",
					"servers": []map[string]interface{}{
						{
							"server_id":   "${alicloud_instance.default5DqP8f.id}",
							"server_type": "Ecs",
						},
					},
					"server_group_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scheduler":         "5TCH",
						"protocol":          "GENEVE",
						"server_group_type": "Instance",
						"vpc_id":            CHECKSET,
						"dry_run":           "false",
						"servers.#":         "1",
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
					"servers": []map[string]interface{}{
						{
							"server_type": "Eni",
							"server_id":   "${alicloud_instance.default5DqP8f.network_interface_id}",
							"server_ip":   "${alicloud_instance.default5DqP8f.primary_ip_address}",
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

variable "zone_id1" {
  default = "cn-wulanchabu-b"
}

resource "alicloud_vpc" "defaultEaxcvb" {
  cidr_block = "10.0.0.0/8"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "defaultc3uVID" {
  vpc_id       = alicloud_vpc.defaultEaxcvb.id
  zone_id      = var.zone_id1
  cidr_block   = "10.0.0.0/24"
  vswitch_name = format("%%s3", var.name)
}

resource "alicloud_security_group" "default" {
  name = "tf-test"
  description = "New security group"
  vpc_id = alicloud_vpc.defaultEaxcvb.id
}

resource "alicloud_instance" "default5DqP8f" {
	vswitch_id = alicloud_vswitch.defaultc3uVID.id
	image_id = "aliyun_2_1903_x64_20G_alibase_20231221.vhd"
	
	instance_type = "ecs.g6.large"
	system_disk_category = "cloud_efficiency"
	
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 5
	instance_name = format("%%s4", var.name)
    description   = "tf-test-ecs"
    security_groups = [alicloud_security_group.default.id]

  	availability_zone          = alicloud_vswitch.defaultc3uVID.zone_id
  	instance_charge_type       = "PostPaid"
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
	name := fmt.Sprintf("tf-testacc%sgwlbservergroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGwlbServerGroupBasicDependence8564)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-wulanchabu"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": "5TCH",
					"health_check_config": []map[string]interface{}{
						{
							"health_check_protocol":        "TCP",
							"health_check_connect_port":    "80",
							"health_check_connect_timeout": "5",
							"health_check_domain":          "www.domain.com",
							"health_check_enabled":         "true",
							"health_check_http_code": []string{
								"http_2xx", "http_4xx", "http_3xx"},
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
							"health_check_domain":          "www.domain.com",
							"health_check_enabled":         "true",
							"health_check_http_code":       []string{"http_5xx"},
							"health_check_interval":        "10",
							"health_check_path":            "/health-check",
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
  vpc_name   = var.name
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
	name := fmt.Sprintf("tf-testacc%sgwlbservergroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGwlbServerGroupBasicDependence8500)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-wulanchabu"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"scheduler": "5TCH",
					"health_check_config": []map[string]interface{}{
						{
							"health_check_protocol":        "TCP",
							"health_check_connect_port":    "80",
							"health_check_connect_timeout": "5",
							"health_check_domain":          "www.domain.com",
							"health_check_enabled":         "true",
							"health_check_http_code": []string{
								"http_2xx", "http_4xx", "http_3xx"},
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
							"health_check_domain":          "www.domain.com",
							"health_check_enabled":         "true",
							"health_check_http_code":       []string{"http_5xx"},
							"health_check_interval":        "10",
							"health_check_path":            "/health-check",
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
  vpc_name   = var.name
}


`, name)
}

// Test Gwlb ServerGroup. <<< Resource test cases, automatically generated.
