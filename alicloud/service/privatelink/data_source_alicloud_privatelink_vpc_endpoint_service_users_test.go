package privatelink_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudPrivatelinkVpcEndpointServiceUsersDataSource(t *testing.T) {
	resourceId := "data.alicloud_privatelink_vpc_endpoint_service_users.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccPrivatelinkVpcEndpointServiceUsers%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourcePrivatelinkVpcEndpointServiceUsersDependence)

	userIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"service_id": "${alicloud_privatelink_vpc_endpoint_service_user.default.service_id}",
			"user_id":    "${alicloud_privatelink_vpc_endpoint_service_user.default.user_id}",
		}),
		FakeConfig: "",
	}

	var existPrivatelinkVpcEndpointServiceUsersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":           "1",
			"ids.0":           acctest.CHECKSET,
			"users.#":         "1",
			"users.0.id":      acctest.CHECKSET,
			"users.0.user_id": acctest.CHECKSET,
		}
	}

	var fakePrivatelinkVpcEndpointServiceUsersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}

	var PrivatelinkVpcEndpointServiceUsersInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existPrivatelinkVpcEndpointServiceUsersMapFunc,
		FakeMapFunc:  fakePrivatelinkVpcEndpointServiceUsersMapFunc,
	}

	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.PrivateLinkRegions)
	}

	PrivatelinkVpcEndpointServiceUsersInfo.DataSourceTestCheckWithPreCheck(t, 0, preCheck, userIdConf)
}

func dataSourcePrivatelinkVpcEndpointServiceUsersDependence(name string) string {
	return fmt.Sprintf(`
	resource "alicloud_ram_user" "user" {
	    name         = "%[1]s"
	    display_name = "user_display_name"
	    mobile       = "86-18688888888"
	    email        = "hello.uuu@aaa.com"
	    comments     = "yoyoyo"
	}
	data "alicloud_vpcs" "default" {
	    name_regex = "^default-NODELETING$"
	}
	resource "alicloud_security_group" "default" {
	    name        = "tf-testAcc-for-privatelink"
	    description = "privatelink test security group"
	    vpc_id      = data.alicloud_vpcs.default.ids.0
	}
	resource "alicloud_privatelink_vpc_endpoint_service" "default" {
		service_description = "%[1]s"
		connect_bandwidth = 103
		auto_accept_connection = false
	}
    resource "alicloud_privatelink_vpc_endpoint_service_user" "default" {
        service_id = alicloud_privatelink_vpc_endpoint_service.default.id
        user_id = alicloud_ram_user.user.id
	}
	`, name)
}
