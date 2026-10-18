package cms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudCmsEntityStore_basic covers the singleton EntityStore
// lifecycle: create under a fresh workspace, read back the identity and
// region, import by workspace name, and destroy.
func TestAccAliCloudCmsEntityStore_basic(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", r)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCmsEntityStoreConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_entity_store.default", "workspace_name", name),
					resource.TestCheckResourceAttrSet("alicloud_cms_entity_store.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_cms_entity_store.default", "region_id"),
				),
			},
			{
				ResourceName:            "alicloud_cms_entity_store.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccCmsEntityStoreConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%[1]s"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = var.name
  sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_entity_store" "default" {
  workspace_name = alicloud_cms_workspace.default.workspace_name
}
`, name)
}
