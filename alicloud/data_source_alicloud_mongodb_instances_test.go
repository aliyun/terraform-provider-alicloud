package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudMongoDBInstancesDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mongodb_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mongodb_instance.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_mongodb_instance.default.name}"`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_mongodb_instance.default.name}_fake"`,
		}),
	}

	instanceClassConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_mongodb_instance.default.id}"]`,
			"instance_class": `"${alicloud_mongodb_instance.default.db_instance_class}"`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"instance_class": `"${alicloud_mongodb_instance.default.db_instance_class}_fake"`,
		}),
	}

	availabilityZoneConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_mongodb_instance.default.id}"]`,
			"availability_zone": `"${alicloud_mongodb_instance.default.zone_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"availability_zone": `"${alicloud_mongodb_instance.default.zone_id}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_mongodb_instance.default.id}"]`,
			"status": `"Running"`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"status": `"Deleting"`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"tags": `{
						Created = "TF"
						For 	= "Instance"
					  }`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"tags": `{
						Created = "TF_Update"
						For 	= "Instance_Update"
					  }`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_mongodb_instance.default.id}"]`,
			"name_regex":        `"${alicloud_mongodb_instance.default.name}"`,
			"instance_class":    `"${alicloud_mongodb_instance.default.db_instance_class}"`,
			"availability_zone": `"${alicloud_mongodb_instance.default.zone_id}"`,
			"status":            `"Running"`,
			"tags": `{
						Created = "TF"
						For 	= "Instance"
					  }`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_mongodb_instance.default.id}_fake"]`,
			"name_regex":        `"${alicloud_mongodb_instance.default.name}_fake"`,
			"instance_class":    `"${alicloud_mongodb_instance.default.db_instance_class}_fake"`,
			"availability_zone": `"${alicloud_mongodb_instance.default.zone_id}_fake"`,
			"status":            `"Deleting"`,
			"tags": `{
						Created = "TF_Update"
						For 	= "Instance_Update"
					  }`,
		}),
	}

	var existAliCloudMongoDBInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"instances.#":                   "1",
			"instances.0.id":                CHECKSET,
			"instances.0.engine":            CHECKSET,
			"instances.0.engine_version":    CHECKSET,
			"instances.0.instance_type":     CHECKSET,
			"instances.0.instance_class":    CHECKSET,
			"instances.0.storage":           CHECKSET,
			"instances.0.network_type":      CHECKSET,
			"instances.0.availability_zone": CHECKSET,
			"instances.0.name":              CHECKSET,
			"instances.0.charge_type":       CHECKSET,
			"instances.0.replication":       CHECKSET,
			"instances.0.lock_mode":         CHECKSET,
			"instances.0.region_id":         CHECKSET,
			"instances.0.creation_time":     CHECKSET,
			"instances.0.expiration_time":   CHECKSET,
			"instances.0.status":            CHECKSET,
			"instances.0.tags.%":            "2",
			"instances.0.tags.Created":      "TF",
			"instances.0.tags.For":          "Instance",
			"instances.0.mongos.#":          "0",
			"instances.0.shards.#":          "0",
			"instances.0.restore_ranges.#":  "0",
		}
	}

	var fakeAliCloudMongoDBInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"instances.#": "0",
		}
	}

	var aliCloudMongoDBInstancesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_mongodb_instances.default",
		existMapFunc: existAliCloudMongoDBInstancesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudMongoDBInstancesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	aliCloudMongoDBInstancesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, instanceClassConf, availabilityZoneConf, statusConf, tagsConf, allConf)
}

