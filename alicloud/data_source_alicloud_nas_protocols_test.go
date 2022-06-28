package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudNASProtocolsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100000, 999999)
	AllConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasDataSource(map[string]string{
			"type":    `"Performance"`,
			"zone_id": `"${data.alicloud_zones.default.zones.0.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudNasDataSource(map[string]string{
			"type":    `"Performance"`,
			"zone_id": `"${data.alicloud_zones.default.zones.0.id}_fake"`,
		}),
	}
	TypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasDataSource(map[string]string{
			"type": `"Performance"`,
		}),
	}
	nasRecordsCheckInfo.dataSourceTestCheck(t, rand, AllConf, TypeConf)
}

func testAccCheckAlicloudNasDataSource(attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
data "alicloud_zones" "default" {} 

data "alicloud_nas_protocols" "default" {
	%s
}`, strings.Join(pairs, "\n	"))
}

var existNasRecordsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"protocols.#": CHECKSET,
	}
}

var fakeNasRecordsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"protocols.#": "0",
	}
}

var nasRecordsCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_nas_protocols.default",
	existMapFunc: existNasRecordsMapFunc,
	fakeMapFunc:  fakeNasRecordsMapFunc,
}
