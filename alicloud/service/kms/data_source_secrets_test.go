package kms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudKmsSecretsDataSource(t *testing.T) {
	resourceId := "data.alicloud_kms_secrets.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKmsSecret-%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceKmsSecretsDependence)

	NameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"fetch_tags": "true",
			"name_regex": "^${alicloud_kms_secret.default.secret_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"fetch_tags": "true",
			"name_regex": "^${alicloud_kms_secret.default.secret_name}-fake",
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"fetch_tags": "true",
			"ids":        []string{"${alicloud_kms_secret.default.secret_name}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"fetch_tags": "true",
			"ids":        []string{"${alicloud_kms_secret.default.secret_name}-fake"},
		}),
	}
	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${alicloud_kms_secret.default.secret_name}",
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "Secrettest",
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${alicloud_kms_secret.default.secret_name}",
			"tags": map[string]interface{}{
				"Created": "TF_fake",
				"For":     "Secrettest",
			},
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"fetch_tags": "true",
			"name_regex": "^${alicloud_kms_secret.default.secret_name}",
			"ids":        []string{"${alicloud_kms_secret.default.secret_name}"},
			"tags": map[string]interface{}{
				"Created": "TF",
				"For":     "Secrettest",
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
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
			"ids.0":                  acctest.CHECKSET,
			"names.#":                "1",
			"names.0":                acctest.CHECKSET,
			"secrets.#":              "1",
			"secrets.0.id":           acctest.CHECKSET,
			"secrets.0.secret_name":  acctest.CHECKSET,
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

	var kmsSecretInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existKmsSecretMapFunc,
		FakeMapFunc:  fakeKmsSecretMapFunc,
	}

	kmsSecretInfo.DataSourceTestCheck(t, 0, NameRegexConf, idsConf, tagsConf, allConf)
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
