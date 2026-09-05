package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudApigPluginAttachment_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_plugin_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudApigPluginAttachmentMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigPluginAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigPluginAttachmentBasicDependence)
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
					"attach_resource_type": "HttpApi",
					"attach_resource_ids":  []string{"${alicloud_apig_http_api.default.id}"},
					"attach_resource_id":   "${alicloud_apig_http_api.default.id}",
					"environment_id":       "${alicloud_apig_gateway.default.environments.0.environment_id}",
					"enable":               "true",
					"plugin_info": []map[string]interface{}{
						{
							"plugin_id":     "${alicloud_apig_plugin.default.id}",
							"gateway_id":    "${alicloud_apig_gateway.default.id}",
							"plugin_config": "eyJ0ZXN0IjoiaGVsbG8ifQ==",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"attach_resource_type":    "HttpApi",
						"attach_resource_ids.#":   "1",
						"attach_resource_id":      CHECKSET,
						"environment_id":          CHECKSET,
						"enable":                  "true",
						"plugin_info.#":           "1",
						"plugin_info.0.plugin_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"attach_resource_type": "HttpApi",
					"attach_resource_ids":  []string{"${alicloud_apig_http_api.default.id}"},
					"attach_resource_id":   "${alicloud_apig_http_api.default.id}",
					"environment_id":       "${alicloud_apig_gateway.default.environments.0.environment_id}",
					"enable":               "false",
					"plugin_info": []map[string]interface{}{
						{
							"plugin_id":     "${alicloud_apig_plugin.default.id}",
							"gateway_id":    "${alicloud_apig_gateway.default.id}",
							"plugin_config": "eyJ0ZXN0Ijoid29ybGQifQ==",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"attach_resource_type": "HttpApi",
					"attach_resource_ids":  []string{"${alicloud_apig_http_api.default2.id}"},
					"attach_resource_id":   "${alicloud_apig_http_api.default2.id}",
					"environment_id":       "${alicloud_apig_gateway.default.environments.0.environment_id}",
					"enable":               "true",
					"plugin_info": []map[string]interface{}{
						{
							"plugin_id":     "${alicloud_apig_plugin.default.id}",
							"gateway_id":    "${alicloud_apig_gateway.default.id}",
							"plugin_config": "eyJ0ZXN0IjoiaGVsbG8ifQ==",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"attach_resource_ids.#": "1",
						"enable":                "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"plugin_info.0.gateway_id"},
			},
		},
	})
}

var AlicloudApigPluginAttachmentMap = map[string]string{
	"plugin_class_info.#": CHECKSET,
}

func AlicloudApigPluginAttachmentBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "default" {
  network_access_config {
    type = "Internet"
  }
  vswitch {
    vswitch_id = data.alicloud_vswitches.default.ids.0
  }
  zone_config {
    select_option = "Auto"
  }
  vpc {
    vpc_id = data.alicloud_vpcs.default.ids.0
  }
  gateway_type = "API"
  payment_type = "PayAsYouGo"
  gateway_name = var.name
  spec         = "apigw.small.x1"
  log_config {
    sls {
      enable = true
    }
  }
}

resource "alicloud_apig_plugin" "default" {
  plugin_class_id = "pls-crpqb35lhtgo800k2m86"
  gateway_id      = alicloud_apig_gateway.default.id
}


resource "alicloud_apig_http_api" "default" {
  http_api_name = var.name
  protocols     = ["HTTP"]
  type          = "Rest"
  base_path     = "/terraform-acc"
}

resource "alicloud_apig_http_api" "default2" {
  http_api_name = "${var.name}-2"
  protocols     = ["HTTP"]
  type          = "Rest"
  base_path     = "/terraform-acc-2"
}
`, name)
}