func TestAccAliCloudMongoDBInstancesDataSource_basic1(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_mongodb_instance.default.id}"]`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_mongodb_instance.default.id}"]`,
			"enable_details": `false`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_mongodb_instance.default.name}"`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_mongodb_instance.default.name}"`,
			"enable_details": `false`,
		}),
	}

	instanceClassConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_mongodb_instance.default.id}"]`,
			"instance_class": `"${alicloud_mongodb_instance.default.db_instance_class}"`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_mongodb_instance.default.id}"]`,
			"instance_class": `"${alicloud_mongodb_instance.default.db_instance_class}"`,
			"enable_details": `false`,
		}),
	}

	availabilityZoneConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_mongodb_instance.default.id}"]`,
			"availability_zone": `"${alicloud_mongodb_instance.default.zone_id}"`,
			"enable_details":    `true`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_mongodb_instance.default.id}"]`,
			"availability_zone": `"${alicloud_mongodb_instance.default.zone_id}"`,
			"enable_details":    `false`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_mongodb_instance.default.id}"]`,
			"status":         `"Running"`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_mongodb_instance.default.id}"]`,
			"status":         `"Running"`,
			"enable_details": `false`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"tags": `{
						Created = "TF"
						For 	= "Instance"
					  }`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"tags": `{
						Created = "TF"
						For 	= "Instance"
					  }`,
			"enable_details": `false`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_mongodb_instance.default.id}"]`,
			"name_regex":        `"${alicloud_mongodb_instance.default.name}"`,
			"instance_class":    `"${alicloud_mongodb_instance.default.db_instance_class}"`,
			"availability_zone": `"${alicloud_mongodb_instance.default.zone_id}"`,
			"status":            `"Running"`,
			"tags": `{
						Created = "TF"
						For 	= "Instance"
					  }`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_mongodb_instance.default.id}"]`,
			"name_regex":        `"${alicloud_mongodb_instance.default.name}"`,
			"instance_class":    `"${alicloud_mongodb_instance.default.db_instance_class}"`,
			"availability_zone": `"${alicloud_mongodb_instance.default.zone_id}"`,
			"status":            `"Running"`,
			"tags": `{
						Created = "TF"
						For 	= "Instance"
					  }`,
			"enable_details": `false`,
		}),
	}

	var existAliCloudMongoDBInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"instances.#":                   "1",
			"instances.0.id":                CHECKSET,
			"instances.0.engine":            CHECKSET,
			"instances.0.engine_version":    CHECKSET,
			"instances.0.instance_type":     CHECKSET,
			"instances.0.instance_class":    CHECKSET,
			"instances.0.storage":           CHECKSET,
			"instances.0.network_type":      CHECKSET,
			"instances.0.availability_zone": CHECKSET,
			"instances.0.name":              CHECKSET,
			"instances.0.charge_type":       CHECKSET,
			"instances.0.replication":       CHECKSET,
			"instances.0.lock_mode":         CHECKSET,
			"instances.0.region_id":         CHECKSET,
			"instances.0.creation_time":     CHECKSET,
			"instances.0.expiration_time":   CHECKSET,
			"instances.0.status":            CHECKSET,
			"instances.0.tags.%":            "2",
			"instances.0.tags.Created":      "TF",
			"instances.0.tags.For":          "Instance",
			"instances.0.mongos.#":          "0",
			"instances.0.shards.#":          "0",
			"instances.0.restore_ranges.#":  CHECKSET,
		}
	}

	var fakeAliCloudMongoDBInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"instances.#":                   "1",
			"instances.0.id":                CHECKSET,
			"instances.0.engine":            CHECKSET,
			"instances.0.engine_version":    CHECKSET,
			"instances.0.instance_type":     CHECKSET,
			"instances.0.instance_class":    CHECKSET,
			"instances.0.storage":           CHECKSET,
			"instances.0.network_type":      CHECKSET,
			"instances.0.availability_zone": CHECKSET,
			"instances.0.name":              CHECKSET,
			"instances.0.charge_type":       CHECKSET,
			"instances.0.replication":       CHECKSET,
			"instances.0.lock_mode":         CHECKSET,
			"instances.0.region_id":         CHECKSET,
			"instances.0.creation_time":     CHECKSET,
			"instances.0.expiration_time":   CHECKSET,
			"instances.0.status":            CHECKSET,
			"instances.0.tags.%":            "2",
			"instances.0.tags.Created":      "TF",
			"instances.0.tags.For":          "Instance",
			"instances.0.mongos.#":          "0",
			"instances.0.shards.#":          "0",
			"instances.0.restore_ranges.#":  "0",
		}
	}

	var aliCloudMongoDBInstancesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_mongodb_instances.default",
		existMapFunc: existAliCloudMongoDBInstancesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudMongoDBInstancesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	aliCloudMongoDBInstancesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, instanceClassConf, availabilityZoneConf, statusConf, tagsConf, allConf)
}

func testAccCheckAliCloudMongoDBInstancesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-MongoDBInstance-%d"
	}

	data "alicloud_mongodb_zones" "default" {
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "default-NODELETING"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_mongodb_zones.default.zones.0.id
	}

	data "alicloud_security_groups" "default" {
  		vpc_id = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_mongodb_instance" "default" {
  		engine_version      = "4.4"
  		db_instance_class   = "mdb.shard.2x.xlarge.d"
  		db_instance_storage = 20
  		vswitch_id          = data.alicloud_vswitches.default.ids.0
  		name                = var.name
  		tags = {
    		Created = "TF"
    		For     = "Instance"
  		}
	}

	data "alicloud_mongodb_instances" "default" {
		%s
	}
