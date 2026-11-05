package privatelink_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudPrivateLinkVpcEndpointServicesDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(1000, 9999)
	resourceId := "data.alicloud_privatelink_vpc_endpoint_services.default"
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourcePrivateLinkVpcEndpointServicesConfig0)

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_privatelink_vpc_endpoint_service.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_privatelink_vpc_endpoint_service.default.id}_fake"},
		}),
	}

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_privatelink_vpc_endpoint_service.default.vpc_endpoint_service_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_privatelink_vpc_endpoint_service.default.vpc_endpoint_service_name}_fake",
		}),
	}

	vpcEndpointServiceNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"vpc_endpoint_service_name": "${alicloud_privatelink_vpc_endpoint_service.default.vpc_endpoint_service_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"vpc_endpoint_service_name": "${alicloud_privatelink_vpc_endpoint_service.default.vpc_endpoint_service_name}_fake",
		}),
	}

	autoAcceptConnectionConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"vpc_endpoint_service_name": "${alicloud_privatelink_vpc_endpoint_service.default.vpc_endpoint_service_name}",
			"auto_accept_connection":    "${alicloud_privatelink_vpc_endpoint_service.default.auto_accept_connection}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"vpc_endpoint_service_name": "${alicloud_privatelink_vpc_endpoint_service.default.vpc_endpoint_service_name}_fake",
			"auto_accept_connection":    "false",
		}),
	}

	serviceBusinessStatusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"vpc_endpoint_service_name": "${alicloud_privatelink_vpc_endpoint_service.default.vpc_endpoint_service_name}",
			"service_business_status":   "${alicloud_privatelink_vpc_endpoint_service.default.service_business_status}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"service_business_status": "FinancialLocked",
		}),
	}

	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"vpc_endpoint_service_name": "${alicloud_privatelink_vpc_endpoint_service.default.vpc_endpoint_service_name}",
			"status":                    "Active",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"status": "Deleting",
		}),
	}

	tagsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"vpc_endpoint_service_name": "${alicloud_privatelink_vpc_endpoint_service.default.vpc_endpoint_service_name}",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "PrivateLinkVpcEndpointService",
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "PrivateLinkVpcEndpointService_Fake",
			},
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
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
		FakeConfig: testAccConfig(map[string]interface{}{
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
			"services.0.id":                        acctest.CHECKSET,
			"services.0.service_id":                acctest.CHECKSET,
			"services.0.vpc_endpoint_service_name": acctest.CHECKSET,
			"services.0.service_description":       acctest.CHECKSET,
			"services.0.service_domain":            acctest.CHECKSET,
			"services.0.connect_bandwidth":         acctest.CHECKSET,
			"services.0.auto_accept_connection":    acctest.CHECKSET,
			"services.0.service_business_status":   acctest.CHECKSET,
			"services.0.status":                    acctest.CHECKSET,
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

	var aliCloudRamRolesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_privatelink_vpc_endpoint_services.default",
		ExistMapFunc: existAliCloudPrivateLinkVpcEndpointServicesMapFunc,
		FakeMapFunc:  fakeAliCloudPrivateLinkVpcEndpointServicesMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
	}

	aliCloudRamRolesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, vpcEndpointServiceNameConf, autoAcceptConnectionConf, serviceBusinessStatusConf, statusConf, tagsConf, allConf)
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
