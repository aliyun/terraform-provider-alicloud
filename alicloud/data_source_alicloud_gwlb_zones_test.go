package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudGwlbZoneDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.EfloSupportRegions)
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudGwlbZoneSourceConfig(rand, map[string]string{
			"ids":             `["cn-wulanchabu-b"]`,
			"accept_language": `"zh-CN"`,
		}),
		fakeConfig: testAccCheckAlicloudGwlbZoneSourceConfig(rand, map[string]string{
			"ids":             `["cn-wulanchabu-a"]`,
			"accept_language": `"zh-CN"`,
		}),
	}

	GwlbZoneCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}

var existGwlbZoneMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"zones.#":            CHECKSET,
		"zones.0.id":         CHECKSET,
		"zones.0.zone_id":    CHECKSET,
		"zones.0.local_name": CHECKSET,
	}
}

var fakeGwlbZoneMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"zones.#": "0",
	}
}

var GwlbZoneCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_gwlb_zones.default",
	existMapFunc: existGwlbZoneMapFunc,
	fakeMapFunc:  fakeGwlbZoneMapFunc,
}

func testAccCheckAlicloudGwlbZoneSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccGwlbZone%d"
}

data "alicloud_gwlb_zones" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
