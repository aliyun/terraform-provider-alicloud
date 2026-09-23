package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudExpressConnectRouterTrAssociationsDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_express_connect_router_tr_associations.default"
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterexpressconnectroutertrassociation%d", defaultRegionToTest, rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceExpressConnectRouterTrAssociationsConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"ids":    []string{"${alicloud_express_connect_router_tr_association.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"ids":    []string{"${alicloud_express_connect_router_tr_association.default.id}_fake"},
		}),
	}

	associationIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id":         "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"association_id": "${alicloud_express_connect_router_tr_association.default.association_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":         "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"association_id": "${alicloud_express_connect_router_tr_association.default.association_id}_fake",
		}),
	}

	transitRouterIdConfConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id":            "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"transit_router_id": "${alicloud_express_connect_router_tr_association.default.transit_router_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":            "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"transit_router_id": "${alicloud_express_connect_router_tr_association.default.transit_router_id}_fake",
		}),
	}

	associationRegionIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id":                "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"association_region_id": "${alicloud_express_connect_router_tr_association.default.association_region_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id":                "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"association_region_id": "${alicloud_express_connect_router_tr_association.default.association_region_id}_fake",
		}),
	}

	cenIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"cen_id": "${alicloud_express_connect_router_tr_association.default.cen_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"cen_id": "${alicloud_express_connect_router_tr_association.default.cen_id}_fake",
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"status": "INACTIVE",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"status": "CREATING",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id":                "${alicloud_express_connect_router_tr_association.default.ecr_id}",
			"ids":                   []string{"${alicloud_express_connect_router_tr_association.default.id}"},
			"association_id":        "${alicloud_express_connect_router_tr_association.default.association_id}",
			"transit_router_id":     "${alicloud_express_connect_router_tr_association.default.transit_router_id}",
			"association_region_id": "${alicloud_express_connect_router_tr_association.default.association_region_id}",
			"cen_id":                "${alicloud_express_connect_router_tr_association.default.cen_id}",
			"status":                "INACTIVE",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
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
			"associations.0.id":                      CHECKSET,
			"associations.0.ecr_id":                  CHECKSET,
			"associations.0.association_id":          CHECKSET,
			"associations.0.transit_router_id":       CHECKSET,
			"associations.0.association_node_type":   CHECKSET,
			"associations.0.transit_router_owner_id": CHECKSET,
			"associations.0.cen_id":                  CHECKSET,
			"associations.0.allowed_prefixes_mode":   CHECKSET,
			"associations.0.status":                  CHECKSET,
			"associations.0.create_time":             CHECKSET,
			"associations.0.modify_time":             CHECKSET,
			"associations.0.allowed_prefixes.#":      CHECKSET,
		}
	}

	var fakeAliCloudExpressConnectRouterTrAssociationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"associations.#": "0",
		}
	}

	var aliCloudExpressConnectRouterTrAssociationsInfo = dataSourceAttr{
		resourceId:   "data.alicloud_express_connect_router_tr_associations.default",
		existMapFunc: existAliCloudExpressConnectRouterTrAssociationsMapFunc,
		fakeMapFunc:  fakeAliCloudExpressConnectRouterTrAssociationsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudExpressConnectRouterTrAssociationsInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, associationIdConf, transitRouterIdConfConf, associationRegionIdConf, cenIdConf, statusConf, allConf)
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
