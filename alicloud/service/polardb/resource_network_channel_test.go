package polardb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAlicloudPolarDBNetworkChannel_basic(t *testing.T) {
	resourceId := "alicloud_polardb_network_channel.default"
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
				Config: testAccAlicloudPolarDBNetworkChannelConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttr(resourceId, "channel_name", name),
					resource.TestCheckResourceAttrSet(resourceId, "db_cluster_id"),
					resource.TestCheckResourceAttrSet(resourceId, "target_ip"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"notes", "channel_infos"},
			},
		},
	})
}

func testAccCheckPolarDBNetworkChannelDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_polardb_network_channel" {
			continue
		}
		parts := strings.Split(rs.Primary.ID, ":")
		if len(parts) < 2 {
			continue
		}
		request := map[string]interface{}{
			"RegionId":    client.RegionId,
			"DBClusterId": parts[0],
			"ChannelName": parts[1],
		}
		_, err := client.RpcPost("polardb", "2017-08-01", "DescribeNetworkChannel", request, nil, true)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("PolarDB Network Channel %s still exists", rs.Primary.ID)
	}
	return nil
}

func testAccPolarDBNetworkChannelDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_vpc" "default" {
  vpc_name = var.name
}

resource "alicloud_vswitch" "default" {
  vswitch_name = "tf_testAccPolarDB"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-hangzhou-j"
  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
}

locals {
  vpc_id      = alicloud_vpc.default.id
  vswitch_id  = concat(alicloud_vswitch.default.*.id, [""])[0]
}

resource "alicloud_polardb_cluster" "default" {
  db_type        = "PostgreSQL"
  db_version     = "14"
  pay_type       = "PostPaid"
  db_node_class  = "polar.pg.x4.medium"
  vswitch_id     = local.vswitch_id
  description    = var.name
}

resource "alicloud_polardb_cluster" "target" {
  db_type        = "PostgreSQL"
  db_version     = "14"
  pay_type       = "PostPaid"
  db_node_class  = "polar.pg.x4.medium"
  vswitch_id     = local.vswitch_id
  description    = "${var.name}-target"
}

resource "alicloud_polardb_network_channel" "default" {
  db_cluster_id        = alicloud_polardb_cluster.default.id
  channel_name         = var.name
  target_db_cluster_id = alicloud_polardb_cluster.target.id
  target_ip            = alicloud_polardb_cluster.target.connection_string
  target_port          = alicloud_polardb_cluster.target.port
  vpc_id               = local.vpc_id
}
`, name)
}

func testAccAlicloudPolarDBNetworkChannelConfig(name string) string {
	return testAccPolarDBNetworkChannelDependence(name)
}
