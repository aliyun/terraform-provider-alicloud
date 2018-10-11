package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDatahubSubscription_importBasic(t *testing.T) {
	resourceName := "alicloud_datahub_subscription.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatahubSubscriptionDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDatahubSubscription,
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
