package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudDatahubSubscription_importBasic(t *testing.T) {
	resourceName := "alicloud_datahub_subscription.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, true, connectivity.DatahubSupportedRegions) },
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
