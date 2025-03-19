package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudEfloInvocation_basic10348(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_invocation.default"
	ra := resourceAttrInit(resourceId, AlicloudEfloInvocationMap10348)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloInvocation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEfloInvocationBasicDependence10348)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":      "test",
					"content_encoding": "Base64",
					"name":             name,
					"repeat_mode":      "Once",
					"parameters": map[string]interface{}{
						"\"name\"": "Jack",
					},
					"node_id_list": []string{
						"e01-cn-rno46i6rdfn"},
					"timeout":          "68",
					"command_content":  "ZWNobyAxMjM=",
					"working_dir":      "/home/",
					"username":         "root",
					"enable_parameter": "false",
					"termination_mode": "ProcessTree",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":      "test",
						"content_encoding": "Base64",
						"name":             name,
						"repeat_mode":      "Once",
						"timeout":          "68",
						"command_content":  "ZWNobyAxMjM=",
						"working_dir":      "/home/",
						"username":         "root",
						"enable_parameter": "false",
						"termination_mode": "ProcessTree",
					}),
				),
			},
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
				ImportStateVerifyIgnore: []string{"command_content", "command_id", "content_encoding", "description", "enable_parameter", "frequency", "launcher", "name", "parameters", "repeat_mode", "termination_mode", "timeout", "username", "working_dir", "node_id_list"},
			},
		},
	})
}

var AlicloudEfloInvocationMap10348 = map[string]string{}

func AlicloudEfloInvocationBasicDependence10348(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}
