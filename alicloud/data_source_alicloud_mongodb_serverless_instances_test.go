package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudMongodbServerlessInstancesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.MongoDBServerlessSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongodbServerlessInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mongodb_serverless_instance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudMongodbServerlessInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mongodb_serverless_instance.default.id}_fake"]`,
		}),
	}
	dBInstanceClassConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongodbServerlessInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_mongodb_serverless_instance.default.id}"]`,
			"db_instance_class": `"dds.serverless.cu"`,
		}),
		fakeConfig: testAccCheckAlicloudMongodbServerlessInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_mongodb_serverless_instance.default.id}"]`,
			"db_instance_class": `"dds.serverless.cu_fake"`,
		}),
	}
	dBInstanceDescriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongodbServerlessInstancesDataSourceName(rand, map[string]string{
			"ids":                     `["${alicloud_mongodb_serverless_instance.default.id}"]`,
			"db_instance_description": `"${alicloud_mongodb_serverless_instance.default.db_instance_description}"`,
		}),
		fakeConfig: testAccCheckAlicloudMongodbServerlessInstancesDataSourceName(rand, map[string]string{
			"ids":                     `["${alicloud_mongodb_serverless_instance.default.id}"]`,
			"db_instance_description": `"${alicloud_mongodb_serverless_instance.default.db_instance_description}_fake"`,
		}),
	}
	networkTypeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongodbServerlessInstancesDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_mongodb_serverless_instance.default.id}"]`,
			"network_type": `"VPC"`,
		}),
		fakeConfig: testAccCheckAlicloudMongodbServerlessInstancesDataSourceName(rand, map[string]string{
			"ids":          `["${alicloud_mongodb_serverless_instance.default.id}"]`,
			"network_type": `"Classic"`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongodbServerlessInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_mongodb_serverless_instance.default.id}"]`,
			"resource_group_id": `"${alicloud_mongodb_serverless_instance.default.resource_group_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudMongodbServerlessInstancesDataSourceName(rand, map[string]string{
			"ids":               `["${alicloud_mongodb_serverless_instance.default.id}"]`,
			"resource_group_id": `"${alicloud_mongodb_serverless_instance.default.resource_group_id}_fake"`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongodbServerlessInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mongodb_serverless_instance.default.id}"]`,
			"tags": `{
				"Created" = "MongodbServerlessInstance"
				"For" = "TF"
			}`,
		}),
		fakeConfig: testAccCheckAlicloudMongodbServerlessInstancesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mongodb_serverless_instance.default.id}"]`,
			"tags": `{
				"Created" = "MongodbServerlessInstance1"
				"For" = "TF1"
			}`,
		}),
	}
	vSwitchIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongodbServerlessInstancesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_mongodb_serverless_instance.default.id}"]`,
			"vswitch_id": `"${alicloud_mongodb_serverless_instance.default.vswitch_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudMongodbServerlessInstancesDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_mongodb_serverless_instance.default.id}"]`,
			"vswitch_id": `"${alicloud_mongodb_serverless_instance.default.vswitch_id}_fake"`,
		}),
	}
	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongodbServerlessInstancesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_mongodb_serverless_instance.default.id}"]`,
			"vpc_id": `"${alicloud_mongodb_serverless_instance.default.vpc_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudMongodbServerlessInstancesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_mongodb_serverless_instance.default.id}"]`,
			"vpc_id": `"${alicloud_mongodb_serverless_instance.default.vpc_id}_fake"`,
		}),
	}
	zoneIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongodbServerlessInstancesDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_mongodb_serverless_instance.default.id}"]`,
			"zone_id": `"${alicloud_mongodb_serverless_instance.default.zone_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudMongodbServerlessInstancesDataSourceName(rand, map[string]string{
			"ids":     `["${alicloud_mongodb_serverless_instance.default.id}"]`,
			"zone_id": `"${alicloud_mongodb_serverless_instance.default.zone_id}_fake"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongodbServerlessInstancesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_mongodb_serverless_instance.default.id}"]`,
			"status": `"Running"`,
		}),
		fakeConfig: testAccCheckAlicloudMongodbServerlessInstancesDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_mongodb_serverless_instance.default.id}"]`,
			"status": `"Creating"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMongodbServerlessInstancesDataSourceName(rand, map[string]string{
			"db_instance_class":       `"dds.serverless.cu"`,
			"db_instance_description": `"${alicloud_mongodb_serverless_instance.default.db_instance_description}"`,
			"ids":                     `["${alicloud_mongodb_serverless_instance.default.id}"]`,
			"network_type":            `"VPC"`,
			"resource_group_id":       `"${alicloud_mongodb_serverless_instance.default.resource_group_id}"`,
			"status":                  `"Running"`,
			"tags": `{
				"Created" = "MongodbServerlessInstance"
				"For" = "TF"
			}`,
			"vpc_id":     `"${alicloud_mongodb_serverless_instance.default.vpc_id}"`,
			"vswitch_id": `"${alicloud_mongodb_serverless_instance.default.vswitch_id}"`,
			"zone_id":    `"${alicloud_mongodb_serverless_instance.default.zone_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudMongodbServerlessInstancesDataSourceName(rand, map[string]string{
			"db_instance_class":       `"dds.serverless.cu_fake"`,
			"db_instance_description": `"${alicloud_mongodb_serverless_instance.default.db_instance_description}_fake"`,
			"ids":                     `["${alicloud_mongodb_serverless_instance.default.id}_fake"]`,
			"network_type":            `"Classic"`,
			"resource_group_id":       `"${alicloud_mongodb_serverless_instance.default.resource_group_id}_fake"`,
			"status":                  `"Creating"`,
			"tags": `{
				"Created" = "MongodbServerlessInstance1"
				"For" = "TF1"
			}`,
			"vpc_id":     `"${alicloud_mongodb_serverless_instance.default.vpc_id}_fake"`,
			"vswitch_id": `"${alicloud_mongodb_serverless_instance.default.vswitch_id}_fake"`,
			"zone_id":    `"${alicloud_mongodb_serverless_instance.default.zone_id}_fake"`,
		}),
	}
	var existAlicloudMongodbServerlessInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                                        "1",
			"instances.#":                                                  "1",
			"instances.0.capacity_unit":                                    "100",
			"instances.0.payment_type":                                     "Subscription",
			"instances.0.db_instance_class":                                CHECKSET,
			"instances.0.db_instance_description":                          fmt.Sprintf("tf-testacc-mongodbserverlessinstance-%d", rand),
			"instances.0.id":                                               CHECKSET,
			"instances.0.db_instance_id":                                   CHECKSET,
			"instances.0.db_instance_release_protection":                   CHECKSET,
			"instances.0.db_instance_storage":                              "5",
			"instances.0.engine":                                           "MongoDB",
			"instances.0.engine_version":                                   "4.2",
			"instances.0.expire_time":                                      CHECKSET,
			"instances.0.kind_code":                                        "0",
			"instances.0.lock_mode":                                        "Unlock",
			"instances.0.maintain_end_time":                                CHECKSET,
			"instances.0.maintain_start_time":                              CHECKSET,
			"instances.0.max_connections":                                  CHECKSET,
			"instances.0.max_iops":                                         CHECKSET,
			"instances.0.network_type":                                     "VPC",
			"instances.0.protocol_type":                                    CHECKSET,
			"instances.0.resource_group_id":                                CHECKSET,
			"instances.0.security_ip_groups.#":                             "1",
			"instances.0.security_ip_groups.0.security_ip_group_attribute": "test",
			"instances.0.security_ip_groups.0.security_ip_group_name":      "test",
			"instances.0.security_ip_groups.0.security_ip_list":            "192.168.0.1",
			"instances.0.status":                                           "Running",
			"instances.0.storage_engine":                                   "WiredTiger",
			"instances.0.tags.%":                                           "2",
			"instances.0.tags.Created":                                     "MongodbServerlessInstance",
			"instances.0.tags.For":                                         "TF",
			"instances.0.vswitch_id":                                       CHECKSET,
			"instances.0.vpc_auth_mode":                                    "",
			"instances.0.vpc_id":                                           CHECKSET,
			"instances.0.zone_id":                                          CHECKSET,
		}
	}
	var fakeAlicloudMongodbServerlessInstancesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudMongodbServerlessInstancesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_mongodb_serverless_instances.default",
		existMapFunc: existAlicloudMongodbServerlessInstancesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudMongodbServerlessInstancesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudMongodbServerlessInstancesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, dBInstanceClassConf, dBInstanceDescriptionConf, networkTypeConf, resourceGroupIdConf, tagsConf, vSwitchIdConf, vpcIdConf, zoneIdConf, statusConf, allConf)
}
func testAccCheckAlicloudMongodbServerlessInstancesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testacc-mongodbserverlessinstance-%d"
}

data "alicloud_mongodb_zones" "default" {}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_mongodb_zones.default.zones.0.id
}

resource "alicloud_vswitch" "default" {
    count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
    vswitch_name = var.name
    vpc_id       = data.alicloud_vpcs.default.ids.0
    zone_id      = data.alicloud_mongodb_zones.default.zones.0.id
    cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
}

resource "alicloud_mongodb_serverless_instance" "default" {
  account_password        = "Abc12345"
  db_instance_description = var.name
  db_instance_storage     = 5
  capacity_unit           = 100
  engine                  = "MongoDB"
  engine_version          = "4.2"
  period                  = 1
  period_price_type       = "Month"
  vpc_id                  = data.alicloud_vpcs.default.ids.0
  zone_id                 = data.alicloud_mongodb_zones.default.zones.0.id
  vswitch_id              = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.default.*.id, [""])[0]
  tags = {
    Created = "MongodbServerlessInstance"
    For     = "TF"
  }
  security_ip_groups {
    security_ip_group_attribute = "test"
    security_ip_group_name      = "test"
    security_ip_list            = "192.168.0.1"
  }
}

data "alicloud_mongodb_serverless_instances" "default" {	
	enable_details = true
	%s	
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
