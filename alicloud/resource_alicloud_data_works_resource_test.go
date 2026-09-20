package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test DataWorks Resource. >>> Resource test cases.
func TestAccAliCloudDataWorksResource_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_resource.default"
	ra := resourceAttrInit(resourceId, AlicloudDataWorksResourceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataWorksServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDataWorksResource")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDataWorksResourceBasicDependence)
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
					"project_id": "${alicloud_data_works_project.default.id}",
					"spec":       "{}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_id":  CHECKSET,
						"resource_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"project_id":    "${alicloud_data_works_project.default.id}",
					"spec":          "{}",
					"resource_name": name,
					"resource_file": "test-resource-file-content",
					"path":          "/test/path",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_id":    CHECKSET,
						"resource_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"project_id": "${alicloud_data_works_project.default.id}",
					"spec":       "{}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"resource_file", "path", "spec"},
			},
		},
	})
}

var AlicloudDataWorksResourceMap = map[string]string{}

func AlicloudDataWorksResourceBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_data_works_project" "default" {
    project_name           = var.name
    display_name           = var.name
    pai_task_enabled       = false
    dev_environment_enabled = true
}

`, name)
}

// Test DataWorks Resource. <<< Resource test cases.
