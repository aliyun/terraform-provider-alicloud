package polardbx_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudDrdsPolardbxInstancesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_drds_polardbx_instances.default"
	name := fmt.Sprintf("tfaccdrdsds%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceDrdsPolardbxInstancesConfigDependence)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_drds_polardbx_instance.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_drds_polardbx_instance.default.id}-fake"},
		}),
	}

	descriptionRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_drds_polardbx_instance.default.description}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_drds_polardbx_instance.default.description}-fake",
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_drds_polardbx_instance.default.id}"},
			"status": "Running",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_drds_polardbx_instance.default.id}"},
			"status": "Deleting",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_drds_polardbx_instance.default.id}"},
			"description_regex": "${alicloud_drds_polardbx_instance.default.description}",
			"status":            "Running",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_drds_polardbx_instance.default.id}"},
			"description_regex": "${alicloud_drds_polardbx_instance.default.description}-fake",
			"status":            "Running",
		}),
	}

	var existMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"descriptions.#":                   "1",
			"ids.0":                            acctest.CHECKSET,
			"instances.#":                      "1",
			"instances.0.id":                   acctest.CHECKSET,
			"instances.0.polardbx_instance_id": acctest.CHECKSET,
			"instances.0.description":          strings.ToLower(name),
			"instances.0.cn_class":             acctest.CHECKSET,
			"instances.0.dn_class":             acctest.CHECKSET,
			"instances.0.cn_node_count":        acctest.CHECKSET,
			"instances.0.dn_node_count":        acctest.CHECKSET,
			"instances.0.topology_type":        acctest.CHECKSET,
			"instances.0.status":               "Running",
			"instances.0.primary_zone":         acctest.CHECKSET,
			"instances.0.vpc_id":               acctest.CHECKSET,
			"instances.0.create_time":          acctest.CHECKSET,
			"instances.0.region_id":            acctest.CHECKSET,
		}
	}

	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"descriptions.#": "0",
			"instances.#":    "0",
		}
	}

	var checkInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existMapFunc,
		FakeMapFunc:  fakeMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
	}

	checkInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, descriptionRegexConf, statusConf, allConf)
}

func dataSourceDrdsPolardbxInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-beijing-f"
}

resource "alicloud_drds_polardbx_instance" "default" {
  topology_type = "1azone"
  vswitch_id    = data.alicloud_vswitches.default.ids.0
  primary_zone  = "cn-beijing-f"
  cn_node_count = 2
  dn_class      = "mysql.n4.medium.25"
  cn_class      = "polarx.x4.medium.2e"
  dn_node_count = 2
  vpc_id        = data.alicloud_vpcs.default.ids.0
  description   = lower(var.name)
}
`, name)
}
