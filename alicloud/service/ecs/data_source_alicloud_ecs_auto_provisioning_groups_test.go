package ecs_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The existing alicloud_auto_provisioning_group resource reads the group
// immediately after creation and can observe it as absent while the service is
// still converging, so it is not used as an acceptance fixture here. This test
// exercises the live query, the complete filter request and the empty-result
// semantics; the non-empty object mapping is covered by the data source unit
// tests.
func TestAccAliCloudECSAutoProvisioningGroupsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccautoprovisioninggroupds-%d", rand)
	resourceId := "data.alicloud_ecs_auto_provisioning_groups.default"
	nameFilter := fmt.Sprintf("  auto_provisioning_group_name = %q\n", name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: autoProvisioningGroupsDataBlock(nameFilter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "groups.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "ids.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "names.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config: autoProvisioningGroupsDataBlock(
					"  ids               = [\"apg-does-not-exist\"]\n  resource_group_id = \"rg-does-not-exist\"\n  tags = {\n    env = \"tf-not-matching\"\n  }\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "groups.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config: autoProvisioningGroupsDataBlock(
					nameFilter + "  status      = [\"deleted\"]\n  types       = [\"maintain\"]\n  page_number = 1\n  page_size   = 10\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "groups.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config:      autoProvisioningGroupsDataBlock("  types = [\"not-a-type\"]\n"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match|expected value to be one of`),
			},
			{
				Config:      autoProvisioningGroupsDataBlock("  status = [\"not-a-status\"]\n"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match|expected value to be one of`),
			},
		},
	})
}

func autoProvisioningGroupsDataBlock(arguments string) string {
	return fmt.Sprintf(`
data "alicloud_ecs_auto_provisioning_groups" "default" {
%s}
`, arguments)
}
