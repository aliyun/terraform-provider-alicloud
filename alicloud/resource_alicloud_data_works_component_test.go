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
	// spec is a JSON string with embedded double quotes. Passing it through the
	// testAccConfig map helper breaks HCL parsing: valueConvert
	// (service_alicloud_common_test.go) wraps string values in quotes without
	// escaping inner quotes, so the parser sees the inner " as the string end
	// and "nodeType" as a bare argument with no newline. Hand-write the HCL and
	// escape inner quotes with \" — same convention as monitor_contacts_json /
	// monitor_config_json in resource_alicloud_schedulerx_job_test.go.
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: AlicloudDataWorksComponentBasic8904Config(name, `"{\"nodeType\":\"NODE_TYPE_DEFAULT\",\"componentName\":\"tf-test-component\"}"`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_id": "34051",
					}),
				),
			},
			{
				Config: AlicloudDataWorksComponentBasic8904Config(name, `"{\"nodeType\":\"NODE_TYPE_DEFAULT\",\"componentName\":\"tf-test-component-updated\",\"nodeMode\":\"MODE_SYNC\"}"`),
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

// AlicloudDataWorksComponentBasic8904Config builds the HCL for the
// DataWorks Component acceptance test. spec is a JSON string with embedded
// double quotes; callers must pass it in HCL-escaped form (inner quotes as \")
// so the HCL parser treats it as a single string. Hand-writing the resource
// block instead of using the testAccConfig map helper avoids valueConvert
// wrapping the value in unescaped quotes (see comment in TestAcc... above).
func AlicloudDataWorksComponentBasic8904Config(name, spec string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_data_works_component" "default" {
  project_id     = "34051"
  spec           = %s
  component_type = "NODE_TYPE_DEFAULT"
  source         = "MANUAL"
}
`, name, spec)
}
