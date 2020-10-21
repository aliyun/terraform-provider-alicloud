package alicloud

import (
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/config"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudConfigConfigurationRecorder_basic(t *testing.T) {
	var v config.ConfigurationRecorder
	resourceId := "alicloud_config_configuration_recorder.default"
	ra := resourceAttrInit(resourceId, ConfigConfigurationRecorderMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ConfigService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeConfigConfigurationRecorder")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", ConfigConfigurationRecorderBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.CloudConfigSupportedRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enterprise_edition"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_types": []string{"ACS::ECS::Instance"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_types.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_types": []string{"ACS::ECS::Instance", "ACS::ECS::Disk"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_types.#": "2",
					}),
				),
			},
		},
	})
}

var ConfigConfigurationRecorderMap = map[string]string{
	"enterprise_edition":         "false",
	"organization_enable_status": CHECKSET,
	"organization_master_id":     CHECKSET,
	"status":                     CHECKSET,
}

func ConfigConfigurationRecorderBasicdependence(name string) string {
	return ""
}
