// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Dms Airflow. >>> Resource test cases, automatically generated.
// Case Airflow_tests 11207
func TestAccAliCloudDmsAirflow_basic11207(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dms_airflow.default"
	ra := resourceAttrInit(resourceId, AlicloudDmsAirflowMap11207)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDmsAirflow")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDmsAirflowBasicDependence11207)
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
					"worker_serverless_replicas": "0",
					"description":                "zhenyuan-test",
					"zone_id":                    "cn-hangzhou-h",
					"workspace_id":               "${alicloud_dms_enterprise_workspace.workspace.id}",
					"vpc_id":                     "${data.alicloud_vpcs.default.ids.0}",
					"oss_path":                   "/",
					"app_spec":                   "SMALL",
					"oss_bucket_name":            "hansheng",
					"airflow_name":               name,
					"security_group_id":          "${alicloud_security_group.security_group.id}",
					"vswitch_id":                 "${data.alicloud_vswitches.default.ids.0}",
					"requirement_file":           "default/requirements.txt",
					"plugins_dir":                "default/plugins",
					"dags_dir":                   "default/dags",
					"startup_file":               "default/startup.sh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"worker_serverless_replicas": "0",
						"description":                "zhenyuan-test",
						"zone_id":                    "cn-hangzhou-h",
						"workspace_id":               CHECKSET,
						"vpc_id":                     CHECKSET,
						"oss_path":                   "/",
						"app_spec":                   "SMALL",
						"oss_bucket_name":            "hansheng",
						"airflow_name":               name,
						"security_group_id":          CHECKSET,
						"vswitch_id":                 CHECKSET,
						"requirement_file":           "default/requirements.txt",
						"plugins_dir":                "default/plugins",
						"dags_dir":                   "default/dags",
						"startup_file":               "default/startup.sh",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"worker_serverless_replicas": "1",
					"description":                "11111",
					"app_spec":                   "MEDIUM",
					"airflow_name":               name + "_update",
					"requirement_file":           "default/requirements2.txt",
					"plugins_dir":                "default/plugins2",
					"dags_dir":                   "default/dags2",
					"startup_file":               "default/startup2.sh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"worker_serverless_replicas": "1",
						"description":                CHECKSET,
						"app_spec":                   "MEDIUM",
						"airflow_name":               name + "_update",
						"requirement_file":           "default/requirements2.txt",
						"plugins_dir":                "default/plugins2",
						"dags_dir":                   "default/dags2",
						"startup_file":               "default/startup2.sh",
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

var AlicloudDmsAirflowMap11207 = map[string]string{
	"airflow_id": CHECKSET,
	"region_id":  CHECKSET,
}

func AlicloudDmsAirflowBasicDependence11207(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-hangzhou-h"
}

resource "alicloud_security_group" "security_group" {
  description         = "zhenyuan_test_group"
  security_group_name = "zhenyuan_test_group"
  vpc_id              = data.alicloud_vpcs.default.ids.0
  security_group_type = "normal"
  inner_access_policy = "Accept"
}

resource "alicloud_dms_enterprise_workspace" "workspace" {
  description    = "zhenyuan-test"
  vpc_id         = data.alicloud_vpcs.default.ids.0
  workspace_name = "zhenyuan-test"
}


`, name)
}

// Test Dms Airflow. <<< Resource test cases, automatically generated.
