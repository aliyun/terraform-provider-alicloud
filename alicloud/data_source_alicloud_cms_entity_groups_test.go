package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCmsEntityGroupsDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	dataSourceId := "data.alicloud_cms_entity_groups.default"

	dependence := AliCloudCmsEntityGroupBasicDependence(name)
	resourceBlock := fmt.Sprintf(`
resource "alicloud_cms_entity_group" "default" {
  entity_group_name = "%s"
  description       = "tf test datasource entity group"
  workspace         = alicloud_cms_workspace.default.id

  entity_rules {
    entity_types = ["ECS"]
  }
}
`, name)
	dataSourceBlock := fmt.Sprintf(`
data "alicloud_cms_entity_groups" "default" {
  workspace = alicloud_cms_workspace.default.id
  name      = "%s"
}
`, name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudCmsEntityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: dependence + resourceBlock + dataSourceBlock,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceId, "groups.#", "1"),
					resource.TestCheckResourceAttr(dataSourceId, "groups.0.entity_group_name", name),
					resource.TestCheckResourceAttr(dataSourceId, "groups.0.entity_rules.#", "1"),
				),
			},
		},
	})
}

func TestAccAliCloudCmsEntityGroupsDataSource_filterByType(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%d", rand)
	dataSourceId := "data.alicloud_cms_entity_groups.default"

	dependence := AliCloudCmsEntityGroupBasicDependence(name)
	resourceBlock := fmt.Sprintf(`
resource "alicloud_cms_entity_group" "default" {
  entity_group_name = "%s"
  description       = "tf test datasource filter"
  workspace         = alicloud_cms_workspace.default.id

  entity_rules {
    entity_types = ["ECS"]
  }
}
`, name)
	dataSourceBlock := `
data "alicloud_cms_entity_groups" "default" {
  workspace         = alicloud_cms_workspace.default.id
  entity_group_type = "ECS"
}
`

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudCmsEntityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: dependence + resourceBlock + dataSourceBlock,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceId, "groups.#", "1"),
				),
			},
		},
	})
}
