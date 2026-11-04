package oss_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	osssvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/oss"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAliCloudOssBucketLifecycle_basic exercises a single rule with a days
// based expiration, then updates it to a transition and re-checks import. The
// lifecycle_rule_allow_same_action_overlap flag is plan-only (the API never
// echoes it back), so it is ignored during ImportStateVerify, exactly as the
// embedded bucket implementation ignores it.
func TestAccAliCloudOssBucketLifecycle_basic(t *testing.T) {
	resourceId := "alicloud_oss_bucket_lifecycle.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketlc%d", acctest.DefaultRegionToTest, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckOssBucketLifecycleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudOssBucketLifecycleBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "bucket", name),
					resource.TestCheckResourceAttr(resourceId, "id", name),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.id", "rule1"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.prefix", "path1/"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.expiration.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.expiration.0.days", "365"),
				),
			},
			{
				Config: testAccAliCloudOssBucketLifecycleBasicUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.id", "rule1"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.transitions.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.transitions.0.days", "30"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.transitions.0.storage_class", "IA"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lifecycle_rule_allow_same_action_overlap"},
			},
		},
	})
}

func TestAccAliCloudOssBucketLifecycle_expiration(t *testing.T) {
	resourceId := "alicloud_oss_bucket_lifecycle.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketlcexp%d", acctest.DefaultRegionToTest, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckOssBucketLifecycleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudOssBucketLifecycleExpirationDate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.expiration.0.date", "2025-01-12"),
				),
			},
			{
				Config: testAccAliCloudOssBucketLifecycleExpirationCreatedBefore(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.expiration.0.created_before_date", "2025-01-12"),
				),
			},
			{
				Config: testAccAliCloudOssBucketLifecycleExpirationDeleteMarker(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.expiration.0.expired_object_delete_marker", "true"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lifecycle_rule_allow_same_action_overlap"},
			},
		},
	})
}

func TestAccAliCloudOssBucketLifecycle_abortMultipartUpload(t *testing.T) {
	resourceId := "alicloud_oss_bucket_lifecycle.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketlcamu%d", acctest.DefaultRegionToTest, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckOssBucketLifecycleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudOssBucketLifecycleAbortMultipartDays(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.abort_multipart_upload.0.days", "10"),
				),
			},
			{
				Config: testAccAliCloudOssBucketLifecycleAbortMultipartCreatedBefore(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.abort_multipart_upload.0.created_before_date", "2025-01-22"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lifecycle_rule_allow_same_action_overlap"},
			},
		},
	})
}

func TestAccAliCloudOssBucketLifecycle_noncurrentVersion(t *testing.T) {
	resourceId := "alicloud_oss_bucket_lifecycle.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketlcnv%d", acctest.DefaultRegionToTest, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckOssBucketLifecycleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudOssBucketLifecycleNoncurrentVersion(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.noncurrent_version_expiration.0.days", "30"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.noncurrent_version_transition.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.noncurrent_version_transition.0.days", "10"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.noncurrent_version_transition.0.storage_class", "Archive"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lifecycle_rule_allow_same_action_overlap"},
			},
		},
	})
}

func TestAccAliCloudOssBucketLifecycle_filter(t *testing.T) {
	resourceId := "alicloud_oss_bucket_lifecycle.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketlcflt%d", acctest.DefaultRegionToTest, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckOssBucketLifecycleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudOssBucketLifecycleFilter(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.transitions.0.days", "30"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.transitions.0.storage_class", "IA"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.filter.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.filter.0.object_size_greater_than", "100"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.filter.0.object_size_less_than", "1024"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.filter.0.not.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.filter.0.not.0.prefix", "path1/exclude/"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.filter.0.not.0.tag.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.filter.0.not.0.tag.0.key", "env"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.filter.0.not.0.tag.0.value", "prod"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lifecycle_rule_allow_same_action_overlap"},
			},
		},
	})
}

func TestAccAliCloudOssBucketLifecycle_tags(t *testing.T) {
	resourceId := "alicloud_oss_bucket_lifecycle.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketlctag%d", acctest.DefaultRegionToTest, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckOssBucketLifecycleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudOssBucketLifecycleTags(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.tags.%", "2"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.tags.key2", "value2"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lifecycle_rule_allow_same_action_overlap"},
			},
		},
	})
}

