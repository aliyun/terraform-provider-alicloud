// Package alicloud
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func AliCloudCmsContextStoresBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

	provider "alicloud" {
		alias = "hz"
		region = "cn-hangzhou"
	}
	resource "alicloud_log_project" "default" {
		provider = alicloud.hz
		project_name = var.name
	}
	resource "alicloud_cms_workspace" "default" {
		provider = alicloud.hz
		sls_project = alicloud_log_project.default.project_name
		workspace_name = var.name
	}
	resource "alicloud_cms_context_store" "default" {
		provider = alicloud.hz
		workspace = alicloud_cms_workspace.default.workspace_name
		context_store_name = var.name
		context_type = "private"
		description = var.name
	}
`, name)
}

func TestAccAliCloudCmsContextStores_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_context_store.default"
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsContextStore")
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc("data.alicloud_cms_context_stores.default", name, AliCloudCmsContextStoresBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rc.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"provider":  "alicloud.hz",
					"workspace": "${alicloud_cms_workspace.default.workspace_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_cms_context_stores.default", "context_stores.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_cms_context_stores.default", "ids.#"),
				),
			},
		},
	})
}
