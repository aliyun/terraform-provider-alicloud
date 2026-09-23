package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEbsDiskReplicaGroupsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.EBSSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEbsDiskReplicaGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ebs_disk_replica_group.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEbsDiskReplicaGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ebs_disk_replica_group.default.id}_fake"]`,
		}),
	}
	var existAlicloudEbsDiskReplicaGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"groups.#":                       "1",
			"groups.0.id":                    CHECKSET,
			"groups.0.description":           fmt.Sprintf("tf-testAccEbsDiskReplicaGroup%d", rand),
			"groups.0.destination_region_id": CHECKSET,
			"groups.0.destination_zone_id":   CHECKSET,
			"groups.0.group_name":            fmt.Sprintf("tf-testAccEbsDiskReplicaGroup%d", rand),
			"groups.0.primary_region":        CHECKSET,
			"groups.0.primary_zone":          CHECKSET,
			"groups.0.site":                  CHECKSET,
			"groups.0.rpo":                   "900",
			"groups.0.replica_group_id":      CHECKSET,
			"groups.0.source_region_id":      CHECKSET,
			"groups.0.source_zone_id":        CHECKSET,
			"groups.0.standby_region":        CHECKSET,
			"groups.0.standby_zone":          CHECKSET,
			"groups.0.status":                CHECKSET,
		}
	}
	var fakeAlicloudEbsDiskReplicaGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"groups.#": "0",
		}
	}
	var alicloudEbsDiskReplicaGroupsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ebs_disk_replica_groups.default",
		existMapFunc: existAlicloudEbsDiskReplicaGroupsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudEbsDiskReplicaGroupsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudEbsDiskReplicaGroupsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf)
}
func testAccCheckAlicloudEbsDiskReplicaGroupsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccEbsDiskReplicaGroup%d"
}

variable "region" {
  default = "%s"
}

data "alicloud_ebs_regions" "default" {
  region_id = var.region
}

resource "alicloud_ebs_disk_replica_group" "default" {
  source_region_id      = var.region
  source_zone_id        = data.alicloud_ebs_regions.default.regions[0].zones[0].zone_id
  destination_region_id = var.region
  destination_zone_id   = data.alicloud_ebs_regions.default.regions[0].zones[1].zone_id
  group_name            = var.name
  description           = var.name
  rpo                   = 900
}

data "alicloud_ebs_disk_replica_groups" "default" {	
	%s
}

`, rand, defaultRegionToTest, strings.Join(pairs, "\n   "))
	return config
}
