package alicloud

import (
	"testing"
)

func TestAccAlicloudKmsPlaintextDataSource(t *testing.T) {
	resourceId := "data.alicloud_kms_plaintext.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceKmsPlaintextDependence)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ciphertext_blob": alicloud_kms_ciphertext.default.ciphertext_blob,
		}),
	}

	var existKmsPlaintextMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"plaintext": "plaintext",
		}
	}

	var fakeKmsPlaintextMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ciphertext_blob": NOSET,
		}
	}

	var kmsPlainCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existKmsPlaintextMapFunc,
		fakeMapFunc:  fakeKmsPlaintextMapFunc,
	}

	kmsPlainCheckInfo.dataSourceTestCheck(t, 0, allConf)
}

func dataSourceKmsPlaintextDependence(name string) string {
	return `
	resource "alicloud_kms_key" "default" {
    	is_enabled = true
	}

	resource "alicloud_kms_ciphertext" "default" {
		key_id = alicloud_kms_key.default.id
		plaintext = "plaintext"
	}

	`
}
