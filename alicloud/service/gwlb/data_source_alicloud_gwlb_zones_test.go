package gwlb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudGwlbZoneDataSource(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.EfloSupportRegions)
	rand := tfacctest.RandIntRange(1000000, 9999999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudGwlbZoneSourceConfig(rand, map[string]string{
			"ids":             `["cn-wulanchabu-b"]`,
			"accept_language": `"zh-CN"`,
		}),
		FakeConfig: testAccCheckAlicloudGwlbZoneSourceConfig(rand, map[string]string{
			"ids":             `["cn-wulanchabu-a"]`,
			"accept_language": `"zh-CN"`,
		}),
	}

	GwlbZoneCheckInfo.DataSourceTestCheck(t, rand, idsConf)
}

var existGwlbZoneMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"zones.#":            acctest.CHECKSET,
		"zones.0.id":         acctest.CHECKSET,
		"zones.0.zone_id":    acctest.CHECKSET,
		"zones.0.local_name": acctest.CHECKSET,
	}
}

var fakeGwlbZoneMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"zones.#": "0",
	}
}

var GwlbZoneCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_gwlb_zones.default",
	ExistMapFunc: existGwlbZoneMapFunc,
	FakeMapFunc:  fakeGwlbZoneMapFunc,
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
