package alicloud

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	ossv2 "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAlicloudOssBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":       "${alicloud_oss_bucket.default.bucket}",
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
					"bucket":                 "${alicloud_oss_bucket.default.bucket}",
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
	//      (Step 1 and Step 2) actually reaches the OSS API — otherwise
	//      OSS rejects with "The retain until date must be in the future."
	//   2. They must have expired by the time the test framework destroys
	//      the object — otherwise DeleteObject is blocked by COMPLIANCE.
	// 60s/90s gives Step 0 (bucket create) and Step 1/2 (apply) ample
	// time, and the final Check below sleeps past extendedUntil so the
	// destroy phase can succeed without retries.
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
				// Pre-flight: only create the bucket (with versioning, a
				// prerequisite for object-level worm). Object-level worm must
				// be enabled on the bucket before any worm header is accepted
				// on PutObject; we enable it in the next step's PreConfig.
				Config: resourceOssBucketObjectVersioningDependence(name),
			},
			{
				PreConfig: func() { enableBucketObjectWorm(t, name) },
				Config: testAccConfig(map[string]interface{}{
					"bucket":                        "${alicloud_oss_bucket.default.bucket}",
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
					"bucket":                        "${alicloud_oss_bucket.default.bucket}",
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

// enableBucketObjectWorm turns on object-level worm on the given bucket so
// that subsequent PutObject calls accept x-oss-object-worm-* headers. No
// default retention is configured; per-object retention is supplied via
// the object's own worm headers in the test config.
func enableBucketObjectWorm(t *testing.T, bucket string) {
	t.Helper()
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	_, err := client.WithOssClientV2(func(c *ossv2.Client) (interface{}, error) {
		return c.PutBucketObjectWormConfiguration(context.Background(), &ossv2.PutBucketObjectWormConfigurationRequest{
			Bucket: ossv2.Ptr(bucket),
			ObjectWormConfiguration: &ossv2.ObjectWormConfiguration{
				ObjectWormEnabled: ossv2.Ptr("Enabled"),
			},
		})
	})
	if err != nil {
		t.Fatalf("PutBucketObjectWormConfiguration on %s failed: %v", bucket, err)
	}
}

func resourceOssBucketObjectConfigDependence(name string) string {

	return fmt.Sprintf(`
resource "alicloud_oss_bucket" "default" {
	bucket = "%s"
}
data "alicloud_kms_keys" "enabled" {
	status = "Enabled"
}
`, name)
}

// resourceOssBucketObjectVersioningDependence is like
// resourceOssBucketObjectConfigDependence but with bucket versioning enabled,
// which is a prerequisite for object-level worm. force_destroy = true is
// required because DeleteObject on a versioned bucket only writes a delete
// marker — the original versions (including the one previously held by the
// worm retention) remain, and DeleteBucket would otherwise fail with
// BucketNotEmpty.
func resourceOssBucketObjectVersioningDependence(name string) string {
	return fmt.Sprintf(`
resource "alicloud_oss_bucket" "default" {
	bucket = "%s"
	force_destroy = true
	versioning {
		status = "Enabled"
	}
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
