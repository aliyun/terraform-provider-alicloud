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

// Test Oss BucketStyle. >>> Resource test cases, automatically generated.
// Case Style指定Category 6688
func TestAccAliCloudOssBucketStyle_basic6688(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_style.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOssBucketStyleMap6688)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oss.NewOssServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOssBucketStyle")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccoss%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOssBucketStyleBasicDependence6688)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":     "${alicloud_oss_bucket.CreateBucket.id}",
					"style_name": "style-771",
					"content":    "image/resize,p_75,w_75",
					"category":   "document",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":     acctest.CHECKSET,
						"style_name": acctest.CHECKSET,
						"content":    "image/resize,p_75,w_75",
						"category":   "document",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content":  "image/resize,p_75,w_70",
					"category": "video",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content":  "image/resize,p_75,w_70",
						"category": "video",
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

var AlicloudOssBucketStyleMap6688 = map[string]string{
	"create_time": acctest.CHECKSET,
}

func AlicloudOssBucketStyleBasicDependence6688(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
}


`, name)
}

// Case BucketStyle测试 6687
func TestAccAliCloudOssBucketStyle_basic6687(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_style.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOssBucketStyleMap6687)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oss.NewOssServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOssBucketStyle")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccoss%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOssBucketStyleBasicDependence6687)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket":     "${alicloud_oss_bucket.CreateBucket.id}",
					"style_name": "style-140",
					"content":    "image/resize,p_75,w_75",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket":     acctest.CHECKSET,
						"style_name": acctest.CHECKSET,
						"content":    "image/resize,p_75,w_75",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content":  "image/resize,p_75,w_70",
					"category": "image",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content":  "image/resize,p_75,w_70",
						"category": "image",
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

var AlicloudOssBucketStyleMap6687 = map[string]string{
	"create_time": acctest.CHECKSET,
}

func AlicloudOssBucketStyleBasicDependence6687(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_oss_bucket" "CreateBucket" {
  storage_class = "Standard"
}


`, name)
}

// Test Oss BucketStyle. <<< Resource test cases, automatically generated.
