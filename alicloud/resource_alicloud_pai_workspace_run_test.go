package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test PaiWorkspace Run. >>> Resource test cases, automatically generated.
// Case Run测试_副本1732014483368 9013
func TestAccAliCloudPaiWorkspaceRun_basic9013(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_run.default"
	ra := resourceAttrInit(resourceId, AlicloudPaiWorkspaceRunMap9013)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceRun")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPaiWorkspaceRunBasicDependence9013)
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
					"source_type":   "TrainingService",
					"source_id":     "945",
					"run_name":      name,
					"experiment_id": "${alicloud_pai_workspace_experiment.defaultQRwWbv.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":   "TrainingService",
						"source_id":     CHECKSET,
						"run_name":      name,
						"experiment_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"run_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"run_name": name + "_update",
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

var AlicloudPaiWorkspaceRunMap9013 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudPaiWorkspaceRunBasicDependence9013(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_pai_workspace_workspace" "defaultCAFUa9" {
  description    = "440"
  display_name   = "test_pop_run_530"
  workspace_name = var.name
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_experiment" "defaultQRwWbv" {
  accessibility   = "PRIVATE"
  artifact_uri    = "oss://test.oss-cn-hangzhou.aliyuncs.com/test/"
  experiment_name = format("%%s1", var.name)
  workspace_id    = alicloud_pai_workspace_workspace.defaultCAFUa9.id
}


`, name)
}

// Case Run测试 8781
func TestAccAliCloudPaiWorkspaceRun_basic8781(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_run.default"
	ra := resourceAttrInit(resourceId, AlicloudPaiWorkspaceRunMap8781)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceRun")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPaiWorkspaceRunBasicDependence8781)
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
					"source_type":   "TrainingService",
					"source_id":     "28",
					"run_name":      name,
					"experiment_id": "${alicloud_pai_workspace_experiment.defaultQRwWbv.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"source_type":   "TrainingService",
						"source_id":     CHECKSET,
						"run_name":      name,
						"experiment_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"run_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"run_name": name + "_update",
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

var AlicloudPaiWorkspaceRunMap8781 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudPaiWorkspaceRunBasicDependence8781(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_pai_workspace_workspace" "defaultCAFUa9" {
  description    = "931"
  display_name   = "test_pop_run_674"
  workspace_name = var.name
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_experiment" "defaultQRwWbv" {
  artifact_uri    = "oss://test.oss-cn-hangzhou.aliyuncs.com/test/"
  accessibility   = "PRIVATE"
  experiment_name = format("%%s1", var.name)
  workspace_id    = alicloud_pai_workspace_workspace.defaultCAFUa9.id
}


`, name)
}

// Test PaiWorkspace Run. <<< Resource test cases, automatically generated.
