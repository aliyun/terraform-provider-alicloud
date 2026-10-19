package ens_test

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// alicloud_ens_bucket_acl manages the ACL of an EXISTING ENS bucket, so the
// test must pre-create the bucket (PutBucket) and tear it down (DeleteBucket)
// itself. The provider's Meta is not populated at PreCheck time (it is set only
// after the first test step applies a provider config), so these helpers build a
// client straight from the environment, mirroring sharedClientForRegion but
// using only exported symbols. They poll with a plain time loop and deliberately
// avoid the SDK v2 helper/resource package, which would re-register the -sweep
// flag and collide with plugin-testing once a test binary links both.

func ensTestClient(t *testing.T) *connectivity.AliyunClient {
	t.Helper()
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-beijing"
	}
	accessKey := os.Getenv("ALICLOUD_ACCESS_KEY")
	secretKey := os.Getenv("ALICLOUD_SECRET_KEY")
	if accessKey == "" || secretKey == "" {
		t.Fatal("ALICLOUD_ACCESS_KEY and ALICLOUD_SECRET_KEY must be set for acceptance tests")
	}

	conf := connectivity.Config{
		Region:      connectivity.Region(region),
		RegionId:    region,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		Protocol:    "HTTPS",
		Endpoints:   &sync.Map{},
		SignVersion: &sync.Map{},
	}
	if token := os.Getenv("ALICLOUD_SECURITY_TOKEN"); token != "" {
		conf.SecurityToken = token
	}
	if accountId := firstNonEmpty(os.Getenv("ALICLOUD_ACCOUNT_ID"), os.Getenv("ALIBABA_CLOUD_ACCOUNT_ID")); accountId != "" {
		conf.AccountId = accountId
	}
	if accountType := firstNonEmpty(os.Getenv("ALICLOUD_ACCOUNT_TYPE"), os.Getenv("ALIBABA_CLOUD_ACCOUNT_TYPE")); accountType != "" {
		conf.AccountType = accountType
	}

	credCfg := new(credentials.Config).SetType("access_key").SetAccessKeyId(accessKey).SetAccessKeySecret(secretKey)
	if token := strings.TrimSpace(os.Getenv("ALICLOUD_SECURITY_TOKEN")); token != "" {
		credCfg.SetType("sts").SetSecurityToken(token)
	}
	cred, err := credentials.NewCredential(credCfg)
	if err != nil {
		t.Fatalf("build ENS test client: credentials: %v", err)
	}
	conf.Credential = cred

	client, err := conf.Client()
	if err != nil {
		t.Fatalf("build ENS test client for region %s: %v", region, err)
	}
	return client
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func createEnsBucketForTest(t *testing.T, bucketName string) {
	t.Helper()
	client := ensTestClient(t)

	request := map[string]interface{}{"BucketName": bucketName}
	if err := ensRetry(2*time.Minute, func() error {
		_, err := client.RpcPost("Ens", "2017-11-10", "PutBucket", nil, request, true)
		return err
	}); err != nil {
		t.Fatalf("PutBucket %s: %v", bucketName, err)
	}

	// Poll GetBucketInfo until the bucket is queryable; during the creation
	// window the service may report NoSuchBucket.
	query := map[string]interface{}{"BucketName": bucketName}
	if err := ensRetry(3*time.Minute, func() error {
		_, e := client.RpcPost("Ens", "2017-11-10", "GetBucketInfo", query, nil, true)
		if e == nil {
			return nil
		}
		if isEnsNotFound(e) || isTransient(e) {
			return errRetryable(e)
		}
		return e
	}); err != nil {
		t.Fatalf("GetBucketInfo ready check for %s: %v", bucketName, err)
	}
}

