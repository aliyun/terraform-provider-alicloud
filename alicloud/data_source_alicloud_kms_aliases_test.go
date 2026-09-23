package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudKMSAliasesDataSource(t *testing.T) {
	resourceId := "data.alicloud_kms_aliases.default"
	rand := acctest.RandIntRange(1000000, 9999999)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, fmt.Sprintf("alias/tf_testacc_%d", rand), dataSourceKmsAliasesDependence)

	NameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${alicloud_kms_alias.this.alias_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${alicloud_kms_alias.this.alias_name}-fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_kms_alias.this.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_kms_alias.this.id}-fake"},
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${alicloud_kms_alias.this.alias_name}",
			"ids":        []string{"${alicloud_kms_alias.this.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${alicloud_kms_alias.this.alias_name}-fake",
			"ids":        []string{"${alicloud_kms_alias.this.id}-fake"},
		}),
	}
	var existKmsCiphertextMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                "1",
			"ids.0":                CHECKSET,
			"names.#":              "1",
			"names.0":              CHECKSET,
			"aliases.#":            "1",
			"aliases.0.id":         CHECKSET,
			"aliases.0.alias_name": CHECKSET,
			"aliases.0.key_id":     CHECKSET,
		}
	}

	var fakeKmsCiphertextMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"aliases.#": "0",
		}
	}

	var kmsCipherCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existKmsCiphertextMapFunc,
		fakeMapFunc:  fakeKmsCiphertextMapFunc,
	}

	kmsCipherCheckInfo.dataSourceTestCheck(t, 0, NameRegexConf, idsConf, allConf)
}

func dataSourceKmsAliasesDependence(name string) string {
	return fmt.Sprintf(`
    resource "alicloud_kms_key" "this" {
		description = "tf-testacckmskeyforaliasdatasource"
		pending_window_in_days = 7
	}

	resource "alicloud_kms_alias" "this" {
  		alias_name = "%s"
  		key_id = "${alicloud_kms_key.this.id}"
	}
	`, name)
}
