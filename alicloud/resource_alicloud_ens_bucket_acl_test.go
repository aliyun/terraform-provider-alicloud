package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudEnsBucketAcl_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ens_bucket_acl.default"
	ra := resourceAttrInit(resourceId, AlicloudEnsBucketAclMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EnsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEnsBucketAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	bucketName := fmt.Sprintf("tf-testacc-ens-ba-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, bucketName, AlicloudEnsBucketAclBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			createEnsBucketForLifecycleTest(t, bucketName)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			func(*terraform.State) error {
				deleteEnsBucketForLifecycleTest(bucketName)
				return nil
			},
		),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name": bucketName,
					"bucket_acl":  "public-read",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name": bucketName,
						"bucket_acl":  "public-read",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name": bucketName,
					"bucket_acl":  "private",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name": bucketName,
						"bucket_acl":  "private",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket_name": bucketName,
					"bucket_acl":  "public-read-write",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket_name": bucketName,
						"bucket_acl":  "public-read-write",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudEnsBucketAclMap = map[string]string{}

func AlicloudEnsBucketAclBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}
