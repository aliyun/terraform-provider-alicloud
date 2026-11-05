// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package sls_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudSlsEtlDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(1000000, 9999999)

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSlsEtlSourceConfig(rand, map[string]string{
			"ids":      `["${alicloud_sls_etl.default.id}"]`,
			"project":  `"${alicloud_log_project.defaulthhAPo6.id}"`,
			"logstore": `"${alicloud_log_store.defaultzWKLkp.name}"`,
		}),
		FakeConfig: testAccCheckAlicloudSlsEtlSourceConfig(rand, map[string]string{
			"ids":      `["${alicloud_sls_etl.default.id}_fake"]`,
			"project":  `"${alicloud_log_project.defaulthhAPo6.id}"`,
			"logstore": `"${alicloud_log_store.defaultzWKLkp.name}"`,
		}),
	}

	SlsEtlCheckInfo.DataSourceTestCheck(t, rand, allConf)
}

var existSlsEtlMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"etls.#":                    "1",
		"etls.0.status":             acctest.CHECKSET,
		"etls.0.description":        acctest.CHECKSET,
		"etls.0.configuration.#":    acctest.CHECKSET,
		"etls.0.create_time":        acctest.CHECKSET,
		"etls.0.job_name":           acctest.CHECKSET,
		"etls.0.schedule_id":        acctest.CHECKSET,
		"etls.0.last_modified_time": acctest.CHECKSET,
	}
}

var fakeSlsEtlMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"etls.#": "0",
	}
}

var SlsEtlCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_sls_etls.default",
	ExistMapFunc: existSlsEtlMapFunc,
	FakeMapFunc:  fakeSlsEtlMapFunc,
}

func testAccCheckAlicloudSlsEtlSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testaccsls%d"
}
resource "alicloud_log_project" "defaulthhAPo6" {
  description = "terraform-etl-test-312"
  name        = var.name
}

resource "alicloud_log_store" "defaultzWKLkp" {
  hot_ttl          = "8"
  retention_period = "30"
  shard_count      = "2"
  project          = alicloud_log_project.defaulthhAPo6.id
  name             = var.name
}


resource "alicloud_sls_etl" "default" {
  job_name = "etl-1740472705-185721"
  display_name = "etl-1740472705-185721"
  project = "${alicloud_log_project.defaulthhAPo6.id}"
  configuration {
    from_time = "1706771697"
    to_time = "1738394097"
    script = "* | extend a=1"
    lang = "SPL"
    role_arn = "test-role-arn"
      sink {
      role_arn = "test-role-arn"
      name = "11111"
      endpoint = "cn-hangzhou-intranet.log.aliyuncs.com"
      project = "gy-hangzhou-huolang-1"
      logstore = "gy-rm2"
      datasets = [
                   "__UNNAMED__"
                 ]
    }
    
    logstore = "${alicloud_log_store.defaultzWKLkp.name}"
  }
  
}

data "alicloud_sls_etls" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
