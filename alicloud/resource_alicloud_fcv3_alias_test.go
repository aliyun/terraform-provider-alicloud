package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Fcv3 Alias. >>> Resource test cases, automatically generated.
// Case Alias_Base_online 7304
func TestAccAliCloudFcv3Alias_basic7304(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_alias.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3AliasMap7304)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3Alias")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv3alias%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3AliasBasicDependence7304)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"version_id":    "1",
					"function_name": "flask-3xdg",
					"description":   "create alias",
					"alias_name":    name,
					"additional_version_weight": map[string]interface{}{
						"\"2\"": "0.5",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version_id":    "1",
						"function_name": "flask-3xdg",
						"description":   "create alias",
						"alias_name":    name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"version_id":  "2",
					"description": "update",
					"additional_version_weight": map[string]interface{}{
						"\"1\"": "0.5",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version_id":  "2",
						"description": "update",
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

var AlicloudFcv3AliasMap7304 = map[string]string{
	"alias_name":  CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudFcv3AliasBasicDependence7304(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "function_name" {
  default = "TestAlias_Base"
}


`, name)
}

// Case TestAlias_Base 7214
func TestAccAliCloudFcv3Alias_basic7214(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_fcv3_alias.default"
	ra := resourceAttrInit(resourceId, AlicloudFcv3AliasMap7214)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Fcv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeFcv3Alias")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sfcv3alias%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudFcv3AliasBasicDependence7214)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"version_id":    "1",
					"function_name": "flask-3xdg",
					"description":   "create alias",
					"alias_name":    name,
					"additional_version_weight": map[string]interface{}{
						"\"2\"": "0.5",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version_id":    "1",
						"function_name": "flask-3xdg",
						"description":   "create alias",
						"alias_name":    name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"version_id":  "2",
					"description": "update",
					"additional_version_weight": map[string]interface{}{
						"\"1\"": "0.5",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version_id":  "2",
						"description": "update",
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

var AlicloudFcv3AliasMap7214 = map[string]string{
	"alias_name":  CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudFcv3AliasBasicDependence7214(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

`, name)
}

// Test Fcv3 Alias. <<< Resource test cases, automatically generated.
