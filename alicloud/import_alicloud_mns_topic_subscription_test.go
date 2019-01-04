package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudMnsTopicSubscription_importBasic(t *testing.T) {
	resourceName := "alicloud_mns_topic_subscription.subscription"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMNSTopicDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccMNSTopicSubscriptionConfig(acctest.RandIntRange(10000, 999999)),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
