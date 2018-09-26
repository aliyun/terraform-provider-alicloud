package alicloud

import (
	"testing"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudMNSTopicSubscription_importBasic(t *testing.T) {
	resourceName := "alicloud_mns_topic_subscription.subscription"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMNSTopicDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccMNSTopicSubscriptionConfig,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}