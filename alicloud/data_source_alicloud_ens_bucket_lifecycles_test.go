package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAlicloudEnsBucketLifecyclesDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	bucketName := fmt.Sprintf("tf-testacc-ens-bl-ds-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc("data.alicloud_ens_bucket_lifecycles.default", bucketName, func(name string) string {
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
`, name, bucketName)
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			createEnsBucketForLifecycleTest(t, bucketName)
		},
		IDRefreshName: "data.alicloud_ens_bucket_lifecycles.default",
		Providers:     testAccProviders,
		CheckDestroy: func(*terraform.State) error {
			deleteEnsBucketForLifecycleTest(bucketName)
			return nil
		},
		Steps: []resource.TestStep{
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
		},
	})
}

func TestAccAlicloudEnsBucketLifecyclesDataSource_ruleIdFilter(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	bucketName := fmt.Sprintf("tf-testacc-ens-bl-dsf-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc("data.alicloud_ens_bucket_lifecycles.default", bucketName, func(name string) string {
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
`, name, bucketName)
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			createEnsBucketForLifecycleTest(t, bucketName)
		},
		IDRefreshName: "data.alicloud_ens_bucket_lifecycles.default",
		Providers:     testAccProviders,
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
