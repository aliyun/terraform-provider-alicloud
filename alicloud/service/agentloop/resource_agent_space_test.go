package agentloop_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/agentloop"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Agentloop AgentSpace. >>> Resource test cases, automatically generated.
// Case AgentSpace全生命周期 11075
func TestAccAliCloudAgentloopAgentSpace_basic11075(t *testing.T) {
	resourceId := "alicloud_agentloop_agent_space.default"
	var v map[string]interface{}
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAgentloopAgentSpaceMap11075)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return agentloop.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAgentSpace")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestaccagentspace%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAgentloopAgentSpaceBasicDependence11075)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// The recommended usage: cms_workspace / mse_namespace_id are
				// NOT set, so the backend auto-creates the CMS workspace, the
				// MSE namespace and the SLS project together with the AgentSpace.
				Config: testAccConfig(map[string]interface{}{
					"agent_space": name,
					"description": "terraform自动化测试描述",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"agent_space": name,
						"description": "terraform自动化测试描述",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "terraform自动化测试描述更新了",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "terraform自动化测试描述更新了",
					}),
				),
			},
			{
				// AutoCreated -> UserSelected: rebinding the auto-created CMS
				// workspace to a user-selected one flips the bind type
				// asynchronously; the Update must wait until the settled bind
				// type reads UserSelected, otherwise a destroy right after the
				// update could misread the stale bind type.
				Config: testAccConfig(map[string]interface{}{
					"cms_workspace": "${alicloud_cms_workspace.default.workspace_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cms_workspace":           name + "-cms",
						"cms_workspace_bind_type": "UserSelected",
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

var AlicloudAgentloopAgentSpaceMap11075 = map[string]string{
	"region_id":        acctest.CHECKSET,
	"create_time":      acctest.CHECKSET,
	"cms_workspace":    acctest.CHECKSET,
	"mse_namespace_id": acctest.CHECKSET,
	"sls_project":      acctest.CHECKSET,
	// The bind type is assigned asynchronously by the backend provisioning
	// workflow, so right after creation it may still read "UserSelected"
	// even for an auto-created workspace. Only verify it is set here; the
	// exact enum values are asserted in case 11076 (UserSelected).
	"cms_workspace_bind_type": acctest.CHECKSET,
}

func AlicloudAgentloopAgentSpaceBasicDependence11075(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_log_project" "default" {
  project_name = "${var.name}-log"
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = "${var.name}-cms"
  sls_project    = alicloud_log_project.default.project_name
}


`, name)
}

// Case AgentSpace创建参数覆盖 11076
func TestAccAliCloudAgentloopAgentSpace_basic11076(t *testing.T) {
	resourceId := "alicloud_agentloop_agent_space.default"
	var v map[string]interface{}
	ra := acctest.ResourceAttrInit(resourceId, AlicloudAgentloopAgentSpaceMap11076)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return agentloop.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAgentSpace")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestaccagentspace%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudAgentloopAgentSpaceBasicDependence11076)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// cms_workspace is covered here with user-selected workspaces
				// (11075 covers the recommended auto-created mode): bind one
				// workspace and rebind to another to exercise the update path.
				Config: testAccConfig(map[string]interface{}{
					"agent_space":              name,
					"description":              "terraform自动化测试描述",
					"trajectory_store_enabled": "true",
					"cms_workspace":            "${alicloud_cms_workspace.default.workspace_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"agent_space":              name,
						"description":              "terraform自动化测试描述",
						"trajectory_store_enabled": "true",
						"cms_workspace":            name + "-cms-a",
						"cms_workspace_bind_type":  "UserSelected",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cms_workspace": "${alicloud_cms_workspace.update.workspace_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cms_workspace": name + "-cms-b",
					}),
				),
			},
			{
				// mse_namespace_id requires a pre-existing, unbound MSE namespace, which
				// cannot be created by this provider. Reusing the shared agent space's
				// namespace verifies that the parameter is passed through and validated
				// by the API (MSENamespaceAlreadyBound).
				Config: testAccConfig(map[string]interface{}{
					"mse_namespace_id": "${alicloud_agentloop_agent_space.shared.mse_namespace_id}",
				}),
				ExpectError: regexp.MustCompile(`MSENamespaceAlreadyBound`),
			},
			// No import step here: the ExpectError step above triggers a ForceNew
			// replacement of alicloud_agentloop_agent_space.default whose create is
			// rejected by the API, so the resource is absent from the state
			// afterwards. Import coverage is provided by case 11075.
		},
	})
}

var AlicloudAgentloopAgentSpaceMap11076 = map[string]string{
	"region_id":        acctest.CHECKSET,
	"create_time":      acctest.CHECKSET,
	"mse_namespace_id": acctest.CHECKSET,
}

func AlicloudAgentloopAgentSpaceBasicDependence11076(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_agentloop_agent_space" "shared" {
  agent_space = "${var.name}-shared"
}

resource "alicloud_log_project" "default" {
  project_name = "${var.name}-log-a"
}

resource "alicloud_log_project" "update" {
  project_name = "${var.name}-log-b"
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = "${var.name}-cms-a"
  sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_workspace" "update" {
  workspace_name = "${var.name}-cms-b"
  sls_project    = alicloud_log_project.update.project_name
}`, name)
}

// Test Agentloop AgentSpace. <<< Resource test cases, automatically generated.
