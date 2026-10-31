package sae_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudSAEInstanceSpecificationsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudSaeInstanceSpecificationsDataSourceName(rand, map[string]string{}),
		FakeConfig: testAccCheckAlicloudSaeInstanceSpecificationsDataSourceName(rand, map[string]string{
			"ids": `["fake"]`,
		}),
	}

	var existAlicloudSaeInstanceSpecificationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      acctest.CHECKSET,
			"specifications.#":           acctest.CHECKSET,
			"specifications.0.enable":    acctest.CHECKSET,
			"specifications.0.cpu":       acctest.CHECKSET,
			"specifications.0.memory":    acctest.CHECKSET,
			"specifications.0.spec_info": acctest.CHECKSET,
			"specifications.0.version":   acctest.CHECKSET,
		}
	}
	var fakeAlicloudSaeInstanceSpecificationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudSaeInstanceSpecificationsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_sae_instance_specifications.default",
		ExistMapFunc: existAlicloudSaeInstanceSpecificationsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudSaeInstanceSpecificationsDataSourceNameMapFunc,
	}
	alicloudSaeInstanceSpecificationsCheckInfo.DataSourceTestCheck(t, rand, idsConf)
}
func testAccCheckAlicloudSaeInstanceSpecificationsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
data "alicloud_sae_instance_specifications" "default" {	
	%s
}
`, strings.Join(pairs, " \n "))
	return config
}
