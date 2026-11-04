package kms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudKmsSecretVersionsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_kms_secret_versions.default"
	name := fmt.Sprintf("tf_testAccKmsSecretsVersionsDataSource_%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceKmsSecretsVersionsConfigDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"secret_name":    "${alicloud_kms_secret.default.secret_name}",
			"enable_details": "true",
			"ids":            []string{"${alicloud_kms_secret.default.version_id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"secret_name":    "${alicloud_kms_secret.default.secret_name}",
			"enable_details": "true",
			"ids":            []string{"${alicloud_kms_secret.default.version_id}-fake"},
		}),
	}

	VersionStageConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"secret_name":    "${alicloud_kms_secret.default.secret_name}",
			"enable_details": "true",
			"version_stage":  "ACSCurrent",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"secret_name":    "${alicloud_kms_secret.default.secret_name}",
			"enable_details": "true",
			"version_stage":  "ACSCurrent-fake",
		}),
	}

	var existKmsSecretVersionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"ids.0":                       acctest.CHECKSET,
			"versions.#":                  "1",
			"versions.0.secret_name":      name,
			"versions.0.version_id":       acctest.CHECKSET,
			"versions.0.version_stages.#": "1",
			"versions.0.secret_data":      acctest.CHECKSET,
			"versions.0.secret_data_type": acctest.CHECKSET,
		}
	}

	var fakeKmsSecretVersionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"versions.#": "0",
		}
	}

	var kmsSecretVersionsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existKmsSecretVersionsMapFunc,
		FakeMapFunc:  fakeKmsSecretVersionsMapFunc,
	}

	kmsSecretVersionsCheckInfo.DataSourceTestCheck(t, rand, idsConf, VersionStageConf)
}

func dataSourceKmsSecretsVersionsConfigDependence(name string) string {
	return fmt.Sprintf(`
		resource "alicloud_kms_secret" "default" {
		  secret_name = "%s"
		  secret_data = "user:root:passwd:1234"
		  version_id = "v001"
		}
	`, name)
}
