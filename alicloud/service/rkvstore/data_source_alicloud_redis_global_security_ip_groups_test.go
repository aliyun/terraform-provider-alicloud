// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package rkvstore_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudRedisGlobalSecurityIpGroupDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRedisGlobalSecurityIpGroupSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_redis_global_security_ip_group.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudRedisGlobalSecurityIpGroupSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_redis_global_security_ip_group.default.id}_fake"]`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRedisGlobalSecurityIpGroupSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_redis_global_security_ip_group.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudRedisGlobalSecurityIpGroupSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_redis_global_security_ip_group.default.id}_fake"]`,
		}),
	}

	RedisGlobalSecurityIpGroupCheckInfo.DataSourceTestCheck(t, rand, idsConf, allConf)
}

var existRedisGlobalSecurityIpGroupMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#":                          "1",
		"groups.0.global_ip_group_name":     acctest.CHECKSET,
		"groups.0.global_ip_list":           acctest.CHECKSET,
		"groups.0.global_security_group_id": acctest.CHECKSET,
		"groups.0.region_id":                acctest.CHECKSET,
	}
}

var fakeRedisGlobalSecurityIpGroupMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#": "0",
	}
}

var RedisGlobalSecurityIpGroupCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_redis_global_security_ip_groups.default",
	ExistMapFunc: existRedisGlobalSecurityIpGroupMapFunc,
	FakeMapFunc:  fakeRedisGlobalSecurityIpGroupMapFunc,
}

func testAccCheckAlicloudRedisGlobalSecurityIpGroupSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccRedisGlobalSecurityIpGroup%d"
}
variable "zone_id" {
  default = "cn-beijing-h"
}

variable "region_id" {
  default = "cn-beijing"
}



resource "alicloud_redis_global_security_ip_group" "default" {
  global_ip_group_name = "ggn_test_create"
  global_ip_list = "192.168.0.1,10.10.10.10,172.16.0.1"
}

data "alicloud_redis_global_security_ip_groups" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
