package expressconnectrouter_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudExpressConnectRouterTrAssociationsDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_express_connect_router_tr_associations.default"
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectroutertrassociation%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceExpressConnectRouterTrAssociationsConfig)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"ids":    []string{"${alicloud_express_connect_router_tr_association.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"ids":    []string{"${alicloud_express_connect_router_tr_association.default.id}_fake"},
		}),
	}

	associationIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id":         "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"association_id": "${alicloud_express_connect_router_tr_association.default.association_id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":         "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"association_id": "${alicloud_express_connect_router_tr_association.default.association_id}_fake",
		}),
	}

	transitRouterIdConfConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id":            "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"transit_router_id": "${alicloud_express_connect_router_tr_association.default.transit_router_id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":            "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"transit_router_id": "${alicloud_express_connect_router_tr_association.default.transit_router_id}_fake",
		}),
	}

	associationRegionIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id":                "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"association_region_id": "${alicloud_express_connect_router_tr_association.default.association_region_id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":                "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"association_region_id": "${alicloud_express_connect_router_tr_association.default.association_region_id}_fake",
		}),
	}

	cenIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"cen_id": "${alicloud_express_connect_router_tr_association.default.cen_id}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"cen_id": "${alicloud_express_connect_router_tr_association.default.cen_id}_fake",
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"status": "INACTIVE",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"status": "CREATING",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ecr_id":                "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"ids":                   []string{"${alicloud_express_connect_router_tr_association.default.id}"},
			"association_id":        "${alicloud_express_connect_router_tr_association.default.association_id}",
			"transit_router_id":     "${alicloud_express_connect_router_tr_association.default.transit_router_id}",
			"association_region_id": "${alicloud_express_connect_router_tr_association.default.association_region_id}",
			"cen_id":                "${alicloud_express_connect_router_tr_association.default.cen_id}",
			"status":                "INACTIVE",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":                "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"ids":                   []string{"${alicloud_express_connect_router_tr_association.default.id}_fake"},
			"association_id":        "${alicloud_express_connect_router_tr_association.default.association_id}_fake",
			"transit_router_id":     "${alicloud_express_connect_router_tr_association.default.transit_router_id}_fake",
			"association_region_id": "${alicloud_express_connect_router_tr_association.default.association_region_id}_fake",
			"cen_id":                "${alicloud_express_connect_router_tr_association.default.cen_id}_fake",
			"status":                "CREATING",
		}),
	}

	var existAliCloudExpressConnectRouterTrAssociationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"associations.#":                         "1",
			"associations.0.id":                      acctest.CHECKSET,
			"associations.0.ecr_id":                  acctest.CHECKSET,
			"associations.0.association_id":          acctest.CHECKSET,
			"associations.0.transit_router_id":       acctest.CHECKSET,
			"associations.0.association_node_type":   acctest.CHECKSET,
			"associations.0.transit_router_owner_id": acctest.CHECKSET,
			"associations.0.cen_id":                  acctest.CHECKSET,
			"associations.0.allowed_prefixes_mode":   acctest.CHECKSET,
			"associations.0.status":                  acctest.CHECKSET,
			"associations.0.create_time":             acctest.CHECKSET,
			"associations.0.modify_time":             acctest.CHECKSET,
			"associations.0.allowed_prefixes.#":      acctest.CHECKSET,
		}
	}

	var fakeAliCloudExpressConnectRouterTrAssociationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"associations.#": "0",
		}
	}

	var aliCloudExpressConnectRouterTrAssociationsInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_express_connect_router_tr_associations.default",
		ExistMapFunc: existAliCloudExpressConnectRouterTrAssociationsMapFunc,
		FakeMapFunc:  fakeAliCloudExpressConnectRouterTrAssociationsMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudExpressConnectRouterTrAssociationsInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, associationIdConf, transitRouterIdConfConf, associationRegionIdConf, cenIdConf, statusConf, allConf)
}

func dataSourceExpressConnectRouterTrAssociationsConfig(name string) string {
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

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "default" {
  cen_id = alicloud_cen_instance.default.id
}

resource "alicloud_express_connect_router_tr_association" "default" {
  ecr_id                  = alicloud_express_connect_router_express_connect_router.default.id
  transit_router_id       = alicloud_cen_transit_router.default.transit_router_id
  cen_id                  = alicloud_cen_transit_router.default.cen_id
  transit_router_owner_id = data.alicloud_account.default.id
  association_region_id   = data.alicloud_regions.default.regions.0.id
  allowed_prefixes        = ["10.0.0.0/24", "10.0.1.0/24", "10.0.2.0/24"]
}
`, name)
}
