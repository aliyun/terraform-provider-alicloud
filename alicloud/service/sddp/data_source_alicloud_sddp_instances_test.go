package sddp_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudSddpInstancesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSddpInstanceDataSourceConfig(rand, map[string]string{}),
		FakeConfig:  "",
	}

	var existSddpInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#":              "1",
			"instances.0.id":           acctest.CHECKSET,
			"instances.0.instance_id":  acctest.CHECKSET,
			"instances.0.status":       acctest.CHECKSET,
			"instances.0.payment_type": "POSTPAY",
		}
	}

	var fakeSddpInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var sddpInstancesRecordsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_sddp_instances.default",
		ExistMapFunc: existSddpInstanceMapFunc,
		FakeMapFunc:  fakeSddpInstanceMapFunc,
	}

	var preCheck = func() {
		acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-zhangjiakou"})
	}

	sddpInstancesRecordsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)

}

func testAccCheckAlicloudSddpInstanceDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
data "alicloud_sddp_instances" "default" {
  %s
}
`, strings.Join(pairs, "\n  "))
	return config
}
