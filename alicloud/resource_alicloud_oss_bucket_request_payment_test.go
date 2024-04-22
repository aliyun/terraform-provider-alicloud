package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Oss BucketRequestPayment. >>> Resource test cases, automatically generated.
// Case RequestPayment测试 6448
func TestAccAliCloudOssBucketRequestPayment_basic6448(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_request_payment.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketRequestPaymentMap6448)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketRequestPayment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketrequestpayment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketRequestPaymentBasicDependence6448)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
					"payer":  "Requester",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bucket": CHECKSET,
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
						"bucket": CHECKSET,
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
func TestAccAliCloudOssBucketRequestPayment_basic6448_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_oss_bucket_request_payment.default"
	ra := resourceAttrInit(resourceId, AlicloudOssBucketRequestPaymentMap6448)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &OssServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeOssBucketRequestPayment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sossbucketrequestpayment%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudOssBucketRequestPaymentBasicDependence6448)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payer":  "BucketOwner",
					"bucket": "${alicloud_oss_bucket.CreateBucket.bucket}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payer":  "BucketOwner",
						"bucket": CHECKSET,
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
