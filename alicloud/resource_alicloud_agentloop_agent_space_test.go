package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test AgentLoop AgentSpace. >>> Resource test cases.

// AgentLoop currently only provides openAPI endpoints in cn-hangzhou and
// ap-southeast-1; pin the test region to cn-hangzhou unless ALICLOUD_REGION
// is set explicitly. The pin must happen before testAccPreCheck, which
// defaults an empty ALICLOUD_REGION to cn-beijing.
func testAccPreCheckAgentLoopRegion(t *testing.T) {
	if v := os.Getenv("ALICLOUD_REGION"); v == "" {
		os.Setenv("ALICLOUD_REGION", "cn-hangzhou")
	}
	testAccPreCheck(t)
}

var AliCloudAgentLoopAgentSpaceMap = map[string]string{
	"id":          CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
	"update_time": CHECKSET,
}

func TestAccAliCloudAgentLoopAgentSpace_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_agentloop_agent_space.default"
	ra := resourceAttrInit(resourceId, AliCloudAgentLoopAgentSpaceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AgentLoopServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAgentLoopAgentSpace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccagentloop%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudAgentLoopAgentSpaceBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckAgentLoopRegion(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"agent_space":   name,
					"description":   name,
					"cms_workspace": "",
					"mse_workspace": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"agent_space": name,
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"agent_space": name,
					"description": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "-update",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func AliCloudAgentLoopAgentSpaceBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}

// Test AgentLoop AgentSpace. <<< Resource test cases.
