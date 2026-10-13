package alicloud

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUnitFlattenPolarDBGatewayAttribute(t *testing.T) {
	endpoints := flattenPolarDBGatewayEndpoints(map[string]interface{}{
		"Endpoint": []interface{}{
			map[string]interface{}{
				"Address": "gw.example.com", "EndpointId": "gwe-1", "GwClusterId": "gw-1",
				"Port": "3306", "TunnelId": "tunnel-1", "VpcId": "vpc-1", "NetType": "Private",
			},
			map[string]interface{}{
				"Address": "gw2.example.com", "EndpointId": "gwe-2", "GwClusterId": "gw-1",
				"Port": float64(8443), "TunnelId": "tunnel-2", "VpcId": "vpc-1", "NetType": "Public",
			},
		},
	})
	wantEndpoints := []map[string]interface{}{
		{
			"address": "gw.example.com", "endpoint_id": "gwe-1", "gateway_id": "gw-1",
			"port": 3306, "tunnel_id": "tunnel-1", "vpc_id": "vpc-1", "network_type": "Private",
		},
		{
			"address": "gw2.example.com", "endpoint_id": "gwe-2", "gateway_id": "gw-1",
			"port": 8443, "tunnel_id": "tunnel-2", "vpc_id": "vpc-1", "network_type": "Public",
		},
	}
	if !reflect.DeepEqual(endpoints, wantEndpoints) {
		t.Fatalf("unexpected endpoints: %#v", endpoints)
	}

	securityIPArrays := flattenPolarDBGatewaySecurityIPArrays(map[string]interface{}{
		"SecurityIPArray": []interface{}{map[string]interface{}{
			"SecurityIPArrayName": "default", "SecurityIPList": "127.0.0.1",
		}},
	})
	wantSecurityIPArrays := []map[string]interface{}{{"name": "default", "ip_list": "127.0.0.1"}}
	if !reflect.DeepEqual(securityIPArrays, wantSecurityIPArrays) {
		t.Fatalf("unexpected security IP arrays: %#v", securityIPArrays)
	}
}

func TestAccAliCloudPolarDBGateway_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_polardb_gateway.default"
	ra := resourceAttrInit(resourceId, map[string]string{
		"status": CHECKSET,
	})
	serviceFunc := func() interface{} {
		return &PolarDBService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribePolarDBGatewayAttribute")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccPolarDBGateway-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourcePolarDBGatewayConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Beijing})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":           "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
					"db_cluster_class":  "polar.app.g2.small",
					"pay_type":          "Prepaid",
					"auto_renew":        "false",
					"period":            "Month",
					"used_time":         "1",
					"vpc_id":            "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"db_type":           "PostgreSQL",
				}),
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":           "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
					"db_cluster_class":  "polar.app.g2.small",
					"pay_type":          "Postpaid",
					"auto_renew":        "false",
					"vpc_id":            "${data.alicloud_vpcs.default.ids.0}",
					"vswitch_id":        "${data.alicloud_vswitches.default.ids.0}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"db_type":           "PostgreSQL",
				}),
				Check: resource.ComposeTestCheckFunc(
					rac.resourceAttrMapUpdateSet()(map[string]string{
						"region_id": CHECKSET, "zone_id": CHECKSET, "db_cluster_class": "polar.app.g2.small",
						"pay_type": "Postpaid", "auto_renew": "false", "vpc_id": CHECKSET,
						"vswitch_id": CHECKSET, "security_group_id": CHECKSET, "db_type": "PostgreSQL",
						"status": CHECKSET, "create_time": CHECKSET, "expired": CHECKSET,
						"endpoints.#": "2", "endpoints.0.address": CHECKSET, "endpoints.0.port": CHECKSET,
						"security_ip_arrays.#": "0",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"zone_id", "db_cluster_class", "pay_type", "auto_renew", "period", "used_time",
					"vpc_id", "vswitch_id", "security_group_id", "db_type", "modify_time",
				},
			},
		},
	})
}

func resourcePolarDBGatewayConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_vpcs" "default" {
  is_default = true
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-beijing-l"
}

resource "alicloud_security_group" "default" {
  security_group_name = var.name
  vpc_id              = data.alicloud_vpcs.default.ids.0
}
`, name)
}
