// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test PolarDb GlobalSecurityIpGroup. >>> Resource test cases, automatically generated.
// Case IP白名单模板用例 11862
func TestAccAliCloudPolarDbGlobalSecurityIpGroup_basic11862(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polar_db_global_security_ip_group.default"
	ra := resourceAttrInit(resourceId, AlicloudPolarDbGlobalSecurityIpGroupMap11862)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PolarDbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePolarDbGlobalSecurityIpGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccpolardb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudPolarDbGlobalSecurityIpGroupBasicDependence11862)
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
					"global_ip_list":       "192.168.0.1",
					"global_ip_group_name": "test_template",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_ip_list":       "192.168.0.1",
						"global_ip_group_name": "test_template",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_ip_list":       "192.168.2.3",
					"global_ip_group_name": "test_name_new",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_ip_list":       "192.168.2.3",
						"global_ip_group_name": "test_name_new",
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

var AlicloudPolarDbGlobalSecurityIpGroupMap11862 = map[string]string{
	"region_id": CHECKSET,
}

func AlicloudPolarDbGlobalSecurityIpGroupBasicDependence11862(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test PolarDb GlobalSecurityIpGroup. <<< Resource test cases, automatically generated.
