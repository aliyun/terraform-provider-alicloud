package kms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudKmsKeysDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_kms_keys.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccKmsKeysDataSource_%d", rand),
		dataSourceKmsKeysConfigDependence)

	descriptionRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"description_regex": "^${alicloud_kms_key.default.description}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "^${alicloud_kms_key.default.description}-fake",
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"description_regex": "^${alicloud_kms_key.default.description}",
			"status":            "Enabled",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "^${alicloud_kms_key.default.description}",
			"status":            "Disabled",
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_kms_key.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_kms_key.default.id}-fake"},
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"description_regex": "^${alicloud_kms_key.default.description}",
			"status":            "Enabled",
			"ids":               []string{"${alicloud_kms_key.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "^${alicloud_kms_key.default.description}-fake",
			"status":            "Enabled",
			"ids":               []string{"${alicloud_kms_key.default.id}"},
		}),
	}

	var existKmsKeysMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                "1",
			"ids.0":                acctest.CHECKSET,
			"keys.#":               "1",
			"keys.0.id":            acctest.CHECKSET,
			"keys.0.arn":           acctest.CHECKSET,
			"keys.0.status":        "Enabled",
			"keys.0.description":   fmt.Sprintf("tf_testAccKmsKeysDataSource_%d", rand),
			"keys.0.creation_date": acctest.CHECKSET,
			"keys.0.creator":       acctest.CHECKSET,
		}
	}

	var fakeKmsKeysMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":  "0",
			"keys.#": "0",
		}
	}

	var kmsKeysCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existKmsKeysMapFunc,
		FakeMapFunc:  fakeKmsKeysMapFunc,
	}

	kmsKeysCheckInfo.DataSourceTestCheck(t, rand, descriptionRegexConf, statusConf, idsConf, allConf)
}

func dataSourceKmsKeysConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_kms_key" "default" {
    description = "%s"
    pending_window_in_days = 7
}
`, name)
}
