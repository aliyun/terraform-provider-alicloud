package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceAlicloudEnsBucketAcls_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "data.alicloud_ens_bucket_acls.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsBucketAclsMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsBucketAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	bucketName := fmt.Sprintf("tf-testacc-ens-ba-ds-%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			createEnsBucketForLifecycleTest(t, bucketName)
		},
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			func(*terraform.State) error {
				deleteEnsBucketForLifecycleTest(bucketName)
				return nil
			},
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "alicloud_ens_bucket_acl" "default" {
  bucket_name = "%[1]s"
  bucket_acl  = "public-read"
}

data "alicloud_ens_bucket_acls" "default" {
  bucket_name = alicloud_ens_bucket_acl.default.bucket_name
}
`, bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name": bucketName,
					}),
					resource.TestCheckResourceAttr("data.alicloud_ens_bucket_acls.default", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ens_bucket_acls.default", "bucket_acls.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ens_bucket_acls.default", "bucket_acls.0.bucket_name", bucketName),
				),
			},
		},
	})
}

var AlicloudEnsBucketAclsMap = map[string]string{}
