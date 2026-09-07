package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudOtsTunnelsDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_ots_tunnels.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("testAcc%d", rand),
		dataSourceOtsTunnelsConfigDependence)

	instanceNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_tunnel.default.instance_name}",
			"table_name":    "${alicloud_ots_tunnel.default.table_name}",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_tunnel.default.instance_name}",
			"table_name":    "${alicloud_ots_tunnel.default.table_name}",
			"ids":           []string{"${alicloud_ots_tunnel.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_tunnel.default.instance_name}",
			"table_name":    "${alicloud_ots_tunnel.default.table_name}",
			"ids":           []string{"${alicloud_ots_tunnel.default.id}-fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_tunnel.default.instance_name}",
			"table_name":    "${alicloud_ots_tunnel.default.table_name}",
			"name_regex":    "${alicloud_ots_tunnel.default.tunnel_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_tunnel.default.instance_name}",
			"table_name":    "${alicloud_ots_tunnel.default.table_name}",
			"name_regex":    "${alicloud_ots_tunnel.default.tunnel_name}-fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_tunnel.default.instance_name}",
			"table_name":    "${alicloud_ots_tunnel.default.table_name}",
			"ids":           []string{"${alicloud_ots_tunnel.default.id}"},
			"name_regex":    "${alicloud_ots_tunnel.default.tunnel_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_name": "${alicloud_ots_tunnel.default.instance_name}",
			"table_name":    "${alicloud_ots_tunnel.default.table_name}",
			"ids":           []string{"${alicloud_ots_tunnel.default.id}"},
			"name_regex":    "${alicloud_ots_tunnel.default.tunnel_name}-fake",
		}),
	}

	var existOtsTunnelsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                             "1",
			"names.0":                             CHECKSET,
			"tunnels.#":                           "1",
			"tunnels.0.instance_name":             CHECKSET,
			"tunnels.0.table_name":                CHECKSET,
			"tunnels.0.tunnel_name":               CHECKSET,
			"tunnels.0.tunnel_id":                 CHECKSET,
			"tunnels.0.tunnel_rpo":                "0",
			"tunnels.0.tunnel_type":               "BaseAndStream",
			"tunnels.0.tunnel_stage":              "ProcessBaseData",
			"tunnels.0.create_time":               CHECKSET,
			"tunnels.0.expired":                   "false",
			"tunnels.0.channels.#":                "2",
			"tunnels.0.channels.0.channel_id":     CHECKSET,
			"tunnels.0.channels.0.channel_type":   CHECKSET,
			"tunnels.0.channels.0.channel_status": CHECKSET,
			"tunnels.0.channels.0.client_id":      "",
			"tunnels.0.channels.0.channel_rpo":    "0",
			"tunnels.0.channels.1.channel_id":     CHECKSET,
			"tunnels.0.channels.1.channel_type":   CHECKSET,
			"tunnels.0.channels.1.channel_status": CHECKSET,
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

	var otsTunnelsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existOtsTunnelsMapFunc,
		fakeMapFunc:  fakeOtsTunnelsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
	}
	otsTunnelsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, instanceNameConf, idsConf, nameRegexConf, allConf)
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
