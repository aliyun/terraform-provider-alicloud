package mongodb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudMongoDBInstancesDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 9999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mongodb_instance.default.id}"]`,
		}),
		FakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mongodb_instance.default.id}_fake"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_mongodb_instance.default.name}"`,
		}),
		FakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_mongodb_instance.default.name}_fake"`,
		}),
	}

	instanceClassConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_mongodb_instance.default.id}"]`,
			"instance_class": `"${alicloud_mongodb_instance.default.db_instance_class}"`,
		}),
		FakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"instance_class": `"${alicloud_mongodb_instance.default.db_instance_class}_fake"`,
		}),
	}

	availabilityZoneConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_mongodb_instance.default.id}"]`,
			"availability_zone": `"${alicloud_mongodb_instance.default.zone_id}"`,
		}),
		FakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"availability_zone": `"${alicloud_mongodb_instance.default.zone_id}_fake"`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_mongodb_instance.default.id}"]`,
			"status": `"Running"`,
		}),
		FakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"status": `"Deleting"`,
		}),
	}

	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"tags": `{
						Created = "TF"
						For 	= "Instance"
					  }`,
		}),
		FakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"tags": `{
						Created = "TF_Update"
						For 	= "Instance_Update"
					  }`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
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
		FakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
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
			"instances.0.id":                acctest.CHECKSET,
			"instances.0.engine":            acctest.CHECKSET,
			"instances.0.engine_version":    acctest.CHECKSET,
			"instances.0.instance_type":     acctest.CHECKSET,
			"instances.0.instance_class":    acctest.CHECKSET,
			"instances.0.storage":           acctest.CHECKSET,
			"instances.0.network_type":      acctest.CHECKSET,
			"instances.0.availability_zone": acctest.CHECKSET,
			"instances.0.name":              acctest.CHECKSET,
			"instances.0.charge_type":       acctest.CHECKSET,
			"instances.0.replication":       acctest.CHECKSET,
			"instances.0.lock_mode":         acctest.CHECKSET,
			"instances.0.region_id":         acctest.CHECKSET,
			"instances.0.creation_time":     acctest.CHECKSET,
			"instances.0.expiration_time":   acctest.CHECKSET,
			"instances.0.status":            acctest.CHECKSET,
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

	var aliCloudMongoDBInstancesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_mongodb_instances.default",
		ExistMapFunc: existAliCloudMongoDBInstancesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAliCloudMongoDBInstancesDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
	}

	aliCloudMongoDBInstancesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, instanceClassConf, availabilityZoneConf, statusConf, tagsConf, allConf)
}

