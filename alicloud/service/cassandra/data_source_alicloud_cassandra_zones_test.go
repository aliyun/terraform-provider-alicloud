package cassandra_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func SkipTestAccAlicloudCassandraZonesDataSource_basic(t *testing.T) {
	// Cassandra has been offline
	t.Skip("Cassandra has been offline")
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_cassandra_zones.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, "", dataSourceCassandraZonesConfigDependence)

	multiConfig := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"multi": "true",
		}),
	}

	var existCassandraZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   acctest.CHECKSET,
			"zones.#": acctest.CHECKSET,
		}
	}

	var fakeCassandraZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var cassandraZonesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existCassandraZonesMapFunc,
		FakeMapFunc:  fakeCassandraZonesMapFunc,
	}

	cassandraZonesCheckInfo.DataSourceTestCheck(t, rand, multiConfig)
}

func dataSourceCassandraZonesConfigDependence(name string) string {
	return ""
}
