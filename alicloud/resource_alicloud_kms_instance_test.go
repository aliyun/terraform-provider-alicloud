package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Kms Instance. >>> Resource test cases, automatically generated.
// Case 3873
func TestAccAlicloudKmsInstance_basic3873(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsInstanceMap3873)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmsinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsInstanceBasicDependence3873)
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
					"vpc_num":      "1",
					"key_num":      "1000",
					"secret_num":   "0",
					"spec":         "1000",
					"product_type": "kms_ddi_public_cn",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_num":      "1",
						"key_num":      "1000",
						"secret_num":   "0",
						"spec":         "1000",
						"product_type": "kms_ddi_public_cn",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_num": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_num": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_num": "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_num": "1000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secret_num": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secret_num": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spec": "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec": "1000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_num": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_num": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"key_num": "2000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"key_num": "2000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"secret_num": "1000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"secret_num": "1000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"spec": "2000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"spec": "2000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpc_num":         "1",
					"key_num":         "1000",
					"secret_num":      "0",
					"spec":            "1000",
					"renew_status":    "AutoRenewal",
					"product_type":    "kms_ddi_public_cn",
					"product_version": "3",
					"renew_period":    "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_num":         "1",
						"key_num":         "1000",
						"secret_num":      "0",
						"spec":            "1000",
						"renew_status":    "AutoRenewal",
						"product_type":    "kms_ddi_public_cn",
						"product_version": "3",
						"renew_period":    "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_type", "product_version", "renew_period", "renew_status"},
			},
		},
	})
}

var AlicloudKmsInstanceMap3873 = map[string]string{
	"create_time":   CHECKSET,
	"instance_name": CHECKSET,
	"payment_type":  CHECKSET,
}

func AlicloudKmsInstanceBasicDependence3873(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case 3873  twin
func TestAccAlicloudKmsInstance_basic3873_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsInstanceMap3873)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%skmsinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsInstanceBasicDependence3873)
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
					"vpc_num":         "2",
					"key_num":         "2000",
					"secret_num":      "1000",
					"spec":            "2000",
					"renew_status":    "AutoRenewal",
					"product_type":    "kms_ddi_public_cn",
					"product_version": "3",
					"renew_period":    "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_num":         "2",
						"key_num":         "2000",
						"secret_num":      "1000",
						"spec":            "2000",
						"renew_status":    "AutoRenewal",
						"product_type":    "kms_ddi_public_cn",
						"product_version": "3",
						"renew_period":    "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"product_type", "product_version", "renew_period", "renew_status"},
			},
		},
	})
}

// Test Kms Instance. <<< Resource test cases, automatically generated.
