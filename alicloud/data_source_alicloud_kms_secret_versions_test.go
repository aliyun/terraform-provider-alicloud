package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudKmsSecretVersionsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_kms_secret_versions.default"
	name := fmt.Sprintf("tf_testAccKmsSecretsVersionsDataSource_%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceKmsSecretsVersionsConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"secret_name":    "${alicloud_kms_secret.default.secret_name}",
			"enable_details": "true",
			"ids":            []string{"${alicloud_kms_secret.default.version_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"secret_name":    "${alicloud_kms_secret.default.secret_name}",
			"enable_details": "true",
			"ids":            []string{"${alicloud_kms_secret.default.version_id}-fake"},
		}),
	}

	var existKmsSecretVersionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"ids.0":                       CHECKSET,
			"versions.#":                  "1",
			"versions.0.secret_name":      name,
			"versions.0.version_id":       CHECKSET,
			"versions.0.version_stages.#": "1",
			"versions.0.secret_data":      CHECKSET,
			"versions.0.secret_data_type": CHECKSET,
		}
	}

	var fakeKmsSecretVersionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"versions.#": "0",
		}
	}

	var kmsSecretVersionsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existKmsSecretVersionsMapFunc,
		fakeMapFunc:  fakeKmsSecretVersionsMapFunc,
	}

	kmsSecretVersionsCheckInfo.dataSourceTestCheck(t, rand, idsConf)
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
