package eflo_test

import (
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudEfloErRouteEntriesDataSource creates an Eflo Er and reads its
// route entries via the alicloud_eflo_er_route_entries data source. The Er
// resource is SDKv2 and the data source is Framework; both are served through
// the single ProtoV5 muxed provider, so ProtoV5ProviderFactories covers both.
//
// Eflo ER is only exercised in cn-hangzhou (master_zone_id cn-hangzhou-a),
// mirroring the existing alicloud_eflo_er resource test's region handling: the
// root-package test gates with testAccPreCheckWithRegions(t, true, [Hangzhou]),
// which switches ALICLOUD_REGION to cn-hangzhou rather than skipping when the
// runner's default region is unsupported. That helper is unexported in package
// alicloud, so this package-eflo_test replicates the switch. An all-SKIP run is
// a failure for the remote ACC runner, so the test must actually execute.
func TestAccAliCloudEfloErRouteEntriesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudEfloErRouteEntriesDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_eflo_er_route_entries.default", "id"),
					resource.TestCheckResourceAttrSet("data.alicloud_eflo_er_route_entries.default", "er_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_eflo_er_route_entries.default", "entries.#"),
				),
			},
		},
	})
}

func testAccAliCloudEfloErRouteEntriesDataSourceConfig() string {
	return `
variable "name" {
  default = "tf-acc-eflo-er-route-entries"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_eflo_er" "default" {
  er_name        = var.name
  master_zone_id = "cn-hangzhou-a"
  description    = "tf-acc data source test"
}

data "alicloud_eflo_er_route_entries" "default" {
  er_id = alicloud_eflo_er.default.id
}
`
}
