package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudDdoscooPortsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1, 65535)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDdoscooPortsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ddoscoo_port.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudDdoscooPortsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ddoscoo_port.default.id}_fake"]`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudDdoscooPortsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_ddoscoo_port.default.id}"]`,
			"frontend_port":     fmt.Sprintf(`"%d"`, rand+1),
			"frontend_protocol": `"tcp"`,
		}),
		fakeConfig: testAccCheckAlicloudDdoscooPortsDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_ddoscoo_port.default.id}_fake"]`,
			"frontend_port":     fmt.Sprintf(`"%d"`, rand+1),
			"frontend_protocol": `"tcp"`,
		}),
	}
	var existAlicloudDdoscooPortsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"ports.#":                   "1",
			"ports.0.backend_port":      fmt.Sprintf("%d", rand),
			"ports.0.frontend_port":     fmt.Sprintf("%d", rand+1),
			"ports.0.frontend_protocol": `tcp`,
			"ports.0.instance_id":       CHECKSET,
			"ports.0.id":                CHECKSET,
			"ports.0.real_servers.#":    `1`,
		}
	}
	var fakeAlicloudDdoscooPortsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudDdoscooPortsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ddoscoo_ports.default",
		existMapFunc: existAlicloudDdoscooPortsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudDdoscooPortsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.DdoscooSupportedRegions)
	}
	alicloudDdoscooPortsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, allConf)
}
func testAccCheckAlicloudDdoscooPortsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccPort-%d"
}

data "alicloud_ddoscoo_instances" "default" {}

resource "alicloud_ddoscoo_port" "default" {
	backend_port = "%d"
	frontend_port = "%d"
	instance_id = data.alicloud_ddoscoo_instances.default.ids.0
	frontend_protocol = "tcp"
	real_servers = ["192.168.0.1"]
}

data "alicloud_ddoscoo_ports" "default" {	
	instance_id = alicloud_ddoscoo_port.default.instance_id
	%s
}
`, rand, rand, rand+1, strings.Join(pairs, " \n "))
	return config
}
