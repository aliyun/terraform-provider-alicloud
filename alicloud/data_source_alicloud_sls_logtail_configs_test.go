// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSlsLogtailConfigDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	ProjectNameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlsLogtailConfigSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_sls_logtail_config.default.id}"]`,
			"project_name":  `"${var.project_name}"`,
			"logstore_name": `"test"`,
			"offset":        `"0"`,
			"size":          `"10"`,
		}),
		fakeConfig: testAccCheckAlicloudSlsLogtailConfigSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_sls_logtail_config.default.id}_fake"]`,
			"project_name":  `"${var.project_name}"`,
			"logstore_name": `"test"`,
			"offset":        `"0"`,
			"size":          `"10"`,
		}),
	}

	AllConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlsLogtailConfigSourceConfig(rand, map[string]string{
			"ids":                 `["${alicloud_sls_logtail_config.default.id}"]`,
			"project_name":        `"${var.project_name}"`,
			"logtail_config_name": `"tfaccsls12378"`,
			"logstore_name":       `"test"`,
			"offset":              `"0"`,
			"size":                `"10"`,
		}),
		fakeConfig: testAccCheckAlicloudSlsLogtailConfigSourceConfig(rand, map[string]string{
			"ids":                 `["${alicloud_sls_logtail_config.default.id}_fake"]`,
			"project_name":        `"${var.project_name}"`,
			"logtail_config_name": `"tfaccsls12378"`,
			"logstore_name":       `"test"`,
			"offset":              `"0"`,
			"size":                `"10"`,
		}),
	}

	SlsLogtailConfigCheckInfo.dataSourceTestCheck(t, rand, ProjectNameConf, AllConf)
}

var existSlsLogtailConfigMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"configs.#":                     "1",
		"configs.0.logtail_config_name": CHECKSET,
	}
}

var fakeSlsLogtailConfigMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"configs.#": "0",
	}
}

var SlsLogtailConfigCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_sls_logtail_configs.default",
	existMapFunc: existSlsLogtailConfigMapFunc,
	fakeMapFunc:  fakeSlsLogtailConfigMapFunc,
}

func testAccCheckAlicloudSlsLogtailConfigSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSlsLogtailConfig%d"
}
variable "project_name" {
  default = "project-for-logtail-terraform"
}

resource "alicloud_log_project" "defaultuA28zS" {
  description = "project for terrafrom test"
  name        = var.project_name
}

resource "alicloud_sls_logtail_config" "default" {
  output_type = "LogService"
  input_detail = "{\"adjustTimezone\":false,\"delayAlarmBytes\":0,\"delaySkipBytes\":0,\"discardNonUtf8\":false,\"discardUnmatch\":true,\"dockerFile\":false,\"enableRawLog\":false,\"enableTag\":false,\"fileEncoding\":\"utf8\",\"filePattern\":\"access*.log\",\"filterKey\":[\"key1\"],\"filterRegex\":[\"regex1\"],\"key\":[\"key1\",\"key2\"],\"localStorage\":true,\"logBeginRegex\":\".*\",\"logPath\":\"/var/log/httpd\",\"logTimezone\":\"\",\"logType\":\"common_reg_log\",\"maxDepth\":1000,\"maxSendRate\":-1,\"mergeType\":\"topic\",\"preserve\":true,\"preserveDepth\":0,\"priority\":0,\"regex\":\"(w+)(s+)\",\"sendRateExpire\":0,\"sensitive_keys\":[],\"tailExisted\":false,\"timeFormat\":\"%%Y/%%m/%%d %%H:%%M:%%S\",\"timeKey\":\"time\",\"topicFormat\":\"none\"}"
  logtail_config_name = "tfaccsls12378"
  input_type = "file"
  project_name = "${alicloud_log_project.defaultuA28zS.project_name}"
  output_detail {
    endpoint = "cn-hangzhou-intranet.log.aliyuncs.com"
    region = "cn-hangzhou"
    logstore_name = "test"
  }
}


data "alicloud_sls_logtail_configs" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
