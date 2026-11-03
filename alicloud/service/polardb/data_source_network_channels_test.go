package polardb_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAlicloudPolarDBNetworkChannels_basic(t *testing.T) {
	resourceId := "alicloud_polardb_network_channel.default"
	dataSourceId := "data.alicloud_polardb_network_channels.default"
	rand := tfacctest.RandStringFromCharSet(10, tfacctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testacc-polardb-nc-%s", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPolarDBNetworkChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudPolarDBNetworkChannelsConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttrSet(dataSourceId, "id"),
					resource.TestCheckResourceAttr(dataSourceId, "channel_infos.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceId, "channel_infos.0.channel_name"),
					resource.TestCheckResourceAttrSet(dataSourceId, "channel_infos.0.db_cluster_id"),
					resource.TestCheckResourceAttrSet(dataSourceId, "channel_infos.0.target_ip"),
					resource.TestCheckResourceAttr(dataSourceId, "ids.#", "1"),
					resource.TestCheckResourceAttr(dataSourceId, "names.#", "1"),
				),
			},
		},
	})
}

func testAccAlicloudPolarDBNetworkChannelsConfig(name string) string {
	return fmt.Sprintf(`
%s

data "alicloud_polardb_network_channels" "default" {
  db_cluster_id = alicloud_polardb_cluster.default.id
}
`, testAccPolarDBNetworkChannelDependence(name))
}
