package alicloud

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_kms_key", &resource.Sweeper{
		Name: "alicloud_kms_key",
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
			if strings.HasPrefix(strings.ToLower(key.Description), strings.ToLower(description)) {
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

func TestAccAlicloudKMSKey_basic(t *testing.T) {
	var v kms.KeyMetadata
	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, KmsKeyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccKmsKey%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, KmsKeyBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":            name,
					"key_spec":               "Aliyun_SM4",
					"protection_level":       "HSM",
					"pending_window_in_days": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":            name,
						"key_spec":               "Aliyun_SM4",
						"protection_level":       "HSM",
						"pending_window_in_days": "7",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pending_window_in_days", "deletion_window_in_days", "is_enabled"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_state": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_state": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_state": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_state": "Enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "from_terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "from_terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"automatic_rotation": "Enabled",
					"rotation_interval":  "2678400s",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"automatic_rotation": "Enabled",
						"rotation_interval":  "2678400s",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":        name,
					"automatic_rotation": "Disabled",
					"rotation_interval":  REMOVEKEY,
				}),

				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":        name,
						"automatic_rotation": "Disabled",
						"rotation_interval":  REMOVEKEY,
					}),
				),
			},
		},
	})
}

var KmsKeyMap = map[string]string{
	"arn":                 CHECKSET,
	"automatic_rotation":  "Disabled",
	"creation_date":       CHECKSET,
	"creator":             CHECKSET,
	"key_state":           "Enabled",
	"key_usage":           "ENCRYPT/DECRYPT",
	"last_rotation_date":  CHECKSET,
	"origin":              "Aliyun_KMS",
	"primary_key_version": CHECKSET,
	"protection_level":    "SOFTWARE",
}

func KmsKeyBasicdependence(name string) string {
	return ""
}
