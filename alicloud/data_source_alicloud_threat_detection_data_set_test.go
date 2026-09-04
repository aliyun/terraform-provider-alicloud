package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudThreatDetectionDataSetsDataSource_nameRegex(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionDataSetsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_threat_detection_data_set.default.data_set_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionDataSetsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_threat_detection_data_set.default.data_set_name}_fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudThreatDetectionDataSetsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_data_set.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudThreatDetectionDataSetsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_threat_detection_data_set.default.id}_fake"]`,
		}),
	}
	var existAlicloudThreatDetectionDataSetsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"names.#":                     "1",
			"data_sets.#":                 "1",
			"data_sets.0.id":              CHECKSET,
			"data_sets.0.data_set_id":     CHECKSET,
			"data_sets.0.data_set_name":   CHECKSET,
			"data_sets.0.data_set_status": CHECKSET,
		}
	}
	var fakeAlicloudThreatDetectionDataSetsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"data_sets.#": "0",
		}
	}
	var alicloudThreatDetectionDataSetsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_threat_detection_data_sets.default",
		existMapFunc: existAlicloudThreatDetectionDataSetsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudThreatDetectionDataSetsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudThreatDetectionDataSetsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf)
}

func testAccCheckAlicloudThreatDetectionDataSetsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccThreatDetectionDataSet-%d"
}

resource "alicloud_threat_detection_data_set" "default" {
  data_set_name           = var.name
  data_set_field_key_name = "ip"
  data_set_file_name      = "example.csv"
  data_set_description    = "description for test"
  data_set_type           = "userDefined"
  data_set_status         = 1
}

data "alicloud_threat_detection_data_sets" "default" {
  %s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
