// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Wafv3 DefenseResourceGroup. >>> Resource test cases, automatically generated.
// Case 防护对象组-20251016-防护对象传入 11639
func TestAccAliCloudWafv3DefenseResourceGroup_basic11639(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_wafv3_defense_resource_group.default"
	ra := resourceAttrInit(resourceId, AlicloudWafv3DefenseResourceGroupMap11639)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Wafv3ServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeWafv3DefenseResourceGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccwafv3%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudWafv3DefenseResourceGroupBasicDependence11639)
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
					"group_name": "testfromTF",
					"resource_list": []string{
						"${alicloud_wafv3_domain.defaultHVcskT.domain_id}"},
					"description": "test",
					"instance_id": "${data.alicloud_wafv3_instances.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":      "testfromTF",
						"resource_list.#": "1",
						"description":     "test",
						"instance_id":     CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_list": []string{
						"${alicloud_wafv3_domain.defaultHVcskT.domain_id}", "${alicloud_wafv3_domain.defaultEH4CwO.domain_id}", "${alicloud_wafv3_domain.defaultY0ge1N.domain_id}"},
					"description": "testabc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_list.#": "3",
						"description":     "testabc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_list": []string{
						"${alicloud_wafv3_domain.defaultEH4CwO.domain_id}"},
					"description": "assss",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_list.#": "1",
						"description":     "assss",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_list": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_list.#": "0",
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

var AlicloudWafv3DefenseResourceGroupMap11639 = map[string]string{}

func AlicloudWafv3DefenseResourceGroupBasicDependence11639(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-hangzhou"
}

data "alicloud_wafv3_instances" "default" {
}

resource "alicloud_wafv3_domain" "defaultHVcskT" {
  instance_id = data.alicloud_wafv3_instances.default.ids.0
  listen {
    http_ports = ["80"]
  }
  redirect {
    loadbalance 	= "iphash"
    backends    	= ["6.36.36.36"]
    connect_timeout = 5
    read_timeout    = 120
    write_timeout   = 120
  }
  domain      = "1511928242963727_1.wafqax.top"
  access_type = "share"
}

resource "alicloud_wafv3_domain" "defaultEH4CwO" {
  instance_id = data.alicloud_wafv3_instances.default.ids.0
  listen {
    http_ports = ["80"]
  }
  redirect {
    loadbalance 	= "iphash"
    backends    	= ["6.36.36.36"]
    connect_timeout = 5
    read_timeout    = 120
    write_timeout   = 120
  }
  domain      = "1511928242963727_2.wafqax.top"
  access_type = "share"
}

resource "alicloud_wafv3_domain" "defaultY0ge1N" {
  instance_id = data.alicloud_wafv3_instances.default.ids.0
  listen {
    http_ports = ["80"]
  }
  redirect {
    loadbalance 	= "iphash"
    backends    	= ["6.36.36.36"]
    connect_timeout = 5
    read_timeout    = 120
    write_timeout   = 120
  }
  domain      = "1511928242963727_3.wafqax.top"
  access_type = "share"
}


`, name)
}

// Test Wafv3 DefenseResourceGroup. <<< Resource test cases, automatically generated.
