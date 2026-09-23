// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudFcv3TriggerDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	FunctionNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFcv3TriggerSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_fcv3_trigger.default.id}"]`,
			"function_name": `"${alicloud_fcv3_trigger.default.function_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudFcv3TriggerSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_fcv3_trigger.default.id}_fake"]`,
			"function_name": `"${alicloud_fcv3_trigger.default.function_name}"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFcv3TriggerSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_fcv3_trigger.default.id}"]`,
			"function_name": `"${alicloud_fcv3_trigger.default.function_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudFcv3TriggerSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_fcv3_trigger.default.id}_fake"]`,
			"function_name": `"${alicloud_fcv3_trigger.default.function_name}"`,
		}),
	}

	Fcv3TriggerCheckInfo.dataSourceTestCheck(t, rand, FunctionNameConf, allConf)
}

var existFcv3TriggerMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"triggers.#":                    "1",
		"triggers.0.trigger_type":       CHECKSET,
		"triggers.0.description":        CHECKSET,
		"triggers.0.create_time":        CHECKSET,
		"triggers.0.trigger_config":     CHECKSET,
		"triggers.0.trigger_name":       CHECKSET,
		"triggers.0.trigger_id":         CHECKSET,
		"triggers.0.last_modified_time": CHECKSET,
		"triggers.0.qualifier":          CHECKSET,
		"triggers.0.http_trigger.#":     CHECKSET,
	}
}

var fakeFcv3TriggerMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"triggers.#": "0",
	}
}

var Fcv3TriggerCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_fcv3_triggers.default",
	existMapFunc: existFcv3TriggerMapFunc,
	fakeMapFunc:  fakeFcv3TriggerMapFunc,
}

func testAccCheckAlicloudFcv3TriggerSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccFcv3Trigger%d"
}
variable "function_name" {
  default = "TestTriggerResourceAPI"
}

variable "trigger_name" {
  default = "TestTrigger_HTTP"
}

resource "alicloud_fcv3_function" "function" {
  memory_size = "512"
  cpu         = 0.5
  handler     = "index.Handler"
  code {
    zip_file = "UEsDBBQACAAIAAAAAAAAAAAAAAAAAAAAAAAIAAAAaW5kZXgucHmEkEFKxEAQRfd9ig9ZTCJOooIwDMwNXLqXnnQlaalUhU5lRj2KZ/FOXkESGR114bJ/P/7jV4b1xRq1hijtFpM1682cuNgPmgysbRulPT0fRxXnMtwrSPyeCdYRokSLnuMLJTTkbUqEvDMbxm1VdcRD6Tk+T1LW2ldB66knsYdA5iNX17ebm6tN2VnPhcswMPmREPuBacb+CiapLarAj9gT6/H97dVlCNScY3mtYvRkxdZlwDKDEnanPWVLdrdkeXEGlFEazVdfPVHaVeHc3N15CUwppwOJXeK7HshAB8NuOU7J6sP4SRXuH/EvbUfMiqMmDqv5M5FNSfAj/wgAAP//UEsHCPl//NYAAQAArwEAAFBLAQIUABQACAAIAAAAAAD5f/zWAAEAAK8BAAAIAAAAAAAAAAAAAAAAAAAAAABpbmRleC5weVBLBQYAAAAAAQABADYAAAA2AQAAAAA="
  }
  function_name = var.name
  runtime       = "python3.9"
  disk_size     = "512"
  log_config {
    log_begin_rule = "None"
  }
}


resource "alicloud_fcv3_trigger" "default" {
  function_name = "${alicloud_fcv3_function.function.function_name}"
  trigger_type = "http"
  trigger_name = "tf-testacceu-central-1fcv3trigger28547"
  description = "create"
  qualifier = "LATEST"
  trigger_config = "{\"methods\":[\"GET\",\"POST\",\"PUT\",\"DELETE\"],\"authType\":\"anonymous\",\"disableURLInternet\":false}"
}

data "alicloud_fcv3_triggers" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
