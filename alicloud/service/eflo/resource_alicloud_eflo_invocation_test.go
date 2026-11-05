package eflo_test

// TypeList Order Coverage: the multi-member TypeList attributes listed below
// are pre-existing stock debt surfaced by this directory-only migration (no
// schema or behavior change). They carry no convergent reorder acceptance
// test in this package; marked per scripts/collection-order/README.md
// stock-marker contract so the collection-order checker skips them.
// collection-order: bypass node_id_list -- pre-existing stock TypeList without convergent reorder acceptance coverage; directory-only migration surfaces it, no schema or behavior change

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/eflo"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// lintignore: AT001
func TestAccAliCloudEfloInvocation_basic10348(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_invocation.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudEfloInvocationMap10348)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return eflo.NewEfloServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEfloInvocation")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudEfloInvocationBasicDependence10348)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
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
