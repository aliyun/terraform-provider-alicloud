package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudDataWorksComponent_basic8904(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_component.default"
	ra := resourceAttrInit(resourceId, AlicloudDataWorksComponentMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataworksPublicService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDataWorksComponent")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdataworkscomponent%d", defaultRegionToTest, rand)
	// spec is a JSON string with embedded double quotes. The testAccConfig map
	// helper's valueConvert wraps string values in quotes without escaping
	// inner quotes, so a single-escaped \" (Go in-memory ") produces invalid
	// HCL: spec = "{"nodeType":...}" — the parser sees the inner " as the
	// string end and reports "argument definition must end with a newline".
	// Double-escape the inner quotes (\\\" in Go source → \" in memory) so
	// valueConvert emits spec = "{\"nodeType\":...}" — valid HCL with escaped
	// inner quotes. Same convention as api_override_json /
	// monitor_config_json in other resource tests.
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDataWorksComponentBasicDependence0)
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
					"project_id":     "34051",
					"spec":           "{\\\"nodeType\\\":\\\"NODE_TYPE_DEFAULT\\\",\\\"componentName\\\":\\\"tf-test-component\\\"}",
					"component_type": "NODE_TYPE_DEFAULT",
					"source":         "MANUAL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_id": "34051",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"project_id":     "34051",
					"spec":           "{\\\"nodeType\\\":\\\"NODE_TYPE_DEFAULT\\\",\\\"componentName\\\":\\\"tf-test-component-updated\\\",\\\"nodeMode\\\":\\\"MODE_SYNC\\\"}",
					"component_type": "NODE_TYPE_DEFAULT",
					"source":         "MANUAL",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_id": "34051",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"spec"},
			},
		},
	})
}

var AlicloudDataWorksComponentMap0 = map[string]string{
	"component_id":   CHECKSET,
	"project_id":     "34051",
	"spec":           NOSET,
	"component_type": NOSET,
	"source":         NOSET,
}

func AlicloudDataWorksComponentBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}
