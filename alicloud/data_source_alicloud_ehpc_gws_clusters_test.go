package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEhpcGwsClustersDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.EhpcSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEhpcGwsClustersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ehpc_gws_cluster.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEhpcGwsClustersDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ehpc_gws_cluster.default.id}_fake"]`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEhpcGwsClustersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ehpc_gws_cluster.default.id}"]`,
			"status": `"running"`,
		}),
		fakeConfig: testAccCheckAlicloudEhpcGwsClustersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ehpc_gws_cluster.default.id}"]`,
			"status": `"deleted"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEhpcGwsClustersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ehpc_gws_cluster.default.id}"]`,
			"status": `"running"`,
		}),
		fakeConfig: testAccCheckAlicloudEhpcGwsClustersDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_ehpc_gws_cluster.default.id}_fake"]`,
			"status": `"deleted"`,
		}),
	}
	var existAlicloudEhpcGwsClustersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"clusters.#":                "1",
			"clusters.0.vpc_id":         CHECKSET,
			"clusters.0.id":             CHECKSET,
			"clusters.0.gws_cluster_id": CHECKSET,
			"clusters.0.create_time":    CHECKSET,
			"clusters.0.status":         "running",
		}
	}
	var fakeAlicloudEhpcGwsClustersDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudEhpcGwsClustersCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ehpc_gws_clusters.default",
		existMapFunc: existAlicloudEhpcGwsClustersDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEhpcGwsClustersDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudEhpcGwsClustersCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, allConf)
}
func testAccCheckAlicloudEhpcGwsClustersDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccGwsCluster-%d"
}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_ehpc_gws_cluster" "default" {
  name         = var.name
  cluster_type = "gws.s1.standard"
  vswitch_id   = data.alicloud_vswitches.default.ids.0
  vpc_id       = data.alicloud_vpcs.default.ids.0
}
data "alicloud_ehpc_gws_clusters" "default" {	
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
