package drds_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudDRDSInstancesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(0, 9999)
	resourceId := "data.alicloud_drds_instances.default"
	name := fmt.Sprintf("tf-testAcc%s%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceDRDSInstancesConfigDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_drds_instance.default.description}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_drds_instance.default.description}-fake",
		}),
	}

	descriptionRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_drds_instance.default.description}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alicloud_drds_instance.default.description}-fake",
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_drds_instance.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_drds_instance.default.id}-fake"},
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "${alicloud_drds_instance.default.description}",
			"description_regex": "${alicloud_drds_instance.default.description}",
			"ids":               []string{"${alicloud_drds_instance.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "${alicloud_drds_instance.default.description}-fake",
			"description_regex": "${alicloud_drds_instance.default.description}-fake",
			"ids":               []string{"${alicloud_drds_instance.default.id}-fake"},
		}),
	}

	var existDRDSInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"descriptions.#":                "1",
			"ids.0":                         acctest.CHECKSET,
			"instances.#":                   "1",
			"instances.0.description":       name,
			"instances.0.type":              "PRIVATE",
			"instances.0.zone_id":           acctest.CHECKSET,
			"instances.0.id":                acctest.CHECKSET,
			"instances.0.network_type":      "VPC",
			"instances.0.create_time":       acctest.CHECKSET,
			"instances.0.connection_string": acctest.CHECKSET,
			"instances.0.port":              acctest.CHECKSET,
		}
	}

	var fakeDRDSInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"descriptions.#": "0",
			"instances.#":    "0",
		}
	}

	var drdsInstancesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existDRDSInstancesMapFunc,
		FakeMapFunc:  fakeDRDSInstancesMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.DrdsSupportedRegions)
	}

	drdsInstancesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, descriptionRegexConf, idsConf, allConf)
}

func dataSourceDRDSInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
 	variable "name" {
		default = "%s"
	}
	# DRDS 1.0 (classic, Type=1/V1) instance creation is only supported in
	# specific zones of cn-beijing. An existing running DRDS 1.0 instance
	# (drds.sn2.4c16g, drdsPost, V1) in this account confirms cn-beijing-h is a
	# valid DRDS 1.0 zone; the first VSwitch-capable zone (cn-beijing-f) is a
	# DRDS 2.0/polardbx zone and returns InternalError on classic CreateDrdsInstance.
	# ACC region is fixed to cn-beijing, so the zone is pinned to cn-beijing-h.
	resource "alicloud_vpc" "default" {
		name       = "tf-testAcc-vpc-${var.name}"
		cidr_block = "172.16.0.0/16"
	}
	resource "alicloud_vswitch" "default" {
		name       = "tf-testAcc-vsw-${var.name}"
		vpc_id     = alicloud_vpc.default.id
		cidr_block = "172.16.0.0/24"
		zone_id    = "cn-beijing-h"
	}
 	resource "alicloud_drds_instance" "default" {
  		description = "${var.name}"
  		zone_id = "${alicloud_vswitch.default.zone_id}"
  		instance_series = "drds.sn2.4c16g"
  		instance_charge_type = "PostPaid"
		vswitch_id = "${alicloud_vswitch.default.id}"
  		specification = "drds.sn2.4c16g.8c32g"
}
 `, name)
}
