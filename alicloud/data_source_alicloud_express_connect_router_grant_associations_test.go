package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAliCloudExpressConnectRouterGrantAssociationsDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_express_connect_router_grant_associations.default"
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouter%d", defaultRegionToTest, rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceExpressConnectRouterGrantAssociationsConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_grant_association.default.ecr_id}",
			"ids":    []string{"${alicloud_express_connect_router_grant_association.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_grant_association.default.ecr_id}",
			"ids":    []string{"${alicloud_express_connect_router_grant_association.default.id}_fake"},
		}),
	}

	instanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id":      "${alicloud_express_connect_router_grant_association.default.ecr_id}",
			"instance_id": "${alicloud_express_connect_router_grant_association.default.instance_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":      "${alicloud_express_connect_router_grant_association.default.ecr_id}",
			"instance_id": "${alicloud_express_connect_router_grant_association.default.instance_id}_fake",
		}),
	}

	instanceRegionIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id":             "${alicloud_express_connect_router_grant_association.default.ecr_id}",
			"instance_region_id": "${alicloud_express_connect_router_grant_association.default.instance_region_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":             "${alicloud_express_connect_router_grant_association.default.ecr_id}",
			"instance_region_id": "${alicloud_express_connect_router_grant_association.default.instance_region_id}_fake",
		}),
	}

	instanceTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id":        "${alicloud_express_connect_router_grant_association.default.ecr_id}",
			"instance_type": "${alicloud_express_connect_router_grant_association.default.instance_type}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":        "${alicloud_express_connect_router_grant_association.default.ecr_id}",
			"instance_type": "VBR",
		}),
	}

	callerTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id":      "${alicloud_express_connect_router_grant_association.default.id}",
			"caller_type": "ECR",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":      "${alicloud_express_connect_router_grant_association.default.id}",
			"caller_type": "OTHER",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id":             "${alicloud_express_connect_router_grant_association.default.ecr_id}",
			"ids":                []string{"${alicloud_express_connect_router_grant_association.default.id}"},
			"instance_id":        "${alicloud_express_connect_router_grant_association.default.instance_id}",
			"instance_region_id": "${alicloud_express_connect_router_grant_association.default.instance_region_id}",
			"instance_type":      "${alicloud_express_connect_router_grant_association.default.instance_type}",
			"caller_type":        "ECR",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":             "${alicloud_express_connect_router_grant_association.default.ecr_id}",
			"ids":                []string{"${alicloud_express_connect_router_grant_association.default.id}_fake"},
			"instance_id":        "${alicloud_express_connect_router_grant_association.default.instance_id}_fake",
			"instance_region_id": "${alicloud_express_connect_router_grant_association.default.instance_region_id}_fake",
			"instance_type":      "VBR",
			"caller_type":        "OTHER",
		}),
	}

	var existAliCloudExpressConnectRouterGrantAssociationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"associations.#":                    "1",
			"associations.0.id":                 CHECKSET,
			"associations.0.ecr_id":             CHECKSET,
			"associations.0.instance_id":        CHECKSET,
			"associations.0.instance_region_id": CHECKSET,
			"associations.0.owner_id":           CHECKSET,
			"associations.0.instance_type":      CHECKSET,
			"associations.0.grant_id":           CHECKSET,
			"associations.0.instance_owner_id":  CHECKSET,
			"associations.0.instance_owner_bid": CHECKSET,
			"associations.0.status":             CHECKSET,
			"associations.0.create_time":        CHECKSET,
			"associations.0.modify_time":        CHECKSET,
		}
	}

	var fakeAliCloudExpressConnectRouterGrantAssociationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"associations.#": "0",
		}
	}

	var aliCloudExpressConnectRouterGrantAssociationsInfo = dataSourceAttr{
		resourceId:   "data.alicloud_express_connect_router_grant_associations.default",
		existMapFunc: existAliCloudExpressConnectRouterGrantAssociationsMapFunc,
		fakeMapFunc:  fakeAliCloudExpressConnectRouterGrantAssociationsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudExpressConnectRouterGrantAssociationsInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, instanceIdConf, instanceRegionIdConf, instanceTypeConf, callerTypeConf, allConf)
}

func dataSourceExpressConnectRouterGrantAssociationsConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

variable "vpc_id" {
  # You need to modify this value to an existing VPC under your account
  default = "vpc-xxx"
}

variable "ecr_owner_uid" {
  # You need to modify this value to ecr owner ali uid
  default = "18xxx"
}

variable "ecr_id" {
  # You need to modify this value to an existing ecr id
  default = "ecr-xxx"
}

variable "region" {
  default = "cn-hangzhou"
}

resource "alicloud_express_connect_router_grant_association" "default" {
  ecr_id             = var.ecr_id
  instance_region_id = var.region
  instance_id        = var.vpc_id
  ecr_owner_ali_uid  = var.ecr_owner_uid
  instance_type      = "VPC"
}
`, name)
}
