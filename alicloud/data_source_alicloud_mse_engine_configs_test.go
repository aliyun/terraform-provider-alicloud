package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudMSEEngineConfigsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.MSESupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMseEngineConfigsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mse_engine_config.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudMseEngineConfigsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mse_engine_config.default.id}_fake"]`,
		}),
	}
	namespaceConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMseEngineConfigsDataSourceNamespace(rand, map[string]string{
			"ids": `["${alicloud_mse_engine_config.default.id}"]`,
		}),
	}

	var existAlicloudMseEngineConfigsDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":              "1",
			"configs.#":          "1",
			"configs.0.type":     CHECKSET,
			"configs.0.app_name": "test",
			//"configs.0.tags":               NOSET,
			"configs.0.md5":     CHECKSET,
			"configs.0.data_id": fmt.Sprintf("tf-testAccEngineConfig-%d", rand),
			"configs.0.content": "test",
			"configs.0.group":   fmt.Sprintf("tf-testAccEngineConfig-%d", rand),
			"configs.0.desc":    "test",
			//"configs.0.encrypted_data_key": REMOVEKEY,
			//"configs.0.beta_ips":         ,
		}
	}
	var fakeAlicloudMseEngineConfigsDataSourceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudMseEngineConfigsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_mse_engine_configs.default",
		existMapFunc: existAlicloudMseEngineConfigsDataSourceMapFunc,
		fakeMapFunc:  fakeAlicloudMseEngineConfigsDataSourceMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudMseEngineConfigsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, namespaceConf)

}
func testAccCheckAlicloudMseEngineConfigsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAccEngineConfig-%d"
	}

	data "alicloud_mse_clusters" "tf" {
	  name_regex = "tf.*"
	}

	resource "alicloud_mse_engine_config" "default" {
      instance_id = data.alicloud_mse_clusters.tf.clusters.0.instance_id
	  data_id =  var.name
      group = var.name
      content ="test"
      app_name="test"
	  desc="test"
      type="text"
	  namespace_id="default"
	}

	data "alicloud_mse_engine_configs" "default"{
      instance_id =  data.alicloud_mse_clusters.tf.clusters.0.instance_id
      enable_details = "true"
      namespace_id="default"
      %s
	 
	}

`, rand, strings.Join(pairs, " \n "))
	return config

}

func testAccCheckAlicloudMseEngineConfigsDataSourceNamespace(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAccEngineConfig-%d"
	}

	data "alicloud_mse_clusters" "tf" {
	  name_regex = "tf.*"
	}

	resource "alicloud_mse_engine_config" "default" {
      instance_id = data.alicloud_mse_clusters.tf.clusters.0.instance_id
	  data_id =  var.name
      group = var.name
      content ="test"
      app_name="test"
	  desc="test"
      type="text"
		 
	}

	data "alicloud_mse_engine_configs" "default"{
      instance_id =  data.alicloud_mse_clusters.tf.clusters.0.instance_id
      enable_details = "true"
      namespace_id=""
      %s
	}

`, rand, strings.Join(pairs, " \n "))
	return config

}
