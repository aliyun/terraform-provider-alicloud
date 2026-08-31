package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func SkipTestAccAlicloudCassandraZonesDataSource_basic(t *testing.T) {
	// Cassandra has been offline
	t.Skip("Cassandra has been offline")
	rand := acctest.RandInt()
	resourceId := "data.alicloud_cassandra_zones.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceCassandraZonesConfigDependence)

	multiConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"multi": "true",
		}),
	}

	var existCassandraZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   CHECKSET,
			"zones.#": CHECKSET,
		}
	}

	var fakeCassandraZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var cassandraZonesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCassandraZonesMapFunc,
		fakeMapFunc:  fakeCassandraZonesMapFunc,
	}

	cassandraZonesCheckInfo.dataSourceTestCheck(t, rand, multiConfig)
}

func dataSourceCassandraZonesConfigDependence(name string) string {
	return ""
}
