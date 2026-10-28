package privatelink_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudPrivatelinkVpcEndpointsDataSource(t *testing.T) {
	resourceId := "data.alicloud_privatelink_vpc_endpoints.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccPrivatelinkVpcEndpoints%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourcePrivatelinkVpcEndpointsDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_privatelink_vpc_endpoint.default.vpc_endpoint_name}",
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_privatelink_vpc_endpoint.default.vpc_endpoint_name}-fake",
			"enable_details": "true",
		}),
	}
	tagConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_privatelink_vpc_endpoint.default.vpc_endpoint_name}",
			"enable_details": "true",
			"tags": map[string]string{
				"Created": "TF",
			},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_privatelink_vpc_endpoint.default.vpc_endpoint_name}-fake",
			"enable_details": "true",
			"tags": map[string]string{
				"Created": "TF-fake",
			},
		}),
	}
	statusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_privatelink_vpc_endpoint.default.id}"},
			"status":         "Active",
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_privatelink_vpc_endpoint.default.id}"},
			"status":         "Creating",
			"enable_details": "true",
		}),
	}
	connectionStatusConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_privatelink_vpc_endpoint.default.id}"},
			"connection_status": "Disconnected",
			"enable_details":    "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_privatelink_vpc_endpoint.default.id}"},
			"connection_status": "Pending",
			"enable_details":    "true",
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_privatelink_vpc_endpoint.default.id}"},
			"enable_details": "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_privatelink_vpc_endpoint.default.id}-fake"},
			"enable_details": "true",
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":        name,
			"status":            "Active",
			"connection_status": "Disconnected",
			"ids":               []string{"${alicloud_privatelink_vpc_endpoint.default.id}"},
			"enable_details":    "true",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":        name + "_fake",
			"status":            "Creating",
			"connection_status": "Pending",
			"ids":               []string{"${alicloud_privatelink_vpc_endpoint.default.id}-fake"},
			"enable_details":    "true",
		}),
	}
	var existPrivatelinkVpcEndpointsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                "1",
			"ids.0":                                acctest.CHECKSET,
			"names.#":                              "1",
			"names.0":                              acctest.CHECKSET,
			"endpoints.#":                          "1",
			"endpoints.0.id":                       acctest.CHECKSET,
			"endpoints.0.bandwidth":                "103",
			"endpoints.0.connection_status":        "Disconnected",
			"endpoints.0.endpoint_business_status": "Normal",
			"endpoints.0.endpoint_description":     "",
			"endpoints.0.endpoint_domain":          acctest.CHECKSET,
			"endpoints.0.endpoint_id":              acctest.CHECKSET,
			"endpoints.0.security_group_ids.#":     "1",
			"endpoints.0.service_id":               acctest.CHECKSET,
			"endpoints.0.service_name":             acctest.CHECKSET,
			"endpoints.0.status":                   "Active",
			"endpoints.0.vpc_endpoint_name":        name,
			"endpoints.0.vpc_id":                   acctest.CHECKSET,
			"endpoints.0.zone.#":                   "0",
			"endpoints.0.tags.%":                   "1",
			"endpoints.0.tags.Created":             "TF",
		}
	}

	var fakePrivatelinkVpcEndpointsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"endpoints.#": "0",
		}
	}

	var PrivatelinkVpcEndpointsInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existPrivatelinkVpcEndpointsMapFunc,
		FakeMapFunc:  fakePrivatelinkVpcEndpointsMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.PrivateLinkRegions)
	}

	PrivatelinkVpcEndpointsInfo.DataSourceTestCheckWithPreCheck(t, 0, preCheck, nameRegexConf, tagConf, statusConf, connectionStatusConf, idsConf, allConf)
}

func dataSourcePrivatelinkVpcEndpointsDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_vpcs" "default" {
	  name_regex = "^default-NODELETING$"
	}
	resource "alicloud_security_group" "default" {
	  name        = "tf-testAcc-for-privatelink"
	  description = "privatelink test security group"
	  vpc_id      = data.alicloud_vpcs.default.ids.0
	}
	resource "alicloud_privatelink_vpc_endpoint_service" "default" {
	service_description = ""
	connect_bandwidth = 103
    auto_accept_connection = false
	}
    resource "alicloud_privatelink_vpc_endpoint" "default" {
	  tags = {
		Created = "TF"
	  }
	 service_id = alicloud_privatelink_vpc_endpoint_service.default.id
	 vpc_id = data.alicloud_vpcs.default.ids.0
     security_group_ids = [alicloud_security_group.default.id]
	 vpc_endpoint_name = "%[1]s"
	 depends_on = [alicloud_privatelink_vpc_endpoint_service.default]
	}
	`, name)
}
