package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudPolarDBApplication_Create(t *testing.T) {
	v := map[string]interface{}{}
	rand := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBApplication-%s", rand)
	resourceId := "alicloud_polardb_application.default"
	regionId := os.Getenv("ALICLOUD_REGION")
	var basicMap = map[string]string{
		"vswitch_id": CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBApplicationAttribute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBApplicationConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":      "${var.name}",
					"application_type": "polarclaw",
					"architecture":     "x86",
					"pay_type":         "PostPaid",
					"region_id":        regionId,
					"vswitch_id":       "${local.vswitch_id}",
					"vpc_id":           "${local.vpc_id}",
					"zone_id":          "${data.alicloud_polardb_node_classes.this.classes.1.zone_id}",
					"model_from":       "bailian",
					"model_base_url":   "https://dashscope.aliyuncs.com/compatible-mode/v1",
					"model_name":       "qwen3.6-plus",
					"components": []map[string]interface{}{
						{
							"component_type":    "polarclaw_comp",
							"component_class":   "polar.app.g2.medium",
							"component_replica": 1,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Activated",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"upgrade_version": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"upgrade_version": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"parameters": []map[string]interface{}{
						{
							"parameter_name":  "secret.dashscope.apiKey",
							"parameter_value": "abc",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Activated",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modify_mode":            "Append",
					"security_ip_array_name": "test",
					"security_ip_list":       "127.0.0.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_array_name": "test",
						"security_ip_list":       "127.0.0.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"modify_mode":            "Delete",
					"security_ip_array_name": "test",
					"security_ip_list":       "127.0.0.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_array_name": "test",
						"security_ip_list":       "127.0.0.1",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudPolarDBApplication_CreateFull(t *testing.T) {
	v := map[string]interface{}{}
	rand := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBApplication-%s", rand)
	resourceId := "alicloud_polardb_application.default"
	regionId := os.Getenv("ALICLOUD_REGION")
	var basicMap = map[string]string{
		"vswitch_id": CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBApplicationAttribute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBApplicationFullConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":      "${var.name}",
					"application_type": "polarclaw",
					"architecture":     "x86",
					"pay_type":         "PostPaid",
					"region_id":        regionId,
					"vswitch_id":       "${local.vswitch_id}",
					"vpc_id":           "${local.vpc_id}",
					"zone_id":          "${data.alicloud_polardb_node_classes.this.classes.1.zone_id}",
					"model_from":       "bailian",
					"model_base_url":   "https://dashscope.aliyuncs.com/compatible-mode/v1",
					"model_name":       "qwen3.6-plus",
					"components": []map[string]interface{}{
						{
							"component_type":    "polarclaw_comp",
							"component_class":   "polar.app.g2.medium",
							"component_replica": 1,
						},
					},
					"auto_renew":        false,
					"period":            1,
					"model_api_key":     "sk-xxx",
					"model_api":         "openai-completions",
					"ai_db_cluster_id":  "${alicloud_polardb_cluster.cluster.id}",
					"db_cluster_id":     "${alicloud_polardb_cluster.cluster.id}",
					"used_time":         1,
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Activated",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"component_id":           "",
					"modify_mode":            "Append",
					"security_ip_array_name": "test",
					"security_ip_list":       "127.0.0.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Activated",
					}),
				),
			},
		},
	})
}

func resourcePolarDBApplicationConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	resource "alicloud_vpc" "default" {
		vpc_name = var.name
	}

	resource "alicloud_vswitch" "default" {
		zone_id  = data.alicloud_polardb_node_classes.this.classes.1.zone_id
		vpc_id   = alicloud_vpc.default.id
		cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
	}

	locals {
		vpc_id     = alicloud_vpc.default.id
		vswitch_id = alicloud_vswitch.default.id

	}

	data "alicloud_polardb_node_classes" "this" {
		db_type    = "MySQL"
		db_version = "8.0"
		pay_type   = "PostPaid"
		category   = "Normal"
	}

`, name)
}

func resourcePolarDBApplicationFullConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	resource "alicloud_vpc" "default" {
		vpc_name = var.name
	}

	resource "alicloud_vswitch" "default" {
		zone_id  = data.alicloud_polardb_node_classes.this.classes.1.zone_id
		vpc_id   = alicloud_vpc.default.id
		cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
	}

	locals {
		vpc_id     = alicloud_vpc.default.id
		vswitch_id = alicloud_vswitch.default.id
		component_id = alicloud_polardb_application.default.component_id

	}

	data "alicloud_polardb_node_classes" "this" {
		db_type    = "MySQL"
		db_version = "8.0"
		pay_type   = "PostPaid"
		category   = "Normal"
	}

	data "alicloud_resource_manager_resource_groups" "default" {
		status = "OK"
	}

	resource "alicloud_polardb_cluster" "cluster" {
		db_type = "MySQL"
		db_version = "8.0"
		pay_type = "PostPaid"
        db_node_count = "2"
		db_node_class = "polar.mysql.x4.large"
		vswitch_id = "${local.vswitch_id}"
		description = "${var.name}"
	}

`, name)
}
