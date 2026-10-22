package agentloop_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/agentloop"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Agentloop ContextStoreApiKey. >>> Resource test cases, automatically generated.
// Case ContextStoreAPIKey全生命周期 11078
func TestAccAliCloudAgentloopContextStoreApiKey_basic11078(t *testing.T) {
	resourceId := "alicloud_agentloop_context_store_api_key.default"
	rac := resourceAttrCheckInit(resourceId, AlicloudAgentloopContextStoreApiKeyMap11078, (*agentloop.Service).DescribeContextStoreApiKey)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tftestaccapikey%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAgentloopContextStoreApiKeyBasicDependence11078)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:            resourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"agent_space":        "${alicloud_agentloop_agent_space.default.agent_space}",
					"context_store_name": "${alicloud_agentloop_context_store.default.context_store_name}",
					"name":               name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"agent_space":        name + "-as",
						"context_store_name": name + "cs",
						"name":               name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// The complete apiKey is only returned at creation time and the Get API
				// only returns a desensitized prefix, so it cannot be verified on import.
				ImportStateVerifyIgnore: []string{"api_key"},
			},
		},
	})
}

var AlicloudAgentloopContextStoreApiKeyMap11078 = map[string]string{
	"api_key":     CHECKSET,
	"region_id":   CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudAgentloopContextStoreApiKeyBasicDependence11078(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_agentloop_agent_space" "default" {
    agent_space = "${var.name}-as"
}

resource "alicloud_agentloop_context_store" "default" {
    agent_space        = alicloud_agentloop_agent_space.default.agent_space
    context_store_name = "${var.name}cs"
    context_type       = "experience"
    config {
        service_names = ["svc-tf-acc"]
    }
}

`, name)
}

// Test Agentloop ContextStoreApiKey. <<< Resource test cases, automatically generated.
