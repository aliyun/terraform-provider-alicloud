package kms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudKmsKeyVersionsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_kms_key_versions.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccKmsKeyVersionsDataSource_%d", rand),
		dataSourceKmsKeyVersionsConfigDependence)

	basicConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"key_id": "${alicloud_kms_key.default.id}",
			"ids":    []string{"${alicloud_kms_key.default.primary_key_version}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"key_id": "${alicloud_kms_key.default.id}",
			"ids":    []string{"${alicloud_kms_key.default.primary_key_version}-fake"},
		}),
	}
	var existKmsKeyVersionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"ids.0":                     acctest.CHECKSET,
			"versions.#":                "1",
			"versions.0.create_time":    acctest.CHECKSET,
			"versions.0.key_version_id": acctest.CHECKSET,
		}
	}

	var fakeKmsKeyVersionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"versions.#": "0",
		}
	}

	var kmsKeyVersionsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existKmsKeyVersionsMapFunc,
		FakeMapFunc:  fakeKmsKeyVersionsMapFunc,
	}

	kmsKeyVersionsCheckInfo.DataSourceTestCheck(t, rand, basicConf)
}

func dataSourceKmsKeyVersionsConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_kms_key" "default" {
  description = "%s"
  pending_window_in_days = 7
}`, name)
}
