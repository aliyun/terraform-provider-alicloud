package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAliCloudKmsCiphertext_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
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

func TestAccAliCloudKmsCiphertext_validate(t *testing.T) {

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
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

func TestAccAliCloudKmsCiphertext_validate_withContext(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
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

func TestAccAliCloudKmsCiphertext_plaintextWo(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudKmsCiphertextConfig_plaintextWo(acctest.RandomWithPrefix("tf-testacc-plaintextwo"), 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_kms_ciphertext.default", "ciphertext_blob"),
					resource.TestCheckResourceAttr("alicloud_kms_ciphertext.default", "plaintext_wo_version", "1"),
					resource.TestCheckResourceAttr("alicloud_kms_ciphertext.default", "has_plaintext_wo", "true"),
				),
			},
			{
				Config: testAccAlicloudKmsCiphertextConfig_plaintextWo(acctest.RandomWithPrefix("tf-testacc-plaintextwo"), 2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_kms_ciphertext.default", "ciphertext_blob"),
					resource.TestCheckResourceAttr("alicloud_kms_ciphertext.default", "plaintext_wo_version", "2"),
					resource.TestCheckResourceAttr("alicloud_kms_ciphertext.default", "has_plaintext_wo", "true"),
				),
			},
			{
				Config: testAccAlicloudKmsCiphertextConfig_basic(acctest.RandomWithPrefix("tf-testacc-plaintextwo")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_kms_ciphertext.default", "ciphertext_blob"),
					resource.TestCheckResourceAttr("alicloud_kms_ciphertext.default", "has_plaintext_wo", "false"),
				),
			},
		},
	})
}

var testAccAlicloudKmsCiphertextConfig_plaintextWo = func(keyId string, version int) string {
	return fmt.Sprintf(`
resource "alicloud_kms_key" "default" {
  	description = "%s"
	is_enabled  = true
	pending_window_in_days = 7
}

resource "alicloud_kms_ciphertext" "default" {
	key_id = "${alicloud_kms_key.default.id}"
	plaintext_wo = "plaintext-wo"
	plaintext_wo_version = %d
}
`, keyId, version)
}
