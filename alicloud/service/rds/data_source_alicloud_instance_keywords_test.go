package rds_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudRdsInstanceKeywordsDatasource(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_instance_keywords.default"
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, "", testAccAlicloudRdsInstanceKeywordsDataSourceConfig)

	AccountKeywordsConfig := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"key": "account",
		}),
	}

	DatabaseKeywordsConfig := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"key": "database",
		}),
	}

	var existRdsInstanceKeywordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      acctest.CHECKSET,
			"ids.0":      acctest.CHECKSET,
			"keywords.#": acctest.CHECKSET,
			"keywords.0": acctest.CHECKSET,
		}
	}

	var fakeRdsInstanceKeywordsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"keywords.#": "0",
		}
	}

	var RdsInstanceKeywordsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existRdsInstanceKeywordsMapFunc,
		FakeMapFunc:  fakeRdsInstanceKeywordsMapFunc,
	}

	RdsInstanceKeywordsCheckInfo.DataSourceTestCheck(t, rand, AccountKeywordsConfig, DatabaseKeywordsConfig)
}

func testAccAlicloudRdsInstanceKeywordsDataSourceConfig(name string) string {
	return ""
}
