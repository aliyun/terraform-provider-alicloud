package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudAgentloopAgentSpaces_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-agentloop-ds-%d", rand)
	testAccConfig := resourceTestAccConfigFunc("alicloud_agentloop_agent_space.default", name, AlicloudAgentloopAgentSpaceBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_agentloop_agent_space.default",
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"agent_space":          name,
					"description":          "description for datasource test",
					"cms_workspace":        fmt.Sprintf("cms-%d", rand),
					"mse_workspace":        fmt.Sprintf("mse-%d", rand),
					"delete_cms_workspace": false,
					"delete_mse_namespace": false,
					"delete_sls_project":   false,
				}) + fmt.Sprintf(`
data "alicloud_agentloop_agent_spaces" "default" {
    agent_space   = alicloud_agentloop_agent_space.default.agent_space
    biz_region_id = "cn-hangzhou"
}

output "spaces_id" {
    value = data.alicloud_agentloop_agent_spaces.default.spaces.0.id
}

output "spaces_agent_space" {
    value = data.alicloud_agentloop_agent_spaces.default.spaces.0.agent_space
}
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_agentloop_agent_spaces.default", "spaces.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_agentloop_agent_spaces.default", "spaces.0.agent_space", name),
				),
			},
		},
	})
}

// silence unused import if connectivity is not directly referenced
var _ connectivity.Region
