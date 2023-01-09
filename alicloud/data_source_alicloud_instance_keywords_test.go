package alicloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"testing"
)

func TestAccAlicloudRdsInstanceKeywordsDatasource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_instance_keywords.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", testAccAlicloudRdsInstanceKeywordsDataSourceConfig)

	AccountKeywordsConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"key": "account",
		}),
	}

	DatabaseKeywordsConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"key": "database",
		}),
	}

	var existRdsInstanceKeywordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      CHECKSET,
			"ids.0":      CHECKSET,
			"keywords.#": CHECKSET,
			"keywords.0": CHECKSET,
		}
	}

	var fakeRdsInstanceKeywordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"keywords.#": "0",
		}
	}

	var RdsInstanceKeywordsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existRdsInstanceKeywordsMapFunc,
		fakeMapFunc:  fakeRdsInstanceKeywordsMapFunc,
	}

	RdsInstanceKeywordsCheckInfo.dataSourceTestCheck(t, rand, AccountKeywordsConfig, DatabaseKeywordsConfig)
}

func testAccAlicloudRdsInstanceKeywordsDataSourceConfig(name string) string {
	return ""
}