// TestAccAliCloudOssBucketLifecycle_multiRule mirrors the multi-rule scenario
// kept in the embedded alicloud_oss_bucket test: an expiration-days rule, a
// transition rule, and an abort_multipart_upload rule side by side.
func TestAccAliCloudOssBucketLifecycle_multiRule(t *testing.T) {
	resourceId := "alicloud_oss_bucket_lifecycle.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketlcmulti%d", acctest.DefaultRegionToTest, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckOssBucketLifecycleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudOssBucketLifecycleMultiRule(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.#", "3"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.id", "rule1"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.0.expiration.0.days", "365"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.1.id", "rule2"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.1.transitions.0.days", "30"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.1.transitions.0.storage_class", "Archive"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.2.id", "rule3"),
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.2.abort_multipart_upload.0.days", "10"),
				),
			},
			{
				// Removing every rule makes the resource delete the lifecycle
				// configuration entirely, which is the same semantics the embedded
				// bucket implementation used when lifecycle_rule was emptied.
				Config: testAccAliCloudOssBucketLifecycleEmpty(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "lifecycle_rule.#", "0"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lifecycle_rule_allow_same_action_overlap"},
			},
		},
	})
}

func testAccCheckOssBucketLifecycleDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_oss_bucket_lifecycle" {
			continue
		}
		bucket := rs.Primary.ID
		service := osssvc.NewOssServiceV2(client)
		_, err := service.DescribeOssBucketLifecycle(bucket)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing oss bucket lifecycle %s failed: %s", bucket, err)
		}
		return fmt.Errorf("oss bucket lifecycle %s still exists", bucket)
	}
	return nil
}

func testAccAliCloudOssBucketLifecycleBucket(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_oss_bucket" "default" {
  storage_class = "Standard"
  bucket        = var.name

  # The standalone alicloud_oss_bucket_lifecycle resource manages the same
  # lifecycle configuration that the embedded lifecycle_rule block would
  # otherwise own. ignore_changes prevents the two resources from fighting
  # over the same API state, matching the cors/website/logging/versioning
  # split pattern.
  lifecycle {
    ignore_changes = [
      lifecycle_rule,
    ]
  }
}
`, name)
}

func testAccAliCloudOssBucketLifecycleBasic(name string) string {
	return fmt.Sprintf(`
%s
resource "alicloud_oss_bucket_lifecycle" "default" {
  bucket = alicloud_oss_bucket.default.bucket

  lifecycle_rule {
    id      = "rule1"
    prefix  = "path1/"
    enabled = true

    expiration {
      days = 365
    }
  }
}
`, testAccAliCloudOssBucketLifecycleBucket(name))
}

func testAccAliCloudOssBucketLifecycleBasicUpdate(name string) string {
	return fmt.Sprintf(`
%s
resource "alicloud_oss_bucket_lifecycle" "default" {
  bucket = alicloud_oss_bucket.default.bucket

  lifecycle_rule {
    id      = "rule1"
    prefix  = "path1/"
    enabled = false

    transitions {
      days          = 30
      storage_class = "IA"
    }
  }
}
`, testAccAliCloudOssBucketLifecycleBucket(name))
}

func testAccAliCloudOssBucketLifecycleExpirationDate(name string) string {
	return fmt.Sprintf(`
%s
resource "alicloud_oss_bucket_lifecycle" "default" {
  bucket = alicloud_oss_bucket.default.bucket

  lifecycle_rule {
    id      = "rule1"
    prefix  = "path1/"
    enabled = true

    expiration {
      date = "2025-01-12"
    }
  }
}
`, testAccAliCloudOssBucketLifecycleBucket(name))
}

func testAccAliCloudOssBucketLifecycleExpirationCreatedBefore(name string) string {
	return fmt.Sprintf(`
%s
resource "alicloud_oss_bucket_lifecycle" "default" {
  bucket = alicloud_oss_bucket.default.bucket

  lifecycle_rule {
    id      = "rule1"
    prefix  = "path1/"
    enabled = true

    expiration {
      created_before_date = "2025-01-12"
    }
  }
}
`, testAccAliCloudOssBucketLifecycleBucket(name))
}

func testAccAliCloudOssBucketLifecycleExpirationDeleteMarker(name string) string {
	return fmt.Sprintf(`
%s
resource "alicloud_oss_bucket_lifecycle" "default" {
  bucket = alicloud_oss_bucket.default.bucket

  lifecycle_rule {
    id      = "rule1"
    prefix  = "path1/"
    enabled = true

    expiration {
      expired_object_delete_marker = true
    }
  }
}
`, testAccAliCloudOssBucketLifecycleBucket(name))
}

func testAccAliCloudOssBucketLifecycleAbortMultipartDays(name string) string {
	return fmt.Sprintf(`
