package nas_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudNASProtocolsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(100000, 999999)
	AllConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNasDataSource(map[string]string{
			"type":    `"Performance"`,
			"zone_id": `"${data.alicloud_zones.default.zones.0.id}"`,
		}),
		FakeConfig: testAccCheckAlicloudNasDataSource(map[string]string{
			"type":    `"Performance"`,
			"zone_id": `"${data.alicloud_zones.default.zones.0.id}_fake"`,
		}),
	}
	TypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudNasDataSource(map[string]string{
			"type": `"Performance"`,
		}),
	}
	nasRecordsCheckInfo.DataSourceTestCheck(t, rand, AllConf, TypeConf)
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
		"protocols.#": acctest.CHECKSET,
	}
}

var fakeNasRecordsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"protocols.#": "0",
	}
}

var nasRecordsCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_nas_protocols.default",
	ExistMapFunc: existNasRecordsMapFunc,
	FakeMapFunc:  fakeNasRecordsMapFunc,
}
