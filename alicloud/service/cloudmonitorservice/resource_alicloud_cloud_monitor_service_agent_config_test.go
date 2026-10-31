// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package cloudmonitorservice_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudmonitorservice"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test CloudMonitorService AgentConfig. >>> Resource test cases, automatically generated.
// Case AgentConfigTest 5607
// lintignore: AT001
func TestAccAliCloudCloudMonitorServiceAgentConfig_basic5607(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_agent_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCloudMonitorServiceAgentConfigMap5607)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCloudMonitorServiceServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCloudMonitorServiceAgentConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCloudMonitorServiceAgentConfigBasicDependence5607)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
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
