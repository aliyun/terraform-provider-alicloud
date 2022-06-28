package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudKMSSecretsDataSource(t *testing.T) {
	resourceId := "data.alicloud_kms_secrets.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKmsSecret-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceKmsSecretsDependence)

	NameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"fetch_tags": "true",
			"name_regex": "^${alicloud_kms_secret.default.secret_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"fetch_tags": "true",
			"name_regex": "^${alicloud_kms_secret.default.secret_name}-fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"fetch_tags": "true",
			"ids":        []string{"${alicloud_kms_secret.default.secret_name}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"fetch_tags": "true",
			"ids":        []string{"${alicloud_kms_secret.default.secret_name}-fake"},
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"fetch_tags": "true",
			"name_regex": "^${alicloud_kms_secret.default.secret_name}",
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "Secrettest",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"fetch_tags": "true",
			"name_regex": "^${alicloud_kms_secret.default.secret_name}",
			"tags": map[string]interface{}{
				"Created": "TF_fake",
				"For":     "Secrettest",
			},
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"fetch_tags": "true",
			"name_regex": "^${alicloud_kms_secret.default.secret_name}",
			"ids":        []string{"${alicloud_kms_secret.default.secret_name}"},
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "Secrettest",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"fetch_tags": "true",
			"name_regex": "^${alicloud_kms_secret.default.secret_name}",
			"ids":        []string{"${alicloud_kms_secret.default.secret_name}-fake"},
			"tags": map[string]interface{}{
				"Created": "TF_fake",
				"For":     "Secrettest",
			},
		}),
	}
	var existKmsSecretMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                  "1",
			"ids.0":                  CHECKSET,
			"names.#":                "1",
			"names.0":                CHECKSET,
			"secrets.#":              "1",
			"secrets.0.id":           CHECKSET,
			"secrets.0.secret_name":  CHECKSET,
			"secrets.0.tags.%":       "2",
			"secrets.0.tags.Created": "TF",
			"secrets.0.tags.For":     "Secrettest",
		}
	}

	var fakeKmsSecretMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"secrets.#": "0",
		}
	}

	var kmsSecretInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existKmsSecretMapFunc,
		fakeMapFunc:  fakeKmsSecretMapFunc,
	}

	kmsSecretInfo.dataSourceTestCheck(t, 0, NameRegexConf, idsConf, tagsConf, allConf)
}

func dataSourceKmsSecretsDependence(name string) string {
	return fmt.Sprintf(`
		resource "alicloud_kms_secret" "default" {
		  secret_name = "%s"
		  secret_data = "user:root:passwd:1234"
		  version_id = "v001"
		  tags = {
			Created = "TF"
			For 	= "Secrettest"
		  }
		}
	`, name)
}
