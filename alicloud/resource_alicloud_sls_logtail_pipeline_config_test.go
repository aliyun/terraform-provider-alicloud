// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Sls LogtailPipelineConfig. >>> Resource test cases, automatically generated.
// Case LogtailPipelineConfigTestPL 12633
func TestAccAliCloudSlsLogtailPipelineConfig_basic12633(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_logtail_pipeline_config.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsLogtailPipelineConfigMap12633)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsLogtailPipelineConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsLogtailPipelineConfigBasicDependence12633)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project":     "terraform-logstore-test-578",
					"config_name": "pl-auto-test",
					"flushers":    []map[string]interface{}{},
					"inputs":      []map[string]interface{}{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project":     "terraform-logstore-test-578",
						"config_name": "pl-auto-test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"inputs": []map[string]interface{}{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"logstore_name"},
			},
		},
	})
}

var AlicloudSlsLogtailPipelineConfigMap12633 = map[string]string{}

func AlicloudSlsLogtailPipelineConfigBasicDependence12633(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Sls LogtailPipelineConfig. <<< Resource test cases, automatically generated.
