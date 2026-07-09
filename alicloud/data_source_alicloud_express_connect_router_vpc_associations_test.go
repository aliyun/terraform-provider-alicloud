package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudExpressConnectRouterVpcAssociationsDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_express_connect_router_vpc_associations.default"
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectroutervpcassociation%d", defaultRegionToTest, rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceExpressConnectRouterVpcAssociationsConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"ids":    []string{"${alicloud_express_connect_router_vpc_association.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"ids":    []string{"${alicloud_express_connect_router_vpc_association.default.id}_fake"},
		}),
	}

	associationIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id":         "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"association_id": "${alicloud_express_connect_router_vpc_association.default.association_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":         "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"association_id": "${alicloud_express_connect_router_vpc_association.default.association_id}_fake",
		}),
	}

	vpcIdConfConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"vpc_id": "${alicloud_express_connect_router_vpc_association.default.vpc_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"vpc_id": "${alicloud_express_connect_router_vpc_association.default.vpc_id}_fake",
		}),
	}

	associationRegionIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id":                "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"association_region_id": "${alicloud_express_connect_router_vpc_association.default.association_region_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":                "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"association_region_id": "${alicloud_express_connect_router_vpc_association.default.association_region_id}_fake",
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"status": "ACTIVE",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"status": "CREATING",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id":                "${alicloud_express_connect_router_vpc_association.default.ecr_id}",
			"ids":                   []string{"${alicloud_express_connect_router_vpc_association.default.id}"},
			"association_id":        "${alicloud_express_connect_router_vpc_association.default.association_id}",
			"vpc_id":                "${alicloud_express_connect_router_vpc_association.default.vpc_id}",
			"association_region_id": "${alicloud_express_connect_router_vpc_association.default.association_region_id}",
			"status":                "ACTIVE",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
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
			"associations.0.id":                    CHECKSET,
			"associations.0.ecr_id":                CHECKSET,
			"associations.0.association_id":        CHECKSET,
			"associations.0.vpc_id":                CHECKSET,
			"associations.0.association_node_type": CHECKSET,
			"associations.0.vpc_owner_id":          CHECKSET,
			"associations.0.allowed_prefixes_mode": CHECKSET,
			"associations.0.status":                CHECKSET,
			"associations.0.create_time":           CHECKSET,
			"associations.0.modify_time":           CHECKSET,
			"associations.0.allowed_prefixes.#":    CHECKSET,
		}
	}

	var fakeAliCloudExpressConnectRouterVpcAssociationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"associations.#": "0",
		}
	}

	var aliCloudExpressConnectRouterVpcAssociationsInfo = dataSourceAttr{
		resourceId:   "data.alicloud_express_connect_router_vpc_associations.default",
		existMapFunc: existAliCloudExpressConnectRouterVpcAssociationsMapFunc,
		fakeMapFunc:  fakeAliCloudExpressConnectRouterVpcAssociationsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudExpressConnectRouterVpcAssociationsInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, associationIdConf, vpcIdConfConf, associationRegionIdConf, statusConf, allConf)
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
