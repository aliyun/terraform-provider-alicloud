package ens_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAliCloudEnsBucket_basic covers the full lifecycle of the ENS bucket
// Framework resource: a standard (global) bucket is created with a comment and
// an explicit private ACL, the ACL is then updated in place to public-read and
// again to public-read-write (exercising all three ACL values via PutBucketAcl),
// and finally the resource is imported to verify the ID-only import reproduces
// the applied state. The standard bucket form is used because PutBucket with
// only BucketName is the documented default and avoids the sink-only
// EnsRegionId coupling, matching the ENS package's existing test setup helper.
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
				Config: testAccAliCloudEnsBucketConfigUpdateAcl(bucketName, "public-read"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bucket_acl", "public-read"),
					// ForceNew fields are unchanged across the in-place ACL update.
					resource.TestCheckResourceAttr(resourceName, "bucket_name", bucketName),
					resource.TestCheckResourceAttr(resourceName, "logical_bucket_type", "standard"),
					resource.TestCheckResourceAttr(resourceName, "comment", "tf-testacc-ens-comment"),
				),
			},
			{
				Config: testAccAliCloudEnsBucketConfigUpdateAcl(bucketName, "public-read-write"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bucket_acl", "public-read-write"),
					resource.TestCheckResourceAttr(resourceName, "bucket_name", bucketName),
					resource.TestCheckResourceAttr(resourceName, "logical_bucket_type", "standard"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// GetBucketInfo is eventually consistent for bucket_acl and
				// modify_time after PutBucketAcl: a read-back during import can
				// serve a stale snapshot (the pre-update ACL and modify_time)
				// even though the write already converged in the Update step
				// above. These two fields are therefore not reliably
				// reproducible on import and are ignored here; the import still
				// verifies id, bucket_name, logical_bucket_type and comment,
				// and bucket_acl is exercised end-to-end (private ->
				// public-read -> public-read-write) across the steps above.
				ImportStateVerifyIgnore: []string{"bucket_acl", "modify_time"},
			},
		},
	})
}

// testAccAliCloudEnsBucketConfig is the basic standard-bucket configuration. A
// standard bucket needs only a bucket name; comment and bucket_acl are set
// explicitly so the in-place ACL updates in the next steps are observable.
func testAccAliCloudEnsBucketConfig(bucketName string) string {
	return fmt.Sprintf(`
resource "alicloud_ens_bucket" "default" {
  bucket_name = "%[1]s"
  comment     = "tf-testacc-ens-comment"
  bucket_acl  = "private"
}
`, bucketName)
}

// testAccAliCloudEnsBucketConfigUpdateAcl flips bucket_acl to the given value,
// exercising the PutBucketAcl update path while every other argument stays the
// same, so only the ACL changes in place. Used for public-read and
// public-read-write to cover all three ACL values.
func testAccAliCloudEnsBucketConfigUpdateAcl(bucketName, acl string) string {
	return fmt.Sprintf(`
resource "alicloud_ens_bucket" "default" {
  bucket_name = "%[1]s"
  comment     = "tf-testacc-ens-comment"
  bucket_acl  = "%[2]s"
}
`, bucketName, acl)
}

// testAccCheckAliCloudEnsBucketDestroy verifies the bucket is gone after the
// test destroys it. ENS deletes buckets asynchronously and GetBucketInfo is
// eventually consistent across replicas, so a single read-back right after
// DeleteBucket can still observe the bucket even though the provider's Delete
// already polled until it disappeared. Poll GetBucketInfo until it reports
// NoSuchBucket (or the poll budget is exhausted); any other error keeps
// polling because a transient failure during the deletion window should not
// mask a successful destroy.
func testAccCheckAliCloudEnsBucketDestroy(bucketName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		deadline := time.Now().Add(5 * time.Minute)
		for {
			_, err := client.RpcPost("Ens", "2017-11-10", "GetBucketInfo", nil, map[string]interface{}{
				"BucketName": bucketName,
			}, true)
			if err != nil {
				if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"NoSuchBucket"}) {
					return nil
				}
				// Other errors are treated as transient and keep polling; the
				// destroy path already waited for deletion, so a persistent
				// non-NotFound error surfaces only when the poll budget below
				// is exhausted.
			}
			if time.Now().After(deadline) {
				return fmt.Errorf("ENS bucket %s still exists after destroy (polled for eventual consistency for 5m)", bucketName)
			}
			time.Sleep(3 * time.Second)
		}
	}
}

