package apig_test

import (
	"fmt"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	apig "github.com/aliyun/terraform-provider-alicloud/alicloud/service/apig"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Apig AiModelProvider. >>> Resource test cases, hand-written.
func TestAccAliCloudApigAiModelProvider_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_ai_model_provider.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudApigAiModelProviderMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return apig.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeApigAiModelProvider")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapigaimp%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudApigAiModelProviderBasicDependence)
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
					"gateway_id":     "${alicloud_apig_gateway.default.id}",
					"model_provider": "openai",
					"display_name":   name,
					"service_ids":    []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"gateway_id":     acctest.CHECKSET,
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
	"model_count": acctest.CHECKSET,
	"source":      acctest.CHECKSET,
	"update_time": acctest.CHECKSET,
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

// collection-order: bypass service_ids -- round-trippable (Read sets service_ids from boundServices; non-ForceNew) but non-empty member-order convergence is unverifiable in-harness: alicloud_apig_service.source_type ValidateFunc only allows DNS/VIP/FC3, none qualify as the AI service the gateway requires (ACC TestAccAliCloudApigAiModelProvider_memberOrder step1 CREATE failed InvalidParameter.WithDetail "AI service not found in gateway"); basic test sidesteps with service_ids=[] + ImportStateVerifyIgnore. Per-property adjudication, not a uniform historical bypass.

// Test Apig AiModelProvider. <<< Resource test cases, hand-written.
