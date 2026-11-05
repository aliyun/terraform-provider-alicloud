package ebs_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudEbsRegionsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	acctest.CheckoutSupportedRegions(t, true, connectivity.EBSSupportRegions)
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEbsRegionsSourceConfig(rand, map[string]string{}),
		FakeConfig:  "",
	}

	var existEbsRegionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"regions.#":           acctest.CHECKSET,
			"regions.0.region_id": acctest.CHECKSET,
			"regions.0.zones.#":   acctest.CHECKSET,
		}
	}

	var fakeEbsRegionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"regions.#": "0",
		}
	}

	var EbsRegionsRecordsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ebs_regions.default",
		ExistMapFunc: existEbsRegionsMapFunc,
		FakeMapFunc:  fakeEbsRegionsMapFunc,
	}

	EbsRegionsRecordsCheckInfo.DataSourceTestCheck(t, rand, allConf)

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
