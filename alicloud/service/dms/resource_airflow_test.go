// Package dms. This file is generated automatically. Please do not modify it manually, thank you!
package dms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dms"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Dms Airflow. >>> Resource test cases, automatically generated.
// Case Airflow_tests 11207
func TestAccAliCloudDmsAirflow_basic11207(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	var v map[string]interface{}
	resourceId := "alicloud_dms_airflow.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDmsAirflowMap11207)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dms.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAirflow")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccdms%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDmsAirflowBasicDependence11207)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
					"requirement_file":           "/Workspace/requirements.txt",
					"plugins_dir":                "/Workspace/plugins",
					"dags_dir":                   "/Workspace/dags",
					"startup_file":               "/Workspace/startup.sh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"worker_serverless_replicas": "0",
						"description":                "zhenyuan-test",
						"zone_id":                    "cn-hangzhou-h",
						"workspace_id":               acctest.CHECKSET,
						"vpc_id":                     acctest.CHECKSET,
						"oss_path":                   "/",
						"app_spec":                   "SMALL",
						"oss_bucket_name":            "hansheng",
						"airflow_name":               name,
						"security_group_id":          acctest.CHECKSET,
						"vswitch_id":                 acctest.CHECKSET,
						"requirement_file":           "/Workspace/requirements.txt",
						"plugins_dir":                "/Workspace/plugins",
						"dags_dir":                   "/Workspace/dags",
						"startup_file":               "/Workspace/startup.sh",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"worker_serverless_replicas": "1",
					"description":                "11111",
					"app_spec":                   "MEDIUM",
					"airflow_name":               name + "_update",
					"requirement_file":           "/Workspace/requirements2.txt",
					"plugins_dir":                "/Workspace/plugins2",
					"dags_dir":                   "/Workspace/dags2",
					"startup_file":               "/Workspace/startup2.sh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"worker_serverless_replicas": "1",
						"description":                acctest.CHECKSET,
						"app_spec":                   "MEDIUM",
						"airflow_name":               name + "_update",
						"requirement_file":           "/Workspace/requirements2.txt",
						"plugins_dir":                "/Workspace/plugins2",
						"dags_dir":                   "/Workspace/dags2",
						"startup_file":               "/Workspace/startup2.sh",
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
	"airflow_id": acctest.CHECKSET,
	"region_id":  acctest.CHECKSET,
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
