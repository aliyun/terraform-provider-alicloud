package ecs_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// A real auto provisioning group is created so the query is proven against a
// matching object: the assertion checks the target ID, the name and the group
// state, and a second query with a non-existent ID proves the ID filter is
// applied instead of returning every group.
func TestAccAliCloudECSAutoProvisioningGroupsDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 999999)
	name := fmt.Sprintf("tf-testAccautoprovisioninggroupds-%d", rand)
	resourceId := "data.alicloud_ecs_auto_provisioning_groups.default"
	groupResourceId := "alicloud_auto_provisioning_group.default"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: autoProvisioningGroupsConfig(name, "  ids = [alicloud_auto_provisioning_group.default.id]\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "groups.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceId, "ids.0", groupResourceId, "id"),
					resource.TestCheckResourceAttrPair(resourceId, "groups.0.auto_provisioning_group_id", groupResourceId, "id"),
					resource.TestCheckResourceAttrPair(resourceId, "groups.0.auto_provisioning_group_name", groupResourceId, "auto_provisioning_group_name"),
					resource.TestCheckResourceAttr(resourceId, "groups.0.auto_provisioning_group_type", "maintain"),
					resource.TestCheckResourceAttrSet(resourceId, "groups.0.status"),
					resource.TestCheckResourceAttr(resourceId, "names.#", "1"),
					resource.TestCheckResourceAttrPair(resourceId, "names.0", groupResourceId, "auto_provisioning_group_name"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "1"),
				),
			},
			{
				Config: autoProvisioningGroupsConfig(name, fmt.Sprintf("  auto_provisioning_group_name = %q\n", name)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "groups.#", "1"),
					resource.TestCheckResourceAttrPair(resourceId, "ids.0", groupResourceId, "id"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "1"),
				),
			},
			{
				Config: autoProvisioningGroupsConfig(name, `  ids = ["apg-does-not-exist"]
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "groups.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "ids.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config: autoProvisioningGroupsConfig(name, fmt.Sprintf("  auto_provisioning_group_name = \"no-such-group-%d\"\n", rand)),
				Check:  resource.TestCheckResourceAttr(resourceId, "groups.#", "0"),
			},
			{
				Config:      autoProvisioningGroupsConfig(name, "  types = [\"not-a-type\"]\n"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match|expected value to be one of`),
			},
			{
				Config:      autoProvisioningGroupsConfig(name, "  status = [\"not-a-status\"]\n"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match|expected value to be one of`),
			},
			{
				Config:      autoProvisioningGroupsConfig(name, "  tags = { \"env\" = \"https://example.com\" }\n"),
				ExpectError: regexp.MustCompile(`Invalid Tag Value|must not contain a URL`),
			},
			{
				Config:      autoProvisioningGroupsConfig(name, "  tags = { \"aliyun-env\" = \"test\" }\n"),
				ExpectError: regexp.MustCompile(`Invalid Tag Key|must not start with`),
			},
			{
				// The destroy phase re-reads the last configuration, so a valid
				// step must follow the expected validation failures.
				Config: autoProvisioningGroupsConfig(name, "  ids = [alicloud_auto_provisioning_group.default.id]\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "groups.#", "1"),
				),
			},
		},
	})
}

func autoProvisioningGroupsConfig(name, dataArguments string) string {
	return fmt.Sprintf(`
variable "name" {
  default = %q
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
  cpu_core_count    = 2
  memory_size       = 4
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu"
  most_recent = true
  owners      = "system"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_launch_template" "default" {
  name              = var.name
  image_id          = data.alicloud_images.default.images.0.id
  instance_type     = data.alicloud_instance_types.default.instance_types.0.id
  security_group_id = alicloud_security_group.default.id
  vswitch_id        = alicloud_vswitch.default.id
}

resource "alicloud_auto_provisioning_group" "default" {
  launch_template_id            = alicloud_launch_template.default.id
  auto_provisioning_group_name  = var.name
  auto_provisioning_group_type  = "maintain"
  total_target_capacity         = "1"
  spot_target_capacity          = "1"
  default_target_capacity_type  = "Spot"

  launch_template_config {
    vswitch_id        = alicloud_vswitch.default.id
    weighted_capacity = "1"
    max_price         = "2"
    instance_type     = data.alicloud_instance_types.default.instance_types.0.id
  }
}

data "alicloud_ecs_auto_provisioning_groups" "default" {
%s}
`, name, dataArguments)
}
