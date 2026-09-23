package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudApigSourceDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApigSourceSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_source.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudApigSourceSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_apig_source.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApigSourceSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_apig_source.default.source_name}"`,
		}),
		fakeConfig: testAccCheckAliCloudApigSourceSourceConfig(rand, map[string]string{
			"name_regex": `"^tf-fake-name-not-exist$"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudApigSourceSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_apig_source.default.id}"]`,
			"name_regex":        `"${alicloud_apig_source.default.source_name}"`,
			"gateway_id":        `"${alicloud_apig_gateway.default.id}"`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
			"type":              `"K8S"`,
			"enable_details":    `true`,
		}),
		fakeConfig: testAccCheckAliCloudApigSourceSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_apig_source.default.id}_fake"]`,
			"name_regex":        `"^tf-fake-name-not-exist$"`,
			"gateway_id":        `"${alicloud_apig_gateway.default.id}"`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
			"type":              `"K8S"`,
			"enable_details":    `true`,
		}),
	}

	ApigSourceCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

var existApigSourceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                 "1",
		"sources.#":             "1",
		"sources.0.id":          CHECKSET,
		"sources.0.source_id":   CHECKSET,
		"sources.0.source_name": CHECKSET,
		"sources.0.create_time": CHECKSET,
		"sources.0.update_time": CHECKSET,
	}
}

var fakeApigSourceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"sources.#": "0",
		"ids.#":     "0",
	}
}

var ApigSourceCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_apig_sources.default",
	existMapFunc: existApigSourceMapFunc,
	fakeMapFunc:  fakeApigSourceMapFunc,
}

func testAccCheckAliCloudApigSourceSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tfaccapigsource%d"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_cs_managed_kubernetes" "default" {
    name                 = var.name
    cluster_spec         = "ack.pro.small"
    worker_vswitch_ids   = [alicloud_vswitch.default.id]
    pod_cidr             = "10.95.0.0/16"
    service_cidr         = "172.23.0.0/16"
    slb_internet_enabled = true
    new_nat_gateway      = true
    deletion_protection  = false
}

resource "alicloud_vpc" "default" {
    vpc_name   = var.name
    cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
    vswitch_name = var.name
    vpc_id       = alicloud_vpc.default.id
    zone_id      = "cn-hangzhou-i"
    cidr_block   = "192.168.8.0/24"
}

resource "alicloud_apig_gateway" "default" {
    gateway_name    = var.name
    spec            = "apigw.small.x1"
    gateway_edition = "Professional"
    gateway_type    = "API"
    payment_type    = "PayAsYouGo"
    vpc {
        vpc_id = alicloud_vpc.default.id
    }
    vswitch {
        vswitch_id = alicloud_vswitch.default.id
    }
    network_access_config {
        type = "Internet"
    }
    zone_config {
        select_option = "Auto"
    }
    log_config {
        sls {
            enable = false
        }
    }
}

resource "alicloud_apig_source" "default" {
    type              = "K8S"
    gateway_id        = alicloud_apig_gateway.default.id
    resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
    k8s_source_info {
        cluster_id = alicloud_cs_managed_kubernetes.default.id
    }
}

data "alicloud_apig_sources" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
