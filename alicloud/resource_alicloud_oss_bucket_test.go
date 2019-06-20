package alicloud

import (
	"fmt"
	"log"
	"testing"

	"strings"
	"time"

	"strconv"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_oss_bucket", &resource.Sweeper{
		Name: "alicloud_oss_bucket",
		F:    testSweepOSSBuckets,
	})
}

func testSweepOSSBuckets(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testacc",
		"tf-test-",
		"test-bucket-",
		"tf-oss-test-",
		"tf-object-test-",
		"test-acc-alicloud-",
	}

	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		return ossClient.ListBuckets()
	})
	if err != nil {
		return fmt.Errorf("Error retrieving OSS buckets: %s", err)
	}
	resp, _ := raw.(oss.ListBucketsResult)
	sweeped := false

	for _, v := range resp.Buckets {
		name := v.Name
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping OSS bucket: %s", name)
			continue
		}
		sweeped = true
		raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			return ossClient.Bucket(name)
		})
		if err != nil {
			return fmt.Errorf("Error getting bucket (%s): %#v", name, err)
		}
		bucket, _ := raw.(*oss.Bucket)
		if objects, err := bucket.ListObjects(); err != nil {
			log.Printf("[ERROR] Failed to list objects: %s", err)
		} else if len(objects.Objects) > 0 {
			for _, o := range objects.Objects {
				if err := bucket.DeleteObject(o.Key); err != nil {
					log.Printf("[ERROR] Failed to delete object (%s): %s.", o.Key, err)
				}
			}

		}

		log.Printf("[INFO] Deleting OSS bucket: %s", name)

		_, err = client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			return nil, ossClient.DeleteBucket(name)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete OSS bucket (%s): %s", name, err)
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudOssBucketBasic(t *testing.T) {
	var bucket oss.BucketInfo
	resourceName := "alicloud_oss_bucket.basic"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOssBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudOssBucketBasicConfig(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssBucketExists(
						"alicloud_oss_bucket.basic", &bucket),
					resource.TestCheckResourceAttrSet(
						"alicloud_oss_bucket.basic",
						"location"),
					resource.TestCheckResourceAttr(
						"alicloud_oss_bucket.basic",
						"acl",
						"public-read"),
					resource.TestCheckResourceAttr("alicloud_oss_bucket.basic", "lifecycle_rule.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})

}

func TestAccAlicloudOssBucketCors(t *testing.T) {
	var bucket oss.BucketInfo
	resourceName := "alicloud_oss_bucket.cors"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOssBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudOssBucketCorsConfig(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssBucketExists(
						"alicloud_oss_bucket.cors", &bucket),
					resource.TestCheckResourceAttr(
						"alicloud_oss_bucket.cors",
						"cors_rule.#",
						"2"),
					resource.TestCheckResourceAttr(
						"alicloud_oss_bucket.cors",
						"cors_rule.0.allowed_headers.0",
						"authorization"),
					resource.TestCheckResourceAttr("alicloud_oss_bucket.cors", "lifecycle_rule.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccAlicloudOssBucketWebsite(t *testing.T) {
	var bucket oss.BucketInfo
	resourceName := "alicloud_oss_bucket.website"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOssBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudOssBucketWebsiteConfig(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssBucketExists(
						"alicloud_oss_bucket.website", &bucket),
					resource.TestCheckResourceAttr(
						"alicloud_oss_bucket.website",
						"website.#",
						"1"),
					resource.TestCheckResourceAttr("alicloud_oss_bucket.website", "lifecycle_rule.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}
func TestAccAlicloudOssBucketLogging(t *testing.T) {
	var bucket oss.BucketInfo
	resourceName := "alicloud_oss_bucket.logging"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOssBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudOssBucketLoggingConfig(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssBucketExists(
						"alicloud_oss_bucket.target", &bucket),
					testAccCheckOssBucketExists(
						"alicloud_oss_bucket.logging", &bucket),
					resource.TestCheckResourceAttr(
						"alicloud_oss_bucket.logging",
						"logging.#",
						"1"),
					resource.TestCheckResourceAttr("alicloud_oss_bucket.logging", "lifecycle_rule.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccAlicloudOssBucketReferer(t *testing.T) {
	var bucket oss.BucketInfo
	resourceName := "alicloud_oss_bucket.referer"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_oss_bucket.referer",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOssBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudOssBucketRefererConfig(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssBucketExists(
						"alicloud_oss_bucket.referer", &bucket),
					resource.TestCheckResourceAttr(
						"alicloud_oss_bucket.referer",
						"referer_config.#",
						"1"),
					resource.TestCheckResourceAttr("alicloud_oss_bucket.referer", "lifecycle_rule.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}
func TestAccAlicloudOssBucketLifecycle(t *testing.T) {
	var bucket oss.BucketInfo
	resourceName := "alicloud_oss_bucket.lifecycle"

	hashcode1 := strconv.Itoa(expirationHash(map[string]interface{}{
		"days": 365,
		"date": "",
	}))
	hashcode2 := strconv.Itoa(expirationHash(map[string]interface{}{
		"days": 0,
		"date": "2018-01-12",
	}))
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOssBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudOssBucketLifecycleConfig(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssBucketExists(
						"alicloud_oss_bucket.lifecycle", &bucket),
					resource.TestCheckResourceAttr("alicloud_oss_bucket.lifecycle", "lifecycle_rule.#", "2"),
					resource.TestCheckResourceAttr("alicloud_oss_bucket.lifecycle", "lifecycle_rule.0.id", "rule1"),
					resource.TestCheckResourceAttr("alicloud_oss_bucket.lifecycle", "lifecycle_rule.0.prefix", "path1/"),
					resource.TestCheckResourceAttr("alicloud_oss_bucket.lifecycle", "lifecycle_rule.0.enabled", "true"),
					resource.TestCheckResourceAttr("alicloud_oss_bucket.lifecycle", "lifecycle_rule.0.expiration."+hashcode1+".days", "365"),
					resource.TestCheckResourceAttr("alicloud_oss_bucket.lifecycle", "lifecycle_rule.1.id", "rule2"),
					resource.TestCheckResourceAttr("alicloud_oss_bucket.lifecycle", "lifecycle_rule.1.prefix", "path2/"),
					resource.TestCheckResourceAttr("alicloud_oss_bucket.lifecycle", "lifecycle_rule.1.enabled", "true"),
					resource.TestCheckResourceAttr("alicloud_oss_bucket.lifecycle", "lifecycle_rule.1.expiration."+hashcode2+".date", "2018-01-12"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccAlicloudOssBucketPolicy(t *testing.T) {
	var bucket oss.BucketInfo
	resourceName := "alicloud_oss_bucket.policy"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOssBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudOssBucketPolicyConfig(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssBucketExists(
						"alicloud_oss_bucket.policy", &bucket),
					resource.TestCheckResourceAttr(
						"alicloud_oss_bucket.policy",
						"policy",
						"{\"Statement\":[{\"Action\":[\"oss:*\"],\"Effect\":\"Allow\",\"Resource\":[\"acs:oss:*:*:*\"]}],\"Version\":\"1\"}"),
					resource.TestCheckResourceAttr("alicloud_oss_bucket.policy", "lifecycle_rule.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccAlicloudOssBucketStorageClass(t *testing.T) {
	var bucket oss.BucketInfo
	resourceName := "alicloud_oss_bucket.storage"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOssBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudOssBucketStorageClassConfig(acctest.RandInt()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssBucketExists(
						"alicloud_oss_bucket.storage", &bucket),
					resource.TestCheckResourceAttr(
						"alicloud_oss_bucket.storage",
						"storage_class",
						"IA"),
					resource.TestCheckResourceAttr("alicloud_oss_bucket.storage", "lifecycle_rule.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccAlicloudOssBucketSseRule(t *testing.T) {
	var bucket oss.BucketInfo
	resourceName := "alicloud_oss_bucket.sserule"

	rand := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.OssSseSupportedRegions)
		},

		// module name
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOssBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudOssBucketSseRuleConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssBucketExists(
						"alicloud_oss_bucket.sserule", &bucket),
					resource.TestCheckResourceAttr(
						"alicloud_oss_bucket.sserule",
						"server_side_encryption_rule.0.sse_algorithm",
						"AES256"),
					resource.TestCheckResourceAttr("alicloud_oss_bucket.sserule", "lifecycle_rule.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccAlicloudOssBucketUpdateSseRuleConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssBucketExists(
						"alicloud_oss_bucket.sserule", &bucket),
					resource.TestCheckResourceAttr(
						"alicloud_oss_bucket.sserule",
						"server_side_encryption_rule.0.sse_algorithm",
						"KMS"),
				),
			},
			{
				Config: testAccAlicloudOssBucketDeleteSseRuleConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssBucketExists(
						"alicloud_oss_bucket.sserule", &bucket),
					resource.TestCheckResourceAttr(
						"alicloud_oss_bucket.sserule",
						"server_side_encryption_rule.#",
						"0"),
				),
			},
		},
	})
}

func TestAccAlicloudOssBucketTags(t *testing.T) {
	var bucket oss.BucketInfo
	resourceName := "alicloud_oss_bucket.tags"

	rand := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOssBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudOssBucketTagsConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssBucketExists(
						"alicloud_oss_bucket.tags", &bucket),
					resource.TestCheckResourceAttr(
						"alicloud_oss_bucket.tags",
						"tags.%",
						"2"),
					resource.TestCheckResourceAttr(
						"alicloud_oss_bucket.tags",
						"tags.key1",
						"value1"),
					resource.TestCheckResourceAttr(
						"alicloud_oss_bucket.tags",
						"tags.key2",
						"value2"),
					resource.TestCheckResourceAttr("alicloud_oss_bucket.tags", "lifecycle_rule.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccAlicloudOssBucketUpdateTagsConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssBucketExists(
						"alicloud_oss_bucket.tags", &bucket),
					resource.TestCheckResourceAttr(
						"alicloud_oss_bucket.tags",
						"tags.%",
						"3"),
					resource.TestCheckResourceAttr(
						"alicloud_oss_bucket.tags",
						"tags.key1-update",
						"value1-update"),
					resource.TestCheckResourceAttr(
						"alicloud_oss_bucket.tags",
						"tags.key2-update",
						"value2-update"),
					resource.TestCheckResourceAttr(
						"alicloud_oss_bucket.tags",
						"tags.key3-new",
						"value3-new"),
				),
			},
			{
				Config: testAccAlicloudOssBucketDeleteTagsConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssBucketExists(
						"alicloud_oss_bucket.tags", &bucket),
					resource.TestCheckResourceAttr(
						"alicloud_oss_bucket.tags",
						"tags.%",
						"0"),
				),
			},
		},
	})
}

func TestAccAlicloudOssBucketVersioning(t *testing.T) {
	var bucket oss.BucketInfo
	resourceName := "alicloud_oss_bucket.versioning"
	rand := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.OssVersioningSupportedRegions)
		},

		// module name
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOssBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudOssBucketVersioningConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssBucketExists(
						"alicloud_oss_bucket.versioning", &bucket),
					resource.TestCheckResourceAttr(
						"alicloud_oss_bucket.versioning",
						"versioning.0.status",
						"Enabled"),
					resource.TestCheckResourceAttr("alicloud_oss_bucket.versioning", "lifecycle_rule.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccAlicloudOssBucketUpdateVersioningConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOssBucketExists(
						"alicloud_oss_bucket.versioning", &bucket),
					resource.TestCheckResourceAttr(
						"alicloud_oss_bucket.versioning",
						"versioning.0.status",
						"Suspended"),
				),
			},
		},
	})
}

