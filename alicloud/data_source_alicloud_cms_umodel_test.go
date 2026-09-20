package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCmsUmodelDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccumodelds%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc("data.alicloud_cms_umodel.default", name, AliCloudCmsUmodelDataSourceBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "data.alicloud_cms_umodel.default",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"workspace": "${alicloud_cms_umodel.default.workspace}",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_cms_umodel.default", "region_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_cms_umodel.default", "id"),
				),
			},
		},
	})
}

func TestAccAliCloudCmsUmodelDataSource_disappear(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccumodeldsdisp%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc("data.alicloud_cms_umodel.default", name, AliCloudCmsUmodelDataSourceBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "data.alicloud_cms_umodel.default",
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"workspace": "${alicloud_cms_umodel.default.workspace}",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_cms_umodel.default", "id"),
				),
			},
		},
	})
}

func AliCloudCmsUmodelDataSourceBasicDependence(name string) string {
	return fmt.Sprintf(`
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "%s"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  sls_project    = alicloud_log_project.default.project_name
  workspace_name = var.name
}

resource "alicloud_cms_umodel" "default" {
  workspace   = alicloud_cms_workspace.default.workspace_name
  description = var.name
}
`, name)
}
