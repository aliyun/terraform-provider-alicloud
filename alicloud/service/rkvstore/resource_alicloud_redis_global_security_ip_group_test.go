// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package rkvstore_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/rkvstore"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Redis GlobalSecurityIpGroup. >>> Resource test cases, automatically generated.
// Case 白名单模板模型测试 12846
func TestAccAliCloudRedisGlobalSecurityIpGroup_basic12846(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_global_security_ip_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudRedisGlobalSecurityIpGroupMap12846)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rkvstore.NewRedisServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRedisGlobalSecurityIpGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccredis%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudRedisGlobalSecurityIpGroupBasicDependence12846)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"global_ip_group_name": name,
					"global_ip_list":       "192.168.1.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_ip_group_name": name,
						"global_ip_list":       "192.168.1.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_ip_group_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_ip_group_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_ip_list": "192.168.1.1,192.168.1.2,192.168.1.3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_ip_list": "192.168.1.1,192.168.1.2,192.168.1.3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_ip_list": "192.168.1.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_ip_list": "192.168.1.1",
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

func TestAccAliCloudRedisGlobalSecurityIpGroup_basic12846_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_global_security_ip_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudRedisGlobalSecurityIpGroupMap12846)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rkvstore.NewRedisServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRedisGlobalSecurityIpGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccredis%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudRedisGlobalSecurityIpGroupBasicDependence12846)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"global_ip_group_name": name,
					"global_ip_list":       "192.168.1.1,192.168.1.2,192.168.1.3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_ip_group_name": name,
						"global_ip_list":       "192.168.1.1,192.168.1.2,192.168.1.3",
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

func TestAccAliCloudRedisGlobalSecurityIpGroup_multi(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_global_security_ip_group.default.5"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudRedisGlobalSecurityIpGroupMap12846)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rkvstore.NewRedisServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRedisGlobalSecurityIpGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccredis%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudRedisGlobalSecurityIpGroupBasicDependence12846)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"count":                "6",
					"global_ip_group_name": name + "${count.index}",
					"global_ip_list":       "192.168.1.1,192.168.1.2,192.168.1.3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_ip_group_name": name + fmt.Sprint(5),
						"global_ip_list":       "192.168.1.1,192.168.1.2,192.168.1.3",
					}),
				),
			},
		},
	})
}

var AliCloudRedisGlobalSecurityIpGroupMap12846 = map[string]string{}

func AliCloudRedisGlobalSecurityIpGroupBasicDependence12846(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}

// Test Redis GlobalSecurityIpGroup. <<< Resource test cases, automatically generated.

// DescribeGlobalSecurityIPGroup is a list query that returns every global
// security IP group, so Read must locate the matching group by
// GlobalSecurityGroupId rather than blindly taking list element [0]. This
// regression test creates several groups with count > 1 and then asserts a
// clean plan, which fails when Read writes the first list entry into every
// instance's state and produces a persistent diff.
func TestAccAliCloudRedisGlobalSecurityIpGroup_count(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_global_security_ip_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudRedisGlobalSecurityIpGroupMap12846)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return rkvstore.NewRedisServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRedisGlobalSecurityIpGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccredis%d", rand)
	testConfig := fmt.Sprintf(`
resource "alicloud_redis_global_security_ip_group" "default" {
  count                = 3
  global_ip_group_name = "%s${count.index}"
  global_ip_list       = "192.168.${count.index}.1,10.10.10.1${count.index}"
}
`, name)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		Providers:    acctest.SDKProviders,
		CheckDestroy: rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_redis_global_security_ip_group.default.0", "global_ip_group_name", name+"0"),
					resource.TestCheckResourceAttr("alicloud_redis_global_security_ip_group.default.0", "global_ip_list", "192.168.0.1,10.10.10.10"),
					resource.TestCheckResourceAttr("alicloud_redis_global_security_ip_group.default.1", "global_ip_group_name", name+"1"),
					resource.TestCheckResourceAttr("alicloud_redis_global_security_ip_group.default.1", "global_ip_list", "192.168.1.1,10.10.10.11"),
					resource.TestCheckResourceAttr("alicloud_redis_global_security_ip_group.default.2", "global_ip_group_name", name+"2"),
					resource.TestCheckResourceAttr("alicloud_redis_global_security_ip_group.default.2", "global_ip_list", "192.168.2.1,10.10.10.12"),
				),
			},
			{
				Config:             testConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
