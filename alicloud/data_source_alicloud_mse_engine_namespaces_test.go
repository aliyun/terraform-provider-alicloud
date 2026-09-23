package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudMseEngineNamespacesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.MSESupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMseEngineNamespacesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mse_engine_namespace.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudMseEngineNamespacesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mse_engine_namespace.default.id}_fake"]`,
		}),
	}
	var existAlicloudMseEngineNamespacesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"namespaces.#":                     "1",
			"namespaces.0.namespace_desc":      "",
			"namespaces.0.namespace_show_name": "public",
			"namespaces.0.namespace_id":        "",
			"namespaces.0.service_count":       CHECKSET,
			"namespaces.0.quota":               CHECKSET,
			"namespaces.0.type":                CHECKSET,
			"namespaces.0.config_count":        CHECKSET,
			"namespaces.0.id":                  CHECKSET,
		}
	}
	var fakeAlicloudMseEngineNamespacesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "2",
		}
	}
	var alicloudMseEngineNamespacesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_mse_engine_namespaces.default",
		existMapFunc: existAlicloudMseEngineNamespacesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudMseEngineNamespacesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudMseEngineNamespacesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf)
}
func testAccCheckAlicloudMseEngineNamespacesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

	variable "name" {	
			default = "tf-testAccEngineNamespace-%d"
	}
	
	data "alicloud_zones" "default" {
	  available_resource_creation = "VSwitch"
	}
	
	resource "alicloud_vpc" "default" {
	  vpc_name   = "default"
	  cidr_block = "172.17.3.0/24"
	}
	
	resource "alicloud_vswitch" "default" {
	  vswitch_name = "default"
	  cidr_block   = "172.17.3.0/24"
	  vpc_id       = alicloud_vpc.default.id
	  zone_id      = data.alicloud_zones.default.zones.0.id
	}
	
	resource "alicloud_mse_cluster" "default" {
	  cluster_specification = "MSE_SC_1_2_60_c"
	  cluster_type          = "Nacos-Ans"
	  cluster_version       = "NACOS_2_0_0"
	  instance_count        = 3
	  net_type              = "privatenet"
	  vswitch_id            = alicloud_vswitch.default.id
	  connection_type       = "slb"
	  pub_network_flow      = "1"
	  mse_version           = "mse_pro"
	  vpc_id                = alicloud_vpc.default.id
	}

	resource "alicloud_mse_engine_namespace" "default" {
		instance_id = alicloud_mse_cluster.default.id
		namespace_show_name = var.name
		namespace_id = var.name
	}
	
	data "alicloud_mse_engine_namespaces" "default" {
	   instance_id = alicloud_mse_cluster.default.id
	}

`, rand)
	return config
}
