package alicloud

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudMnsTopicDataSource(t *testing.T) {
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudMNSTopicDataSourceConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_mns_topics.topics"),
					resource.TestCheckResourceAttr("data.alicloud_mns_topics.topics", "topics.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_mns_topics.topics", "topics.0.name", fmt.Sprintf("tf-testAccMNSTopicConfig-%d", rand)),
					resource.TestCheckResourceAttr("data.alicloud_mns_topics.topics", "topics.0.maximum_message_size", "12357"),
					resource.TestCheckResourceAttr("data.alicloud_mns_topics.topics", "topics.0.logging_enabled", "true"),
				),
			},
		},
	})
}

func TestAccAlicloudMnsTopicDataSourceEmpty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudMNSTopicDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_mns_topics.topics"),
					resource.TestCheckResourceAttr("data.alicloud_mns_topics.topics", "topics.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_mns_topics.topics", "topics.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_mns_topics.topics", "topics.0.maximum_message_size"),
					resource.TestCheckNoResourceAttr("data.alicloud_mns_topics.topics", "topics.0.logging_enabled"),
				),
			},
		},
	})
}

func testAccCheckAlicloudMNSTopicDataSourceConfig(rand int) string {
	return fmt.Sprintf(`
	data "alicloud_mns_topics" "topics" {
	  name_prefix = "${alicloud_mns_topic.topic.name}"
	}

	resource "alicloud_mns_topic" "topic"{
		name="tf-testAccMNSTopicConfig-%d"
		maximum_message_size=12357
		logging_enabled=true
	}
	`, rand)
}

const testAccCheckAlicloudMNSTopicDataSourceEmpty = `
data "alicloud_mns_topics" "topics" {
  name_prefix = "tf-testacc-fake-name"
}

`