func TestAccAliCloudMongoDBInstancesDataSource_basic1(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 9999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_mongodb_instance.default.id}"]`,
			"enable_details": `true`,
		}),
		FakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_mongodb_instance.default.id}"]`,
			"enable_details": `false`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_mongodb_instance.default.name}"`,
			"enable_details": `true`,
		}),
		FakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_mongodb_instance.default.name}"`,
			"enable_details": `false`,
		}),
	}

	instanceClassConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_mongodb_instance.default.id}"]`,
			"instance_class": `"${alicloud_mongodb_instance.default.db_instance_class}"`,
			"enable_details": `true`,
		}),
		FakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_mongodb_instance.default.id}"]`,
			"instance_class": `"${alicloud_mongodb_instance.default.db_instance_class}"`,
			"enable_details": `false`,
		}),
	}

	availabilityZoneConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_mongodb_instance.default.id}"]`,
			"availability_zone": `"${alicloud_mongodb_instance.default.zone_id}"`,
			"enable_details":    `true`,
		}),
		FakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_mongodb_instance.default.id}"]`,
			"availability_zone": `"${alicloud_mongodb_instance.default.zone_id}"`,
			"enable_details":    `false`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_mongodb_instance.default.id}"]`,
			"status":         `"Running"`,
			"enable_details": `true`,
		}),
		FakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_mongodb_instance.default.id}"]`,
			"status":         `"Running"`,
			"enable_details": `false`,
		}),
	}

	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"tags": `{
						Created = "TF"
						For 	= "Instance"
					  }`,
			"enable_details": `true`,
		}),
		FakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
			"tags": `{
						Created = "TF"
						For 	= "Instance"
					  }`,
			"enable_details": `false`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
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
		FakeConfig: testAccCheckAliCloudMongoDBInstancesDataSourceName(rand, map[string]string{
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
			"instances.0.id":                acctest.CHECKSET,
			"instances.0.engine":            acctest.CHECKSET,
			"instances.0.engine_version":    acctest.CHECKSET,
			"instances.0.instance_type":     acctest.CHECKSET,
			"instances.0.instance_class":    acctest.CHECKSET,
			"instances.0.storage":           acctest.CHECKSET,
			"instances.0.network_type":      acctest.CHECKSET,
			"instances.0.availability_zone": acctest.CHECKSET,
			"instances.0.name":              acctest.CHECKSET,
			"instances.0.charge_type":       acctest.CHECKSET,
			"instances.0.replication":       acctest.CHECKSET,
			"instances.0.lock_mode":         acctest.CHECKSET,
			"instances.0.region_id":         acctest.CHECKSET,
			"instances.0.creation_time":     acctest.CHECKSET,
			"instances.0.expiration_time":   acctest.CHECKSET,
			"instances.0.status":            acctest.CHECKSET,
			"instances.0.tags.%":            "2",
			"instances.0.tags.Created":      "TF",
			"instances.0.tags.For":          "Instance",
			"instances.0.mongos.#":          "0",
			"instances.0.shards.#":          "0",
			"instances.0.restore_ranges.#":  acctest.CHECKSET,
		}
	}

	var fakeAliCloudMongoDBInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"instances.#":                   "1",
			"instances.0.id":                acctest.CHECKSET,
			"instances.0.engine":            acctest.CHECKSET,
			"instances.0.engine_version":    acctest.CHECKSET,
			"instances.0.instance_type":     acctest.CHECKSET,
			"instances.0.instance_class":    acctest.CHECKSET,
			"instances.0.storage":           acctest.CHECKSET,
			"instances.0.network_type":      acctest.CHECKSET,
			"instances.0.availability_zone": acctest.CHECKSET,
			"instances.0.name":              acctest.CHECKSET,
			"instances.0.charge_type":       acctest.CHECKSET,
			"instances.0.replication":       acctest.CHECKSET,
			"instances.0.lock_mode":         acctest.CHECKSET,
			"instances.0.region_id":         acctest.CHECKSET,
			"instances.0.creation_time":     acctest.CHECKSET,
			"instances.0.expiration_time":   acctest.CHECKSET,
			"instances.0.status":            acctest.CHECKSET,
			"instances.0.tags.%":            "2",
			"instances.0.tags.Created":      "TF",
			"instances.0.tags.For":          "Instance",
			"instances.0.mongos.#":          "0",
			"instances.0.shards.#":          "0",
			"instances.0.restore_ranges.#":  "0",
		}
	}

	var aliCloudMongoDBInstancesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_mongodb_instances.default",
		ExistMapFunc: existAliCloudMongoDBInstancesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAliCloudMongoDBInstancesDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
	}

	aliCloudMongoDBInstancesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, instanceClassConf, availabilityZoneConf, statusConf, tagsConf, allConf)
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
	rand := tfacctest.RandIntRange(1000, 9999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mongodb_sharding_instance.default.id}"]`,
		}),
		FakeConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mongodb_sharding_instance.default.id}_fake"]`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_mongodb_sharding_instance.default.name}"`,
		}),
		FakeConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_mongodb_sharding_instance.default.name}_fake"`,
		}),
	}

	availabilityZoneConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_mongodb_sharding_instance.default.id}"]`,
			"availability_zone": `"${alicloud_mongodb_sharding_instance.default.zone_id}"`,
		}),
		FakeConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"availability_zone": `"${alicloud_mongodb_sharding_instance.default.zone_id}_fake"`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_mongodb_sharding_instance.default.id}"]`,
			"status": `"Running"`,
		}),
		FakeConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"status": `"Deleting"`,
		}),
	}

	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"tags": `{
						Created = "TF"
						For 	= "ShardingInstance"
					  }`,
		}),
		FakeConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"tags": `{
						Created = "TF_Update"
						For 	= "ShardingInstance_Update"
					  }`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_mongodb_sharding_instance.default.id}"]`,
			"name_regex":        `"${alicloud_mongodb_sharding_instance.default.name}"`,
			"availability_zone": `"${alicloud_mongodb_sharding_instance.default.zone_id}"`,
			"status":            `"Running"`,
			"tags": `{
						Created = "TF"
						For 	= "ShardingInstance"
					  }`,
		}),
		FakeConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
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
			"instances.0.id":                acctest.CHECKSET,
			"instances.0.engine":            acctest.CHECKSET,
			"instances.0.engine_version":    acctest.CHECKSET,
			"instances.0.instance_type":     acctest.CHECKSET,
			"instances.0.instance_class":    acctest.CHECKSET,
			"instances.0.storage":           acctest.CHECKSET,
			"instances.0.network_type":      acctest.CHECKSET,
			"instances.0.availability_zone": acctest.CHECKSET,
			"instances.0.name":              acctest.CHECKSET,
			"instances.0.charge_type":       acctest.CHECKSET,
			"instances.0.lock_mode":         acctest.CHECKSET,
			"instances.0.region_id":         acctest.CHECKSET,
			"instances.0.creation_time":     acctest.CHECKSET,
			"instances.0.expiration_time":   acctest.CHECKSET,
			"instances.0.status":            acctest.CHECKSET,
			"instances.0.tags.%":            "2",
			"instances.0.tags.Created":      "TF",
			"instances.0.tags.For":          "ShardingInstance",
			"instances.0.mongos.#":          acctest.CHECKSET,
			"instances.0.shards.#":          acctest.CHECKSET,
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

	var aliCloudMongoDBShardingInstancesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_mongodb_instances.default",
		ExistMapFunc: existAliCloudMongoDBShardingInstancesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAliCloudMongoDBShardingInstancesDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
	}

	aliCloudMongoDBShardingInstancesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, availabilityZoneConf, statusConf, tagsConf, allConf)
}

