package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudGPDBAccount_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_account.default"
	checkoutSupportedRegions(t, true, connectivity.GPDBAccountSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudGPDBAccountMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbAccount")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tftest%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGPDBAccountBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					//"db_instance_id":      "${alicloud_gpdb_instance.default.id}",
					"db_instance_id":      "gp-bp1q1m5z33131uf05",
					"account_name":        name,
					"account_password":    "TFTest123",
					"account_description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_name":        name,
						"account_description": name,
						"db_instance_id":      CHECKSET,
					}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"account_description": name + "update",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"account_description": name + "update",
			//		}),
			//	),
			//},
			{
				Config: testAccConfig(map[string]interface{}{
					"account_password": "TFTest123" + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"account_description": name,
			//		"account_password":    "TFTest123",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"account_description": name,
			//		}),
			//	),
			//},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_password"},
			},
		},
	})
}

var AlicloudGPDBAccountMap0 = map[string]string{}

func AlicloudGPDBAccountBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

`, name)
}
