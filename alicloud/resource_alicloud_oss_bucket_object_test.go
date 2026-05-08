package alicloud

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlicloudOssBucketObject_basic(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "tf-oss-object-test-acc-source")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// first write some data to the tempfile just so it's not 0 bytes.
	err = ioutil.WriteFile(tmpFile.Name(), []byte("{anything will do }"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	var v http.Header
	resourceId := "alicloud_oss_bucket_object.default"
	ra := resourceAttrInit(resourceId, ossBucketObjectBasicMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-object-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOssBucketObjectConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  testAccCheckAlicloudOssBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":       "${alicloud_oss_bucket_public_access_block.default.bucket}",
					"key":          "test-object-source-key",
					"source":       tmpFile.Name(),
					"content_type": "binary/octet-stream",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudOssBucketObjectExists(
						"alicloud_oss_bucket_object.default", name, v),
					testAccCheck(map[string]string{
						"bucket": name,
						"source": tmpFile.Name(),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"source":  REMOVEKEY,
					"content": "some words for test oss object content",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudOssBucketObjectExists(
						"alicloud_oss_bucket_object.default", name, v),
					testAccCheck(map[string]string{
						"source":  REMOVEKEY,
						"content": "some words for test oss object content",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"acl": "public-read",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"acl": "public-read",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_side_encryption": "KMS",
					"kms_key_id":             "${data.alicloud_kms_keys.enabled.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":                 "${alicloud_oss_bucket_public_access_block.default.bucket}",
					"server_side_encryption": "AES256",
					"kms_key_id":             REMOVEKEY,
					"key":                    "test-object-source-key",
					"content":                REMOVEKEY,
					"source":                 tmpFile.Name(),
					"content_type":           "binary/octet-stream",
					"acl":                    REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudOssBucketObjectExists(
						"alicloud_oss_bucket_object.default", name, v),
					testAccCheck(map[string]string{
						"bucket":       name,
						"key":          "test-object-source-key",
						"content":      REMOVEKEY,
						"source":       tmpFile.Name(),
						"content_type": "binary/octet-stream",
						"acl":          "private",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudOssBucketObject_worm(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "tf-oss-object-worm-source")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if err := ioutil.WriteFile(tmpFile.Name(), []byte("worm content"), 0644); err != nil {
		t.Fatal(err)
	}

	var v http.Header
	resourceId := "alicloud_oss_bucket_object.default"
	ra := resourceAttrInit(resourceId, ossBucketObjectBasicMap)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-object-worm-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOssBucketObjectVersioningDependence)

	// Retention windows are computed at test start. They must satisfy two
	// conflicting constraints:
	//   1. They must still be in the future when the corresponding Apply
	//      reaches the OSS API — otherwise OSS rejects with
	//      "The retain until date must be in the future."
	//   2. They must have expired by the time the test framework destroys
	//      the object — otherwise DeleteObject is blocked by COMPLIANCE.
	// 60s/90s gives Step 0 (bucket + worm config + object create) and
	// Step 1 (extend) ample time, and the final Check below sleeps past
	// extendedUntil so destroy can succeed without retries.
	//
	// The OSS retain-until-date parameter is ISO8601 with millisecond
	// precision (e.g. 2026-09-30T00:00:00.000Z). Use the same layout so
	// testCheck's literal string compare matches what OSS echoes back.
	const iso8601Ms = "2006-01-02T15:04:05.000Z"
	initialUntil := time.Now().UTC().Add(60 * time.Second).Format(iso8601Ms)
	extendedUntilTime := time.Now().UTC().Add(90 * time.Second)
	extendedUntil := extendedUntilTime.Format(iso8601Ms)

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAlicloudOssBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				// The bucket field references
				// alicloud_oss_bucket_object_worm_configuration so that the
				// object is created only after object-level worm has been
				// enabled on the bucket — without this, PutObject would be
				// rejected because the worm header is not yet accepted.
				Config: testAccConfig(map[string]interface{}{
					"bucket":                        "${alicloud_oss_bucket_object_worm_configuration.default.bucket_name}",
					"key":                           "test-object-worm-key",
					"source":                        tmpFile.Name(),
					"content_type":                  "binary/octet-stream",
					"object_worm_mode":              "COMPLIANCE",
					"object_worm_retain_until_date": initialUntil,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudOssBucketObjectExists(resourceId, name, v),
					testAccCheck(map[string]string{
						"bucket":                        name,
						"key":                           "test-object-worm-key",
						"object_worm_mode":              "COMPLIANCE",
						"object_worm_retain_until_date": initialUntil,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":                        "${alicloud_oss_bucket_object_worm_configuration.default.bucket_name}",
					"key":                           "test-object-worm-key",
					"source":                        tmpFile.Name(),
					"content_type":                  "binary/octet-stream",
					"object_worm_mode":              "COMPLIANCE",
					"object_worm_retain_until_date": extendedUntil,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudOssBucketObjectExists(resourceId, name, v),
					testAccCheck(map[string]string{
						"object_worm_retain_until_date": extendedUntil,
					}),
					// Block until the retention has expired so the framework
					// can DeleteObject during destroy without hitting the
					// COMPLIANCE-mode worm.
					func(s *terraform.State) error {
						wait := time.Until(extendedUntilTime) + 10*time.Second
						if wait > 0 {
							t.Logf("waiting %v for object worm retention to expire before destroy", wait)
							time.Sleep(wait)
						}
						return nil
					},
				),
			},
		},
	})
}

func resourceOssBucketObjectConfigDependence(name string) string {

	return fmt.Sprintf(`
resource "alicloud_oss_bucket" "default" {
	bucket = "%s"
}
resource "alicloud_oss_bucket_public_access_block" "default" {
	bucket              = alicloud_oss_bucket.default.bucket
	block_public_access = false
}
data "alicloud_kms_keys" "enabled" {
	status = "Enabled"
}
`, name)
}

// resourceOssBucketObjectVersioningDependence is like
// resourceOssBucketObjectConfigDependence but tailored for the worm test:
// bucket versioning is enabled (a prerequisite for object-level worm),
// force_destroy = true so DeleteBucket can clean up the original versions
// (DeleteObject on a versioned bucket only writes a delete marker), and an
// alicloud_oss_bucket_object_worm_configuration resource turns on
// object-level worm so PutObject accepts x-oss-object-worm-* headers.
func resourceOssBucketObjectVersioningDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_oss_bucket" "default" {
	bucket = "%s"
	force_destroy = true
	versioning {
		status = "Enabled"
	}
}
resource "alicloud_oss_bucket_object_worm_configuration" "default" {
	bucket_name         = alicloud_oss_bucket.default.bucket
	object_worm_enabled = "Enabled"
}
data "alicloud_kms_keys" "enabled" {
	status = "Enabled"
}
`, name)
}

var ossBucketObjectBasicMap = map[string]string{
	"bucket":       CHECKSET,
	"key":          "test-object-source-key",
	"source":       CHECKSET,
	"content_type": "binary/octet-stream",
	"acl":          "private",
}

func testAccCheckAlicloudOssBucketObjectExists(n string, bucket string, obj http.Header) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckOssBucketObjectExistsWithProviders(n, bucket, obj, &providers)
}
func testAccCheckOssBucketObjectExistsWithProviders(n string, bucket string, obj http.Header, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		for _, provider := range *providers {
			// Ignore if Meta is empty, this can happen for validation providers
			if provider.Meta() == nil {
				continue
			}
			client := provider.Meta().(*connectivity.AliyunClient)
			raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
				return ossClient.Bucket(bucket)
			})
			buck, _ := raw.(*oss.Bucket)
			if err != nil {
				return fmt.Errorf("Error getting bucket: %#v", err)
			}
			object, err := buck.GetObjectMeta(rs.Primary.ID)
			log.Printf("[WARN]get oss bucket object %#v", bucket)
			if err == nil {
				if object != nil {
					obj = object
					return nil
				}
				continue
			} else if err != nil {
				return err

			}
		}

		return fmt.Errorf("Bucket not found")
	}
}
func testAccCheckAlicloudOssBucketObjectDestroy(s *terraform.State) error {
	return testAccCheckOssBucketObjectDestroyWithProvider(s, testAccProvider)
}

func testAccCheckOssBucketObjectDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	client := provider.Meta().(*connectivity.AliyunClient)
	var bucket *oss.Bucket
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_oss_bucket" {
			continue
		}
		raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			return ossClient.Bucket(rs.Primary.ID)
		})
		if err != nil {
			return fmt.Errorf("Error getting bucket: %#v", err)
		}
		bucket, _ = raw.(*oss.Bucket)
	}
	if bucket == nil {
		return nil
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_oss_bucket_object" {
			continue
		}

		// Try to find the resource
		exist, err := bucket.IsObjectExist(rs.Primary.ID)
		if err != nil {
			if IsExpectedErrors(err, []string{"NoSuchBucket"}) {
				return nil
			}
			return fmt.Errorf("IsObjectExist got an error: %#v", err)
		}

		if !exist {
			return nil
		}

		return fmt.Errorf("Found oss object: %s", rs.Primary.ID)
	}

	return nil
}
