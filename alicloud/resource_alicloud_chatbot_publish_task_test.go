package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlicloudChatbotPublishTask_basic2099(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.ChatbotSupportRegions)
	resourceId := "alicloud_chatbot_publish_task.default"
	ra := resourceAttrInit(resourceId, AlicloudChatbotPublishTaskMap2099)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ChatbotService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeChatbotPublishTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sChatbotPublishTask%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudChatbotPublishTaskBasicDependence2099)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"agent_key": "${data.alicloud_chatbot_agents.default.agents.0.agent_key}",
					"biz_type":  "faq",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"biz_type": "faq",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				// Chatbot PublishTask is asynchronous on the server side: after Create
				// the task moves through FE_RUNNING -> FE_SUCCESS and modify_time
				// updates accordingly. ImportStateVerify reads twice with a small
				// gap, so the snapshot can land on different state/time values.
				// Skip those two attributes for the equality compare.
				ImportStateVerifyIgnore: []string{"agent_key", "data_id_list", "status", "modify_time"},
			},
		},
	})
}

var AlicloudChatbotPublishTaskMap2099 = map[string]string{}

func AlicloudChatbotPublishTaskBasicDependence2099(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_chatbot_agents" "default" {}

`, name)
}
