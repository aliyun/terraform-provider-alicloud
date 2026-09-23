package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
		Providers:     testAccProviders,
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
				ImportStateVerifyIgnore: []string{"agent_key", "data_id_list"},
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
