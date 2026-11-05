package oss_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// callback and callback_var are Base64-encoded JSON documents, as the OSS
// PutBucketCallbackPolicy API requires. The API performs no canonical
// re-encoding, so the values are stored verbatim; the per-step attribute checks
// and ImportStateVerify are the primary correctness assertions.

const (
	// {"callbackUrl":"test.com/test.php", "callbackBody":"bucket=${bucket}&object=${object}&a=${x:a}&b=${x:b}"}
	cbCallbackTest = "eyJjYWxsYmFja1VybCI6InRlc3QuY29tL3Rlc3QucGhwIiwgImNhbGxiYWNrQm9keSI6ImJ1Y2tldD0ke2J1Y2tldH0mb2JqZWN0PSR7b2JqZWN0fSZhPSR7eDphfSZiPSR7eDpifSJ9"
	// {"callbackUrl":"test2.com/test.php", "callbackBody":"bucket=${bucket}&object=${object}&a=${x:a}&b=${x:b}"}
	cbCallbackTest2 = "eyJjYWxsYmFja1VybCI6InRlc3QyLmNvbS90ZXN0LnBocCIsICJjYWxsYmFja0JvZHkiOiJidWNrZXQ9JHtidWNrZXR9Jm9iamVjdD0ke29iamVjdH0mYT0ke3g6YX0mYj0ke3g6Yn0ifQ=="
	// {"callbackUrl":"test1.com/test.php", "callbackBody":"bucket=${bucket}&object=${object}&a=${x:a}&b=${x:b}"}
	cbCallbackTest1 = "eyJjYWxsYmFja1VybCI6InRlc3QxLmNvbS90ZXN0LnBocCIsICJjYWxsYmFja0JvZHkiOiJidWNrZXQ9JHtidWNrZXR9Jm9iamVjdD0ke29iamVjdH0mYT0ke3g6YX0mYj0ke3g6Yn0ifQ=="
	// {"x:a":"a","x:b":"b"}
	cbVarB = "eyJ4OmEiOiJhIiwgIng6YiI6ImIifQ=="
	// {"x:a":"a","x:b":"d"}
	cbVarD = "eyJ4OmEiOiJhIiwgIng6YiI6ImQifQ=="
	// {"x:a":"a","x:b":"c"}
	cbVarC = "eyJ4OmEiOiJhIiwgIng6YiI6ImMifQ=="
)

func TestAccAliCloudOssBucketCallbackPolicy_basic(t *testing.T) {
	resourceId := "alicloud_oss_bucket_callback_policy.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	region := string(connectivity.Hangzhou)
	// OSS bucket names must be lowercase, 3-63 chars and globally unique.
	name := fmt.Sprintf("tf-testacc-osscbpol-%s-%d", region, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckOssBucketCallbackPolicy(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckOssBucketCallbackPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudOssBucketCallbackPolicyConfig(name, region, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "bucket", name),
					resource.TestCheckResourceAttr(resourceId, "id", name),
					resource.TestCheckResourceAttr(resourceId, "policy_item.#", "3"),
					resource.TestCheckResourceAttr(resourceId, "policy_item.0.policy_name", "name1"),
					resource.TestCheckResourceAttr(resourceId, "policy_item.0.callback", cbCallbackTest),
					resource.TestCheckResourceAttr(resourceId, "policy_item.0.callback_var", cbVarB),
					resource.TestCheckResourceAttr(resourceId, "policy_item.1.policy_name", "name2"),
					resource.TestCheckResourceAttr(resourceId, "policy_item.1.callback", cbCallbackTest2),
					resource.TestCheckResourceAttr(resourceId, "policy_item.1.callback_var", cbVarD),
					resource.TestCheckResourceAttr(resourceId, "policy_item.2.policy_name", "name3"),
					resource.TestCheckResourceAttr(resourceId, "policy_item.2.callback", cbCallbackTest2),
					resource.TestCheckResourceAttr(resourceId, "policy_item.2.callback_var", cbVarD),
				),
			},
			{
				// Update to a single policy item: the Put replaces the whole
				// policy, so the prior three items must all be gone.
				Config: testAccAliCloudOssBucketCallbackPolicyConfig(name, region, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "policy_item.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "policy_item.0.policy_name", "name4"),
					resource.TestCheckResourceAttr(resourceId, "policy_item.0.callback", cbCallbackTest1),
					resource.TestCheckResourceAttr(resourceId, "policy_item.0.callback_var", cbVarC),
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

func testAccPreCheckOssBucketCallbackPolicy(t *testing.T) {
	acctest.PreCheck(t)
	// Callback policy is served from cn-hangzhou; pin the region so the case
	// does not depend on the caller's configured region.
	os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
}

func testAccCheckOssBucketCallbackPolicyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_oss_bucket_callback_policy" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		var response map[string]interface{}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err := retry.Retry(1*time.Minute, func() *retry.RetryError {
			var requestErr error
			response, requestErr = client.Do("Oss", helper.XmlParam("GET", "2019-05-17", "GetBucketCallbackPolicy", "/?policy&comp=callback"),
				map[string]*string{}, nil, nil, map[string]*string{"bucket": helper.StringPointer(rs.Primary.ID)}, false)
			if requestErr != nil {
				if helper.NeedRetry(requestErr) {
					wait()
					return retry.RetryableError(requestErr)
				}
				return retry.NonRetryableError(requestErr)
			}
			return nil
		})
		if err != nil {
			// A gone bucket (the test stack tears the bucket down too) or a gone
			// callback policy both count as deleted.
			if errs.IsExpectedErrors(err, []string{"NoSuchBucket", "BucketCallbackPolicyNotExist"}) || errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing bucket callback policy %s failed: %s", rs.Primary.ID, err)
		}
		if response != nil {
			return fmt.Errorf("bucket callback policy %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccAliCloudOssBucketCallbackPolicyConfig(name, region string, items int) string {
	var block string
	switch items {
	case 3:
		block = fmt.Sprintf(`
  policy_item {
    policy_name  = "name1"
    callback     = %q
    callback_var = %q
  }
  policy_item {
    policy_name  = "name2"
    callback     = %q
    callback_var = %q
  }
  policy_item {
    policy_name  = "name3"
    callback     = %q
    callback_var = %q
  }`,
			cbCallbackTest, cbVarB,
			cbCallbackTest2, cbVarD,
			cbCallbackTest2, cbVarD,
		)
	case 1:
		block = fmt.Sprintf(`
  policy_item {
    policy_name  = "name4"
    callback     = %q
    callback_var = %q
  }`,
			cbCallbackTest1, cbVarC,
		)
	}
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "region" {
  default = "%s"
}

resource "alicloud_oss_bucket" "default" {
  bucket        = var.name
  storage_class = "Standard"
}

resource "alicloud_oss_bucket_callback_policy" "default" {
  bucket = alicloud_oss_bucket.default.bucket
%s
}
`, name, region, block)
}
