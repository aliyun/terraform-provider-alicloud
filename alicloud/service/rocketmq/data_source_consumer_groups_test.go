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
	name := fmt.Sprintf("tf-testAcc%srocketmqds%d", acctest.DefaultRegionToTest, rand)
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
  depends_on  = [alicloud_rocketmq_consumer_group.first, alicloud_rocketmq_consumer_group.second]
  ids         = [var.name, "${var.name}-second"]
}
`,
				Check: resource.ComposeTestCheckFunc(
					rocketmqDataSourceChecks(resourceId, map[string]string{
						"groups.#":                   "2",
						"groups.0.consumer_group_id": name,
						"groups.0.status":            acctest.CHECKSET,
						"groups.0.remark":            "example",
						"groups.0.instance_id":       acctest.CHECKSET,
						"groups.0.region_id":         acctest.CHECKSET,
						"groups.0.create_time":       acctest.CHECKSET,
						"groups.0.update_time":       acctest.CHECKSET,
						"groups.0.max_receive_tps":   "1000",
						// ListConsumerGroups returns messageModel only for lite
						// consumer groups; a standard group maps it to an empty
						// value while the attribute stays exported. The mapping
						// itself is asserted by TestUnitRocketmqConsumerGroupMapping.
						"groups.0.message_model":     "",
						"groups.1.consumer_group_id": name + "-second",
						"ids.#":                      "2",
						"names.#":                    "2",
						"names.0":                    name,
						"id":                         acctest.CHECKSET,
					}),
					resource.TestCheckResourceAttrPair(resourceId, "groups.0.instance_id", "alicloud_rocketmq_instance.default", "id"),
				),
			},
			{
				Config: dependence + `
data "alicloud_rocketmq_consumer_groups" "default" {
  instance_id = alicloud_rocketmq_instance.default.id
  depends_on  = [alicloud_rocketmq_consumer_group.first, alicloud_rocketmq_consumer_group.second]
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
  depends_on  = [alicloud_rocketmq_consumer_group.first, alicloud_rocketmq_consumer_group.second]
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
  depends_on  = [alicloud_rocketmq_consumer_group.first, alicloud_rocketmq_consumer_group.second]
  ids         = ["not-exist-group"]
}
`,
				Check: rocketmqDataSourceChecks(resourceId, map[string]string{
					"groups.#": "0",
					"ids.#":    "0",
					"names.#":  "0",
				}),
			},
			{
				Config: dependence + `
data "alicloud_rocketmq_consumer_groups" "default" {
  instance_id = alicloud_rocketmq_instance.default.id
  depends_on  = [alicloud_rocketmq_consumer_group.first, alicloud_rocketmq_consumer_group.second]
  ids         = [var.name, "${var.name}-second"]
}

data "alicloud_rocketmq_instances" "default" {
  ids = [alicloud_rocketmq_instance.default.id]
}
`,
				Check: resource.ComposeTestCheckFunc(
					rocketmqDataSourceChecks(resourceId, map[string]string{
						"groups.#": "2",
					}),
					resource.TestCheckResourceAttrWith(
						"data.alicloud_rocketmq_instances.default",
						"instances.0.group_count",
						rocketmqCheckAtLeastInt(2),
					),
				),
			},
			{
				// instance_id is Required: a config without it must fail before
				// any ListConsumerGroups request, with a diagnostic naming the
				// missing argument.
				Config: dependence + `
data "alicloud_rocketmq_consumer_groups" "default" {
  name_regex = "^${var.name}$"
}
`,
				ExpectError: rocketmqMissingInstanceIdError(),
			},
			{
				// More consumer groups than one page (pageSize=100) so the walk
				// must read the second page; the filtered result proves the
				// later page is included and the primary key is not duplicated.
				Config: dependence + `
resource "alicloud_rocketmq_consumer_group" "page" {
  count               = 101
  instance_id         = alicloud_rocketmq_instance.default.id
  consumer_group_id   = format("%s-p%03d", var.name, count.index + 1)
  delivery_order_type = "Concurrently"
  consume_retry_policy {
    retry_policy = "DefaultRetryPolicy"
  }
}

data "alicloud_rocketmq_consumer_groups" "default" {
  instance_id = alicloud_rocketmq_instance.default.id
  depends_on  = [alicloud_rocketmq_consumer_group.page]
  name_regex  = "^${var.name}-p"
}
`,
				Check: rocketmqDataSourceChecks(resourceId, map[string]string{
					"groups.#":                     "101",
					"groups.0.consumer_group_id":   fmt.Sprintf("%s-p001", name),
					"groups.100.consumer_group_id": fmt.Sprintf("%s-p101", name),
					"names.#":                      "101",
					"ids.#":                        "101",
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
  consume_retry_policy {
    retry_policy = "DefaultRetryPolicy"
  }
}

resource "alicloud_rocketmq_consumer_group" "second" {
  instance_id         = alicloud_rocketmq_instance.default.id
  consumer_group_id   = "${var.name}-second"
  delivery_order_type = "Concurrently"
  remark              = "example"
  consume_retry_policy {
    retry_policy = "DefaultRetryPolicy"
  }
}
`
}
