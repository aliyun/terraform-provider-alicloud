// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSlsIndexDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlsIndexSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_sls_index.default.id}"]`,
			"project_name":  `"${alicloud_log_project.defaultdCM1bA.project_name}"`,
			"logstore_name": `"${alicloud_log_store.default7MW26R.logstore_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlsIndexSourceConfig(rand, map[string]string{
			"ids":           `["${alicloud_sls_index.default.id}_fake"]`,
			"project_name":  `"${alicloud_log_project.defaultdCM1bA.project_name}"`,
			"logstore_name": `"${alicloud_log_store.default7MW26R.logstore_name}"`,
		}),
	}

	SlsIndexCheckInfo.dataSourceTestCheck(t, rand, allConf)
}

var existSlsIndexMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"indexs.#":                         "1",
		"indexs.0.line.#":                  "1",
		"indexs.0.max_text_len":            CHECKSET,
		"indexs.0.log_reduce_black_list.#": CHECKSET,
		"indexs.0.log_reduce_white_list.#": CHECKSET,
		"indexs.0.keys":                    CHECKSET,
		"indexs.0.ttl":                     CHECKSET,
	}
}

var fakeSlsIndexMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"indexs.#": "0",
	}
}

var SlsIndexCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_sls_indexs.default",
	existMapFunc: existSlsIndexMapFunc,
	fakeMapFunc:  fakeSlsIndexMapFunc,
}

func testAccCheckAlicloudSlsIndexSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSlsIndex%d"
}
variable "logstore_name" {
  default = "logstore-test-v6"
}

variable "project_name" {
  default = "project-for-index-terraform-test-v6"
}

resource "alicloud_log_project" "defaultdCM1bA" {
  description = "terrafrom test"
  name        = var.project_name
}

resource "alicloud_log_store" "default7MW26R" {
  hot_ttl          = "7"
  retention_period = "30"
  shard_count      = "2"
  project          = alicloud_log_project.defaultdCM1bA.project_name
  name             = var.logstore_name
}



resource "alicloud_sls_index" "default" {
  line {
    chn            = true
    case_sensitive = true
    token          = ["a"]
    exclude_keys   = ["t", "tt"]
  }
  logstore_name = alicloud_log_store.default7MW26R.logstore_name
  project_name  = alicloud_log_project.defaultdCM1bA.project_name
  max_text_len          = "2048"
  log_reduce_black_list = ["test"]
  log_reduce_white_list = ["name"]
  log_reduce            = true
  keys = jsonencode(
    {
      "example" : {
        "caseSensitive" : false,
        "token" : [
          "\n",
          "\t",
          ",",
          " ",
          ";",
          "\"",
          "'",
          "(",
          ")",
          "{",
          "}",
          "[",
          "]",
          "<",
          ">",
          "?",
          "/",
          "#",
          ":"
        ],
        "type" : "text",
        "doc_value" : false,
        "alias" : "",
        "chn" : false
      }
    }
  )
}

data "alicloud_sls_indexs" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
