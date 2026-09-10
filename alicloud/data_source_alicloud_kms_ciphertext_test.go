package alicloud

import (
	"testing"
)

func TestAccAlicloudKMSCiphertextDataSource(t *testing.T) {
	resourceId := "data.alicloud_kms_ciphertext.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceKmsCiphertextDependence)

	plaintextConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"key_id":    "${alicloud_kms_key.default.id}",
			"plaintext": "plaintext",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"key_id":    "${alicloud_kms_key.default.id}",
			"plaintext": "plaintext",
			"encryption_context": map[string]string{
				"key": "value",
			},
		}),
	}

	var existKmsCiphertextMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ciphertext_blob": CHECKSET,
		}
	}

	var fakeKmsCiphertextMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ciphertext_blob": NOSET,
		}
	}

	var kmsCipherCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existKmsCiphertextMapFunc,
		fakeMapFunc:  fakeKmsCiphertextMapFunc,
	}

	kmsCipherCheckInfo.dataSourceTestCheck(t, 0, plaintextConf, allConf)
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