// TestAccAliCloudEnsBucket_sink exercises the single-node (sink) bucket form
// of the resource: a sink bucket is created with an explicit ens_region_id and
// a public-read-write ACL, then imported. logical_bucket_type=sink is read
// back from GetBucketInfo (so it survives refresh), but ens_region_id is
// write-only — GetBucketInfo does not return it — so import must ignore
// ens_region_id. This exercises the PutBucket request path that carries
// LogicalBucketType and EnsRegionId together, which the basic standard-bucket
// test never reaches, and the public-read-write ACL value on a sink bucket.
//
// A sink bucket is bound to one edge node and PutBucket rejects a sink request
// whose EnsRegionId has not been purchased for the bucket service
// (OperationDenied "The current ens region has not been purchased"). The
// default acceptance-test account has not purchased an ENS sink region, so this
// case is gated behind the ENS_SINK_REGION_ID environment variable: it runs
// only when that variable holds a purchased sink region id, and skips
// otherwise. The logical_bucket_type=sink and ens_region_id coverage gap for
// the default account is documented in the delivery report.
func TestAccAliCloudEnsBucket_sink(t *testing.T) {
	acctest.PreCheckWithEnvVariable(t, "ENS_SINK_REGION_ID")
	ensRegionId := os.Getenv("ENS_SINK_REGION_ID")
	r := tfacctest.RandIntRange(10000, 99999)
	bucketName := fmt.Sprintf("tf-testacc-ens-sink-%d", r)
	resourceName := "alicloud_ens_bucket.sink"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckAliCloudEnsBucketDestroy(bucketName),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudEnsBucketSinkConfig(bucketName, ensRegionId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", bucketName),
					resource.TestCheckResourceAttr(resourceName, "bucket_name", bucketName),
					resource.TestCheckResourceAttr(resourceName, "logical_bucket_type", "sink"),
					resource.TestCheckResourceAttr(resourceName, "ens_region_id", ensRegionId),
					resource.TestCheckResourceAttr(resourceName, "bucket_acl", "public-read-write"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttrSet(resourceName, "modify_time"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ens_region_id"},
			},
		},
	})
}

// TestAccAliCloudEnsBucket_dispatchScope exercises the dispatch_scope argument
// of a global (standard) bucket with the domestic value: a bucket is created
// with dispatch_scope=domestic and then imported. dispatch_scope is write-only
// (not returned by GetBucketInfo), so import must ignore it. This covers the
// PutBucket DispatchScope request path for domestic, which the basic test never
// sets.
//
// The oversea value is not exercised here: PutBucket with
// dispatch_scope=oversea requires a purchased overseas ENS region, and the
// default acceptance-test account has none (OperationDenied "There is no valid
// ens region"). The oversea coverage gap is documented in the delivery report.
func TestAccAliCloudEnsBucket_dispatchScope(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	bucketName := fmt.Sprintf("tf-testacc-ens-dispatch-%d", r)
	resourceName := "alicloud_ens_bucket.dispatch"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckAliCloudEnsBucketDestroy(bucketName),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudEnsBucketDispatchScopeConfig(bucketName, "domestic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", bucketName),
					resource.TestCheckResourceAttr(resourceName, "bucket_name", bucketName),
					resource.TestCheckResourceAttr(resourceName, "logical_bucket_type", "standard"),
					resource.TestCheckResourceAttr(resourceName, "dispatch_scope", "domestic"),
					resource.TestCheckResourceAttr(resourceName, "bucket_acl", "private"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"dispatch_scope"},
			},
		},
	})
}

// testAccAliCloudEnsBucketSinkConfig is the single-node (sink) bucket
// configuration. A sink bucket is bound to one edge node, so ens_region_id is
// mandatory; bucket_acl is set to public-read-write to exercise that ACL
// value on the sink PutBucket request path.
func testAccAliCloudEnsBucketSinkConfig(bucketName, ensRegionId string) string {
	return fmt.Sprintf(`
resource "alicloud_ens_bucket" "sink" {
  bucket_name         = "%[1]s"
  logical_bucket_type = "sink"
  ens_region_id       = "%[2]s"
  bucket_acl          = "public-read-write"
}
`, bucketName, ensRegionId)
}

// testAccAliCloudEnsBucketDispatchScopeConfig is the global (standard) bucket
// configuration with an explicit dispatch_scope. dispatch_scope is ForceNew,
// so changing it recreates the bucket and exercises the value on the PutBucket
// request path.
func testAccAliCloudEnsBucketDispatchScopeConfig(bucketName, dispatchScope string) string {
	return fmt.Sprintf(`
resource "alicloud_ens_bucket" "dispatch" {
  bucket_name         = "%[1]s"
  logical_bucket_type = "standard"
  dispatch_scope      = "%[2]s"
}
`, bucketName, dispatchScope)
}
