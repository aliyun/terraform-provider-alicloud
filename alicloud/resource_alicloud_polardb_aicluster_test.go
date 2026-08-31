package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudPolarDBAICluster_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_aicluster.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"status": CHECKSET,
	})
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBAIClusterAttribute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccPolarDBAICluster-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBAIClusterConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"region_id":              "cn-beijing",
					"zone_id":                "${data.alicloud_polardb_node_classes.this.classes.1.zone_id}",
					"db_node_class":          "polar.mysql.g8.4xlarge.gu50",
					"db_cluster_description": name,
					"pay_type":               "Postpaid",
					"vswitch_id":             "${local.vswitch_id}",
					"vpc_id":                 "${local.vpc_id}",
					"kube_type":              "ainode",
					"model_name":             "GLM-4.7",
					"extension":              "maas",
					"inference_engine":       "sglang",
					"db_cluster_id":          "${alicloud_polardb_cluster.cluster.id}",
					"security_group_id":      "${alicloud_security_group.default.id}",
					"auto_use_coupon":        "true",
					"auto_renew":             "false",
					"period":                 "Month",
					"used_time":              "1",
					"promotion_code":         "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"region_id":              CHECKSET,
						"zone_id":                CHECKSET,
						"db_node_class":          "polar.mysql.g8.4xlarge.gu50",
						"db_cluster_description": name,
						"pay_type":               "Postpaid",
						"vpc_id":                 CHECKSET,
						"vswitch_id":             CHECKSET,
						"kube_type":              "ainode",
						"model_name":             "GLM-4.7",
						"extension":              "maas",
						"inference_engine":       "sglang",
						"status":                 CHECKSET,
						"auto_use_coupon":        "true",
						"auto_renew":             "false",
						"period":                 "Month",
						"used_time":              "1",
						"promotion_code":         "",
						"connection_string":      CHECKSET,
						"api_key":                CHECKSET,
					}),
				),
			},
		},
	})
}

func resourcePolarDBAIClusterConfigDependence(name string) string {
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
	
	data "alicloud_polardb_node_classes" "this" {
		db_type    = "MySQL"
		db_version = "8.0"
		pay_type   = "PostPaid"
		category   = "Normal"
	}

	locals {
		vpc_id     = alicloud_vpc.default.id
		vswitch_id = alicloud_vswitch.default.id
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

	resource "alicloud_security_group" "default" {
		name   = var.name
		vpc_id = "${local.vpc_id}"
	}
`, name)
}
