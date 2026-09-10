package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Apig Source. >>> Resource test cases, hand-written.
func TestAccAliCloudApigSource_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_source.default"
	ra := resourceAttrInit(resourceId, AlicloudApigSourceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapigsource%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigSourceBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"type":              "K8S",
					"gateway_id":        "${alicloud_apig_gateway.default.id}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"k8s_source_info": []map[string]interface{}{
						{
							"cluster_id": "${alicloud_cs_managed_kubernetes.default.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":              "K8S",
						"gateway_id":        CHECKSET,
						"resource_group_id": CHECKSET,
						"k8s_source_info.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
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

var AlicloudApigSourceMap = map[string]string{
	"create_time": CHECKSET,
	"source_name": CHECKSET,
	"update_time": CHECKSET,
}

func AlicloudApigSourceBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "default" {
    vpc_name   = var.name
    cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
    vswitch_name = var.name
    vpc_id       = alicloud_vpc.default.id
    zone_id      = "cn-hangzhou-i"
    cidr_block   = "192.168.8.0/24"
}

resource "alicloud_cs_managed_kubernetes" "default" {
    name                 = var.name
    cluster_spec         = "ack.pro.small"
    worker_vswitch_ids   = [alicloud_vswitch.default.id]
    pod_cidr             = "10.95.0.0/16"
    service_cidr         = "172.23.0.0/16"
    slb_internet_enabled = true
    new_nat_gateway      = true
    deletion_protection  = false
}

resource "alicloud_apig_gateway" "default" {
    gateway_name    = var.name
    spec            = "apigw.small.x1"
    gateway_edition = "Professional"
    gateway_type    = "API"
    payment_type    = "PayAsYouGo"
    vpc {
        vpc_id = alicloud_vpc.default.id
    }
    vswitch {
        vswitch_id = alicloud_vswitch.default.id
    }
    network_access_config {
        type = "Internet"
    }
    zone_config {
        select_option = "Auto"
    }
    log_config {
        sls {
            enable = false
        }
    }
}

`, name)
}

func TestAccAliCloudApigSource_nacos(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_source.default"
	ra := resourceAttrInit(resourceId, AlicloudApigSourceNacosMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigSource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapigsource%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigSourceNacosDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"type":              "MSE_NACOS",
					"gateway_id":        "${alicloud_apig_gateway.default.id}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"nacos_source_info": []map[string]interface{}{
						{
							"instance_id": "${alicloud_mse_cluster.default.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":                "MSE_NACOS",
						"gateway_id":          CHECKSET,
						"resource_group_id":   CHECKSET,
						"nacos_source_info.#": "1",
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

var AlicloudApigSourceNacosMap = map[string]string{
	"create_time": CHECKSET,
	"source_name": CHECKSET,
	"update_time": CHECKSET,
}

func AlicloudApigSourceNacosDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "default" {
    vpc_name   = var.name
    cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
    vswitch_name = var.name
    vpc_id       = alicloud_vpc.default.id
    zone_id      = "cn-hangzhou-i"
    cidr_block   = "192.168.8.0/24"
}

resource "alicloud_mse_cluster" "default" {
    cluster_specification = "MSE_SC_1_2_60_c"
    cluster_type          = "Nacos-Ans"
    cluster_version       = "NACOS_2_0_0"
    instance_count        = 3
    net_type              = "privatenet"
    vswitch_id            = alicloud_vswitch.default.id
    connection_type       = "slb"
    pub_network_flow      = "1"
    mse_version           = "mse_pro"
    vpc_id                = alicloud_vpc.default.id
}

resource "alicloud_apig_gateway" "default" {
    gateway_name    = var.name
    spec            = "apigw.small.x1"
    gateway_edition = "Professional"
    gateway_type    = "API"
    payment_type    = "PayAsYouGo"
    vpc {
        vpc_id = alicloud_vpc.default.id
    }
    vswitch {
        vswitch_id = alicloud_vswitch.default.id
    }
    network_access_config {
        type = "Internet"
    }
    zone_config {
        select_option = "Auto"
    }
    log_config {
        sls {
            enable = false
        }
    }
}

`, name)
}

// Test Apig Source. <<< Resource test cases, hand-written.
