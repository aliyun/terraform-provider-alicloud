package paifs_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudPaiFeatureStoreDatasources_basic creates a Datasource and then
// lists the instance's datasources, verifying that the newly created datasource
// is returned by the collection Read. Membership is asserted through a contains
// output rather than an exact count, so the test stays robust whether or not the
// optional Name filter narrows the list to a single entry.
func TestAccAliCloudPaiFeatureStoreDatasources_basic(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	instanceId, workspaceId, dsType, uri, config := paifsDatasourceFixtures(t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudPaiFeatureStoreDatasourcesConfig(r, instanceId, workspaceId, dsType, uri, config),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_pai_feature_store_datasources.default", "instance_id", instanceId),
					resource.TestCheckResourceAttrSet("data.alicloud_pai_feature_store_datasources.default", "ids.#"),
					resource.TestCheckResourceAttrSet("data.alicloud_pai_feature_store_datasources.default", "total_count"),
					resource.TestCheckResourceAttrSet("data.alicloud_pai_feature_store_datasources.default", "datasources.#"),
					resource.TestCheckResourceAttr("data.alicloud_pai_feature_store_datasources.default", "found", "true"),
				),
			},
		},
	})
}

func testAccAliCloudPaiFeatureStoreDatasourcesConfig(r int, instanceId, workspaceId, dsType, uri, config string) string {
	name := fmt.Sprintf("tf-testacc%d", r)
	return fmt.Sprintf(`
variable "instance_id" {
  type    = string
  default = "%[2]s"
}
variable "workspace_id" {
  type    = string
  default = "%[3]s"
}

resource "alicloud_pai_feature_store_datasource" "default" {
  instance_id     = var.instance_id
  workspace_id    = var.workspace_id
  type            = "%[4]s"
  datasource_name = "%[6]s"
  uri             = "%[5]s"
%[7]s}

data "alicloud_pai_feature_store_datasources" "default" {
  instance_id = var.instance_id
  name        = "%[6]s"
}

output "found" {
  value = contains(data.alicloud_pai_feature_store_datasources.default.ids, alicloud_pai_feature_store_datasource.default.datasource_id)
}
`, r, instanceId, workspaceId, dsType, uri, name, paifsConfigBlock(config))
}
