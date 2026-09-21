package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test AgentLoop AgentSpace. >>> Resource test cases, hand-written.
func TestAccAliCloudAgentloopAgentSpace_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_agentloop_agent_space.default"
	ra := resourceAttrInit(resourceId, AlicloudAgentloopAgentSpaceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AgentLoopServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAgentLoopAgentSpace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-agentloop-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAgentloopAgentSpaceBasicDependence)
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
					"agent_space":          name,
					"description":          "description for test",
					"cms_workspace":        fmt.Sprintf("cms-%d", rand),
					"mse_workspace":        fmt.Sprintf("mse-%d", rand),
					"delete_cms_workspace": false,
					"delete_mse_namespace": false,
					"delete_sls_project":   false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"agent_space":   name,
						"description":   "description for test",
						"cms_workspace": fmt.Sprintf("cms-%d", rand),
						"mse_workspace": fmt.Sprintf("mse-%d", rand),
						"create_time":   CHECKSET,
						"region_id":     CHECKSET,
						"sls_project":   CHECKSET,
						"update_time":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"agent_space":          name,
					"description":          "description for test updated",
					"cms_workspace":        fmt.Sprintf("cms-%d", rand),
					"mse_workspace":        fmt.Sprintf("mse-%d", rand),
					"delete_cms_workspace": false,
					"delete_mse_namespace": false,
					"delete_sls_project":   false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "description for test updated",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_cms_workspace", "delete_mse_namespace", "delete_sls_project"},
			},
		},
	})
}

var AlicloudAgentloopAgentSpaceMap = map[string]string{
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
	"sls_project": CHECKSET,
	"update_time": CHECKSET,
}

func AlicloudAgentloopAgentSpaceBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}
