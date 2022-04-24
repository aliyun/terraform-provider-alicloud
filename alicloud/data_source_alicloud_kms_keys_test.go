package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudKMSKeysDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_kms_keys.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccKmsKeysDataSource_%d", rand),
		dataSourceKmsKeysConfigDependence)

	descriptionRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "^${alicloud_kms_key.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "^${alicloud_kms_key.default.description}-fake",
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "^${alicloud_kms_key.default.description}",
			"status":            "Enabled",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "^${alicloud_kms_key.default.description}",
			"status":            "Disabled",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_kms_key.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_kms_key.default.id}-fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "^${alicloud_kms_key.default.description}",
			"status":            "Enabled",
			"ids":               []string{"${alicloud_kms_key.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "^${alicloud_kms_key.default.description}-fake",
			"status":            "Enabled",
			"ids":               []string{"${alicloud_kms_key.default.id}"},
		}),
	}

	var existKmsKeysMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                "1",
			"ids.0":                CHECKSET,
			"keys.#":               "1",
			"keys.0.id":            CHECKSET,
			"keys.0.arn":           CHECKSET,
			"keys.0.status":        "Enabled",
			"keys.0.description":   fmt.Sprintf("tf_testAccKmsKeysDataSource_%d", rand),
			"keys.0.creation_date": CHECKSET,
			"keys.0.creator":       CHECKSET,
		}
	}

	var fakeKmsKeysMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":  "0",
			"keys.#": "0",
		}
	}

	var kmsKeysCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existKmsKeysMapFunc,
		fakeMapFunc:  fakeKmsKeysMapFunc,
	}

	kmsKeysCheckInfo.dataSourceTestCheck(t, rand, descriptionRegexConf, statusConf, idsConf, allConf)
}

func dataSourceKmsKeysConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_kms_key" "default" {
    description = "%s"
    pending_window_in_days = 7
}
`, name)
}
