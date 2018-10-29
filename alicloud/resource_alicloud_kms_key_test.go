package alicloud

import (
	"fmt"
	"testing"

	"github.com/denverdino/aliyungo/kms"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudKmsKey_basic(t *testing.T) {
	var keyBefore, keyAfter kms.KeyMetadata

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudKmsKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAlicloudKmsKeyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudKmsKeyExists("alicloud_kms_key.key", &keyBefore),
					resource.TestCheckResourceAttr("alicloud_kms_key.key", "is_enabled", "true"),
				),
			},
			{
				Config: testAlicloudKmsKeyUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudKmsKeyExists("alicloud_kms_key.key", &keyAfter),
					resource.TestCheckResourceAttr("alicloud_kms_key.key", "is_enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckAlicloudKmsKeyExists(name string, key *kms.KeyMetadata) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No KMS Key ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.DescribeKey(rs.Primary.ID)
		})
		if err != nil {
			return err
		}
		o, _ := raw.(*kms.DescribeKeyResponse)
		meta := o.KeyMetadata
		key = &meta

		return nil
	}
}

func testAccCheckAlicloudKmsKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_kms_key" {
			continue
		}

		raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.DescribeKey(rs.Primary.ID)
		})

		if err != nil && !IsExceptedError(err, ForbiddenKeyNotFound) {
			return err
		}
		out, _ := raw.(*kms.DescribeKeyResponse)
		if KeyState(out.KeyMetadata.KeyState) == PendingDeletion {
			return nil
		}

		return fmt.Errorf("KMS key still exists:\n%#v", out.KeyMetadata)
	}

	return nil
}

const testAlicloudKmsKeyBasic = `
resource "alicloud_kms_key" "key" {
    description = "Terraform acc test"
    deletion_window_in_days = 7
}`

const testAlicloudKmsKeyUpdate = `
resource "alicloud_kms_key" "key" {
    description = "Terraform acc test"
    deletion_window_in_days = 7
    is_enabled = false
}`
