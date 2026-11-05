package ecs_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudECSAutoProvisioningGroupsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccautoprovisioninggroupds-%d", rand)
	resourceId := "data.alicloud_ecs_auto_provisioning_groups.default"
	baseConfig := resourceAutoProvisioningGroupConfigDependence(name)
	groupConfig := fmt.Sprintf(`
resource "alicloud_auto_provisioning_group" "default" {
  launch_template_id           = alicloud_launch_template.template.id
  total_target_capacity        = "1"
  spot_target_capacity         = "1"
  auto_provisioning_group_name = "%s"
  launch_template_config {
    instance_type     = "ecs.sn1ne.large"
    vswitch_id        = alicloud_vswitch.default.id
    weighted_capacity = "1"
    max_price         = "2"
  }
}
`, name)
	nameFilter := fmt.Sprintf("  auto_provisioning_group_name = %q\n", name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: baseConfig + groupConfig + autoProvisioningGroupsDataBlock(nameFilter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "groups.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "ids.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "names.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "1"),
					resource.TestCheckResourceAttrPair(resourceId, "ids.0", "alicloud_auto_provisioning_group.default", "id"),
					resource.TestCheckResourceAttrPair(resourceId, "groups.0.auto_provisioning_group_id", "alicloud_auto_provisioning_group.default", "id"),
					resource.TestCheckResourceAttr(resourceId, "groups.0.auto_provisioning_group_name", name),
					resource.TestCheckResourceAttr(resourceId, "groups.0.auto_provisioning_group_type", "maintain"),
					resource.TestCheckResourceAttrSet(resourceId, "groups.0.status"),
					resource.TestCheckResourceAttrSet(resourceId, "groups.0.launch_template_id"),
					resource.TestCheckResourceAttr(resourceId, "groups.0.launch_template_configs.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "groups.0.launch_template_configs.0.instance_type", "ecs.sn1ne.large"),
					resource.TestCheckResourceAttr(resourceId, "groups.0.launch_template_configs.0.max_price", "2"),
					resource.TestCheckResourceAttr(resourceId, "groups.0.launch_template_configs.0.weighted_capacity", "1"),
					resource.TestCheckResourceAttr(resourceId, "groups.0.target_capacity_specification.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "groups.0.target_capacity_specification.0.total_target_capacity", "1"),
					resource.TestCheckResourceAttr(resourceId, "groups.0.spot_options.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "groups.0.spot_options.0.allocation_strategy", "lowest-price"),
				),
			},
			{
				Config: baseConfig + groupConfig + autoProvisioningGroupsDataBlock(
					"  ids = [alicloud_auto_provisioning_group.default.id]\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "groups.#", "1"),
					resource.TestCheckResourceAttrPair(resourceId, "groups.0.auto_provisioning_group_id", "alicloud_auto_provisioning_group.default", "id"),
				),
			},
			{
				Config: baseConfig + groupConfig + autoProvisioningGroupsDataBlock(
					"  ids               = [\"apg-does-not-exist\"]\n  resource_group_id = \"rg-does-not-exist\"\n  tags = {\n    env = \"tf-not-matching\"\n  }\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "groups.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config: baseConfig + groupConfig + autoProvisioningGroupsDataBlock(
					nameFilter+"  types = [\"maintain\"]\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "groups.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "groups.0.auto_provisioning_group_type", "maintain"),
				),
			},
			{
				Config: baseConfig + groupConfig + autoProvisioningGroupsDataBlock(
					nameFilter+"  status      = [\"deleted\"]\n  page_number = 1\n  page_size   = 10\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "groups.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config: baseConfig + groupConfig + autoProvisioningGroupsDataBlock(
					fmt.Sprintf("  auto_provisioning_group_name = %q\n", name+"-fake")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "groups.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "ids.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "names.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config:      baseConfig + groupConfig + autoProvisioningGroupsDataBlock("  types = [\"not-a-type\"]\n"),
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
