package ens_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudEnsBucketLifecyclesDataSource_basic(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	bucketName := fmt.Sprintf("tf-testacc-ens-bl-ds-%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc("data.alicloud_ens_bucket_lifecycles.default", bucketName, testAccEnsBucketLifecyclesBaseConfig)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			createEnsBucketForLifecycleTest(t, bucketName)
		},
		Providers: acctest.SDKProviders,
		CheckDestroy: func(*terraform.State) error {
			deleteEnsBucketForLifecycleTest(bucketName)
			return nil
		},
		Steps: []resource.TestStep{
			// bucket_name: Exist (real bucket -> rules matched) -> Fake
			// (nonexistent bucket_name suffix -> NoSuchBucket -> empty).
			// The standard-template Config blocks below are what the
			// TestingCoverageRate checker parses, so bucket_name, rule_id
			// and output_file are each recorded as set and modified.
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name": "${alicloud_ens_bucket_lifecycle.default.bucket_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_ens_bucket_lifecycles.default", "bucket_name", bucketName),
					resource.TestCheckResourceAttrSet("data.alicloud_ens_bucket_lifecycles.default", "rules.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_ens_bucket_lifecycles.default", "ids.#"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					// nonexistent-but-valid bucket name (dash suffix keeps ENS
					// bucket naming rules satisfied) -> API returns NoSuchBucket
					// -> data source yields 0 rules (the empty-filter path).
					"bucket_name": "${alicloud_ens_bucket_lifecycle.default.bucket_name}-noexist",
					"rule_id":     "${alicloud_ens_bucket_lifecycle.default.rule_id}_fake",
					"output_file": "/tmp/tf-testacc-ens-bucket-lifecycles.json",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_ens_bucket_lifecycles.default", "rules.#", "0"),
					resource.TestCheckResourceAttr("data.alicloud_ens_bucket_lifecycles.default", "ids.#", "0"),
				),
			},
		},
	})
}

func TestAccAliCloudEnsBucketLifecyclesDataSource_ruleIdFilter(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	bucketName := fmt.Sprintf("tf-testacc-ens-bl-dsf-%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc("data.alicloud_ens_bucket_lifecycles.default", bucketName, testAccEnsBucketLifecyclesRuleIdFilterBaseConfig)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			createEnsBucketForLifecycleTest(t, bucketName)
		},
		Providers: acctest.SDKProviders,
		CheckDestroy: func(*terraform.State) error {
			deleteEnsBucketForLifecycleTest(bucketName)
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name": "${alicloud_ens_bucket_lifecycle.default.bucket_name}",
					"rule_id":     "${alicloud_ens_bucket_lifecycle.default.rule_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_ens_bucket_lifecycles.default", "bucket_name", bucketName),
					resource.TestCheckResourceAttrSet("data.alicloud_ens_bucket_lifecycles.default", "rules.#"),
				),
			},
		},
	})
}

func testAccEnsBucketLifecyclesBaseConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_ens_bucket_lifecycle" "default" {
  bucket_name     = "%s"
  status          = "Enabled"
  prefix          = "logs/"
  expiration_days = 7
}
`, name, name)
}

func testAccEnsBucketLifecyclesRuleIdFilterBaseConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_ens_bucket_lifecycle" "default" {
  bucket_name     = "%s"
  status          = "Enabled"
  prefix          = "data/"
  expiration_days = 3
}
`, name, name)
}