`, rand, strings.Join(pairs, " \n "))

	return config
}

func TestAccAliCloudMongoDBShardingInstancesDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mongodb_sharding_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mongodb_sharding_instance.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_mongodb_sharding_instance.default.name}"`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_mongodb_sharding_instance.default.name}_fake"`,
		}),
	}

	availabilityZoneConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_mongodb_sharding_instance.default.id}"]`,
			"availability_zone": `"${alicloud_mongodb_sharding_instance.default.zone_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"availability_zone": `"${alicloud_mongodb_sharding_instance.default.zone_id}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_mongodb_sharding_instance.default.id}"]`,
			"status": `"Running"`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"status": `"Deleting"`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"tags": `{
						Created = "TF"
						For 	= "ShardingInstance"
					  }`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"tags": `{
						Created = "TF_Update"
						For 	= "ShardingInstance_Update"
					  }`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_mongodb_sharding_instance.default.id}"]`,
			"name_regex":        `"${alicloud_mongodb_sharding_instance.default.name}"`,
			"availability_zone": `"${alicloud_mongodb_sharding_instance.default.zone_id}"`,
			"status":            `"Running"`,
			"tags": `{
						Created = "TF"
						For 	= "ShardingInstance"
					  }`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_mongodb_sharding_instance.default.id}_fake"]`,
			"name_regex":        `"${alicloud_mongodb_sharding_instance.default.name}_fake"`,
			"availability_zone": `"${alicloud_mongodb_sharding_instance.default.zone_id}_fake"`,
			"status":            `"Deleting"`,
			"tags": `{
						Created = "TF_Update"
						For 	= "ShardingInstance_Update"
					  }`,
		}),
	}

	var existAliCloudMongoDBShardingInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"instances.#":                   "1",
			"instances.0.id":                CHECKSET,
			"instances.0.engine":            CHECKSET,
			"instances.0.engine_version":    CHECKSET,
			"instances.0.instance_type":     CHECKSET,
			"instances.0.instance_class":    CHECKSET,
			"instances.0.storage":           CHECKSET,
			"instances.0.network_type":      CHECKSET,
			"instances.0.availability_zone": CHECKSET,
			"instances.0.name":              CHECKSET,
			"instances.0.charge_type":       CHECKSET,
			"instances.0.lock_mode":         CHECKSET,
			"instances.0.region_id":         CHECKSET,
			"instances.0.creation_time":     CHECKSET,
			"instances.0.expiration_time":   CHECKSET,
			"instances.0.status":            CHECKSET,
			"instances.0.tags.%":            "2",
			"instances.0.tags.Created":      "TF",
			"instances.0.tags.For":          "ShardingInstance",
			"instances.0.mongos.#":          CHECKSET,
			"instances.0.shards.#":          CHECKSET,
			"instances.0.restore_ranges.#":  "0",
		}
	}

	var fakeAliCloudMongoDBShardingInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"instances.#": "0",
		}
	}

	var aliCloudMongoDBShardingInstancesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_mongodb_instances.default",
		existMapFunc: existAliCloudMongoDBShardingInstancesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudMongoDBShardingInstancesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	aliCloudMongoDBShardingInstancesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, availabilityZoneConf, statusConf, tagsConf, allConf)
}

