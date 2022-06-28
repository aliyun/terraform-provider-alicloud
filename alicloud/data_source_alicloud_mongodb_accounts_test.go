package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudMongodbAccountsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.MongoDBSupportRegions)
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongodbAccountsDataSourceName(rand, map[string]string{
			"instance_id":  `"${alicloud_mongodb_account.default.instance_id}"`,
			"account_name": `"${alicloud_mongodb_account.default.account_name}"`,
		}),
		fakeConfig: "",
	}
	var existAlicloudMongodbAccountsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"accounts.#":                     "1",
			"accounts.0.id":                  CHECKSET,
			"accounts.0.account_name":        "root",
			"accounts.0.instance_id":         CHECKSET,
			"accounts.0.account_description": fmt.Sprintf("tf-testAccAccount-%d", rand),
			"accounts.0.character_type":      CHECKSET,
			"accounts.0.status":              "Available",
		}
	}
	var fakeAlicloudMongodbAccountsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"accounts.#": "0",
		}
	}
	var alicloudMongodbAccountsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_mongodb_accounts.default",
		existMapFunc: existAlicloudMongodbAccountsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudMongodbAccountsDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudMongodbAccountsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)
}
func testAccCheckAlicloudMongodbAccountsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccAccount-%d"
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

resource "alicloud_mongodb_account" "default" {
  account_name = "root"
  account_password = "YourPassword+12345"
  instance_id = alicloud_mongodb_instance.default.id
  account_description = var.name
}

data "alicloud_mongodb_accounts" "default" {
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
