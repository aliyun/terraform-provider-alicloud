package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSAEInstanceSpecificationsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSaeInstanceSpecificationsDataSourceName(rand, map[string]string{}),
		fakeConfig: testAccCheckAlicloudSaeInstanceSpecificationsDataSourceName(rand, map[string]string{
			"ids": `["fake"]`,
		}),
	}

	var existAlicloudSaeInstanceSpecificationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      CHECKSET,
			"specifications.#":           CHECKSET,
			"specifications.0.enable":    CHECKSET,
			"specifications.0.cpu":       CHECKSET,
			"specifications.0.memory":    CHECKSET,
			"specifications.0.spec_info": CHECKSET,
			"specifications.0.version":   CHECKSET,
		}
	}
	var fakeAlicloudSaeInstanceSpecificationsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudSaeInstanceSpecificationsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_sae_instance_specifications.default",
		existMapFunc: existAlicloudSaeInstanceSpecificationsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudSaeInstanceSpecificationsDataSourceNameMapFunc,
	}
	alicloudSaeInstanceSpecificationsCheckInfo.dataSourceTestCheck(t, rand, idsConf)
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
