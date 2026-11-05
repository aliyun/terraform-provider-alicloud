package rocketmq_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudRocketmqTopics_basic(t *testing.T) {
	resourceId := "data.alicloud_rocketmq_topics.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqds%d", acctest.DefaultRegionToTest, rand)
	dependence := testAccRocketmqTopicsDataSourceDependence(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheckWithRegions(t, true, connectivity.RocketMQSupportRegions)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: dependence + `
data "alicloud_rocketmq_topics" "default" {
  instance_id = alicloud_rocketmq_instance.default.id
  depends_on  = [alicloud_rocketmq_topic.first, alicloud_rocketmq_topic.second]
  ids         = [var.name, "${var.name}-second"]
}
`,
				Check: rocketmqDataSourceChecks(resourceId, map[string]string{
					"topics.#":              "2",
					"topics.0.topic_name":   name,
					"topics.0.message_type": "NORMAL",
					"topics.0.status":       acctest.CHECKSET,
					"topics.0.instance_id":  acctest.CHECKSET,
					"topics.0.region_id":    acctest.CHECKSET,
					"topics.0.create_time":  acctest.CHECKSET,
					"topics.0.max_send_tps": "1500",
					"topics.1.topic_name":   name + "-second",
					"topics.1.message_type": "FIFO",
					"ids.#":                 "2",
					"names.#":               "2",
					"names.0":               name,
					"id":                    acctest.CHECKSET,
				}),
			},
			{
				Config: dependence + `
data "alicloud_rocketmq_topics" "default" {
  instance_id = alicloud_rocketmq_instance.default.id
  depends_on  = [alicloud_rocketmq_topic.first, alicloud_rocketmq_topic.second]
  name_regex  = "^${var.name}$"
}
`,
				Check: rocketmqDataSourceChecks(resourceId, map[string]string{
					"topics.#":            "1",
					"topics.0.topic_name": name,
				}),
			},
			{
				Config: dependence + `
data "alicloud_rocketmq_topics" "default" {
  instance_id   = alicloud_rocketmq_instance.default.id
  depends_on    = [alicloud_rocketmq_topic.first, alicloud_rocketmq_topic.second]
  name_regex    = "^${var.name}"
  message_types = ["FIFO"]
}
`,
				Check: rocketmqDataSourceChecks(resourceId, map[string]string{
					"topics.#":              "1",
					"topics.0.topic_name":   name + "-second",
					"topics.0.message_type": "FIFO",
				}),
			},
			{
				Config: dependence + `
data "alicloud_rocketmq_topics" "default" {
  instance_id = alicloud_rocketmq_instance.default.id
  depends_on  = [alicloud_rocketmq_topic.first, alicloud_rocketmq_topic.second]
  filter      = "${var.name}-second"
}
`,
				Check: rocketmqDataSourceChecks(resourceId, map[string]string{
					"topics.#":            "1",
					"topics.0.topic_name": name + "-second",
				}),
			},
			{
				Config: dependence + `
data "alicloud_rocketmq_topics" "default" {
  instance_id = alicloud_rocketmq_instance.default.id
  depends_on  = [alicloud_rocketmq_topic.first, alicloud_rocketmq_topic.second]
  ids         = ["not-exist-topic"]
}
`,
				Check: rocketmqDataSourceChecks(resourceId, map[string]string{
					"topics.#": "0",
					"ids.#":    "0",
					"names.#":  "0",
				}),
			},
			{
				Config: dependence + `
data "alicloud_rocketmq_topics" "default" {
  instance_id = alicloud_rocketmq_instance.default.id
  depends_on  = [alicloud_rocketmq_topic.first, alicloud_rocketmq_topic.second]
  ids         = [var.name, "${var.name}-second"]
}

data "alicloud_rocketmq_instances" "default" {
  ids = [alicloud_rocketmq_instance.default.id]
}
`,
				Check: resource.ComposeTestCheckFunc(
					rocketmqDataSourceChecks(resourceId, map[string]string{
						"topics.#": "2",
					}),
					resource.TestCheckResourceAttrWith(
						"data.alicloud_rocketmq_instances.default",
						"instances.0.topic_count",
						rocketmqCheckAtLeastInt(2),
					),
				),
			},
		},
	})
}

func testAccRocketmqTopicsDataSourceDependence(name string) string {
	return testAccRocketmqDataSourceInstance(name) + `
resource "alicloud_rocketmq_topic" "first" {
  instance_id  = alicloud_rocketmq_instance.default.id
  topic_name   = var.name
  message_type = "NORMAL"
  remark       = "example"
  max_send_tps = 1500
}

resource "alicloud_rocketmq_topic" "second" {
  instance_id  = alicloud_rocketmq_instance.default.id
  topic_name   = "${var.name}-second"
  message_type = "FIFO"
  remark       = "example"
}
`
}
