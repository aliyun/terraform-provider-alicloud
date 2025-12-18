// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Kms ValueAddedService. >>> Resource test cases, automatically generated.
// Case 默认密钥增值服务 11636
func TestAccAliCloudKmsValueAddedService_basic11636(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_kms_value_added_service.default"
	ra := resourceAttrInit(resourceId, AlicloudKmsValueAddedServiceMap11636)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &KmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeKmsValueAddedService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacckms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudKmsValueAddedServiceBasicDependence11636)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"value_added_service": "2",
					"payment_type":        "Subscription",
					"period":              "1",
					"renew_period":        "1",
					"renew_status":        "AutoRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"value_added_service": CHECKSET,
						"payment_type":        "Subscription",
						"period":              "1",
						"renew_period":        "1",
						"renew_status":        "AutoRenewal",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renew_period": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renew_period": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"renew_status": "ManualRenewal",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"renew_status": "ManualRenewal",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period", "value_added_service"},
			},
		},
	})
}

var AlicloudKmsValueAddedServiceMap11636 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudKmsValueAddedServiceBasicDependence11636(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Kms ValueAddedService. <<< Resource test cases, automatically generated.
