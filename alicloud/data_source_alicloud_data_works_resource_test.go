package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudDataWorksResourceDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	resourceId := "data.alicloud_data_works_resource.default"
	testAccDataSourceConfig := dataSourceTestAccConfigFunc(resourceId, name, AlicloudDataWorksResourceDataSourceBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConfig(map[string]interface{}{
					"project_id":  "${alicloud_data_works_project.default.id}",
					"output_file": "output.txt",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "resources.#"),
				),
			},
		},
	})
}

func AlicloudDataWorksResourceDataSourceBasicDependence(name string) string {
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
