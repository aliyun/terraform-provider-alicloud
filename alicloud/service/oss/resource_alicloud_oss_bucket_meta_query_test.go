package oss_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/oss"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Oss BucketMetaQuery. >>> Resource test cases, automatically generated.
// Case MetaQuery测试 6694
func TestAccAliCloudOssBucketMetaQuery_basic6694(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_meta_query.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOssBucketMetaQueryMap6694)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oss.NewOssServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOssBucketMetaQuery")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketmetaquery%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOssBucketMetaQueryBasicDependence6694)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status"},
			},
		},
	})
}

var AlicloudOssBucketMetaQueryMap6694 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

func AlicloudOssBucketMetaQueryBasicDependence6694(name string) string {
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

// Case MetaQuery测试 6694  raw
func TestAccAliCloudOssBucketMetaQuery_basic6694_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_meta_query.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOssBucketMetaQueryMap6694)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oss.NewOssServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOssBucketMetaQuery")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketmetaquery%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOssBucketMetaQueryBasicDependence6694)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": acctest.CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"status"},
			},
		},
	})
}

// Test Oss BucketMetaQuery. <<< Resource test cases, automatically generated.
