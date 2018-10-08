package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudMNSQueueDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudMNSQueueDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_mns_queues.queues"),
					resource.TestCheckResourceAttr("data.alicloud_mns_queues.queues", "queues.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_mns_queues.queues", "queues.0.name", "tf-testAccMNSQueueConfig1"),
					resource.TestCheckResourceAttr("data.alicloud_mns_queues.queues", "queues.0.delay_seconds", "60478"),
					resource.TestCheckResourceAttr("data.alicloud_mns_queues.queues", "queues.0.maximum_message_size", "12357"),
					resource.TestCheckResourceAttr("data.alicloud_mns_queues.queues", "queues.0.message_retention_period", "256000"),
					resource.TestCheckResourceAttr("data.alicloud_mns_queues.queues", "queues.0.visibility_timeouts", "30"),
					resource.TestCheckResourceAttr("data.alicloud_mns_queues.queues", "queues.0.polling_wait_seconds", "3"),
				),
			},
		},
	})
}

const testAccCheckAlicloudMNSQueueDataSourceConfig = `
data "alicloud_mns_queues" "queues" {
  name_prefix = "${alicloud_mns_queue.queue.name}"
}

resource "alicloud_mns_queue" "queue"{
	name="tf-testAccMNSQueueConfig1"
	delay_seconds=60478
	maximum_message_size=12357
	message_retention_period=256000
	visibility_timeout=30
	polling_wait_seconds=3
}
`
