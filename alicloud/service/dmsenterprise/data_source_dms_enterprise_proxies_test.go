package dmsenterprise_test

import (
	"fmt"
	"strings"
	"testing"

	tfacctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudDmsEnterpriseProxiesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	tfacctest.CheckoutSupportedRegions(t, true, connectivity.DMSEnterpriseSupportRegions)
	idsConf := tfacctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDmsEnterpriseProxiesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dms_enterprise_proxy.default.id}"]`,
		}),
		FakeConfig: testAccCheckAlicloudDmsEnterpriseProxiesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dms_enterprise_proxy.default.id}_fake"]`,
		}),
	}
	allConf := tfacctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudDmsEnterpriseProxiesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dms_enterprise_proxy.default.id}"]`,
			"tid": `"${data.alicloud_dms_user_tenants.default.ids.0}"`,
		}),
		FakeConfig: testAccCheckAlicloudDmsEnterpriseProxiesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_dms_enterprise_proxy.default.id}_fake"]`,
		}),
	}
	var existAlicloudDmsEnterpriseProxiesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"proxies.#":                "1",
			"proxies.0.instance_id":    tfacctest.CHECKSET,
			"proxies.0.creator_id":     tfacctest.CHECKSET,
			"proxies.0.creator_name":   tfacctest.CHECKSET,
			"proxies.0.https_port":     tfacctest.CHECKSET,
			"proxies.0.private_enable": tfacctest.CHECKSET,
			"proxies.0.private_host":   tfacctest.CHECKSET,
			"proxies.0.protocol_port":  tfacctest.CHECKSET,
			"proxies.0.protocol_type":  tfacctest.CHECKSET,
			"proxies.0.id":             tfacctest.CHECKSET,
			"proxies.0.proxy_id":       tfacctest.CHECKSET,
			"proxies.0.public_enable":  tfacctest.CHECKSET,
			"proxies.0.public_host":    tfacctest.CHECKSET,
		}
	}
	var fakeAlicloudDmsEnterpriseProxiesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudDmsEnterpriseProxiesCheckInfo = tfacctest.DataSourceAttr{
		ResourceId:   "data.alicloud_dms_enterprise_proxies.default",
		ExistMapFunc: existAlicloudDmsEnterpriseProxiesDataSourceNameMapFunc,
		FakeMapFunc:  fakeAlicloudDmsEnterpriseProxiesDataSourceNameMapFunc,
	}

	preCheck := func() {
		tfacctest.PreCheck(t)
	}
	alicloudDmsEnterpriseProxiesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, allConf)
}
func testAccCheckAlicloudDmsEnterpriseProxiesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccProxy-%d"
}

data "alicloud_account" "current" {}

data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "8.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
	zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "8.0"
	category = "HighAvailability"
	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
	name = var.name
	vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_db_instance" "instance" {
	engine = "MySQL"
	engine_version = "8.0"
	db_instance_storage_type = "cloud_essd"
	instance_type = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
	vswitch_id       = data.alicloud_vswitches.default.ids.0
	instance_name    = var.name
	security_ips     = ["100.104.5.0/24","192.168.0.6"]
	tags = {
		"key1" = "value1"
		"key2" = "value2"
	}
}

resource "alicloud_db_account" "account" {
	instance_id = "${alicloud_db_instance.instance.id}"
	name        = "tftestnormal"
	password    = "Test12345"
	type        = "Normal"
}

data "alicloud_dms_user_tenants" "default" {
	status = "ACTIVE"
}
resource "alicloud_dms_enterprise_instance" "default" {
  dba_uid           =  tonumber(data.alicloud_account.current.id)
  host              =  "${alicloud_db_instance.instance.connection_string}"
  port              =  "3306"
  network_type      =	 "VPC"
  safe_rule         =	"自由操作"
  tid               =  data.alicloud_dms_user_tenants.default.ids.0
  instance_type     =	 "mysql"
  instance_source   =	 "RDS"
  env_type          =	 "test"
  database_user     =	 alicloud_db_account.account.name
  database_password =	 alicloud_db_account.account.password
  instance_alias    =	 var.name
  query_timeout     =	 "70"
  export_timeout    =	 "2000"
  ecs_region        =	 "%s"
  sell_trust        =	 "true"
}

resource "alicloud_dms_enterprise_proxy" "default" {
	instance_id = alicloud_dms_enterprise_instance.default.instance_id
	password = "Test12345"
	username = "tftestnormal"
	tid = data.alicloud_dms_user_tenants.default.ids.0
}

data "alicloud_dms_enterprise_proxies" "default" {	
	%s
}
`, rand, tfacctest.DefaultRegionToTest, strings.Join(pairs, " \n "))
	return config
}
