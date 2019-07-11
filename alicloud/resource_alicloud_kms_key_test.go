package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
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
	var v *kms.DescribeKeyResponse

	resourceId := "alicloud_kms_key.default"
	ra := resourceAttrInit(resourceId, kmsKeyBasicMap)

	serviceFunc := func() interface{} {
		return &KmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf_testKmsKey_%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceKmsKeyConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"key_usage":   "ENCRYPT/DECRYPT",
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_usage":   "ENCRYPT/DECRYPT",
						"description": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_window_in_days"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"is_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"is_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"is_enabled": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"is_enabled": "true",
					}),
				),
			},
		},
	})
}

func resourceKmsKeyConfigDependence(name string) string {
	return ""
}

var kmsKeyBasicMap = map[string]string{
	"description": CHECKSET,
	"key_usage":   "ENCRYPT/DECRYPT",
	"is_enabled":  "true",
	"arn":         CHECKSET,
}
