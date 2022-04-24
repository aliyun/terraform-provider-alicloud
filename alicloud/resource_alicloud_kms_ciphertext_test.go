package alicloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudKMSCiphertext_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudKmsCiphertextConfig_basic(acctest.RandomWithPrefix("tf-testacc-basic")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"alicloud_kms_ciphertext.default", "ciphertext_blob"),
				),
			},
		},
	})
}

func TestAccAlicloudKMSCiphertext_validate(t *testing.T) {

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudKmsCiphertextConfig_validate(acctest.RandomWithPrefix("tf-testacc-validate")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_kms_ciphertext.default", "ciphertext_blob"),
					resource.TestCheckResourceAttrPair("alicloud_kms_ciphertext.default", "plaintext", "data.alicloud_kms_plaintext.default", "plaintext"),
				),
			},
		},
	})
}

func TestAccAlicloudKMSCiphertext_validate_withContext(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudKmsCiphertextConfig_validate_withContext(acctest.RandomWithPrefix("tf-testacc-validate-withcontext")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_kms_ciphertext.default", "ciphertext_blob"),
					resource.TestCheckResourceAttrPair("alicloud_kms_ciphertext.default", "plaintext", "data.alicloud_kms_plaintext.default", "plaintext"),
				),
			},
		},
	})
}

var testAccAlicloudKmsCiphertextConfig_basic = func(keyId string) string {
	return fmt.Sprintf(`
resource "alicloud_kms_key" "default" {
  	description = "%s"
	is_enabled  = true
	pending_window_in_days = 7
}

resource "alicloud_kms_ciphertext" "default" {
	key_id = "${alicloud_kms_key.default.id}"
	plaintext = "plaintext"
}
`, keyId)
}

var testAccAlicloudKmsCiphertextConfig_validate = func(keyId string) string {
	return fmt.Sprintf(`
	resource "alicloud_kms_key" "default" {
        description = "%s"
		is_enabled  = true
		pending_window_in_days = 7
	}
	
	resource "alicloud_kms_ciphertext" "default" {
		key_id = "${alicloud_kms_key.default.id}"
		plaintext = "plaintext"
	}
	
	data "alicloud_kms_plaintext" "default" {
	  ciphertext_blob = "${alicloud_kms_ciphertext.default.ciphertext_blob}"
	}
	`, keyId)
}

var testAccAlicloudKmsCiphertextConfig_validate_withContext = func(keyId string) string {
	return fmt.Sprintf(`
	resource "alicloud_kms_key" "default" {
        description = "%s"
		is_enabled  = true
		pending_window_in_days = 7
	}
	
	resource "alicloud_kms_ciphertext" "default" {
		key_id = "${alicloud_kms_key.default.id}"
		plaintext = "plaintext"
        encryption_context = {
    		name = "value"
  		}
	}
	
	data "alicloud_kms_plaintext" "default" {
	  ciphertext_blob = "${alicloud_kms_ciphertext.default.ciphertext_blob}"
	  encryption_context = {
		name = "value"
	  }
	}
	`, keyId)
}
