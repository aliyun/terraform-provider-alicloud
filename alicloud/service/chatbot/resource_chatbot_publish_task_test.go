package chatbot_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	chatbot "github.com/aliyun/terraform-provider-alicloud/alicloud/service/chatbot"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudChatbotPublishTask_basic2099(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.ChatbotSupportRegions)
	resourceId := "alicloud_chatbot_publish_task.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudChatbotPublishTaskMap2099)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return chatbot.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeChatbotPublishTask")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sChatbotPublishTask%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudChatbotPublishTaskBasicDependence2099)
	// CreatePublishTask only accepts DataIdList when BizType is "faq", and the
	// values are FAQ category ids. testAccChatbotFaqCategoryIds returns the
	// real ids published on the caller's default business space; two distinct
	// ids let step1 -> step2 swap the member order and exercise the
	// data_id_list TypeList reorder coverage required after relocation to the
	// service package. A fabricated placeholder would make the static gate
	// green without proving order-independence, so the ids come from the API.
	// GetPublishTaskState does not return DataIdList, so Read never sets it;
	// step3 PlanOnly then proves whether the reordered config converges on a
	// clean plan after step2's apply (the SDK writes the config value into
	// state for the non-computed attribute even though Read does not Set it).
	faqIds := testAccChatbotFaqCategoryIds(t)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"agent_key":    "${data.alicloud_chatbot_agents.default.agents.0.agent_key}",
					"biz_type":     "faq",
					"data_id_list": []interface{}{faqIds[0], faqIds[1]},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"biz_type": "faq",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"agent_key":    "${data.alicloud_chatbot_agents.default.agents.0.agent_key}",
					"biz_type":     "faq",
					"data_id_list": []interface{}{faqIds[1], faqIds[0]},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"biz_type": "faq",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"agent_key":    "${data.alicloud_chatbot_agents.default.agents.0.agent_key}",
					"biz_type":     "faq",
					"data_id_list": []interface{}{faqIds[1], faqIds[0]},
				}),
				PlanOnly: true,
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
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
