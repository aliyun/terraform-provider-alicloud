// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudElasticsearchLogstashDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudElasticsearchLogstashSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_elasticsearch_logstash.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudElasticsearchLogstashSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_elasticsearch_logstash.default.id}_fake"]`,
		}),
	}

	ResourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudElasticsearchLogstashSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_elasticsearch_logstash.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
		}),
		fakeConfig: testAccCheckAlicloudElasticsearchLogstashSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_elasticsearch_logstash.default.id}_fake"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
		}),
	}
	DescriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudElasticsearchLogstashSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_elasticsearch_logstash.default.id}"]`,
			"description": `"tf-acc-create-test-1"`,
		}),
		fakeConfig: testAccCheckAlicloudElasticsearchLogstashSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_elasticsearch_logstash.default.id}_fake"]`,
			"description": `"tf-acc-create-test-1_fake"`,
		}),
	}
	VersionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudElasticsearchLogstashSourceConfig(rand, map[string]string{
			"ids":     `["${alicloud_elasticsearch_logstash.default.id}"]`,
			"version": `"7.4_with_X-Pack"`,
		}),
		fakeConfig: testAccCheckAlicloudElasticsearchLogstashSourceConfig(rand, map[string]string{
			"ids":     `["${alicloud_elasticsearch_logstash.default.id}_fake"]`,
			"version": `"7.4_with_X-Pack_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudElasticsearchLogstashSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_elasticsearch_logstash.default.id}"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
			"description":       `"tf-acc-create-test-1"`,
			"version":           `"7.4_with_X-Pack"`,
		}),
		fakeConfig: testAccCheckAlicloudElasticsearchLogstashSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_elasticsearch_logstash.default.id}_fake"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
			"description":       `"tf-acc-create-test-1_fake"`,
			"version":           `"7.4_with_X-Pack_fake"`,
		}),
	}

	ElasticsearchLogstashCheckInfo.dataSourceTestCheck(t, rand, idsConf, ResourceGroupIdConf, DescriptionConf, VersionConf, allConf)
}

var existElasticsearchLogstashMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"logstashes.#":                   "1",
		"logstashes.0.status":            CHECKSET,
		"logstashes.0.description":       CHECKSET,
		"logstashes.0.resource_group_id": CHECKSET,
		"logstashes.0.instance_id":       CHECKSET,
		"logstashes.0.create_time":       CHECKSET,
		"logstashes.0.network_config.#":  CHECKSET,
		"logstashes.0.node_amount":       CHECKSET,
		"logstashes.0.updated_at":        CHECKSET,
		"logstashes.0.version":           CHECKSET,
		"logstashes.0.node_spec.#":       CHECKSET,
		"logstashes.0.payment_type":      CHECKSET,
		"logstashes.0.tags.%":            CHECKSET,
	}
}

var fakeElasticsearchLogstashMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"logstashes.#": "0",
	}
}

var ElasticsearchLogstashCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_elasticsearch_logstashes.default",
	existMapFunc: existElasticsearchLogstashMapFunc,
	fakeMapFunc:  fakeElasticsearchLogstashMapFunc,
}

func testAccCheckAlicloudElasticsearchLogstashSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccElasticsearchLogstash%d"
}

data "alicloud_elasticsearch_zones" "default" {}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.1.0.0/16"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_elasticsearch_zones.default.zones.0.id
}

resource "alicloud_elasticsearch_logstash" "default" {
  description       = "tf-acc-create-test-1"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  version           = "7.4_with_X-Pack"
  node_spec {
    disk_type = "cloud_efficiency"
    spec      = "elasticsearch.sn1ne.large"
    disk      = "20"
  }
  network_config {
    type       = "vpc"
    vpc_id     = alicloud_vpc.default.id
    vswitch_id = alicloud_vswitch.default.id
    vs_area    = data.alicloud_elasticsearch_zones.default.zones.0.id
  }
  payment_type = "Subscription"
  node_amount  = "1"
}

data "alicloud_elasticsearch_logstashes" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
