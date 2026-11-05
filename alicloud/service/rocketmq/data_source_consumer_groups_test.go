package rocketmq_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudRocketmqConsumerGroups_basic(t *testing.T) {
	resourceId := "data.alicloud_rocketmq_consumer_groups.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqds%d", acctest.DefaultRegionToTest, rand)
	dependence := testAccRocketmqConsumerGroupsDataSourceDependence(name)

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
data "alicloud_rocketmq_consumer_groups" "default" {
  instance_id = alicloud_rocketmq_instance.default.id
  ids         = [var.name, "${var.name}-second"]
}
`,
				Check: rocketmqDataSourceChecks(resourceId, map[string]string{
					"groups.#":                   "2",
					"groups.0.consumer_group_id": name,
					"groups.0.status":            acctest.CHECKSET,
					"groups.0.instance_id":       acctest.CHECKSET,
					"groups.0.region_id":         acctest.CHECKSET,
					"groups.0.message_model":     acctest.CHECKSET,
					"groups.0.create_time":       acctest.CHECKSET,
					"groups.0.max_receive_tps":   "1000",
					"groups.1.consumer_group_id": name + "-second",
					"ids.#":                      "2",
					"names.#":                    "2",
					"names.0":                    name,
					"id":                         acctest.CHECKSET,
				}),
			},
			{
				Config: dependence + `
data "alicloud_rocketmq_consumer_groups" "default" {
  instance_id = alicloud_rocketmq_instance.default.id
  name_regex  = "^${var.name}$"
}
`,
				Check: rocketmqDataSourceChecks(resourceId, map[string]string{
					"groups.#":                   "1",
					"groups.0.consumer_group_id": name,
				}),
			},
			{
				Config: dependence + `
data "alicloud_rocketmq_consumer_groups" "default" {
  instance_id = alicloud_rocketmq_instance.default.id
  filter      = "${var.name}-second"
}
`,
				Check: rocketmqDataSourceChecks(resourceId, map[string]string{
					"groups.#":                   "1",
					"groups.0.consumer_group_id": name + "-second",
				}),
			},
			{
				Config: dependence + `
data "alicloud_rocketmq_consumer_groups" "default" {
  instance_id = alicloud_rocketmq_instance.default.id
  ids         = ["not-exist-group"]
}
`,
				Check: rocketmqDataSourceChecks(resourceId, map[string]string{
					"groups.#": "0",
					"ids.#":    "0",
					"names.#":  "0",
				}),
			},
		},
	})
}

func testAccRocketmqConsumerGroupsDataSourceDependence(name string) string {
	return testAccRocketmqDataSourceInstance(name) + `
resource "alicloud_rocketmq_consumer_group" "first" {
  instance_id         = alicloud_rocketmq_instance.default.id
  consumer_group_id   = var.name
  delivery_order_type = "Concurrently"
  remark              = "example"
  max_receive_tps     = 1000
}

resource "alicloud_rocketmq_consumer_group" "second" {
  instance_id         = alicloud_rocketmq_instance.default.id
  consumer_group_id   = "${var.name}-second"
  delivery_order_type = "Orderly"
  remark              = "example"
}
`
}
