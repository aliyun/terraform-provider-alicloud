package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"strings"
	"testing"
)

func TestAccAlicloudRdsInstanceKeywordsDataSource(t *testing.T) {
	rand := acctest.RandInt()

	testConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRdsInstanceKeywordsDataSourceName(rand, map[string]string{
			"key": `"account"`,
		}),
		fakeConfig: testAccCheckAlicloudRdsInstanceKeywordsDataSourceName(rand, map[string]string{
			"key": `"database"`,
		}),
	}
	var existAlicloudRdsKeywordsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"words.#": CHECKSET,
		}
	}

	var fakeAlicloudRdsKeywordsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"words.#": CHECKSET,
		}
	}

	var alicloudRdsKeywordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_rds_instance_keywords.default",
		existMapFunc: existAlicloudRdsKeywordsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudRdsKeywordsDataSourceNameMapFunc,
	}
	alicloudRdsKeywordsCheckInfo.dataSourceTestCheck(t, rand, testConf)
}

func testAccCheckAlicloudRdsInstanceKeywordsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
data "alicloud_rds_instance_keywords" "default" {
  %s
}`, strings.Join(pairs, "\n"))
}