func TestAccAliCloudMongoDBShardingInstancesDataSource_basic1(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 9999)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_mongodb_sharding_instance.default.id}"]`,
			"enable_details": `true`,
		}),
		FakeConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_mongodb_sharding_instance.default.id}"]`,
			"enable_details": `false`,
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_mongodb_sharding_instance.default.name}"`,
			"enable_details": `true`,
		}),
		FakeConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_mongodb_sharding_instance.default.name}"`,
			"enable_details": `false`,
		}),
	}

	availabilityZoneConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_mongodb_sharding_instance.default.id}"]`,
			"availability_zone": `"${alicloud_mongodb_sharding_instance.default.zone_id}"`,
			"enable_details":    `true`,
		}),
		FakeConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_mongodb_sharding_instance.default.id}"]`,
			"availability_zone": `"${alicloud_mongodb_sharding_instance.default.zone_id}"`,
			"enable_details":    `false`,
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_mongodb_sharding_instance.default.id}"]`,
			"status":         `"Running"`,
			"enable_details": `true`,
		}),
		FakeConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_mongodb_sharding_instance.default.id}"]`,
			"status":         `"Running"`,
			"enable_details": `false`,
		}),
	}

	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"tags": `{
						Created = "TF"
						For 	= "ShardingInstance"
					  }`,
			"enable_details": `true`,
		}),
		FakeConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
			"tags": `{
						Created = "TF"
						For 	= "ShardingInstance"
					  }`,
			"enable_details": `false`,
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{
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
		FakeConfig: testAccCheckAliCloudMongoDBShardingInstancesDataSourceName(rand, map[string]string{

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
			"instances.0.id":                acctest.CHECKSET,
			"instances.0.engine":            acctest.CHECKSET,
			"instances.0.engine_version":    acctest.CHECKSET,
			"instances.0.instance_type":     acctest.CHECKSET,
			"instances.0.instance_class":    acctest.CHECKSET,
			"instances.0.storage":           acctest.CHECKSET,
			"instances.0.network_type":      acctest.CHECKSET,
			"instances.0.availability_zone": acctest.CHECKSET,
			"instances.0.name":              acctest.CHECKSET,
			"instances.0.charge_type":       acctest.CHECKSET,
			"instances.0.lock_mode":         acctest.CHECKSET,
			"instances.0.region_id":         acctest.CHECKSET,
			"instances.0.creation_time":     acctest.CHECKSET,
			"instances.0.expiration_time":   acctest.CHECKSET,
			"instances.0.status":            acctest.CHECKSET,
			"instances.0.tags.%":            "2",
			"instances.0.tags.Created":      "TF",
			"instances.0.tags.For":          "ShardingInstance",
			"instances.0.mongos.#":          acctest.CHECKSET,
			"instances.0.shards.#":          acctest.CHECKSET,
			"instances.0.restore_ranges.#":  acctest.CHECKSET,
		}
	}

	var fakeAliCloudMongoDBShardingInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"instances.#":                   "1",
			"instances.0.id":                acctest.CHECKSET,
			"instances.0.engine":            acctest.CHECKSET,
			"instances.0.engine_version":    acctest.CHECKSET,
			"instances.0.instance_type":     acctest.CHECKSET,
			"instances.0.instance_class":    acctest.CHECKSET,
			"instances.0.storage":           acctest.CHECKSET,
			"instances.0.network_type":      acctest.CHECKSET,
			"instances.0.availability_zone": acctest.CHECKSET,
			"instances.0.name":              acctest.CHECKSET,
			"instances.0.charge_type":       acctest.CHECKSET,
			"instances.0.lock_mode":         acctest.CHECKSET,
			"instances.0.region_id":         acctest.CHECKSET,
			"instances.0.creation_time":     acctest.CHECKSET,
			"instances.0.expiration_time":   acctest.CHECKSET,
			"instances.0.status":            acctest.CHECKSET,
			"instances.0.tags.%":            "2",
			"instances.0.tags.Created":      "TF",
			"instances.0.tags.For":          "ShardingInstance",
			"instances.0.mongos.#":          acctest.CHECKSET,
			"instances.0.shards.#":          acctest.CHECKSET,
			"instances.0.restore_ranges.#":  "0",
		}
	}

	var aliCloudMongoDBShardingInstancesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_mongodb_instances.default",
		ExistMapFunc: existAliCloudMongoDBShardingInstancesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAliCloudMongoDBShardingInstancesDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
	}

	aliCloudMongoDBShardingInstancesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, availabilityZoneConf, statusConf, tagsConf, allConf)
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
