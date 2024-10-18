package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudPrivatelinkVpcEndpointServicesDataSource(t *testing.T) {
	resourceId := "data.alicloud_privatelink_vpc_endpoint_services.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccPrivatelinkVpcEndpointServices%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourcePrivatelinkVpcEndpointServicesDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "com.aliyuncs.privatelink.eu-central-1." + "${alicloud_privatelink_vpc_endpoint_service.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "com.aliyuncs.privatelink.eu-central-1." + "${alicloud_privatelink_vpc_endpoint_service.default.id}-fake",
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_privatelink_vpc_endpoint_service.default.id}"},
			"status": "Active",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":    []string{"${alicloud_privatelink_vpc_endpoint_service.default.id}"},
			"status": "Creating",
		}),
	}
	serviceBusinessStatusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":                     []string{"${alicloud_privatelink_vpc_endpoint_service.default.id}"},
			"service_business_status": "Normal",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":                     []string{"${alicloud_privatelink_vpc_endpoint_service.default.id}"},
			"service_business_status": "FinancialLocked",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_privatelink_vpc_endpoint_service.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_privatelink_vpc_endpoint_service.default.id}-fake"},
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{"Created": "TF", "For": "Test"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{"Created": "TF", "For": "Test-fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":              "com.aliyuncs.privatelink.eu-central-1." + "${alicloud_privatelink_vpc_endpoint_service.default.id}",
			"status":                  "Active",
			"service_business_status": "Normal",
			"ids":                     []string{"${alicloud_privatelink_vpc_endpoint_service.default.id}"}}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":              "com.aliyuncs.privatelink.eu-central-1." + "${alicloud_privatelink_vpc_endpoint_service.default.id}-fake",
			"status":                  "Creating",
			"service_business_status": "FinancialLocked",
			"ids":                     []string{"${alicloud_privatelink_vpc_endpoint_service.default.id}-fake"},
		}),
	}
	var existPrivatelinkVpcEndpointServicesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                "1",
			"ids.0":                                CHECKSET,
			"names.#":                              "1",
			"names.0":                              CHECKSET,
			"services.#":                           "1",
			"services.0.id":                        CHECKSET,
			"services.0.auto_accept_connection":    "false",
			"services.0.connect_bandwidth":         "103",
			"services.0.service_business_status":   CHECKSET,
			"services.0.service_description":       name,
			"services.0.service_domain":            CHECKSET,
			"services.0.service_id":                CHECKSET,
			"services.0.status":                    "Active",
			"services.0.vpc_endpoint_service_name": CHECKSET,
		}
	}

	var fakePrivatelinkVpcEndpointServicesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"services.#": "0",
		}
	}

	var PrivatelinkVpcEndpointServicesInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existPrivatelinkVpcEndpointServicesMapFunc,
		fakeMapFunc:  fakePrivatelinkVpcEndpointServicesMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.PrivateLinkRegions)
	}

	PrivatelinkVpcEndpointServicesInfo.dataSourceTestCheckWithPreCheck(t, 0, preCheck, nameRegexConf, statusConf, serviceBusinessStatusConf, idsConf, tagsConf, allConf)
}

func dataSourcePrivatelinkVpcEndpointServicesDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_privatelink_vpc_endpoint_service" "default" {
	  service_description = "%s"
	  connect_bandwidth = 103
      auto_accept_connection = false
	  tags = {
		Created = "TF",
		For     = "Test",
	  }
	}
	`, name)
}
