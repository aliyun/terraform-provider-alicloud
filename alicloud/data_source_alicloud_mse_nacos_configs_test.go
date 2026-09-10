package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudMseNacosConfigsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.MSESupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMseNacosConfigsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mse_nacos_config.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudMseNacosConfigsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mse_nacos_config.default.id}_fake"]`,
		}),
	}
	var existAlicloudMseNacosConfigsDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":              "1",
			"configs.#":          "1",
			"configs.0.type":     CHECKSET,
			"configs.0.app_name": "test",
			"configs.0.md5":      CHECKSET,
			"configs.0.data_id":  fmt.Sprintf("tf-testAccNacosConfig-%d", rand),
			"configs.0.content":  "test",
			"configs.0.group":    fmt.Sprintf("tf-testAccNacosConfig-%d", rand),
			"configs.0.desc":     "test",
		}
	}
	var fakeAlicloudMseNacosConfigsDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var AlicloudMseNacosConfigsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_mse_nacos_configs.default",
		existMapFunc: existAlicloudMseNacosConfigsDataSourceMapFunc,
		fakeMapFunc:  fakeAlicloudMseNacosConfigsDataSourceMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	AlicloudMseNacosConfigsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf)

}
func testAccCheckAlicloudMseNacosConfigsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAccNacosConfig-%d"
	}

	data "alicloud_zones" "example" {
	  available_resource_creation = "VSwitch"
	}

	resource "alicloud_vpc" "example" {
	  vpc_name   = "terraform-example"
	  cidr_block = "172.17.3.0/24"
	}

	resource "alicloud_vswitch" "example" {
	  vswitch_name = "terraform-example"
	  cidr_block   = "172.17.3.0/24"
	  vpc_id       = alicloud_vpc.example.id
	  zone_id      = data.alicloud_zones.example.zones.0.id
	}

	resource "alicloud_mse_cluster" "example" {
	  connection_type       = "slb"
	  net_type              = "privatenet"
	  vswitch_id            = alicloud_vswitch.example.id
	  cluster_specification = "MSE_SC_1_2_60_c"
	  cluster_version       = "NACOS_2_0_0"
	  instance_count        = "3"
	  pub_network_flow      = "1"
	  cluster_alias_name    = var.name
	  mse_version           = "mse_pro"
	  cluster_type          = "Nacos-Ans"
	}

	resource "alicloud_mse_engine_namespace" "example" {
	  instance_id          = alicloud_mse_cluster.example.id
	  namespace_show_name  = var.name
	  namespace_id         = var.name
	}

	resource "alicloud_mse_nacos_config" "default" {
      instance_id = alicloud_mse_engine_namespace.example.instance_id
	  data_id =  var.name
      group = var.name
      namespace_id=alicloud_mse_engine_namespace.example.namespace_id
      content ="test"
      app_name="test"
	  desc="test"
      type="text"
	}

	data "alicloud_mse_nacos_configs" "default"{
      instance_id =  alicloud_mse_nacos_config.default.instance_id
      enable_details = "true"
      namespace_id = alicloud_mse_nacos_config.default.namespace_id
      %s
	}

`, rand, strings.Join(pairs, " \n "))
	return config

}
