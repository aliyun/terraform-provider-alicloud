package alicloud

import (
	"fmt"
	"testing"

	ossclient "github.com/alibabacloud-go/alibabacloud-gateway-oss-util/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAliCloudOssBucketAcl_pointerResponseUnit(t *testing.T) {
	apiName := "GetBucketAcl"
	responseXML := "<AccessControlPolicy><AccessControlList><Grant>private</Grant></AccessControlList></AccessControlPolicy>"
	parsed, err := ossclient.ParseXml(&responseXML, &apiName)
	if err != nil {
		t.Fatalf("parse OSS ACL response: %s", err)
	}

	acl, err := parseOssBucketAclResponse(tea.ToMap(parsed))
	if err != nil {
		t.Fatalf("parse pointer-backed OSS ACL response: %s", err)
	}
	if got := acl["Grant"]; got != "private" {
		t.Fatalf("Grant = %#v, want private", got)
	}
}

func TestAccAliCloudOssBucketAcl_pointerResponse(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_acl.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketAclMap6192)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketacl%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketAclBasicDependence6192)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy:  rac.checkResourceDestroy(),
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
					"acl":    "private",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": CHECKSET,
						"acl":    "private",
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

// Test Oss BucketAcl. >>> Resource test cases, automatically generated.
// Case 测试BucketAcl 6192
// lintignore: AT001
func TestAccAliCloudOssBucketAcl_basic6192(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_acl.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketAclMap6192)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketacl%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketAclBasicDependence6192)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
					"acl":    "private",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": CHECKSET,
						"acl":    "private",
					}),
				),
			},
			// public acl is not allowed
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"acl": "public-read-write",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"acl": "public-read-write",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
					"acl":    "private",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": CHECKSET,
						"acl":    "private",
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

var AlicloudOssBucketAclMap6192 = map[string]string{}

func AlicloudOssBucketAclBasicDependence6192(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
  bucket        = var.name
}


`, name)
}

// Case 测试BucketAcl 6192  twin
// lintignore: AT001
func TestAccAliCloudOssBucketAcl_basic6192_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_acl.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketAclMap6192)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketAcl")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketacl%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketAclBasicDependence6192)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
					"acl":    "private",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": CHECKSET,
						"acl":    "private",
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

// Test Oss BucketAcl. <<< Resource test cases, automatically generated.
