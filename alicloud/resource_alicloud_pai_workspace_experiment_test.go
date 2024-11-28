package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test PaiWorkspace Experiment. >>> Resource test cases, automatically generated.
// Case 实验测试_副本1732006685057 9003
func TestAccAliCloudPaiWorkspaceExperiment_basic9003(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_experiment.default"
	ra := resourceAttrInit(resourceId, AlicloudPaiWorkspaceExperimentMap9003)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceExperiment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPaiWorkspaceExperimentBasicDependence9003)
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
					"accessibility":   "PRIVATE",
					"artifact_uri":    "oss://yyt-409262.oss-cn-hangzhou.aliyuncs.com/test/",
					"experiment_name": name,
					"workspace_id":    "${alicloud_pai_workspace_workspace.defaultDI9fsL.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accessibility":   "PRIVATE",
						"artifact_uri":    "oss://yyt-409262.oss-cn-hangzhou.aliyuncs.com/test/",
						"experiment_name": name,
						"workspace_id":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accessibility":   "PUBLIC",
					"experiment_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accessibility":   "PUBLIC",
						"experiment_name": name + "_update",
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

var AlicloudPaiWorkspaceExperimentMap9003 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudPaiWorkspaceExperimentBasicDependence9003(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_pai_workspace_workspace" "defaultDI9fsL" {
  description    = "902"
  display_name   = "test_pop_800"
  workspace_name = var.name
  env_types      = ["prod"]
}


`, name)
}

// Case 实验测试 6043
func TestAccAliCloudPaiWorkspaceExperiment_basic6043(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_experiment.default"
	ra := resourceAttrInit(resourceId, AlicloudPaiWorkspaceExperimentMap6043)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceExperiment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPaiWorkspaceExperimentBasicDependence6043)
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
					"accessibility":   "PRIVATE",
					"artifact_uri":    "oss://yyt-409262.oss-cn-hangzhou.aliyuncs.com/test/",
					"experiment_name": name,
					"workspace_id":    "${alicloud_pai_workspace_workspace.defaultDI9fsL.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accessibility":   "PRIVATE",
						"artifact_uri":    "oss://yyt-409262.oss-cn-hangzhou.aliyuncs.com/test/",
						"experiment_name": name,
						"workspace_id":    CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"accessibility":   "PUBLIC",
					"experiment_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accessibility":   "PUBLIC",
						"experiment_name": name + "_update",
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

var AlicloudPaiWorkspaceExperimentMap6043 = map[string]string{
	"create_time": CHECKSET,
}

func AlicloudPaiWorkspaceExperimentBasicDependence6043(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_pai_workspace_workspace" "defaultDI9fsL" {
  description    = "418"
  display_name   = "test_pop_923"
  workspace_name = var.name
  env_types      = ["prod"]
}


`, name)
}

// Test PaiWorkspace Experiment. <<< Resource test cases, automatically generated.