%s
resource "alicloud_oss_bucket_lifecycle" "default" {
  bucket = alicloud_oss_bucket.default.bucket

  lifecycle_rule {
    id      = "rule1"
    prefix  = "path1/"
    enabled = true

    abort_multipart_upload {
      days = 10
    }
  }
}
`, testAccAliCloudOssBucketLifecycleBucket(name))
}

func testAccAliCloudOssBucketLifecycleAbortMultipartCreatedBefore(name string) string {
	return fmt.Sprintf(`
%s
resource "alicloud_oss_bucket_lifecycle" "default" {
  bucket = alicloud_oss_bucket.default.bucket

  lifecycle_rule {
    id      = "rule1"
    prefix  = "path1/"
    enabled = true

    abort_multipart_upload {
      created_before_date = "2025-01-22"
    }
  }
}
`, testAccAliCloudOssBucketLifecycleBucket(name))
}

func testAccAliCloudOssBucketLifecycleNoncurrentVersion(name string) string {
	return fmt.Sprintf(`
%s
resource "alicloud_oss_bucket_lifecycle" "default" {
  bucket = alicloud_oss_bucket.default.bucket

  lifecycle_rule {
    id      = "rule1"
    prefix  = "path1/"
    enabled = true

    noncurrent_version_expiration {
      days = 30
    }

    noncurrent_version_transition {
      days          = 10
      storage_class = "Archive"
    }
  }
}
`, testAccAliCloudOssBucketLifecycleBucket(name))
}

func testAccAliCloudOssBucketLifecycleFilter(name string) string {
	return fmt.Sprintf(`
%s
resource "alicloud_oss_bucket_lifecycle" "default" {
  bucket = alicloud_oss_bucket.default.bucket

  lifecycle_rule {
    id      = "rule1"
    prefix  = "path1/"
    enabled = true

    # OSS requires every lifecycle rule to carry at least one action; a rule
    # exercising only filter/tags still needs an action block to be accepted.
    # Use transitions to match the proven SDKv2 filter test pattern.
    #
    # OSS API constraint (Ec=0014-00000035): when a Prefix is set on the Rule,
    # the Prefix inside Filter/Not MUST start with that same Rule prefix. Here
    # the Rule prefix is "path1/", so the Not prefix must start with "path1/".
    transitions {
      days          = 30
      storage_class = "IA"
    }

    filter {
      object_size_greater_than = 100
      object_size_less_than     = 1024

      not {
        prefix = "path1/exclude/"

        tag {
          key   = "env"
          value = "prod"
        }
      }
    }
  }
}
`, testAccAliCloudOssBucketLifecycleBucket(name))
}

func testAccAliCloudOssBucketLifecycleTags(name string) string {
	return fmt.Sprintf(`
%s
resource "alicloud_oss_bucket_lifecycle" "default" {
  bucket = alicloud_oss_bucket.default.bucket

  lifecycle_rule {
    id      = "rule1"
    prefix  = "path1/"
    enabled = true

    # OSS requires every lifecycle rule to carry at least one action; a rule
    # exercising only filter/tags still needs an action block to be accepted.
    expiration {
      days = 365
    }

    tags = {
      key1 = "value1"
      key2 = "value2"
    }
  }
}
`, testAccAliCloudOssBucketLifecycleBucket(name))
}

func testAccAliCloudOssBucketLifecycleMultiRule(name string) string {
	return fmt.Sprintf(`
%s
resource "alicloud_oss_bucket_lifecycle" "default" {
  bucket = alicloud_oss_bucket.default.bucket

  lifecycle_rule {
    id      = "rule1"
    prefix  = "path1/"
    enabled = true

    expiration {
      days = 365
    }
  }

  lifecycle_rule {
    id      = "rule2"
    prefix  = "path2/"
    enabled = true

    transitions {
      days          = 30
      storage_class = "Archive"
    }
  }

  lifecycle_rule {
    id      = "rule3"
    prefix  = "path3/"
    enabled = true

    abort_multipart_upload {
      days = 10
    }
  }
}
`, testAccAliCloudOssBucketLifecycleBucket(name))
}

func testAccAliCloudOssBucketLifecycleEmpty(name string) string {
	return fmt.Sprintf(`
%s
resource "alicloud_oss_bucket_lifecycle" "default" {
  bucket = alicloud_oss_bucket.default.bucket
}
`, testAccAliCloudOssBucketLifecycleBucket(name))
}