func TestAccAliCloudMongoDBShardingInstancesDataSource_basic1(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_mongodb_sharding_instance.default.id}"]`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_mongodb_sharding_instance.default.id}"]`,
			"enable_details": `false`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_mongodb_sharding_instance.default.name}"`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_mongodb_sharding_instance.default.name}"`,
			"enable_details": `false`,
		}),
	}

	availabilityZoneConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_mongodb_sharding_instance.default.id}"]`,
			"availability_zone": `"${alicloud_mongodb_sharding_instance.default.zone_id}"`,
			"enable_details":    `true`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_mongodb_sharding_instance.default.id}"]`,
			"availability_zone": `"${alicloud_mongodb_sharding_instance.default.zone_id}"`,
			"enable_details":    `false`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_mongodb_sharding_instance.default.id}"]`,
			"status":         `"Running"`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_mongodb_sharding_instance.default.id}"]`,
			"status":         `"Running"`,
			"enable_details": `false`,
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"tags": `{
						Created = "TF"
						For 	= "ShardingInstance"
					  }`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"tags": `{
						Created = "TF"
						For 	= "ShardingInstance"
					  }`,
			"enable_details": `false`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_mongodb_sharding_instance.default.id}"]`,
			"name_regex":        `"${alicloud_mongodb_sharding_instance.default.name}"`,
			"availability_zone": `"${alicloud_mongodb_sharding_instance.default.zone_id}"`,
			"status":            `"Running"`,
			"tags": `{
						Created = "TF"
						For 	= "ShardingInstance"
					  }`,
			"enable_details": `true`,
		}),
		fakeConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{

			"ids":               `["${alicloud_mongodb_sharding_instance.default.id}"]`,
			"name_regex":        `"${alicloud_mongodb_sharding_instance.default.name}"`,
			"availability_zone": `"${alicloud_mongodb_sharding_instance.default.zone_id}"`,
			"status":            `"Running"`,
			"tags": `{
						Created = "TF"
						For 	= "ShardingInstance"
					  }`,
			"enable_details": `false`,
		}),
	}

	var existAliCloudMongoDBShardingInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"instances.#":                   "1",
			"instances.0.id":                CHECKSET,
			"instances.0.engine":            CHECKSET,
			"instances.0.engine_version":    CHECKSET,
			"instances.0.instance_type":     CHECKSET,
			"instances.0.instance_class":    CHECKSET,
			"instances.0.storage":           CHECKSET,
			"instances.0.network_type":      CHECKSET,
			"instances.0.availability_zone": CHECKSET,
			"instances.0.name":              CHECKSET,
			"instances.0.charge_type":       CHECKSET,
			"instances.0.lock_mode":         CHECKSET,
			"instances.0.region_id":         CHECKSET,
			"instances.0.creation_time":     CHECKSET,
			"instances.0.expiration_time":   CHECKSET,
			"instances.0.status":            CHECKSET,
			"instances.0.tags.%":            "2",
			"instances.0.tags.Created":      "TF",
			"instances.0.tags.For":          "ShardingInstance",
			"instances.0.mongos.#":          CHECKSET,
			"instances.0.shards.#":          CHECKSET,
			"instances.0.restore_ranges.#":  CHECKSET,
		}
	}

	var fakeAliCloudMongoDBShardingInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"instances.#":                   "1",
			"instances.0.id":                CHECKSET,
			"instances.0.engine":            CHECKSET,
			"instances.0.engine_version":    CHECKSET,
			"instances.0.instance_type":     CHECKSET,
			"instances.0.instance_class":    CHECKSET,
			"instances.0.storage":           CHECKSET,
			"instances.0.network_type":      CHECKSET,
			"instances.0.availability_zone": CHECKSET,
			"instances.0.name":              CHECKSET,
			"instances.0.charge_type":       CHECKSET,
			"instances.0.lock_mode":         CHECKSET,
			"instances.0.region_id":         CHECKSET,
			"instances.0.creation_time":     CHECKSET,
			"instances.0.expiration_time":   CHECKSET,
			"instances.0.status":            CHECKSET,
			"instances.0.tags.%":            "2",
			"instances.0.tags.Created":      "TF",
			"instances.0.tags.For":          "ShardingInstance",
			"instances.0.mongos.#":          CHECKSET,
			"instances.0.shards.#":          CHECKSET,
			"instances.0.restore_ranges.#":  "0",
		}
	}

	var aliCloudMongoDBShardingInstancesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_mongodb_instances.default",
		existMapFunc: existAliCloudMongoDBShardingInstancesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudMongoDBShardingInstancesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	aliCloudMongoDBShardingInstancesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, availabilityZoneConf, statusConf, tagsConf, allConf)
}

func testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-MongoDBShardingInstance-%d"
	}

	data "alicloud_mongodb_zones" "default" {
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "172.17.3.0/24"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		cidr_block   = "172.17.3.0/24"
  		vpc_id       = alicloud_vpc.default.id
  		zone_id      = data.alicloud_mongodb_zones.default.zones.1.id
	}

	resource "alicloud_mongodb_sharding_instance" "default" {
  		engine_version = "4.2"
  		vswitch_id     = alicloud_vswitch.default.id
  		zone_id        = alicloud_vswitch.default.zone_id
  		name           = var.name
  		mongo_list {
    		node_class = "dds.mongos.mid"
  		}
  		mongo_list {
    		node_class = "dds.mongos.mid"
  		}
  		shard_list {
    		node_class   = "dds.shard.mid"
    		node_storage = "10"
  		}
  		shard_list {
    		node_class        = "dds.shard.standard"
    		node_storage      = "20"
    		readonly_replicas = "1"
  		}
  		tags = {
    		Created = "TF"
    		For     = "ShardingInstance"
  		}
	}

	data "alicloud_mongodb_instances" "default" {
  		instance_type = "sharding"
  		%s
	}
`, rand, strings.Join(pairs, " \n "))

	return config
}
