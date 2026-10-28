package ots_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudOtsTunnelsDataSource_basic(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_ots_tunnels.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("testAcc%d", rand),
		dataSourceOtsTunnelsConfigDependence)

	instanceNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_tunnel.default.instance_name}",
			"table_name":    "${alicloud_ots_tunnel.default.table_name}",
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_tunnel.default.instance_name}",
			"table_name":    "${alicloud_ots_tunnel.default.table_name}",
			"ids":           []string{"${alicloud_ots_tunnel.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_tunnel.default.instance_name}",
			"table_name":    "${alicloud_ots_tunnel.default.table_name}",
			"ids":           []string{"${alicloud_ots_tunnel.default.id}-fake"},
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_tunnel.default.instance_name}",
			"table_name":    "${alicloud_ots_tunnel.default.table_name}",
			"name_regex":    "${alicloud_ots_tunnel.default.tunnel_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_tunnel.default.instance_name}",
			"table_name":    "${alicloud_ots_tunnel.default.table_name}",
			"name_regex":    "${alicloud_ots_tunnel.default.tunnel_name}-fake",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_tunnel.default.instance_name}",
			"table_name":    "${alicloud_ots_tunnel.default.table_name}",
			"ids":           []string{"${alicloud_ots_tunnel.default.id}"},
			"name_regex":    "${alicloud_ots_tunnel.default.tunnel_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_tunnel.default.instance_name}",
			"table_name":    "${alicloud_ots_tunnel.default.table_name}",
			"ids":           []string{"${alicloud_ots_tunnel.default.id}"},
			"name_regex":    "${alicloud_ots_tunnel.default.tunnel_name}-fake",
		}),
	}

	var existOtsTunnelsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                             "1",
			"names.0":                             acctest.CHECKSET,
			"tunnels.#":                           "1",
			"tunnels.0.instance_name":             acctest.CHECKSET,
			"tunnels.0.table_name":                acctest.CHECKSET,
			"tunnels.0.tunnel_name":               acctest.CHECKSET,
			"tunnels.0.tunnel_id":                 acctest.CHECKSET,
			"tunnels.0.tunnel_rpo":                "0",
			"tunnels.0.tunnel_type":               "BaseAndStream",
			"tunnels.0.tunnel_stage":              "ProcessBaseData",
			"tunnels.0.create_time":               acctest.CHECKSET,
			"tunnels.0.expired":                   "false",
			"tunnels.0.channels.#":                "2",
			"tunnels.0.channels.0.channel_id":     acctest.CHECKSET,
			"tunnels.0.channels.0.channel_type":   acctest.CHECKSET,
			"tunnels.0.channels.0.channel_status": acctest.CHECKSET,
			"tunnels.0.channels.0.client_id":      "",
			"tunnels.0.channels.0.channel_rpo":    "0",
			"tunnels.0.channels.1.channel_id":     acctest.CHECKSET,
			"tunnels.0.channels.1.channel_type":   acctest.CHECKSET,
			"tunnels.0.channels.1.channel_status": acctest.CHECKSET,
			"tunnels.0.channels.1.client_id":      "",
			"tunnels.0.channels.1.channel_rpo":    "0",
		}
	}

	var fakeOtsTunnelsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":  "0",
			"tables.#": "0",
		}
	}

	var otsTunnelsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existOtsTunnelsMapFunc,
		FakeMapFunc:  fakeOtsTunnelsMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
	}
	otsTunnelsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, instanceNameConf, idsConf, nameRegexConf, allConf)
}

func dataSourceOtsTunnelsConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}
	resource "alicloud_ots_instance" "default" {
	  name = "tf-${var.name}"
	  description = "${var.name}"
	  accessed_by = "Any"
	  instance_type = "Capacity"
	  tags = {
	    Created = "TF"
	    For = "acceptance test"
	  }
	}

	resource "alicloud_ots_table" "default" {
	  instance_name = "${alicloud_ots_instance.default.name}"
	  table_name = "${var.name}"
	  primary_key {
          name = "pk1"
	      type = "Integer"
	  }
	  primary_key {
          name = "pk2"
          type = "String"
      }
	  time_to_live = -1
	  max_version = 1
	}

	resource "alicloud_ots_tunnel" "default" {
	 instance_name = "${alicloud_ots_instance.default.name}"
	 table_name = "${alicloud_ots_table.default.table_name}"
	 tunnel_name = "${var.name}"
	 tunnel_type = "BaseAndStream"
	}
	`, name)
}
