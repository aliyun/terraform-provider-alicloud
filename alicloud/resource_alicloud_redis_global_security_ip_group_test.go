// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test Redis GlobalSecurityIpGroup. >>> Resource test cases, automatically generated.
// Case 白名单模板模型测试 12846
func TestAccAliCloudRedisGlobalSecurityIpGroup_basic12846(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_global_security_ip_group.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisGlobalSecurityIpGroupMap12846)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisGlobalSecurityIpGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccredis%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRedisGlobalSecurityIpGroupBasicDependence12846)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"global_ip_group_name": "ggn_test_create",
					"global_ip_list":       "192.168.0.1,10.10.10.10,172.16.0.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_ip_group_name": "ggn_test_create",
						"global_ip_list":       "192.168.0.1,10.10.10.10,172.16.0.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_ip_group_name": "ggn_test_update",
					"global_ip_list":       "192.168.0.1,10.10.10.10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_ip_group_name": "ggn_test_update",
						"global_ip_list":       "192.168.0.1,10.10.10.10",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudRedisGlobalSecurityIpGroupMap12846 = map[string]string{}

func AlicloudRedisGlobalSecurityIpGroupBasicDependence12846(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-beijing-h"
}

variable "region_id" {
  default = "cn-beijing"
}


`, name)
}

// Test Redis GlobalSecurityIpGroup. <<< Resource test cases, automatically generated.
