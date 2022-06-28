package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudMongodbAuditPoliciesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.MongoDBSupportRegions)
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongodbAuditPoliciesDataSourceName(rand, map[string]string{
			"db_instance_id": `"${alicloud_mongodb_audit_policy.default.db_instance_id}"`,
		}),
		fakeConfig: "",
	}
	var existAlicloudMongodbAuditPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"policies.#":                "1",
			"policies.0.id":             CHECKSET,
			"policies.0.db_instance_id": CHECKSET,
			"policies.0.audit_status":   "enable",
		}
	}
	var fakeAlicloudMongodbAuditPoliciesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"policies.#": "0",
		}
	}
	var alicloudMongodbAuditPoliciesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_mongodb_audit_policies.default",
		existMapFunc: existAlicloudMongodbAuditPoliciesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudMongodbAuditPoliciesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudMongodbAuditPoliciesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)
}
func testAccCheckAlicloudMongodbAuditPoliciesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccAuditPolicy-%d"
}

data "alicloud_mongodb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_mongodb_zones.default.zones.0.id
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_mongodb_zones.default.zones.0.id
  vswitch_name      = "subnet-for-local-test"
}

resource "alicloud_mongodb_instance" "default" {
  engine_version      = "3.4"
  db_instance_class   = "dds.mongo.mid"
  db_instance_storage = 10
  name                = var.name
  vswitch_id          = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}

resource "alicloud_mongodb_audit_policy" "default" {
  db_instance_id = alicloud_mongodb_instance.default.id
  audit_status   = "enable"
}

data "alicloud_mongodb_audit_policies" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
