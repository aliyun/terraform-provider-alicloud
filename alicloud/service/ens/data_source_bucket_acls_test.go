package ens_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourceAlicloudEnsBucketAcls_basic(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	bucketName := fmt.Sprintf("tf-testacc-ens-ba-ds-%d", r)

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
				Config: testAccDataSourceAlicloudEnsBucketAclsConfig(bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_ens_bucket_acls.default", "id", bucketName),
					resource.TestCheckResourceAttr("data.alicloud_ens_bucket_acls.default", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ens_bucket_acls.default", "bucket_acls.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ens_bucket_acls.default", "bucket_acls.0.bucket_name", bucketName),
					resource.TestCheckResourceAttr("data.alicloud_ens_bucket_acls.default", "bucket_acls.0.bucket_acl", "public-read"),
				),
			},
		},
	})
}

func testAccDataSourceAlicloudEnsBucketAclsConfig(bucketName string) string {
	return fmt.Sprintf(`
resource "alicloud_ens_bucket_acl" "default" {
  bucket_name = "%[1]s"
  bucket_acl  = "public-read"
}

data "alicloud_ens_bucket_acls" "default" {
  bucket_name = alicloud_ens_bucket_acl.default.bucket_name
}
`, bucketName)
}
