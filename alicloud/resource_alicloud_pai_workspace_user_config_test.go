// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test PaiWorkspace UserConfig. >>> Resource test cases, automatically generated.
// Case UserConfig 资源测试01 6821
func TestAccAliCloudPaiWorkspaceUserConfig_basic6821(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_user_config.default"
	ra := resourceAttrInit(resourceId, AliCloudPaiWorkspaceUserConfigMap6821)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceUserConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudPaiWorkspaceUserConfigBasicDependence6821)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"config_key":    "customizePAIAssumedRole",
					"category_name": "DataPrivacyConfig",
					"config_value":  name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_key":    "customizePAIAssumedRole",
						"category_name": "DataPrivacyConfig",
						"config_value":  name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config_value": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_value": name + "update",
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

func TestAccAliCloudPaiWorkspaceUserConfig_basic6821_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_user_config.default"
	ra := resourceAttrInit(resourceId, AliCloudPaiWorkspaceUserConfigMap6821)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceUserConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudPaiWorkspaceUserConfigBasicDependence6821)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"config_key":    "customizePAIAssumedRole",
					"category_name": "DataPrivacyConfig",
					"config_value":  name,
					"scope":         "owner",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_key":    "customizePAIAssumedRole",
						"category_name": "DataPrivacyConfig",
						"config_value":  name,
						"scope":         "owner",
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

func TestAccAliCloudPaiWorkspaceUserConfig_basic6822(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_user_config.default"
	ra := resourceAttrInit(resourceId, AliCloudPaiWorkspaceUserConfigMap6821)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PaiWorkspaceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePaiWorkspaceUserConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudPaiWorkspaceUserConfigBasicDependence6821)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"config_key":    "customizePAIAssumedRole",
					"category_name": "DataPrivacyConfig",
					"config_value":  name,
					"scope":         "subUser",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_key":    "customizePAIAssumedRole",
						"category_name": "DataPrivacyConfig",
						"config_value":  name,
						"scope":         "subUser",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"config_value": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_value": name + "update",
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

var AliCloudPaiWorkspaceUserConfigMap6821 = map[string]string{
	"scope": CHECKSET,
}

func AliCloudPaiWorkspaceUserConfigBasicDependence6821(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test PaiWorkspace UserConfig. <<< Resource test cases, automatically generated.
