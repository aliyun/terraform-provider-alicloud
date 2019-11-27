package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudKmsCiphertextDataSource(t *testing.T) {
	resourceId := "data.alicloud_kms_ciphertext.default"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKmsCiphertextConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						resourceId, "ciphertext_blob",
					),
				),
			},
			{
				ResourceName: resourceId,
			},
		},
	})
}

const testAccDataSourceKmsCiphertextConfigBasic = `
resource "alicloud_kms_key" "default" {
    is_enabled = true
}

data "alicloud_kms_ciphertext" "default" {
	key_id = "${alicloud_kms_key.default.id}"
	plaintext = "plaintext"
}
`
