package ens_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAliCloudEnsBucket_basic covers the full lifecycle of the ENS bucket
// Framework resource: a standard (global) bucket is created with a comment and
// an explicit private ACL, the ACL is then updated in place to public-read via
// PutBucketAcl, and finally the resource is imported to verify the ID-only
// import reproduces the applied state. The standard bucket form is used because
// PutBucket with only BucketName is the documented default and avoids the
// sink-only EnsRegionId coupling, matching the ENS package's existing test
// setup helper.
func TestAccAliCloudEnsBucket_basic(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	bucketName := fmt.Sprintf("tf-testacc-ens-bucket-%d", r)
	resourceName := "alicloud_ens_bucket.default"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckAliCloudEnsBucketDestroy(bucketName),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudEnsBucketConfig(bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", bucketName),
					resource.TestCheckResourceAttr(resourceName, "bucket_name", bucketName),
					resource.TestCheckResourceAttr(resourceName, "logical_bucket_type", "standard"),
					resource.TestCheckResourceAttr(resourceName, "bucket_acl", "private"),
					resource.TestCheckResourceAttr(resourceName, "comment", "tf-testacc-ens-comment"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttrSet(resourceName, "modify_time"),
				),
			},
			{
				Config: testAccAliCloudEnsBucketConfigUpdateAcl(bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bucket_acl", "public-read"),
					// ForceNew fields are unchanged across the in-place ACL update.
					resource.TestCheckResourceAttr(resourceName, "bucket_name", bucketName),
					resource.TestCheckResourceAttr(resourceName, "logical_bucket_type", "standard"),
					resource.TestCheckResourceAttr(resourceName, "comment", "tf-testacc-ens-comment"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// testAccAliCloudEnsBucketConfig is the basic standard-bucket configuration. A
// standard bucket needs only a bucket name; comment and bucket_acl are set
// explicitly so the in-place ACL update in the next step is observable.
func testAccAliCloudEnsBucketConfig(bucketName string) string {
	return fmt.Sprintf(`
resource "alicloud_ens_bucket" "default" {
  bucket_name = "%[1]s"
  comment     = "tf-testacc-ens-comment"
  bucket_acl  = "private"
}
`, bucketName)
}

// testAccAliCloudEnsBucketConfigUpdateAcl flips bucket_acl to public-read,
// exercising the PutBucketAcl update path while every other argument stays
// the same, so only the ACL changes in place.
func testAccAliCloudEnsBucketConfigUpdateAcl(bucketName string) string {
	return fmt.Sprintf(`
resource "alicloud_ens_bucket" "default" {
  bucket_name = "%[1]s"
  comment     = "tf-testacc-ens-comment"
  bucket_acl  = "public-read"
}
`, bucketName)
}

// testAccCheckAliCloudEnsBucketDestroy verifies the bucket is gone after the
// test destroys it. GetBucketInfo returns NoSuchBucket once the bucket is
// deleted; any other error is surfaced so a failed destroy is not hidden.
func testAccCheckAliCloudEnsBucketDestroy(bucketName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		_, err := client.RpcPost("Ens", "2017-11-10", "GetBucketInfo", nil, map[string]interface{}{
			"BucketName": bucketName,
		}, true)
		if err != nil {
			if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"NoSuchBucket"}) {
				return nil
			}
			return fmt.Errorf("checking ENS bucket %s destroy: %w", bucketName, err)
		}
		return fmt.Errorf("ENS bucket %s still exists after destroy", bucketName)
	}
}
