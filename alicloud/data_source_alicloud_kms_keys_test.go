package alicloud

import (
	//"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudKmsKeyDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudKmsKeyDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_kms_keys.keys"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.0.description", "testAccCheckAlicloudKmsKeyDataSourceBasic"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.0.status", "Enabled"),
				),
			},
		},
	})
}

func TestAccAlicloudKmsKeyDataSource_enpty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudKmsKeyDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_kms_keys.keys"),
					resource.TestCheckResourceAttr("data.alicloud_kms_keys.keys", "keys.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_kms_keys.keys", "keys.0.description"),
					resource.TestCheckNoResourceAttr("data.alicloud_kms_keys.keys", "keys.0.status"),
				),
			},
		},
	})
}

const testAccCheckAlicloudKmsKeyDataSourceBasic = `
resource "alicloud_kms_key" "key" {
    description = "testAccCheckAlicloudKmsKeyDataSourceBasic"
    deletion_window_in_days = 7
}

data "alicloud_kms_keys" "keys" {
	description_regex = "testAccCheck*"
	ids = ["${alicloud_kms_key.key.id}"]
	status = "Enabled"
}
`

const testAccCheckAlicloudKmsKeyDataSourceEmpty = `
data "alicloud_kms_keys" "keys" {
	description_regex = "^tf-testAcc-fake-name"
}
`
