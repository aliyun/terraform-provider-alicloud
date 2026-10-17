package paistudio_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudPaiResourceGroup covers the PAI Studio ResourceGroup resource,
// including the governance gap attribute `version` (Create-only, ForceNew).
// The write-only `unbind` attribute (Update-only, no readback) is covered by the
// schema and update path but is not exercised here because it requires a VPC
// already connected to the resource group.
func TestAccAliCloudPaiResourceGroup(t *testing.T) {
	r := tfacctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudPaiResourceGroupConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_resource_group.default", "id"),
					resource.TestCheckResourceAttr("alicloud_pai_resource_group.default", "name", fmt.Sprintf("tf-testacc-%d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_resource_group.default", "description", fmt.Sprintf("tf-testacc-%d-desc", r)),
					resource.TestCheckResourceAttr("alicloud_pai_resource_group.default", "resource_type", "Ecs"),
					resource.TestCheckResourceAttr("alicloud_pai_resource_group.default", "computing_resource_provider", "Ecs"),
					resource.TestCheckResourceAttr("alicloud_pai_resource_group.default", "version", "1.0"),
					resource.TestCheckResourceAttrSet("alicloud_pai_resource_group.default", "status"),
				),
			},
			{
				Config: testAccAliCloudPaiResourceGroupConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pai_resource_group.default", "name", fmt.Sprintf("tf-testacc-%d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_resource_group.default", "description", fmt.Sprintf("tf-testacc-%d-desc-upd", r)),
					resource.TestCheckResourceAttr("alicloud_pai_resource_group.default", "resource_type", "Ecs"),
					resource.TestCheckResourceAttr("alicloud_pai_resource_group.default", "version", "1.0"),
				),
			},
			{
				ResourceName:      "alicloud_pai_resource_group.default",
				ImportState:       true,
				ImportStateVerify: true,
				// ComputingResourceProvider is accepted on Create but GetResourceGroup
				// returns it as an empty string for Ecs resource groups, so it cannot be
				// reconstructed on import. The provider preserves the configured value in
				// state (see echoString in the resource), so it survives refresh but is
				// absent from an import that only knows the id.
				ImportStateVerifyIgnore: []string{"computing_resource_provider"},
			},
		},
	})
}

func testAccAliCloudPaiResourceGroupConfig(r int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testacc-%[1]d"
}

resource "alicloud_pai_resource_group" "default" {
  name                        = var.name
  description                 = "tf-testacc-%[1]d-desc"
  resource_type               = "Ecs"
  computing_resource_provider = "Ecs"
  version                     = "1.0"
}
`, r)
}

func testAccAliCloudPaiResourceGroupConfigUpdate(r int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testacc-%[1]d"
}

resource "alicloud_pai_resource_group" "default" {
  name                        = var.name
  description                 = "tf-testacc-%[1]d-desc-upd"
  resource_type               = "Ecs"
  computing_resource_provider = "Ecs"
  version                     = "1.0"
}
`, r)
}
