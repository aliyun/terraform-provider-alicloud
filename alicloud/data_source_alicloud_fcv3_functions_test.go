// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudFcv3FunctionDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	PrefixIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFcv3FunctionSourceConfig(rand, map[string]string{
			"prefix": `"terraform-example-for"`,
		}),
		fakeConfig: testAccCheckAlicloudFcv3FunctionSourceConfig(rand, map[string]string{
			"prefix": `"terraform-example-for_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFcv3FunctionSourceConfig(rand, map[string]string{
			"ids": `["terraform-example-for-function-alias"]`,
		}),
		fakeConfig: testAccCheckAlicloudFcv3FunctionSourceConfig(rand, map[string]string{
			"ids": `["terraform-example-for-function-alias_fake"]`,
		}),
	}

	ResourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFcv3FunctionSourceConfig(rand, map[string]string{
			"ids":               `["terraform-example-for-function-alias"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
		}),
		fakeConfig: testAccCheckAlicloudFcv3FunctionSourceConfig(rand, map[string]string{
			"ids":               `["terraform-example-for-function-alias_fake"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFcv3FunctionSourceConfig(rand, map[string]string{
			"ids":               `["terraform-example-for-function-alias"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}"`,
		}),
		fakeConfig: testAccCheckAlicloudFcv3FunctionSourceConfig(rand, map[string]string{
			"ids":               `["terraform-example-for-function-alias_fake"]`,
			"resource_group_id": `"${data.alicloud_resource_manager_resource_groups.default.ids.0}_fake"`,
		}),
	}

	Fcv3FunctionCheckInfo.dataSourceTestCheck(t, rand, PrefixIdConf, idsConf, ResourceGroupIdConf, allConf)
}

var existFcv3FunctionMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"functions.#":                           "1",
		"functions.0.id":                        CHECKSET,
		"functions.0.resource_group_id":         CHECKSET,
		"functions.0.vpc_config.#":              CHECKSET,
		"functions.0.tracing_config.#":          CHECKSET,
		"functions.0.cpu":                       "0.5",
		"functions.0.idle_timeout":              CHECKSET,
		"functions.0.function_name":             CHECKSET,
		"functions.0.runtime":                   CHECKSET,
		"functions.0.custom_dns.#":              CHECKSET,
		"functions.0.disk_size":                 "512",
		"functions.0.instance_concurrency":      CHECKSET,
		"functions.0.layers.#":                  CHECKSET,
		"functions.0.nas_config.#":              CHECKSET,
		"functions.0.tags.%":                    CHECKSET,
		"functions.0.function_id":               CHECKSET,
		"functions.0.memory_size":               "512",
		"functions.0.function_arn":              CHECKSET,
		"functions.0.timeout":                   CHECKSET,
		"functions.0.create_time":               CHECKSET,
		"functions.0.handler":                   CHECKSET,
		"functions.0.custom_container_config.#": CHECKSET,
		"functions.0.internet_access":           "true",
		"functions.0.invocation_restriction.#":  CHECKSET,
		"functions.0.code_size":                 CHECKSET,
		"functions.0.custom_runtime_config.#":   CHECKSET,
		"functions.0.gpu_config.#":              CHECKSET,
		"functions.0.oss_mount_config.#":        CHECKSET,
		"functions.0.last_modified_time":        CHECKSET,
		"functions.0.log_config.#":              CHECKSET,
	}
}

var fakeFcv3FunctionMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"functions.#": "0",
	}
}

var Fcv3FunctionCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_fcv3_functions.default",
	existMapFunc: existFcv3FunctionMapFunc,
	fakeMapFunc:  fakeFcv3FunctionMapFunc,
}

func testAccCheckAlicloudFcv3FunctionSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccFcv3Function%d"
}
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_fcv3_function" "default" {
  function_name = "TestNativeRuntime_ResourceGroup"
  memory_size   = "512"
  runtime       = "python3.10"
  timeout       = "3"
  handler       = "index.handler"
  code {
    zip_file = "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA="
  }
  description = "Create"
  cpu               = 0.5
  disk_size         = "512"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  log_config {
    log_begin_rule = "None"
  }
}

data "alicloud_fcv3_functions" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
