// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package mongodb_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/mongodb"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Mongodb GlobalSecurityIPGroup. >>> Resource test cases, automatically generated.
// Case 白名单线上接入 10994
func TestAccAliCloudMongodbGlobalSecurityIPGroup_basic10994(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_global_security_ip_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudMongodbGlobalSecurityIPGroupMap10994)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mongodb.NewMongodbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMongodbGlobalSecurityIPGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmongodb%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudMongodbGlobalSecurityIPGroupBasicDependence10994)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"global_ig_name":          name,
					"global_security_ip_list": "192.168.1.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_ig_name":          name,
						"global_security_ip_list": "192.168.1.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_ig_name":          name + "update",
					"global_security_ip_list": "192.168.1.1,192.168.1.2,192.168.1.3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_ig_name":          name + "update",
						"global_security_ip_list": "192.168.1.1,192.168.1.2,192.168.1.3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_ig_name":          name,
					"global_security_ip_list": "192.168.1.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_ig_name":          name,
						"global_security_ip_list": "192.168.1.1",
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

func TestAccAliCloudMongodbGlobalSecurityIPGroup_basic10994_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_global_security_ip_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudMongodbGlobalSecurityIPGroupMap10994)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mongodb.NewMongodbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMongodbGlobalSecurityIPGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmongodb%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudMongodbGlobalSecurityIPGroupBasicDependence10994)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"global_ig_name":          name,
					"global_security_ip_list": "192.168.1.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_ig_name":          name,
						"global_security_ip_list": "192.168.1.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_ig_name":          name + "update",
					"global_security_ip_list": "192.168.1.1,192.168.1.2,192.168.1.3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_ig_name":          name + "update",
						"global_security_ip_list": "192.168.1.1,192.168.1.2,192.168.1.3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"global_ig_name":          name,
					"global_security_ip_list": "192.168.1.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"global_ig_name":          name,
						"global_security_ip_list": "192.168.1.1",
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

var AliCloudMongodbGlobalSecurityIPGroupMap10994 = map[string]string{
	"region_id": acctest.CHECKSET,
}

func AliCloudMongodbGlobalSecurityIPGroupBasicDependence10994(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

`, name)
}

// Test Mongodb GlobalSecurityIPGroup. <<< Resource test cases, automatically generated.

// DescribeGlobalSecurityIPGroup is a list query that returns every global
// security IP group, so Read must locate the matching group by
// GlobalSecurityGroupId rather than blindly taking list element [0]. This
// regression test creates several groups with count > 1 and then asserts a
// clean plan, which fails when Read writes the first list entry into every
// instance's state and produces a persistent diff.
func TestAccAliCloudMongodbGlobalSecurityIPGroup_count(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_global_security_ip_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudMongodbGlobalSecurityIPGroupMap10994)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mongodb.NewMongodbServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeMongodbGlobalSecurityIPGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmongodb%d", rand)
	testConfig := fmt.Sprintf(`
resource "alicloud_mongodb_global_security_ip_group" "default" {
  count                   = 3
  global_ig_name          = "%s${count.index}"
  global_security_ip_list = "192.168.${count.index}.1"
}
`, name)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_mongodb_global_security_ip_group.default.0", "global_ig_name", name+"0"),
					resource.TestCheckResourceAttr("alicloud_mongodb_global_security_ip_group.default.0", "global_security_ip_list", "192.168.0.1"),
					resource.TestCheckResourceAttr("alicloud_mongodb_global_security_ip_group.default.1", "global_ig_name", name+"1"),
					resource.TestCheckResourceAttr("alicloud_mongodb_global_security_ip_group.default.1", "global_security_ip_list", "192.168.1.1"),
					resource.TestCheckResourceAttr("alicloud_mongodb_global_security_ip_group.default.2", "global_ig_name", name+"2"),
					resource.TestCheckResourceAttr("alicloud_mongodb_global_security_ip_group.default.2", "global_security_ip_list", "192.168.2.1"),
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
