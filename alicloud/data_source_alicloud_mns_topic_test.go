package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudMnsTopicDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudMNSTopicDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_mns_topics.topics"),
					resource.TestCheckResourceAttr("data.alicloud_mns_topics.topics", "topics.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_mns_topics.topics", "topics.0.name", "tf-testAccMNSTopicConfig1"),
					resource.TestCheckResourceAttr("data.alicloud_mns_topics.topics", "topics.0.maximum_message_size", "12357"),
					resource.TestCheckResourceAttr("data.alicloud_mns_topics.topics", "topics.0.logging_enabled", "true"),
				),
			},
		},
	})
}

const testAccCheckAlicloudMNSTopicDataSourceConfig = `
data "alicloud_mns_topics" "topics" {
  name_prefix = "${alicloud_mns_topic.topic.name}"
}

resource "alicloud_mns_topic" "topic"{
	name="tf-testAccMNSTopicConfig1"
	maximum_message_size=12357
	logging_enabled=true
}
`
