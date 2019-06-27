package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

func init() {
	resource.AddTestSweepers("alicloud_kms_keys", &resource.Sweeper{
		Name: "alicloud_kms_keys",
		F:    testSweepKmsKey,
	})
}

func testSweepKmsKey(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
	}

	req := kms.CreateListKeysRequest()
	raw, err := client.WithKmsClient(func(kmsclient *kms.Client) (interface{}, error) {
		return kmsclient.ListKeys(req)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "kms_keys", req.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	keys := raw.(*kms.ListKeysResponse)
	swept := false

	for _, v := range keys.Keys.Key {
		kmsService := &KmsService{client: client}
		key, err := kmsService.DescribeKmsKey(v.KeyId)
		if err != nil {
			if NotFoundError(err) {
				if strings.Contains(err.Error(), "Provider ERROR") {
					continue
				}
				return nil
			}

			return WrapError(err)
		}
		for _, description := range prefixes {
			if strings.HasPrefix(strings.ToLower(key.KeyMetadata.Description), strings.ToLower(description)) {
				req := kms.CreateScheduleKeyDeletionRequest()
				req.KeyId = v.KeyId
				req.PendingWindowInDays = requests.NewInteger(7)
				raw, err = client.WithKmsClient(func(kmsclient *kms.Client) (interface{}, error) {
					return kmsclient.ScheduleKeyDeletion(req)
				})
				swept = true
				if err != nil {
					return WrapErrorf(err, DataDefaultErrorMsg, v.KeyId, req.GetActionName(), AlibabaCloudSdkGoERROR)
				}
				break
			}
		}
	}
	if swept {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudKmsKey_basic(t *testing.T) {
	var keyBefore kms.KeyMetadata
	resourceId := "alicloud_kms_key.key"
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAlicloudKmsKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAlicloudKmsKeyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudKmsKeyExists("alicloud_kms_key.key", &keyBefore),
					resource.TestCheckResourceAttr("alicloud_kms_key.key", "is_enabled", "true"),
					resource.TestCheckResourceAttr("alicloud_kms_key.key", "description", "tf_testaccCheckAlicloudKmsKeyDataSourceBasic"),
					resource.TestCheckResourceAttr("alicloud_kms_key.key", "key_usage", "ENCRYPT/DECRYPT"),
					resource.TestCheckResourceAttr("alicloud_kms_key.key", "deletion_window_in_days", "7"),
					resource.TestCheckResourceAttrSet("alicloud_kms_key.key", "arn"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_window_in_days"},
			},
			{
				Config: testAlicloudKmsKeyUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudKmsKeyExists("alicloud_kms_key.key", &keyBefore),
					resource.TestCheckResourceAttr("alicloud_kms_key.key", "is_enabled", "false"),
					resource.TestCheckResourceAttr("alicloud_kms_key.key", "description", "tf_testaccCheckAlicloudKmsKeyDataSourceBasic"),
					resource.TestCheckResourceAttr("alicloud_kms_key.key", "key_usage", "ENCRYPT/DECRYPT"),
					resource.TestCheckResourceAttr("alicloud_kms_key.key", "deletion_window_in_days", "7"),
					resource.TestCheckResourceAttrSet("alicloud_kms_key.key", "arn"),
				),
			},
		},
	})
}

func testAccCheckAlicloudKmsKeyExists(name string, keymeta *kms.KeyMetadata) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", name))
		}

		if rs.Primary.ID == "" {
			return WrapError(Error("No KMS Key ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		kmsService := &KmsService{client: client}
		key, err := kmsService.DescribeKmsKey(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}

			return WrapError(err)
		}
		meta := key.KeyMetadata
		keymeta = &meta

		return nil
	}
}

func testAccCheckAlicloudKmsKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_kms_key" {
			continue
		}

		kmsService := &KmsService{client: client}
		key, err := kmsService.DescribeKmsKey(rs.Primary.ID)

		if err != nil && !NotFoundError(err) {
			return WrapError(err)
		}
		if strings.Contains(err.Error(), "Provider ERROR") {
			return nil
		}

		return WrapError(fmt.Errorf("KMS key still exists:\n%#v", key.KeyMetadata))
	}

	return nil
}

const testAlicloudKmsKeyBasic = `
resource "alicloud_kms_key" "key" {
    is_enabled = true
    deletion_window_in_days = 7
    description = "tf_testaccCheckAlicloudKmsKeyDataSourceBasic"
}`

const testAlicloudKmsKeyUpdate = `
resource "alicloud_kms_key" "key" {
    is_enabled = false
    deletion_window_in_days = 7 
    description = "tf_testaccCheckAlicloudKmsKeyDataSourceBasic"
}`
