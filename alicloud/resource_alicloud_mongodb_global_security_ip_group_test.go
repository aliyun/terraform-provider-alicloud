// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Mongodb GlobalSecurityIPGroup. >>> Resource test cases, automatically generated.
// Case 白名单线上接入 10994
func TestAccAliCloudMongodbGlobalSecurityIPGroup_basic10994(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mongodb_global_security_ip_group.default"
	ra := resourceAttrInit(resourceId, AliCloudMongodbGlobalSecurityIPGroupMap10994)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongodbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMongodbGlobalSecurityIPGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmongodb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMongodbGlobalSecurityIPGroupBasicDependence10994)
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
	ra := resourceAttrInit(resourceId, AliCloudMongodbGlobalSecurityIPGroupMap10994)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MongodbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMongodbGlobalSecurityIPGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccmongodb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudMongodbGlobalSecurityIPGroupBasicDependence10994)
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
	"region_id": CHECKSET,
}

func AliCloudMongodbGlobalSecurityIPGroupBasicDependence10994(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

`, name)
}

// Test Mongodb GlobalSecurityIPGroup. <<< Resource test cases, automatically generated.
