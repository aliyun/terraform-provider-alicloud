package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func SkipTestAccAlicloudEbsDiskReplicaPairDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEbsDiskReplicaPairSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ebs_disk_replica_pair.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEbsDiskReplicaPairSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ebs_disk_replica_pair.default.id}_fake"]`,
		}),
	}

	siteConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEbsDiskReplicaPairSourceConfig(rand, map[string]string{
			"site": `"production"`,
		}),
		fakeConfig: testAccCheckAlicloudEbsDiskReplicaPairSourceConfig(rand, map[string]string{
			"site": `"backup"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEbsDiskReplicaPairSourceConfig(rand, map[string]string{
			"ids":  `["${alicloud_ebs_disk_replica_pair.default.id}"]`,
			"site": `"production"`,
		}),
		fakeConfig: testAccCheckAlicloudEbsDiskReplicaPairSourceConfig(rand, map[string]string{
			"ids":  `["${alicloud_ebs_disk_replica_pair.default.id}_fake"]`,
			"site": `"backup"`,
		}),
	}

	EbsDiskReplicaPairCheckInfo.dataSourceTestCheck(t, rand, idsConf, siteConf, allConf)
}

var existEbsDiskReplicaPairMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"pairs.#":                       "1",
		"pairs.0.id":                    CHECKSET,
		"pairs.0.bandwidth":             CHECKSET,
		"pairs.0.description":           fmt.Sprintf("tf-testAccEbsDiskReplicaPair%d", rand),
		"pairs.0.destination_disk_id":   CHECKSET,
		"pairs.0.destination_region_id": "cn-hangzhou-onebox-nebula",
		"pairs.0.destination_zone_id":   "cn-hangzhou-onebox-nebula-e",
		"pairs.0.disk_id":               CHECKSET,
		"pairs.0.pair_name":             fmt.Sprintf("tf-testAccEbsDiskReplicaPair%d", rand),
		"pairs.0.payment_type":          CHECKSET,
		"pairs.0.rpo":                   CHECKSET,
		"pairs.0.replica_pair_id":       CHECKSET,
		"pairs.0.resource_group_id":     CHECKSET,
		"pairs.0.source_zone_id":        "cn-hangzhou-onebox-nebula-b",
		"pairs.0.status":                "created",
	}
}

var fakeEbsDiskReplicaPairMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"pairs.#": "0",
	}
}

var EbsDiskReplicaPairCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ebs_disk_replica_pairs.default",
	existMapFunc: existEbsDiskReplicaPairMapFunc,
	fakeMapFunc:  fakeEbsDiskReplicaPairMapFunc,
}

func testAccCheckAlicloudEbsDiskReplicaPairSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccEbsDiskReplicaPair%d"
}


resource "alicloud_ebs_disk_replica_pair" "default" {
  destination_disk_id   = "d-iq8aj6hjpu0m98fyp583"
  destination_region_id = "cn-hangzhou-onebox-nebula"
  bandwidth             = 10240
  destination_zone_id   = "cn-hangzhou-onebox-nebula-e"
  source_zone_id        = "cn-hangzhou-onebox-nebula-b"
  disk_id               = "d-iq89bx99bisov1qetlx8"
  description           = var.name
  pair_name             = var.name
}

data "alicloud_ebs_disk_replica_pairs" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
