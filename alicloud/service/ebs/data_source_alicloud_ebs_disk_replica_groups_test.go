package ebs_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudEbsDiskReplicaGroupsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	acctest.CheckoutSupportedRegions(t, true, connectivity.EBSSupportRegions)
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudEbsDiskReplicaGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ebs_disk_replica_group.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudEbsDiskReplicaGroupsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_ebs_disk_replica_group.default.id}_fake"]`,
		}),
	}
	var existAlicloudEbsDiskReplicaGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          "1",
			"groups.#":                       "1",
			"groups.0.id":                    acctest.CHECKSET,
			"groups.0.description":           fmt.Sprintf("tf-testAccEbsDiskReplicaGroup%d", rand),
			"groups.0.enable_rtc":            acctest.CHECKSET,
			"groups.0.destination_region_id": acctest.CHECKSET,
			"groups.0.destination_zone_id":   acctest.CHECKSET,
			"groups.0.group_name":            fmt.Sprintf("tf-testAccEbsDiskReplicaGroup%d", rand),
			"groups.0.primary_region":        acctest.CHECKSET,
			"groups.0.primary_zone":          acctest.CHECKSET,
			"groups.0.site":                  acctest.CHECKSET,
			"groups.0.rpo":                   "900",
			"groups.0.replica_group_id":      acctest.CHECKSET,
			"groups.0.source_region_id":      acctest.CHECKSET,
			"groups.0.source_zone_id":        acctest.CHECKSET,
			"groups.0.standby_region":        acctest.CHECKSET,
			"groups.0.standby_zone":          acctest.CHECKSET,
			"groups.0.status":                acctest.CHECKSET,
		}
	}
	var fakeAlicloudEbsDiskReplicaGroupsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"groups.#": "0",
		}
	}
	var alicloudEbsDiskReplicaGroupsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_ebs_disk_replica_groups.default",
		ExistMapFunc: existAlicloudEbsDiskReplicaGroupsDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudEbsDiskReplicaGroupsDataSourceNameMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}
	alicloudEbsDiskReplicaGroupsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf)
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

`, rand, acctest.DefaultRegionToTest, strings.Join(pairs, "\n   "))
	return config
}
