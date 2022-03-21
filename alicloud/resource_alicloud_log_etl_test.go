package alicloud

import (
	"fmt"
	"testing"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudLogETL_basic(t *testing.T) {
	var v *sls.ETL
	resourceId := "alicloud_log_etl.default"
	ra := resourceAttrInit(resourceId, logETLMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogossshipper-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogETLConfigDependence)
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
					"etl_name":     name,
					"project":      "${alicloud_log_project.default.name}",
					"display_name": "etl_display",
					"description":  "etl_description",
					"from_time":    "1616486027",
					"to_time":      "1617486027",
					"parameters": map[string]string{
						"test1": "test2",
					},
					"access_key_id":     "access_key_id_test",
					"access_key_secret": "access_key_secret_test",
					"script":            "e_set('new','test')",
					"logstore":          "${alicloud_log_store.default.name}",
					"etl_sinks": []map[string]interface{}{
						{
							"access_key_id":     "test1",
							"access_key_secret": "test2",
							"endpoint":          "cn-hangzhou.log.aliyuncs.com",
							"name":              "target_name",
							"project":           "${alicloud_log_project.default.name}",
							"logstore":          "${alicloud_log_store.default1.name}",
						},
						{
							"access_key_id":     "test1",
							"access_key_secret": "test2",
							"endpoint":          "cn-hangzhou.log.aliyuncs.com",
							"name":              "target_name_2",
							"project":           "${alicloud_log_project.default.name}",
							"logstore":          "${alicloud_log_store.default2.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"etl_name":         name,
						"project":          name,
						"from_time":        "1616486027",
						"to_time":          "1617486027",
						"etl_sinks.#":      "2",
						"parameters.%":     "1",
						"parameters.test1": "test2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"access_key_id", "access_key_secret"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":  "update",
					"display_name": "update_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":  "update",
						"display_name": "update_name",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "STOPPED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "STOPPED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "RUNNING",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "RUNNING",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"script": "e_set('aliyun','aliyun')",
					"parameters": map[string]string{
						"update": "update",
					},
					"etl_sinks": []map[string]interface{}{
						{
							"access_key_id":     "test1",
							"access_key_secret": "test2",
							"endpoint":          "cn-hangzhou.log.aliyuncs.com",
							"name":              "target_name",
							"project":           "${alicloud_log_project.default.name}",
							"logstore":          "${alicloud_log_store.default1.name}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"script":            "e_set('aliyun','aliyun')",
						"parameters.update": "update",
						"parameters.test1":  REMOVEKEY,
						"etl_sinks.#":       "1",
					}),
				),
			},
		},
	})
}

func resourceLogETLConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	    default = "%s"
	}
	resource "alicloud_log_project" "default" {
	    name = "${var.name}"
	    description = "tf unit test"
	}
	resource "alicloud_log_store" "default" {
	    project = "${alicloud_log_project.default.name}"
	    name = "${var.name}"
	    retention_period = "3000"
	    shard_count = 1
	}
	resource "alicloud_log_store" "default1" {
	    project = "${alicloud_log_project.default.name}"
	    name = "${var.name}1"
	    retention_period = "3000"
	    shard_count = 1
	}
	resource "alicloud_log_store" "default2" {
	    project = "${alicloud_log_project.default.name}"
	    name = "${var.name}2"
	    retention_period = "3000"
	    shard_count = 1
	}
	`, name)
}

var logETLMap = map[string]string{
	"project": CHECKSET,
}
