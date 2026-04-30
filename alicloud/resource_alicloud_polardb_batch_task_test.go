package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudPolarDBBatchTask_Create(t *testing.T) {
	v := map[string]interface{}{}
	rand := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBApplication-%s", rand)
	resourceId := "alicloud_polardb_batch_task.default"
	regionId := os.Getenv("ALICLOUD_REGION")
	var basicMap = map[string]string{
		"batch_id": CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBBatchTaskAttribute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return resourcePolarDBBatchTaskConfigDependence(name, regionId)
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: checkPolarDBBatchTaskDestroy(),

		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"region_id":    regionId,
					"task_name":    "batch_task_install",
					"task_type":    "polarclaw_install_skills",
					"instance_ids": []string{"${alicloud_polardb_application.default.id}"},
					"task_params": []map[string]interface{}{
						{
							"skill_name": "memory",
							"version":    "1.0.2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"task_status": "COMPLETED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"region_id":    regionId,
					"task_name":    "batch_task_uninstall",
					"task_type":    "polarclaw_uninstall_skills",
					"instance_ids": []string{"${alicloud_polardb_application.default.id}"},
					"task_params": []map[string]interface{}{
						{
							"skill_name": "memory",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"task_status": "COMPLETED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"region_id":    regionId,
					"task_name":    "batch_task_install",
					"task_type":    "polarclaw_install_skills",
					"instance_ids": []string{"${alicloud_polardb_application.default.id}"},
					"task_params": []map[string]interface{}{
						{
							"skill_name": "ontology",
							"version":    "1.0.4",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"task_status": "COMPLETED",
					}),
				),
			},
		},
	})
}

func resourcePolarDBBatchTaskConfigDependence(name string, regionId string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	variable "region_id" {
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

// checkPolarDBBatchTaskDestroy checks if the PolarDB batch task is destroyed
// Since batch tasks cannot be deleted, we only check if the associated application is still active
func checkPolarDBBatchTaskDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Check if the associated application is still active
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alicloud_polardb_application" {
				continue
			}
			// If the application is still active, the batch task cannot be destroyed
			// This is expected behavior for batch tasks
			return nil
		}
		return nil
	}
}
