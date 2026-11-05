package kms_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
)

func TestAccAliCloudKmsCiphertextDataSource(t *testing.T) {
	resourceId := "data.alicloud_kms_ciphertext.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, "", dataSourceKmsCiphertextDependence)

	plaintextConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"key_id":    "${alicloud_kms_key.default.id}",
			"plaintext": "plaintext",
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"key_id":    "${alicloud_kms_key.default.id}",
			"plaintext": "plaintext",
			"encryption_context": map[string]string{
				"key": "value",
			},
		}),
	}

	var existKmsCiphertextMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ciphertext_blob": acctest.CHECKSET,
		}
	}

	var fakeKmsCiphertextMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ciphertext_blob": acctest.NOSET,
		}
	}

	var kmsCipherCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existKmsCiphertextMapFunc,
		FakeMapFunc:  fakeKmsCiphertextMapFunc,
	}

	kmsCipherCheckInfo.DataSourceTestCheck(t, 0, plaintextConf, allConf)
}

func dataSourceKmsCiphertextDependence(name string) string {
	return `
	resource "alicloud_kms_key" "default" {
		description = "tf-testacckmskeyforCiphertest"
    	is_enabled = true
		pending_window_in_days = 7
	}
	`
}
