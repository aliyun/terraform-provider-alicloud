package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudMongoDBInstancesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongoDBDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_mongodb_instance.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudMongoDBDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_mongodb_instance.default.name}_fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongoDBDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_mongodb_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudMongoDBDataSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_mongodb_instance.default.id}_fake"]`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongoDBDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_mongodb_instance.default.name}"`,
			"tags":       `{Created = "TF"}`,
		}),
		fakeConfig: testAccCheckAlicloudMongoDBDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_mongodb_instance.default.name}"`,
			"tags":       `{Created = "TF1"}`,
		}),
	}

	instanceTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongoDBDataSourceConfig(rand, map[string]string{
			"name_regex":    `"${alicloud_mongodb_instance.default.name}"`,
			"instance_type": `"replicate"`,
		}),
		fakeConfig: testAccCheckAlicloudMongoDBDataSourceConfig(rand, map[string]string{
			"name_regex":    `"${alicloud_mongodb_instance.default.name}"`,
			"instance_type": `"sharding"`,
		}),
	}
	instanceClassConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongoDBDataSourceConfig(rand, map[string]string{
			"name_regex":     `"${alicloud_mongodb_instance.default.name}"`,
			"instance_class": `"dds.mongo.mid"`,
		}),
		fakeConfig: testAccCheckAlicloudMongoDBDataSourceConfig(rand, map[string]string{
			"name_regex":     `"${alicloud_mongodb_instance.default.name}"`,
			"instance_class": `"test.rds.mid"`,
		}),
	}
	availabilityZoneConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongoDBDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_mongodb_instance.default.name}"`,
			"availability_zone": `"${data.alicloud_zones.default.zones.1.id}"`,
		}),
		fakeConfig: testAccCheckAlicloudMongoDBDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_mongodb_instance.default.name}"`,
			"availability_zone": `"test_zone"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongoDBDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_mongodb_instance.default.name}"`,
			"ids":               `["${alicloud_mongodb_instance.default.id}"]`,
			"availability_zone": `"${data.alicloud_zones.default.zones.1.id}"`,
			"instance_type":     `"replicate"`,
			"instance_class":    `"dds.mongo.mid"`,
		}),
		fakeConfig: testAccCheckAlicloudMongoDBDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alicloud_mongodb_instance.default.name}_fake"`,
			"ids":               `["${alicloud_mongodb_instance.default.id}"]`,
			"availability_zone": `"${data.alicloud_zones.default.zones.1.id}"`,
			"instance_type":     `"replicate"`,
			"instance_class":    `"dds.mongo.mid"`,
		}),
	}

	var exisMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#":                 CHECKSET,
			"instances.0.name":            fmt.Sprintf("tf-testAccMongoDBInstance_datasource_%d", rand),
			"instances.0.instance_class":  "dds.mongo.mid",
			"instances.0.engine":          "MongoDB",
			"instances.0.engine_version":  "3.4",
			"instances.0.charge_type":     string(PostPaid),
			"instances.0.storage":         "10",
			"instances.0.instance_type":   "replicate",
			"instances.0.id":              CHECKSET,
			"instances.0.expiration_time": CHECKSET,
			"instances.0.creation_time":   CHECKSET,
			"instances.0.region_id":       CHECKSET,
			"instances.0.status":          CHECKSET,
			"instances.0.network_type":    CHECKSET,
			"instances.0.lock_mode":       CHECKSET,
			"instances.0.tags.%":          "2",
			"instances.0.tags.Created":    "TF",
			"instances.0.tags.For":        "acceptance test",
			"ids.#":                       "1",
			"ids.0":                       CHECKSET,
			"names.#":                     "1",
		}
	}
	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#": "0",
			"ids.#":       "0",
			"names.#":     "0",
		}
	}

	var CheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_mongodb_instances.default",
		existMapFunc: exisMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithNoDefaultVpc(t)
	}
	CheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, tagsConf, instanceTypeConf, instanceClassConf, availabilityZoneConf, allConf)
}

func testAccCheckAlicloudMongoDBDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
data "alicloud_zones" "default" {
  available_resource_creation = "MongoDB"
}

data "alicloud_vpcs" "default" {
	is_default = true
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = "${data.alicloud_zones.default.zones.1.id}"
}

variable "name" {
  default = "tf-testAccMongoDBInstance_datasource_%d"
}

resource "alicloud_mongodb_instance" "default" {
  engine_version      = "3.4"
  db_instance_class   = "dds.mongo.mid"
  db_instance_storage = 10
  name                = "${var.name}"
  vswitch_id          = data.alicloud_vswitches.default.ids.0
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}
data "alicloud_mongodb_instances" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
