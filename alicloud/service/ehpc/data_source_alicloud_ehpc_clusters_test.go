package ehpc_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudEhpcClustersDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.EhpcSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEhpcClustersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ehpc_cluster.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudEhpcClustersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ehpc_cluster.default.id}_fake"]`,
		}),
	}
	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEhpcClustersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ehpc_cluster.default.cluster_name}"`,
		}),
		FakeConfig: testAccCheckAlicloudEhpcClustersDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_ehpc_cluster.default.cluster_name}_fake"`,
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEhpcClustersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ehpc_cluster.default.id}"]`,
			"status": `"running"`,
		}),
		FakeConfig: testAccCheckAlicloudEhpcClustersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ehpc_cluster.default.id}"]`,
			"status": `"stopped"`,
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEhpcClustersDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ehpc_cluster.default.id}"]`,
			"name_regex": `"${alicloud_ehpc_cluster.default.cluster_name}"`,
			"status":     `"running"`,
		}),
		FakeConfig: testAccCheckAlicloudEhpcClustersDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_ehpc_cluster.default.id}_fake"]`,
			"name_regex": `"${alicloud_ehpc_cluster.default.cluster_name}_fake"`,
			"status":     `"stopped"`,
		}),
	}
	var existAlicloudEhpcClustersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"names.#":                          "1",
			"clusters.#":                       "1",
			"clusters.0.account_type":          "nis",
			"clusters.0.client_version":        acctest.CHECKSET,
			"clusters.0.cluster_name":          fmt.Sprintf("tf-testAccCluster-%d", rand),
			"clusters.0.deploy_mode":           "Simple",
			"clusters.0.description":           fmt.Sprintf("tf-testAccCluster-%d", rand),
			"clusters.0.ha_enable":             "false",
			"clusters.0.image_id":              acctest.CHECKSET,
			"clusters.0.image_owner_alias":     "system",
			"clusters.0.compute_count":         "1",
			"clusters.0.compute_instance_type": acctest.CHECKSET,
			"clusters.0.login_count":           "1",
			"clusters.0.login_instance_type":   acctest.CHECKSET,
			"clusters.0.manager_count":         "1",
			"clusters.0.manager_instance_type": acctest.CHECKSET,
			"clusters.0.os_tag":                "CentOS_7.6_64",
			"clusters.0.remote_directory":      acctest.CHECKSET,
			"clusters.0.scc_cluster_id":        "",
			"clusters.0.scheduler_type":        acctest.CHECKSET,
			"clusters.0.security_group_id":     acctest.CHECKSET,
			"clusters.0.volume_id":             acctest.CHECKSET,
			"clusters.0.volume_mountpoint":     acctest.CHECKSET,
			"clusters.0.volume_protocol":       acctest.CHECKSET,
			"clusters.0.volume_type":           acctest.CHECKSET,
			"clusters.0.vswitch_id":            acctest.CHECKSET,
			"clusters.0.vpc_id":                acctest.CHECKSET,
			"clusters.0.zone_id":               acctest.CHECKSET,
			"clusters.0.post_install_script.#": acctest.CHECKSET,
			"clusters.0.application.#":         acctest.CHECKSET,
		}
	}
	var fakeAlicloudEhpcClustersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudEhpcClustersCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ehpc_clusters.default",
		ExistMapFunc: existAlicloudEhpcClustersDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudEhpcClustersDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudEhpcClustersCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, statusConf, allConf)
}
func testAccCheckAlicloudEhpcClustersDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {
  default = "tf-testAccCluster-%d"
}
data "alicloud_zones" default {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}
data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
}
variable "storage_type" {
  default = "Capacity"
}
resource "alicloud_nas_file_system" "default" {
  storage_type  = var.storage_type
  protocol_type = "NFS"
}
resource "alicloud_nas_mount_target" "default" {
  file_system_id    = alicloud_nas_file_system.default.id
  access_group_name = "DEFAULT_VPC_GROUP_NAME"
  vswitch_id        = data.alicloud_vswitches.default.ids.0
}
data "alicloud_images" "default" {
  name_regex = "^centos_7_6_x64*"
  owners     = "system"
}
resource "alicloud_ehpc_cluster" "default" {
  cluster_name          = var.name
  deploy_mode           = "Simple"
  description           = var.name
  ha_enable             = false
  image_id              = data.alicloud_images.default.images.0.id
  image_owner_alias     = "system"
  volume_protocol       = "nfs"
  volume_id             = alicloud_nas_file_system.default.id
  volume_mountpoint     = alicloud_nas_mount_target.default.mount_target_domain
  compute_count         = 1
  compute_instance_type = data.alicloud_instance_types.default.instance_types.0.id
  login_count           = 1
  login_instance_type   = data.alicloud_instance_types.default.instance_types.0.id
  manager_count         = 1
  manager_instance_type = data.alicloud_instance_types.default.instance_types.0.id
  os_tag                = "CentOS_7.6_64"
  scheduler_type        = "pbs"
  password              = "your-password123"
  vswitch_id            = data.alicloud_vswitches.default.ids.0
  vpc_id                = data.alicloud_vpcs.default.ids.0
  zone_id               = data.alicloud_zones.default.zones.0.id
}

data "alicloud_ehpc_clusters" "default" {	
	enable_details = true
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
