package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudMnsQueue_importBasic(t *testing.T) {
	resourceName := "alicloud_mns_queue.queue"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMNSQueueDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccMNSQueueConfig(acctest.RandIntRange(10000, 999999)),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
