package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

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

	request := map[string]interface{}{
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
		"RegionId":   client.RegionId,
	}
	action := "ListKeys"

	var response map[string]interface{}
	conn, err := client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}
	sweeped := false
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_kms_key", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Keys.Key", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Keys.Key", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			if _, ok := item["Description"]; !ok {
				continue
			}
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["Description"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Kms Key: %s", item["Description"].(string))
				continue
			}
			sweeped = true
			action = "ScheduleKeyDeletion"
			request := map[string]interface{}{
				"KeyId":               item["KeyId"],
				"PendingWindowInDays": 7,
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Kms Key (%s): %s", item["Description"].(string), err)
			}
			log.Printf("[INFO] Delete Kms Key success: %s ", item["Description"].(string))
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudKmsKey_basic(t *testing.T) {
	var v map[string]interface{}
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
			// eu-central-1 not support Aliyun_SM4
			testAccPreCheckWithRegions(t, false, connectivity.KmsKeyUnSupportRegions)
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
					"status": "Disabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Disabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Enabled",
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
	"status":              "Enabled",
	"key_usage":           "ENCRYPT/DECRYPT",
	"last_rotation_date":  CHECKSET,
	"origin":              "Aliyun_KMS",
	"primary_key_version": CHECKSET,
	"protection_level":    "SOFTWARE",
}

func KmsKeyBasicdependence(name string) string {
	return ""
}
