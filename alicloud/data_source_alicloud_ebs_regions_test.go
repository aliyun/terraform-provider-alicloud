package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEbsRegionsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEbsRegionsSourceConfig(rand, map[string]string{}),
		fakeConfig:  "",
	}

	var existEbsRegionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"regions.#":           CHECKSET,
			"regions.0.region_id": CHECKSET,
			"regions.0.zones.#":   CHECKSET,
		}
	}

	var fakeEbsRegionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"regions.#": "0",
		}
	}

	var EbsRegionsRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ebs_regions.default",
		existMapFunc: existEbsRegionsMapFunc,
		fakeMapFunc:  fakeEbsRegionsMapFunc,
	}

	EbsRegionsRecordsCheckInfo.dataSourceTestCheck(t, rand, allConf)

}

func testAccCheckAlicloudEbsRegionsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
data "alicloud_ebs_regions" "default"{
%s
}

`, strings.Join(pairs, "\n   "))
	return config
}