func testAccCheckOssBucketExists(n string, b *oss.BucketInfo) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckOssBucketExistsWithProviders(n, b, &providers)
}
func testAccCheckOssBucketExistsWithProviders(n string, b *oss.BucketInfo, providers *[]*schema.Provider) resource.TestCheckFunc {
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
			ossService := OssService{client}
			bucket, err := ossService.QueryOssBucketById(rs.Primary.ID)
			log.Printf("[WARN]get oss bucket %#v", bucket)
			if err == nil && bucket != nil {
				*b = *bucket
				return nil
			}

			// Verify the error is what we want
			e, _ := err.(*oss.ServiceError)
			if e.Code == OssBucketNotFound {
				continue
			}
			if err != nil {
				return err

			}
		}

		return fmt.Errorf("Bucket not found")
	}
}

func TestResourceAlicloudOssBucketAcl_validation(t *testing.T) {
	_, errors := validateOssBucketAcl("incorrect", "acl")
	if len(errors) == 0 {
		t.Fatalf("Expected to trigger a validation error")
	}

	var testCases = []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "public-read",
			ErrCount: 0,
		},
		{
			Value:    "public-read-write",
			ErrCount: 0,
		},
	}

	for _, tc := range testCases {
		_, errors := validateOssBucketAcl(tc.Value, "acl")
		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected not to trigger a validation error")
		}
	}
}

func testAccCheckOssBucketDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_oss_bucket" {
			continue
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ossService := OssService{client}

		// Try to find the resource
		bucket, err := ossService.QueryOssBucketById(rs.Primary.ID)
		if err != nil {
			// Verify the error is what we want
			if IsExceptedErrors(err, []string{OssBucketNotFound}) {
				continue
			}
			return err
		}
		if bucket.Name != "" {
			return fmt.Errorf("Found instance: %s", bucket.Name)
		}
	}

	return nil
}

func testAccAlicloudOssBucketBasicConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_oss_bucket" "basic" {
	bucket = "tf-testacc-bucket-basic-%d"
	acl = "public-read"
}
`, randInt)
}

func testAccAlicloudOssBucketCorsConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_oss_bucket" "cors" {
	bucket = "tf-testacc-bucket-cors-%d"
	cors_rule ={
		allowed_origins=["*"]
		allowed_methods=["PUT","GET"]
		allowed_headers=["authorization"]
	}
	cors_rule ={
		allowed_origins=["http://www.a.com", "http://www.b.com"]
		allowed_methods=["GET"]
		allowed_headers=["authorization"]
		expose_headers=["x-oss-test","x-oss-test1"]
		max_age_seconds=100
	}
}
`, randInt)
}

func testAccAlicloudOssBucketWebsiteConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_oss_bucket" "website"{
	bucket = "tf-testacc-bucket-website-%d"
	website = {
		index_document = "index.html"
		error_document = "error.html"
	}
}
`, randInt)
}

func testAccAlicloudOssBucketLoggingConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_oss_bucket" "target"{
	bucket = "tf-testacc-target-%d"
}
resource "alicloud_oss_bucket" "logging" {
	bucket = "tf-testacc-bucket-logging-%d"
	logging {
		target_bucket = "${alicloud_oss_bucket.target.id}"
		target_prefix = "log/"
	}
}
`, randInt, randInt)
}

func testAccAlicloudOssBucketRefererConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_oss_bucket" "referer" {
	bucket = "tf-testacc-bucket-referer-%d"
	referer_config {
		allow_empty = false
		referers = ["http://www.aliyun.com", "https://www.aliyun.com"]
	}
}
`, randInt)
}

func testAccAlicloudOssBucketLifecycleConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_oss_bucket" "lifecycle"{
	bucket = "tf-testacc-bucket-lifecycle-%d"
	lifecycle_rule {
		id = "rule1"
		prefix = "path1/"
		enabled = true
		expiration {
			days = 365
		}
	}
	lifecycle_rule {
		id = "rule2"
		prefix = "path2/"
		enabled = true
		expiration {
			date = "2018-01-12"
		}
	}
}
`, randInt)
}

func testAccAlicloudOssBucketStorageClassConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_oss_bucket" "storage" {
	bucket = "tf-testacc-bucket-storage-%d"
	storage_class= "IA"
}
`, randInt)
}

func testAccAlicloudOssBucketPolicyConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_oss_bucket" "policy" {
	bucket = "tf-testacc-bucket-policy-%d"
	policy= "{\"Statement\":[{\"Action\":[\"oss:*\"],\"Effect\":\"Allow\",\"Resource\":[\"acs:oss:*:*:*\"]}],\"Version\":\"1\"}"
}
`, randInt)
}

func testAccAlicloudOssBucketSseRuleConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_oss_bucket" "sserule" {
	bucket = "tf-testacc-bucket-sserule-%d"
	server_side_encryption_rule = {
		sse_algorithm = "AES256"
	}
}
`, randInt)
}

func testAccAlicloudOssBucketUpdateSseRuleConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_oss_bucket" "sserule" {
	bucket = "tf-testacc-bucket-sserule-%d"
	server_side_encryption_rule = {
		sse_algorithm = "KMS"
	}
}
`, randInt)
}

func testAccAlicloudOssBucketDeleteSseRuleConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_oss_bucket" "sserule" {
	bucket = "tf-testacc-bucket-sserule-%d"
}
`, randInt)
}

func testAccAlicloudOssBucketTagsConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_oss_bucket" "tags" {
	bucket = "tf-testacc-bucket-tags-%d"
	tags {
		key1 = "value1", 
		key2 = "value2", 
	}
}
`, randInt)
}

func testAccAlicloudOssBucketUpdateTagsConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_oss_bucket" "tags" {
	bucket = "tf-testacc-bucket-tags-%d"
	tags {
		key1-update = "value1-update", 
		key2-update = "value2-update", 
		key3-new = "value3-new", 
	}
}
`, randInt)
}

func testAccAlicloudOssBucketDeleteTagsConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_oss_bucket" "tags" {
	bucket = "tf-testacc-bucket-tags-%d"
}
`, randInt)
}

func testAccAlicloudOssBucketVersioningConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_oss_bucket" "versioning" {
	bucket = "tf-testacc-bucket-version-%d"
	
	versioning = {
		status = "Enabled"
	}
}
`, randInt)
}

func testAccAlicloudOssBucketUpdateVersioningConfig(randInt int) string {
	return fmt.Sprintf(`
resource "alicloud_oss_bucket" "versioning" {
	bucket = "tf-testacc-bucket-version-%d"
	
	versioning = {
		status = "Suspended"
	}
}
`, randInt)
}