func deleteEnsBucketForTest(bucketName string) {
	client, err := ensTestClientSilent()
	if err != nil {
		return
	}
	query := map[string]interface{}{"BucketName": bucketName}
	_ = ensRetry(1*time.Minute, func() error {
		_, e := client.RpcPost("Ens", "2017-11-10", "DeleteBucket", query, nil, true)
		if e == nil || isEnsNotFound(e) {
			return nil
		}
		if isTransient(e) {
			return errRetryable(e)
		}
		return e
	})
}

func ensTestClientSilent() (*connectivity.AliyunClient, error) {
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-beijing"
	}
	accessKey := os.Getenv("ALICLOUD_ACCESS_KEY")
	secretKey := os.Getenv("ALICLOUD_SECRET_KEY")
	if accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("missing credentials")
	}
	conf := connectivity.Config{
		Region:      connectivity.Region(region),
		RegionId:    region,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		Protocol:    "HTTPS",
		Endpoints:   &sync.Map{},
		SignVersion: &sync.Map{},
	}
	if token := os.Getenv("ALICLOUD_SECURITY_TOKEN"); token != "" {
		conf.SecurityToken = token
	}
	credCfg := new(credentials.Config).SetType("access_key").SetAccessKeyId(accessKey).SetAccessKeySecret(secretKey)
	cred, err := credentials.NewCredential(credCfg)
	if err != nil {
		return nil, err
	}
	conf.Credential = cred
	return conf.Client()
}

type retryableError struct{ err error }

func (e *retryableError) Error() string { return e.err.Error() }
func (e *retryableError) Unwrap() error { return e.err }

func errRetryable(err error) error { return &retryableError{err: err} }

func isEnsNotFound(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	for _, frag := range []string{"NoSuchBucket", "InvalidBucketName", "AccessDenied"} {
		if strings.Contains(msg, frag) {
			return true
		}
	}
	return false
}

func isTransient(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	for _, frag := range []string{"Throttling", "RequestTimeout", "ServiceUnavailable", "InternalError", "Busy", "try again", "try it later"} {
		if strings.Contains(msg, frag) {
			return true
		}
	}
	return false
}

func ensRetry(timeout time.Duration, fn func() error) error {
	deadline := time.Now().Add(timeout)
	interval := 3 * time.Second
	step := 5 * time.Second
	var last error
	for {
		err := fn()
		if err == nil {
			return nil
		}
		last = err
		var re *retryableError
		if errors.As(err, &re) || isTransient(err) || isEnsNotFound(err) {
			if time.Now().After(deadline) {
				return last
			}
			time.Sleep(interval)
			interval += step
			continue
		}
		return err
	}
}

func TestAccAliCloudEnsBucketAcl_basic(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	bucketName := fmt.Sprintf("tf-testacc-ens-ba-%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			createEnsBucketForTest(t, bucketName)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy: resource.ComposeTestCheckFunc(
			func(*terraform.State) error {
				deleteEnsBucketForTest(bucketName)
				return nil
			},
		),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudEnsBucketAclConfig(bucketName, "public-read"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_ens_bucket_acl.default", "id", bucketName),
					resource.TestCheckResourceAttr("alicloud_ens_bucket_acl.default", "bucket_name", bucketName),
					resource.TestCheckResourceAttr("alicloud_ens_bucket_acl.default", "bucket_acl", "public-read"),
				),
			},
			{
				Config: testAccAliCloudEnsBucketAclConfig(bucketName, "private"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_ens_bucket_acl.default", "bucket_acl", "private"),
				),
			},
			{
				Config: testAccAliCloudEnsBucketAclConfig(bucketName, "public-read-write"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_ens_bucket_acl.default", "bucket_acl", "public-read-write"),
				),
			},
			{
				ResourceName:      "alicloud_ens_bucket_acl.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAliCloudEnsBucketAclConfig(bucketName, acl string) string {
	return fmt.Sprintf(`
resource "alicloud_ens_bucket_acl" "default" {
  bucket_name = "%[1]s"
  bucket_acl  = "%[2]s"
}
`, bucketName, acl)
}
