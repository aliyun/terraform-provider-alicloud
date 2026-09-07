package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Apig AiModelProvider. >>> Resource test cases, hand-written.
func TestAccAliCloudApigAiModelProvider_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_ai_model_provider.default"
	ra := resourceAttrInit(resourceId, AlicloudApigAiModelProviderMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigAiModelProvider")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapigaimp%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigAiModelProviderBasicDependence)
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
					"gateway_id":     "${alicloud_apig_gateway.default.id}",
					"model_provider": "openai",
					"display_name":   name,
					"service_ids":    []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_id":     CHECKSET,
						"model_provider": "openai",
						"display_name":   name,
						"service_ids.#":  "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": name + "_update",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_ids"},
			},
		},
	})
}

var AlicloudApigAiModelProviderMap = map[string]string{
	"model_count": CHECKSET,
	"source":      CHECKSET,
	"update_time": CHECKSET,
}

func AlicloudApigAiModelProviderBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "j" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-hangzhou-j"
}

data "alicloud_vswitches" "k" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-hangzhou-k"
}

# AI gateway auto-creates a security group whose cleanup lags gateway Delete;
# using the shared NODELETING VPC avoids the framework attempting to delete
# the VPC and hitting DependencyViolation.SecurityGroup during test teardown.
resource "alicloud_apig_gateway" "default" {
  gateway_name = var.name
  spec         = "aigw.small.x1"
  gateway_type = "AI"
  payment_type = "PayAsYouGo"
  vpc {
    vpc_id = data.alicloud_vpcs.default.ids.0
  }
  network_access_config {
    type = "Intranet"
  }
  zone_config {
    select_option = "Manual"
  }
  vswitch {
    vswitch_id = data.alicloud_vswitches.j.ids.0
  }
  log_config {
    sls {
      enable = false
    }
  }
  zones {
    vswitch_id = data.alicloud_vswitches.j.ids.0
    zone_id    = "cn-hangzhou-j"
  }
  zones {
    vswitch_id = data.alicloud_vswitches.k.ids.0
    zone_id    = "cn-hangzhou-k"
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}

`, name)
}

// Test Apig AiModelProvider. <<< Resource test cases, hand-written.
