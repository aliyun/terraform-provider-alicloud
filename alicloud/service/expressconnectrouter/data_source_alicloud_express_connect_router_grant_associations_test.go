package expressconnectrouter_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudExpressConnectRouterGrantAssociationsDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_express_connect_router_grant_associations.default"
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouter%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceExpressConnectRouterGrantAssociationsConfig)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_grant_association.default.ecr_id}",
			"ids":    []string{"${alicloud_express_connect_router_grant_association.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_grant_association.default.ecr_id}",
			"ids":    []string{"${alicloud_express_connect_router_grant_association.default.id}_fake"},
		}),
	}

	instanceIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id":      "${alicloud_express_connect_router_grant_association.default.ecr_id}",
			"instance_id": "${alicloud_express_connect_router_grant_association.default.instance_id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":      "${alicloud_express_connect_router_grant_association.default.ecr_id}",
			"instance_id": "${alicloud_express_connect_router_grant_association.default.instance_id}_fake",
		}),
	}

	instanceRegionIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id":             "${alicloud_express_connect_router_grant_association.default.ecr_id}",
			"instance_region_id": "${alicloud_express_connect_router_grant_association.default.instance_region_id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":             "${alicloud_express_connect_router_grant_association.default.ecr_id}",
			"instance_region_id": "${alicloud_express_connect_router_grant_association.default.instance_region_id}_fake",
		}),
	}

	instanceTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id":        "${alicloud_express_connect_router_grant_association.default.ecr_id}",
			"instance_type": "${alicloud_express_connect_router_grant_association.default.instance_type}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":        "${alicloud_express_connect_router_grant_association.default.ecr_id}",
			"instance_type": "VBR",
		}),
	}

	callerTypeConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id":      "${alicloud_express_connect_router_grant_association.default.id}",
			"caller_type": "ECR",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":      "${alicloud_express_connect_router_grant_association.default.id}",
			"caller_type": "OTHER",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id":             "${alicloud_express_connect_router_grant_association.default.ecr_id}",
			"ids":                []string{"${alicloud_express_connect_router_grant_association.default.id}"},
			"instance_id":        "${alicloud_express_connect_router_grant_association.default.instance_id}",
			"instance_region_id": "${alicloud_express_connect_router_grant_association.default.instance_region_id}",
			"instance_type":      "${alicloud_express_connect_router_grant_association.default.instance_type}",
			"caller_type":        "ECR",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
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
			"associations.0.id":                 acctest.CHECKSET,
			"associations.0.ecr_id":             acctest.CHECKSET,
			"associations.0.instance_id":        acctest.CHECKSET,
			"associations.0.instance_region_id": acctest.CHECKSET,
			"associations.0.owner_id":           acctest.CHECKSET,
			"associations.0.instance_type":      acctest.CHECKSET,
			"associations.0.grant_id":           acctest.CHECKSET,
			"associations.0.instance_owner_id":  acctest.CHECKSET,
			"associations.0.instance_owner_bid": acctest.CHECKSET,
			"associations.0.status":             acctest.CHECKSET,
			"associations.0.create_time":        acctest.CHECKSET,
			"associations.0.modify_time":        acctest.CHECKSET,
		}
	}

	var fakeAliCloudExpressConnectRouterGrantAssociationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"associations.#": "0",
		}
	}

	var aliCloudExpressConnectRouterGrantAssociationsInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_express_connect_router_grant_associations.default",
		ExistMapFunc: existAliCloudExpressConnectRouterGrantAssociationsMapFunc,
		FakeMapFunc:  fakeAliCloudExpressConnectRouterGrantAssociationsMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudExpressConnectRouterGrantAssociationsInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, instanceIdConf, instanceRegionIdConf, instanceTypeConf, callerTypeConf, allConf)
}

func dataSourceExpressConnectRouterGrantAssociationsConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

variable "vpc_id" {
  # Authorized cross-account fixture: matches the pre-provisioned VPC used by
  # resource_alicloud_express_connect_router_grant_association_test (basic9639).
  default = "vpc-7qbkkicdvc69hzhu48x98"
}

variable "ecr_owner_uid" {
  # Authorized ECR owner AliUid (number, TypeInt schema field).
  default = 1891593620094065
}

variable "ecr_id" {
  # Authorized pre-provisioned ECR used by basic9639.
  default = "ecr-g2g2f4ww1vduya9xcn"
}

variable "region" {
  # Region of the pre-provisioned VPC instance.
  default = "cn-wulanchabu-test-5"
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
