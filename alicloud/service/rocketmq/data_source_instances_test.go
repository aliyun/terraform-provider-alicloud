package rocketmq_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudRocketmqInstances_basic(t *testing.T) {
	resourceId := "data.alicloud_rocketmq_instances.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srocketmqds%d", acctest.DefaultRegionToTest, rand)
	dependence := testAccRocketmqDataSourceInstance(name)

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
data "alicloud_rocketmq_instances" "default" {
  ids = [alicloud_rocketmq_instance.default.id]
}
`,
				Check: rocketmqDataSourceChecks(resourceId, map[string]string{
					"instances.#":                   "1",
					"instances.0.instance_id":       acctest.CHECKSET,
					"instances.0.instance_name":     name,
					"instances.0.status":            acctest.CHECKSET,
					"instances.0.region_id":         acctest.CHECKSET,
					"instances.0.series_code":       "ultimate",
					"instances.0.sub_series_code":   "cluster_ha",
					"instances.0.payment_type":      "PayAsYouGo",
					"instances.0.service_code":      "rmq",
					"instances.0.resource_group_id": acctest.CHECKSET,
					"instances.0.topic_count":       acctest.CHECKSET,
					"instances.0.group_count":       acctest.CHECKSET,
					"instances.0.create_time":       acctest.CHECKSET,
					"instances.0.tags.#":            acctest.CHECKSET,
					"ids.#":                         "1",
					"names.#":                       "1",
					"names.0":                       name,
					"id":                            acctest.CHECKSET,
				}),
			},
			{
				Config: dependence + `
data "alicloud_rocketmq_instances" "default" {
  name_regex = "^${var.name}$"
  tags = {
    Created = "TF"
  }
}
`,
				Check: rocketmqDataSourceChecks(resourceId, map[string]string{
					"instances.#":               "1",
					"instances.0.instance_name": name,
					"instances.0.tags.#":        acctest.CHECKSET,
				}),
			},
			{
				Config: dependence + `
data "alicloud_rocketmq_instances" "default" {
  name_regex = "^${var.name}$"
  tags = {
    Created = "NOT_MATCH"
  }
}
`,
				Check: rocketmqDataSourceChecks(resourceId, map[string]string{
					"instances.#": "0",
					"ids.#":       "0",
					"names.#":     "0",
				}),
			},
			{
				Config: dependence + `
data "alicloud_rocketmq_instances" "default" {
  name_regex        = "^${var.name}$"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}
`,
				Check: rocketmqDataSourceChecks(resourceId, map[string]string{
					"instances.#":                   "1",
					"instances.0.instance_name":     name,
					"instances.0.resource_group_id": acctest.CHECKSET,
				}),
			},
			{
				Config: dependence + `
data "alicloud_rocketmq_instances" "default" {
  name_regex   = "^${var.name}$"
  series_codes = ["ultimate"]
}
`,
				Check: rocketmqDataSourceChecks(resourceId, map[string]string{
					"instances.#":               "1",
					"instances.0.instance_name": name,
					"instances.0.series_code":   "ultimate",
				}),
			},
			{
				Config: dependence + `
data "alicloud_rocketmq_instances" "default" {
  name_regex   = "^${var.name}$"
  series_codes = ["standard"]
}
`,
				Check: rocketmqDataSourceChecks(resourceId, map[string]string{
					"instances.#": "0",
				}),
			},
			{
				Config: dependence + `
data "alicloud_rocketmq_instances" "default" {
  filter = alicloud_rocketmq_instance.default.id
}
`,
				Check: rocketmqDataSourceChecks(resourceId, map[string]string{
					"instances.#":             "1",
					"instances.0.instance_id": acctest.CHECKSET,
				}),
			},
			{
				Config: dependence + `
data "alicloud_rocketmq_instances" "default" {
  ids = ["rmq-not-exist-instance"]
}
`,
				Check: rocketmqDataSourceChecks(resourceId, map[string]string{
					"instances.#": "0",
					"ids.#":       "0",
					"names.#":     "0",
				}),
			},
		},
	})
}

// rocketmqDataSourceChecks returns a check set scoped to a single step. The
// accumulating ResourceAttrMapUpdateSet helper would otherwise re-assert keys
// from earlier steps after the result set has shrunk.
func rocketmqDataSourceChecks(resourceId string, checks map[string]string) resource.TestCheckFunc {
	return acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()(checks)
}

// testAccRocketmqDataSourceInstance creates the RocketMQ 5.0 instance that the
// data source tests read. The attributes are the ones the data sources assert.
func testAccRocketmqDataSourceInstance(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  description = "example"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "default" {
  description  = "example"
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_rocketmq_instance" "default" {
  product_info {
    msg_process_spec       = "rmq.u2.10xlarge"
    send_receive_ratio     = "0.3"
    message_retention_time = "70"
  }
  service_code      = "rmq"
  payment_type      = "PayAsYouGo"
  instance_name     = var.name
  sub_series_code   = "cluster_ha"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  remark            = "example"
  tags = {
    Created = "TF"
    For     = "example"
  }
  series_code = "ultimate"
  network_info {
    vpc_info {
      vpc_id = alicloud_vpc.default.id
      vswitches {
        vswitch_id = alicloud_vswitch.default.id
      }
    }
    internet_info {
      internet_spec      = "enable"
      flow_out_type      = "payByBandwidth"
      flow_out_bandwidth = "30"
      ip_whitelist       = []
    }
  }
}
`, name)
}
