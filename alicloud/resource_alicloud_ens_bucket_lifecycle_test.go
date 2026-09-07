package alicloud

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudEnsBucketLifecycle_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_bucket_lifecycle.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsBucketLifecycleMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsBucketLifecycle")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	bucketName := fmt.Sprintf("tf-testacc-ens-bl-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, bucketName, AlicloudEnsBucketLifecycleBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			createEnsBucketForLifecycleTest(t, bucketName)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rac.checkResourceDestroy(),
			func(*terraform.State) error {
				deleteEnsBucketForLifecycleTest(bucketName)
				return nil
			},
		),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name":               bucketName,
					"status":                    "Enabled",
					"prefix":                    "logs/",
					"expiration_days":           7,
					"allow_same_action_overlap": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name":               bucketName,
						"status":                    "Enabled",
						"prefix":                    "logs/",
						"expiration_days":           "7",
						"allow_same_action_overlap": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name":               bucketName,
					"status":                    "Disabled",
					"prefix":                    "logs2/",
					"expiration_days":           7,
					"allow_same_action_overlap": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":          "Disabled",
						"prefix":          "logs2/",
						"expiration_days": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name":               bucketName,
					"status":                    "Enabled",
					"prefix":                    "logs/",
					"expiration_days":           14,
					"allow_same_action_overlap": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":          "Enabled",
						"prefix":          "logs/",
						"expiration_days": "14",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"allow_same_action_overlap"},
			},
		},
	})
}

func TestAccAliCloudEnsBucketLifecycle_ruleId(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_bucket_lifecycle.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsBucketLifecycleMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsBucketLifecycle")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	bucketName := fmt.Sprintf("tf-testacc-ens-bl-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, bucketName, AlicloudEnsBucketLifecycleBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			createEnsBucketForLifecycleTest(t, bucketName)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rac.checkResourceDestroy(),
			func(*terraform.State) error {
				deleteEnsBucketForLifecycleTest(bucketName)
				return nil
			},
		),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name":               bucketName,
					"status":                    "Enabled",
					"prefix":                    "data/",
					"expiration_days":           30,
					"allow_same_action_overlap": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name":     bucketName,
						"rule_id":         CHECKSET,
						"status":          "Enabled",
						"prefix":          "data/",
						"expiration_days": "30",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name":               bucketName,
					"status":                    "Enabled",
					"prefix":                    "data2/",
					"expiration_days":           60,
					"allow_same_action_overlap": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"prefix":          "data2/",
						"expiration_days": "60",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"allow_same_action_overlap"},
			},
		},
	})
}

var AlicloudEnsBucketLifecycleMap = map[string]string{
	"rule_id": CHECKSET,
}

func AlicloudEnsBucketLifecycleBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}

func createEnsBucketForLifecycleTest(t *testing.T, bucketName string) {
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-beijing"
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Fatalf("createEnsBucket sharedClientForRegion: %v", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	action := "PutBucket"
	request := map[string]interface{}{
		"BucketName": bucketName,
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err = client.RpcPost("Ens", "2017-11-10", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("PutBucket %s: %v", bucketName, err)
	}
	// Poll GetBucketInfo until the bucket is queryable; during the creation
	// window the service may report NoSuchBucket.
	readyWait := incrementalWait(5*time.Second, 5*time.Second)
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		getAction := "GetBucketInfo"
		query := map[string]interface{}{"BucketName": bucketName}
		_, e := client.RpcPost("Ens", "2017-11-10", getAction, query, nil, true)
		if e != nil {
			if NotFoundError(e) || IsExpectedErrors(e, []string{"NoSuchBucket"}) {
				readyWait()
				return resource.RetryableError(e)
			}
			return resource.NonRetryableError(e)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("GetBucketInfo ready check for %s: %v", bucketName, err)
	}
}

func deleteEnsBucketForLifecycleTest(bucketName string) {
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-beijing"
	}
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return
	}
	client, ok := rawClient.(*connectivity.AliyunClient)
	if !ok {
		return
	}
	action := "DeleteBucket"
	query := map[string]interface{}{"BucketName": bucketName}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	_ = resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, e := client.RpcPost("Ens", "2017-11-10", action, query, nil, true)
		if e != nil {
			if NeedRetry(e) {
				wait()
				return resource.RetryableError(e)
			}
			return resource.NonRetryableError(e)
		}
		return nil
	})
}
