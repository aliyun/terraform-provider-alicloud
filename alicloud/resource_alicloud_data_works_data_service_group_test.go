package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudDataWorksDataServiceGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_data_service_group.default"
	ra := resourceAttrInit(resourceId, AlicloudDataWorksDataServiceGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataWorksServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDataWorksDataServiceGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_dsg%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDataWorksDataServiceGroupBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		// The DataServiceGroup API does not provide a Delete operation, so Delete only
		// removes the resource from the state file and a post-destroy check can never pass.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project_id":              "${alicloud_data_works_project.default.id}",
					"api_gateway_group_id":    "${alicloud_api_gateway_group.default.id}",
					"data_service_group_name": name,
					"description":             name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_id":              CHECKSET,
						"api_gateway_group_id":    CHECKSET,
						"data_service_group_name": name,
						"description":             name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudDataWorksDataServiceGroupMap = map[string]string{
	"create_time":           CHECKSET,
	"creator_id":            CHECKSET,
	"data_service_group_id": CHECKSET,
	"modified_time":         CHECKSET,
}

func AlicloudDataWorksDataServiceGroupBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_data_works_project" "default" {
  project_name     = var.name
  display_name     = var.name
  description      = var.name
  pai_task_enabled = false
}

resource "alicloud_api_gateway_group" "default" {
  name        = var.name
  description = var.name
}
`, name)
}
