// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudMonitorService AgentConfig. >>> Resource test cases, automatically generated.
// Case AgentConfigTest 5607
// lintignore: AT001
func TestAccAliCloudCloudMonitorServiceAgentConfig_basic5607(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_agent_config.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceAgentConfigMap5607)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceAgentConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceAgentConfigBasicDependence5607)
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
					"enable_install_agent_new_ecs": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_install_agent_new_ecs": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_install_agent_new_ecs": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_install_agent_new_ecs": "true",
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

var AlicloudCloudMonitorServiceAgentConfigMap5607 = map[string]string{}

func AlicloudCloudMonitorServiceAgentConfigBasicDependence5607(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test CloudMonitorService AgentConfig. <<< Resource test cases, automatically generated.
