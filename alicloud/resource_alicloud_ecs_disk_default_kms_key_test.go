package alicloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudEcsDiskDefaultKmsKey_basic12760(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk_default_kms_key.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskDefaultKmsKeyMap12760)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDiskDefaultKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccecs%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskDefaultKmsKeyBasicDependence12760)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEcsDiskDefaultKmsKeyDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"kms_key_id": "${data.alicloud_kms_keys.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kms_key_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"kms_key_id": "${data.alicloud_kms_keys.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kms_key_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudEcsDiskDefaultKmsKey_basic12760_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk_default_kms_key.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsDiskDefaultKmsKeyMap12760)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsDiskDefaultKmsKey")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccecs%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsDiskDefaultKmsKeyBasicDependence12760)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEcsDiskDefaultKmsKeyDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"kms_key_id": "${data.alicloud_kms_keys.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"kms_key_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudEcsDiskDefaultKmsKeyMap12760 = map[string]string{}

func AliCloudEcsDiskDefaultKmsKeyBasicDependence12760(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_kms_keys" "default" {
  filters = "[{\"Key\":\"KeyState\",\"Values\":[\"Enabled\"]}]"
}
`, name)
}

func testAccCheckEcsDiskDefaultKmsKeyDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ecsServiceV2 := EcsServiceV2{client}

		object, err := testAccEcsManagedDefaultKmsKey(client)
		if err != nil {
			return err
		}

		defaultKmsKeyId := fmt.Sprint(object["KeyId"])

		objectRaw, err := ecsServiceV2.DescribeEcsDiskDefaultKmsKey(client.RegionId)
		if err != nil {
			return err
		}

		diskDefaultKmsKeyId := fmt.Sprint(objectRaw["KMSKeyId"])

		// Verify that the default Kms Key Id is now the account's default Kms Key Id.
		if defaultKmsKeyId != diskDefaultKmsKeyId {
			return fmt.Errorf("default CMK (%s) is not the account's default CMK (%s)", defaultKmsKeyId, diskDefaultKmsKeyId)
		}

		return nil
	}
}

// testAccEcsManagedDefaultKmsKey returns' the account's default Kms Key Id.
func testAccEcsManagedDefaultKmsKey(client *connectivity.AliyunClient) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	action := "ListKeys"

	request := map[string]interface{}{
		"Filters":    "[{\"Key\":\"AliasName\", \"Values\":[\"acs/ecs\"]}]",
		"PageSize":   PageSizeLarge,
		"PageNumber": 1,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_disk_default_kms_key", action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.Keys.Key", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, "alicloud_ecs_disk_default_kms_key", "$.Keys.Key", response)
	}

	object = resp.([]interface{})[0].(map[string]interface{})

	return object, nil
}
