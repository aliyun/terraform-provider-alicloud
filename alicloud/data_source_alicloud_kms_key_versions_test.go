package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudKMSKeyVersionsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_kms_key_versions.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccKmsKeyVersionsDataSource_%d", rand),
		dataSourceKmsKeyVersionsConfigDependence)

	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"key_id": "${alicloud_kms_key.default.id}",
			"ids":    []string{"${alicloud_kms_key.default.primary_key_version}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"key_id": "${alicloud_kms_key.default.id}",
			"ids":    []string{"${alicloud_kms_key.default.primary_key_version}-fake"},
		}),
	}
	var existKmsKeyVersionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"ids.0":                     CHECKSET,
			"versions.#":                "1",
			"versions.0.create_time":    CHECKSET,
			"versions.0.key_version_id": CHECKSET,
		}
	}

	var fakeKmsKeyVersionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"versions.#": "0",
		}
	}

	var kmsKeyVersionsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existKmsKeyVersionsMapFunc,
		fakeMapFunc:  fakeKmsKeyVersionsMapFunc,
	}

	kmsKeyVersionsCheckInfo.dataSourceTestCheck(t, rand, basicConf)
}

func dataSourceKmsKeyVersionsConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_kms_key" "default" {
  description = "%s"
  pending_window_in_days = 7
}`, name)
}
