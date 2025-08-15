// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSlsEtlDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSlsEtlSourceConfig(rand, map[string]string{
			"ids":      `["${alicloud_sls_etl.default.id}"]`,
			"project":  `"${alicloud_log_project.defaulthhAPo6.id}"`,
			"logstore": `"${alicloud_log_store.defaultzWKLkp.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSlsEtlSourceConfig(rand, map[string]string{
			"ids":      `["${alicloud_sls_etl.default.id}_fake"]`,
			"project":  `"${alicloud_log_project.defaulthhAPo6.id}"`,
			"logstore": `"${alicloud_log_store.defaultzWKLkp.name}"`,
		}),
	}

	SlsEtlCheckInfo.dataSourceTestCheck(t, rand, allConf)
}

var existSlsEtlMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"etls.#":                    "1",
		"etls.0.status":             CHECKSET,
		"etls.0.description":        CHECKSET,
		"etls.0.configuration.#":    CHECKSET,
		"etls.0.create_time":        CHECKSET,
		"etls.0.job_name":           CHECKSET,
		"etls.0.schedule_id":        CHECKSET,
		"etls.0.last_modified_time": CHECKSET,
	}
}

var fakeSlsEtlMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"etls.#": "0",
	}
}

var SlsEtlCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_sls_etls.default",
	existMapFunc: existSlsEtlMapFunc,
	fakeMapFunc:  fakeSlsEtlMapFunc,
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
