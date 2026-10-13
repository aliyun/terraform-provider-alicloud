package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAliCloudExpressConnectRouterExpressConnectRoutersDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_express_connect_router_express_connect_routers.default"
	name := fmt.Sprintf("tf-testacc%sexpressconnectrouterExpressConnectRouterExpressConnectRouter%d", defaultRegionToTest, rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceExpressConnectRouterExpressConnectRoutersConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_express_connect_router_express_connect_router.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_express_connect_router_express_connect_router.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_express_connect_router_express_connect_router.default.ecr_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_express_connect_router_express_connect_router.default.ecr_name}_fake",
		}),
	}

	ecrIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_express_connect_router.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_express_connect_router.default.id}_fake",
		}),
	}

	ecrNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_name": "${alicloud_express_connect_router_express_connect_router.default.ecr_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_name": "${alicloud_express_connect_router_express_connect_router.default.ecr_name}_fake",
		}),
	}

	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"resource_group_id": "${alicloud_express_connect_router_express_connect_router.default.resource_group_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.2.id}",
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_express_connect_router.default.id}",
			"status": "ACTIVE",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ecr_id": "${alicloud_express_connect_router_express_connect_router.default.id}",
			"status": "UPDATING",
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "ExpressConnectRouter",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "ExpressConnectRouter_Fake",
			},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_express_connect_router_express_connect_router.default.id}"},
			"name_regex":        "${alicloud_express_connect_router_express_connect_router.default.ecr_name}",
			"ecr_id":            "${alicloud_express_connect_router_express_connect_router.default.id}",
			"ecr_name":          "${alicloud_express_connect_router_express_connect_router.default.ecr_name}",
			"resource_group_id": "${alicloud_express_connect_router_express_connect_router.default.resource_group_id}",
			"status":            "ACTIVE",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "ExpressConnectRouter",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_express_connect_router_express_connect_router.default.id}_fake"},
			"name_regex":        "${alicloud_express_connect_router_express_connect_router.default.ecr_name}_fake",
			"ecr_id":            "${alicloud_express_connect_router_express_connect_router.default.id}_fake",
			"ecr_name":          "${alicloud_express_connect_router_express_connect_router.default.ecr_name}_fake",
			"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.2.id}",
			"status":            "UPDATING",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "ExpressConnectRouter_Fake",
			},
		}),
	}

	var existAliCloudExpressConnectRouterExpressConnectRoutersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"names.#":                     "1",
			"routers.#":                   "1",
			"routers.0.id":                CHECKSET,
			"routers.0.ecr_id":            CHECKSET,
			"routers.0.ecr_name":          CHECKSET,
			"routers.0.description":       CHECKSET,
			"routers.0.owner_id":          CHECKSET,
			"routers.0.resource_group_id": CHECKSET,
			"routers.0.alibaba_side_asn":  CHECKSET,
			"routers.0.biz_status":        CHECKSET,
			"routers.0.status":            CHECKSET,
			"routers.0.create_time":       CHECKSET,
			"routers.0.modify_time":       CHECKSET,
			"routers.0.tags.%":            "2",
			"routers.0.tags.Created":      "TF",
			"routers.0.tags.For":          "ExpressConnectRouter",
		}
	}

	var fakeAliCloudExpressConnectRouterExpressConnectRoutersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"routers.#": "0",
		}
	}

	var aliCloudExpressConnectRouterExpressConnectRoutersInfo = dataSourceAttr{
		resourceId:   "data.alicloud_express_connect_router_express_connect_routers.default",
		existMapFunc: existAliCloudExpressConnectRouterExpressConnectRoutersMapFunc,
		fakeMapFunc:  fakeAliCloudExpressConnectRouterExpressConnectRoutersMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudExpressConnectRouterExpressConnectRoutersInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, ecrIdConf, ecrNameConf, resourceGroupIdConf, statusConf, tagsConf, allConf)
}

func dataSourceExpressConnectRouterExpressConnectRoutersConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

data "alicloud_resource_manager_resource_groups" "default" {
}

resource "alicloud_express_connect_router_express_connect_router" "default" {
  alibaba_side_asn  = "65532"
  ecr_name          = var.name
  description       = var.name
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.1.id
  tags = {
    Created = "TF"
    For     = "ExpressConnectRouter"
  }
}
`, name)
}
