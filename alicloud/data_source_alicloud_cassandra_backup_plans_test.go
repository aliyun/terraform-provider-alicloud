package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCassandraBackupPlansDataSource(t *testing.T) {
	// Cassandra has been offline
	t.Skip("Cassandra has been offline")
	rand := acctest.RandInt()

	var existAlicloudCassandraBackupPlanDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                  "1",
			"plans.data_center_id":   CHECKSET,
			"plans.cluster_id":       CHECKSET,
			"plans.backup_time":      "00:10Z",
			"plans.active":           "false",
			"plans.create_time":      CHECKSET,
			"plans.backup_period":    CHECKSET,
			"plans.retention_period": CHECKSET,
		}
	}
	var fakeAlicloudCassandraBackupPlanDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}

	var alicloudCassandraBackupPlanCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cassandra_backup_plans.default_2",
		existMapFunc: existAlicloudCassandraBackupPlanDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudCassandraBackupPlanDataSourceNameMapFunc,
	}

	alicloudCassandraBackupPlanCheckInfo.dataSourceTestCheck(t, rand)
}
