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

// Test Oss BucketRequestPayment. >>> Resource test cases, automatically generated.
// Case RequestPayment测试 6448
// lintignore: AT001
func TestAccAliCloudOssBucketRequestPayment_basic6448(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_request_payment.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOssBucketRequestPaymentMap6448)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oss.NewOssServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOssBucketRequestPayment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketrequestpayment%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOssBucketRequestPaymentBasicDependence6448)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
					"payer":  "Requester",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": acctest.CHECKSET,
						"payer":  "Requester",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payer": "Requester",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payer": "Requester",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payer": "BucketOwner",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payer": "BucketOwner",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payer":  "Requester",
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payer":  "Requester",
						"bucket": acctest.CHECKSET,
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

var AlicloudOssBucketRequestPaymentMap6448 = map[string]string{
	"payer": "BucketOwner",
}

func AlicloudOssBucketRequestPaymentBasicDependence6448(name string) string {
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

// Case RequestPayment测试 6448  twin
// lintignore: AT001
func TestAccAliCloudOssBucketRequestPayment_basic6448_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_request_payment.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudOssBucketRequestPaymentMap6448)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return oss.NewOssServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeOssBucketRequestPayment")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketrequestpayment%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudOssBucketRequestPaymentBasicDependence6448)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payer":  "BucketOwner",
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payer":  "BucketOwner",
						"bucket": acctest.CHECKSET,
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

// Test Oss BucketRequestPayment. <<< Resource test cases, automatically generated.
