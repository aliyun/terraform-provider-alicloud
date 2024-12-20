package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudWafv3InstanceDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.WAFV3SupportRegions)
	rand := acctest.RandIntRange(1000000, 9999999)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudWafv3InstanceSourceConfig(rand, map[string]string{}),
		fakeConfig:  "",
	}

	Wafv3InstanceCheckInfo.dataSourceTestCheck(t, rand, allConf)
}

var existWafv3InstanceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                   "1",
		"instances.#":             "1",
		"instances.0.id":          CHECKSET,
		"instances.0.create_time": CHECKSET,
		"instances.0.instance_id": CHECKSET,
		"instances.0.status":      CHECKSET,
	}
}

var fakeWafv3InstanceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":       "0",
		"instances.#": "0",
	}
}

var Wafv3InstanceCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_wafv3_instances.default",
	existMapFunc: existWafv3InstanceMapFunc,
	fakeMapFunc:  fakeWafv3InstanceMapFunc,
}

func testAccCheckAlicloudWafv3InstanceSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccWafv3Instance%d"
}

resource "alicloud_wafv3_instance" "default" {}

data "alicloud_wafv3_instances" "default" {
    ids = [alicloud_wafv3_instance.default.id]
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
