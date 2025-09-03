// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Sls MachineGroup. >>> Resource test cases, automatically generated.
// Case machine_group_terrafrom 10981
func TestAccAliCloudSlsMachineGroup_basic10981(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_machine_group.default"
	ra := resourceAttrInit(resourceId, AlicloudSlsMachineGroupMap10981)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlsMachineGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccsls%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSlsMachineGroupBasicDependence10981)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-nanjing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"group_name":            "group1",
					"group_type":            "test",
					"project_name":          "${alicloud_log_project.defaultyJqrue.project_name}",
					"machine_identify_type": "ip",
					"group_attribute": []map[string]interface{}{
						{
							"group_topic":   "test",
							"external_name": "test",
						},
					},
					"machine_list": []string{
						"192.168.1.1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":            "group1",
						"group_type":            "test",
						"project_name":          CHECKSET,
						"machine_identify_type": "ip",
						"machine_list.#":        "1",
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

var AlicloudSlsMachineGroupMap10981 = map[string]string{}

func AlicloudSlsMachineGroupBasicDependence10981(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "project_name" {
  default = "project-for-machine-group-terraform"
}

resource "alicloud_log_project" "defaultyJqrue" {
  description = "for terraform test"
  name        = var.project_name
}


`, name)
}

// Test Sls MachineGroup. <<< Resource test cases, automatically generated.
