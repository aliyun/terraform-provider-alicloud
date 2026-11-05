package wafopenapi_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudWafv3InstanceDataSource(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.WAFV3SupportRegions)
	rand := tfacctest.RandIntRange(1000000, 9999999)

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudWafv3InstanceSourceConfig(rand, map[string]string{}),
		FakeConfig:  "",
	}

	Wafv3InstanceCheckInfo.DataSourceTestCheck(t, rand, allConf)
}

var existWafv3InstanceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                   "1",
		"instances.#":             "1",
		"instances.0.id":          acctest.CHECKSET,
		"instances.0.create_time": acctest.CHECKSET,
		"instances.0.instance_id": acctest.CHECKSET,
		"instances.0.status":      acctest.CHECKSET,
	}
}

var fakeWafv3InstanceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":       "0",
		"instances.#": "0",
	}
}

var Wafv3InstanceCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_wafv3_instances.default",
	ExistMapFunc: existWafv3InstanceMapFunc,
	FakeMapFunc:  fakeWafv3InstanceMapFunc,
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
