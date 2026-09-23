package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Oss BucketCnameToken. >>> Resource test cases, automatically generated.
// Case 自定义域名令牌 8382
func TestAccAliCloudOssBucketCnameToken_basic8382(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_cname_token.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketCnameTokenMap8382)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketCnameToken")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketcnametoken%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketCnameTokenBasicDependence8382)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": "${alicloud_oss_bucket.defaultWWM58I.bucket}",
					"domain": "dinary.top",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": CHECKSET,
						"domain": "dinary.top",
					}),
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

var AlicloudOssBucketCnameTokenMap8382 = map[string]string{
	"token": CHECKSET,
}

func AlicloudOssBucketCnameTokenBasicDependence8382(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "defaultWWM58I" {
  bucket        = var.name
  storage_class = "Standard"
}


`, name)
}

// Test Oss BucketCnameToken. <<< Resource test cases, automatically generated.
