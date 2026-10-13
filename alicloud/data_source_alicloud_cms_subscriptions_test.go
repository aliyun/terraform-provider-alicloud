package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCmsSubscriptions_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_subscription.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsSubscriptionMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsSubscription")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-cms-subscription%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCmsSubscriptionsDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_cms_subscriptions.default", "subscriptions.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_cms_subscriptions.default", "subscriptions.0.subscription_id"),
					resource.TestCheckResourceAttr("data.alicloud_cms_subscriptions.default", "subscriptions.0.subscription_name", name),
				),
			},
		},
	})
}

func testAccCmsSubscriptionsDataSourceConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_cms_subscription" "default" {
  subscription_name  = var.name
  description        = var.name
  notify_strategy_id = ""
  workspace          = ""
  enable             = true
  filter_setting {
    conditions {
      field = "product"
      value = "ECS"
      op    = "EQ"
    }
    expression = "product"
    relation   = "AND"
  }
  pushing_setting {
    alert_action_ids   = []
    restore_action_ids = []
    template_uuid      = ""
    response_plan_id   = ""
  }
}

data "alicloud_cms_subscriptions" "default" {
  subscription_name = alicloud_cms_subscription.default.subscription_name
  ids               = [alicloud_cms_subscription.default.id]
}
`, name)
}
