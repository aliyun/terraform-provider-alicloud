package oss_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudOssBucketTags_basic(t *testing.T) {
	resourceId := "alicloud_oss_bucket_tags.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbuckettags%d", acctest.DefaultRegionToTest, rand)
	region := acctest.DefaultRegionToTest

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckOssBucketTagsDestroy,
		Steps: []resource.TestStep{
			{
				// Create with a non-empty tag set.
				Config: testAccAliCloudOssBucketTagsConfig(name, region, `{
  Created = "TF"
  Env     = "test"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "bucket", name),
					resource.TestCheckResourceAttr(resourceId, "id", name),
					resource.TestCheckResourceAttr(resourceId, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceId, "tags.Created", "TF"),
					resource.TestCheckResourceAttr(resourceId, "tags.Env", "test"),
				),
			},
			{
				// Update the tag set: change a value, remove a key, add a key.
				Config: testAccAliCloudOssBucketTagsConfig(name, region, `{
  Created = "TF"
  Env     = "prod"
  Owner   = "team"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceId, "tags.Env", "prod"),
					resource.TestCheckResourceAttr(resourceId, "tags.Owner", "team"),
					resource.TestCheckResourceAttr(resourceId, "tags.Created", "TF"),
				),
			},
			{
				// Clear the tag set with an empty map; DeleteBucketTags runs.
				Config: testAccAliCloudOssBucketTagsConfig(name, region, `{}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "0"),
				),
			},
			{
				// Re-add tags after clearing.
				Config: testAccAliCloudOssBucketTagsConfig(name, region, `{
  Foo = "bar"
}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceId, "tags.Foo", "bar"),
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

func testAccCheckOssBucketTagsDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_oss_bucket_tags" {
			continue
		}
		bucket := rs.Primary.ID
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		_, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
			result, e := ossClient.GetBucketTagging(bucket)
			if e != nil {
				return nil, e
			}
			return result, nil
		})
		if err != nil {
			// NoSuchTagSet (tags cleared) and NoSuchBucket (bucket gone) both
			// mean the managed tag set no longer exists.
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing bucket tags %s failed: %s", bucket, err)
		}
		return fmt.Errorf("bucket tags %s still exist", bucket)
	}
	return nil
}

func testAccAliCloudOssBucketTagsConfig(name, region, tags string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "region" {
  default = "%s"
}

provider "alicloud" {
  region = var.region
}

resource "alicloud_oss_bucket" "default" {
  storage_class = "Standard"
  bucket        = var.name
  lifecycle {
    # alicloud_oss_bucket_tags owns the bucket's tag set; keep the inline
    # tags attribute of alicloud_oss_bucket from fighting this resource.
    ignore_changes = [tags]
  }
}

resource "alicloud_oss_bucket_tags" "default" {
  bucket = alicloud_oss_bucket.default.bucket
  tags   = %s
}
`, name, region, tags)
}
