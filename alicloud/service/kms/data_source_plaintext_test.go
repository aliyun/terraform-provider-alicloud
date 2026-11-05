package kms_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
)

func TestAccAliCloudKmsPlaintextDataSource(t *testing.T) {
	resourceId := "data.alicloud_kms_plaintext.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, "", dataSourceKmsPlaintextDependence)

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ciphertext_blob": "${alicloud_kms_ciphertext.default.ciphertext_blob}",
		}),
	}

	var existKmsPlaintextMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"plaintext": "plaintext",
		}
	}

	var fakeKmsPlaintextMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ciphertext_blob": acctest.NOSET,
		}
	}

	var kmsPlainCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existKmsPlaintextMapFunc,
		FakeMapFunc:  fakeKmsPlaintextMapFunc,
	}

	kmsPlainCheckInfo.DataSourceTestCheck(t, 0, allConf)
}

func dataSourceKmsPlaintextDependence(name string) string {
	return `
	resource "alicloud_kms_key" "default" {
        description = "tf-testacckmskeyforplantexttest"
    	is_enabled = true
		pending_window_in_days = 7
	}

	resource "alicloud_kms_ciphertext" "default" {
		key_id = "${alicloud_kms_key.default.id}"
		plaintext = "plaintext"
	}

	`
}
