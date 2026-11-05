// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package paiworkspace_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/paiworkspace"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test PaiWorkspace UserConfig. >>> Resource test cases, automatically generated.
// Case UserConfig 资源测试01 6821
func TestAccAliCloudPaiWorkspaceUserConfig_basic6821(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_pai_workspace_user_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudPaiWorkspaceUserConfigMap6821)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return paiworkspace.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePaiWorkspaceUserConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudPaiWorkspaceUserConfigBasicDependence6821)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudPaiWorkspaceUserConfigMap6821)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return paiworkspace.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePaiWorkspaceUserConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudPaiWorkspaceUserConfigBasicDependence6821)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudPaiWorkspaceUserConfigMap6821)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return paiworkspace.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribePaiWorkspaceUserConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudPaiWorkspaceUserConfigBasicDependence6821)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
	"scope": acctest.CHECKSET,
}

func AliCloudPaiWorkspaceUserConfigBasicDependence6821(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test PaiWorkspace UserConfig. <<< Resource test cases, automatically generated.
