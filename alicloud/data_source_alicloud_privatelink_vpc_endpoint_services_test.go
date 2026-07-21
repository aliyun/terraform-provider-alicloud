package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudPrivateLinkVpcEndpointServicesDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_privatelink_vpc_endpoint_services.default"
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourcePrivateLinkVpcEndpointServicesConfig0)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_privatelink_vpc_endpoint_service.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_privatelink_vpc_endpoint_service.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_privatelink_vpc_endpoint_service.default.vpc_endpoint_service_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_privatelink_vpc_endpoint_service.default.vpc_endpoint_service_name}_fake",
		}),
	}

	vpcEndpointServiceNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vpc_endpoint_service_name": "${alicloud_privatelink_vpc_endpoint_service.default.vpc_endpoint_service_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"vpc_endpoint_service_name": "${alicloud_privatelink_vpc_endpoint_service.default.vpc_endpoint_service_name}_fake",
		}),
	}

	autoAcceptConnectionConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"auto_accept_connection": "${alicloud_privatelink_vpc_endpoint_service.default.auto_accept_connection}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"vpc_endpoint_service_name": "${alicloud_privatelink_vpc_endpoint_service.default.vpc_endpoint_service_name}_fake",
			"auto_accept_connection":    "false",
		}),
	}

	serviceBusinessStatusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vpc_endpoint_service_name": "${alicloud_privatelink_vpc_endpoint_service.default.vpc_endpoint_service_name}",
			"service_business_status":   "${alicloud_privatelink_vpc_endpoint_service.default.service_business_status}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"service_business_status": "FinancialLocked",
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vpc_endpoint_service_name": "${alicloud_privatelink_vpc_endpoint_service.default.vpc_endpoint_service_name}",
			"status":                    "Active",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"status": "Deleting",
		}),
	}

	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "PrivateLinkVpcEndpointService",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "PrivateLinkVpcEndpointService_Fake",
			},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":                       []string{"${alicloud_privatelink_vpc_endpoint_service.default.id}"},
			"name_regex":                "${alicloud_privatelink_vpc_endpoint_service.default.vpc_endpoint_service_name}",
			"vpc_endpoint_service_name": "${alicloud_privatelink_vpc_endpoint_service.default.vpc_endpoint_service_name}",
			"auto_accept_connection":    "${alicloud_privatelink_vpc_endpoint_service.default.auto_accept_connection}",
			"service_business_status":   "${alicloud_privatelink_vpc_endpoint_service.default.service_business_status}",
			"status":                    "${alicloud_privatelink_vpc_endpoint_service.default.status}",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "PrivateLinkVpcEndpointService",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":                       []string{"${alicloud_privatelink_vpc_endpoint_service.default.id}_fake"},
			"name_regex":                "${alicloud_privatelink_vpc_endpoint_service.default.vpc_endpoint_service_name}_fake",
			"vpc_endpoint_service_name": "${alicloud_privatelink_vpc_endpoint_service.default.vpc_endpoint_service_name}_fake",
			"auto_accept_connection":    "false",
			"service_business_status":   "FinancialLocked",
			"status":                    "Deleting",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "PrivateLinkVpcEndpointService_Fake",
			},
		}),
	}

	var existAliCloudPrivateLinkVpcEndpointServicesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                "1",
			"names.#":                              "1",
			"services.#":                           "1",
			"services.0.id":                        CHECKSET,
			"services.0.service_id":                CHECKSET,
			"services.0.vpc_endpoint_service_name": CHECKSET,
			"services.0.service_description":       CHECKSET,
			"services.0.service_domain":            CHECKSET,
			"services.0.connect_bandwidth":         CHECKSET,
			"services.0.auto_accept_connection":    CHECKSET,
			"services.0.service_business_status":   CHECKSET,
			"services.0.status":                    CHECKSET,
			"services.0.tags.%":                    "2",
			"services.0.tags.Created":              "TF",
			"services.0.tags.For":                  "PrivateLinkVpcEndpointService",
		}
	}

	var fakeAliCloudPrivateLinkVpcEndpointServicesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"services.#": "0",
		}
	}

	var aliCloudRamRolesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_privatelink_vpc_endpoint_services.default",
		existMapFunc: existAliCloudPrivateLinkVpcEndpointServicesMapFunc,
		fakeMapFunc:  fakeAliCloudPrivateLinkVpcEndpointServicesMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}

	aliCloudRamRolesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, vpcEndpointServiceNameConf, autoAcceptConnectionConf, serviceBusinessStatusConf, statusConf, tagsConf, allConf)
}

func dataSourcePrivateLinkVpcEndpointServicesConfig0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_privatelink_vpc_endpoint_service" "default" {
  		service_description    = var.name
  		auto_accept_connection = true
  		tags = {
    		Created = "TF",
    		For     = "PrivateLinkVpcEndpointService",
  		}
	}
`, name)
}
