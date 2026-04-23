package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {}

func TestAccAliCloudPolarDBApplicationEndpoint_Create(t *testing.T) {
	v := map[string]interface{}{}
	rand := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBApplication-%s", rand)
	resourceId := "alicloud_polardb_application_endpoint.default"
	regionId := os.Getenv("ALICLOUD_REGION")
	var basicMap = map[string]string{
		"application_id": CHECKSET,
		"endpoint_id":    CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBApplicationAttribute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return resourcePolarDBApplicationEndpointConfigDependence(name, regionId)
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"application_id": "${alicloud_polardb_application.default.id}",
					"endpoint_id":    "${alicloud_polardb_application.default.id}",
					"net_type":       "Public",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"net_type": "Public",
					}),
				),
			},
		},
	})
}

func resourcePolarDBApplicationEndpointConfigDependence(name string, regionId string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	variable "region_id" {
		default = "%s"
	}

	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
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
		vswitch_id = concat(alicloud_vswitch.default.*.id, [""])[0]

	}

	data "alicloud_polardb_node_classes" "this" {
		db_type    = "MySQL"
		db_version = "8.0"
		pay_type   = "PostPaid"
		category   = "Normal"
	}

    resource "alicloud_polardb_application" "default" {
        description      = "${var.name}"
		application_type = "polarclaw"
		architecture     = "x86"
		pay_type         = "PostPaid"
		region_id        = "${var.region_id}"
		vswitch_id       = "${local.vswitch_id}"
		vpc_id           = "${local.vpc_id}"
		zone_id          = "${data.alicloud_polardb_node_classes.this.classes.1.zone_id}"
		model_from       = "bailian"
		model_base_url   = "https://dashscope.aliyuncs.com/compatible-mode/v1"
		model_name       = "qwen3.6-plus"
		components {
			component_type    = "polarclaw_comp"
			component_class   = "polar.app.g2.medium"
			component_replica = 1
		}
    }

`, name, regionId)
}
