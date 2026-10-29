package expressconnectrouter_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudExpressConnectRouterVpcAssociationsDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_express_connect_router_vpc_associations.default"
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectroutervpcassociation%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceExpressConnectRouterVpcAssociationsConfig)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"ids":    []string{"${alicloud_express_connect_router_vpc_association.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"ids":    []string{"${alicloud_express_connect_router_vpc_association.default.id}_fake"},
		}),
	}

	associationIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id":         "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"association_id": "${alicloud_express_connect_router_vpc_association.default.association_id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":         "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"association_id": "${alicloud_express_connect_router_vpc_association.default.association_id}_fake",
		}),
	}

	vpcIdConfConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"vpc_id": "${alicloud_express_connect_router_vpc_association.default.vpc_id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"vpc_id": "${alicloud_express_connect_router_vpc_association.default.vpc_id}_fake",
		}),
	}

	associationRegionIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id":                "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"association_region_id": "${alicloud_express_connect_router_vpc_association.default.association_region_id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":                "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"association_region_id": "${alicloud_express_connect_router_vpc_association.default.association_region_id}_fake",
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"status": "ACTIVE",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"status": "CREATING",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id":                "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"ids":                   []string{"${alicloud_express_connect_router_vpc_association.default.id}"},
			"association_id":        "${alicloud_express_connect_router_vpc_association.default.association_id}",
			"vpc_id":                "${alicloud_express_connect_router_vpc_association.default.vpc_id}",
			"association_region_id": "${alicloud_express_connect_router_vpc_association.default.association_region_id}",
			"status":                "ACTIVE",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":                "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"ids":                   []string{"${alicloud_express_connect_router_vpc_association.default.id}_fake"},
			"association_id":        "${alicloud_express_connect_router_vpc_association.default.association_id}_fake",
			"vpc_id":                "${alicloud_express_connect_router_vpc_association.default.vpc_id}_fake",
			"association_region_id": "${alicloud_express_connect_router_vpc_association.default.association_region_id}_fake",
			"status":                "CREATING",
		}),
	}

	var existAliCloudExpressConnectRouterVpcAssociationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                "1",
			"associations.#":                       "1",
			"associations.0.id":                    acctest.CHECKSET,
			"associations.0.ecr_id":                acctest.CHECKSET,
			"associations.0.association_id":        acctest.CHECKSET,
			"associations.0.vpc_id":                acctest.CHECKSET,
			"associations.0.association_node_type": acctest.CHECKSET,
			"associations.0.vpc_owner_id":          acctest.CHECKSET,
			"associations.0.allowed_prefixes_mode": acctest.CHECKSET,
			"associations.0.status":                acctest.CHECKSET,
			"associations.0.create_time":           acctest.CHECKSET,
			"associations.0.modify_time":           acctest.CHECKSET,
			"associations.0.allowed_prefixes.#":    acctest.CHECKSET,
		}
	}

	var fakeAliCloudExpressConnectRouterVpcAssociationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"associations.#": "0",
		}
	}

	var aliCloudExpressConnectRouterVpcAssociationsInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_express_connect_router_vpc_associations.default",
		ExistMapFunc: existAliCloudExpressConnectRouterVpcAssociationsMapFunc,
		FakeMapFunc:  fakeAliCloudExpressConnectRouterVpcAssociationsMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudExpressConnectRouterVpcAssociationsInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, associationIdConf, vpcIdConfConf, associationRegionIdConf, statusConf, allConf)
}

func dataSourceExpressConnectRouterVpcAssociationsConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

data "alicloud_account" "default" {
}

data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_express_connect_router_express_connect_router" "default" {
  alibaba_side_asn = "65532"
}

resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_express_connect_router_vpc_association" "default" {
  ecr_id                = alicloud_express_connect_router_express_connect_router.default.id
  vpc_id                = alicloud_vpc.default.id
  association_region_id = data.alicloud_regions.default.regions.0.id
  vpc_owner_id          = data.alicloud_account.default.id
  allowed_prefixes      = ["172.16.1.0/24", "172.16.2.0/24", "172.16.3.0/24"]
}
`, name)
}
