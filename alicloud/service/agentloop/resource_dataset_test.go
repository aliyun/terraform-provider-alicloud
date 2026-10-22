package agentloop_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/agentloop"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Agentloop Dataset. >>> Resource test cases, automatically generated.
// Case Dataset全生命周期 11079
func TestAccAliCloudAgentloopDataset_basic11079(t *testing.T) {
	resourceId := "alicloud_agentloop_dataset.default"
	rac := resourceAttrCheckInit(resourceId, AlicloudAgentloopDatasetMap11079, (*agentloop.Service).DescribeDataset)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestaccdataset%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAgentloopDatasetBasicDependence11079)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"agent_space":  "${alicloud_agentloop_agent_space.default.agent_space}",
					"dataset_name": name,
					"description":  "terraform自动化测试描述",
					"schema": map[string]interface{}{
						"input":  `{\"chn\":true,\"type\":\"text\"}`,
						"output": `{\"chn\":false,\"type\":\"text\"}`,
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"agent_space":  name + "-as",
						"dataset_name": name,
						"description":  "terraform自动化测试描述",
						"schema.%":     "2",
						"schema.input": "{\"chn\":true,\"type\":\"text\"}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "terraform自动化测试描述更新",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "terraform自动化测试描述更新",
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

var AlicloudAgentloopDatasetMap11079 = map[string]string{
	"region_id":   CHECKSET,
	"create_time": CHECKSET,
	"update_time": CHECKSET,
}

func AlicloudAgentloopDatasetBasicDependence11079(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_agentloop_agent_space" "default" {
    agent_space = "${var.name}-as"
}

`, name)
}

// Test Agentloop Dataset. <<< Resource test cases, automatically generated.
