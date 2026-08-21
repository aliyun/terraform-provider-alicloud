package alicloud

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestUnitFlattenPolarDBGatewayAttribute(t *testing.T) {
	endpoints := flattenPolarDBGatewayEndpoints(map[string]interface{}{
		"Endpoint": []interface{}{map[string]interface{}{
			"Address": "gw.example.com", "EndpointId": "gwe-1", "GwClusterId": "gw-1",
			"Port": float64(3306), "TunnelId": "tunnel-1", "VpcId": "vpc-1", "NetType": "Private",
		}},
	})
	wantEndpoints := []map[string]interface{}{{
		"address": "gw.example.com", "endpoint_id": "gwe-1", "gateway_id": "gw-1",
		"port": float64(3306), "tunnel_id": "tunnel-1", "vpc_id": "vpc-1", "network_type": "Private",
	}}
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
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"region_id":         "${data.alicloud_regions.default.regions.0.id}",
					"zone_id":           "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
					"db_cluster_class":  "polar.mysql.x4.large",
					"pay_type":          "Postpaid",
					"auto_renew":        "false",
					"period":            "Month",
					"used_time":         "1",
					"vpc_id":            "${data.alicloud_vswitches.default.vswitches.0.vpc_id}",
					"vswitch_id":        "${data.alicloud_vswitches.default.vswitches.0.id}",
					"security_group_id": "${alicloud_security_group.default.id}",
					"db_type":           "MySQL",
				}),
				Check: resource.ComposeTestCheckFunc(
					rac.resourceAttrMapUpdateSet()(map[string]string{
						"region_id": CHECKSET, "zone_id": CHECKSET, "db_cluster_class": "polar.mysql.x4.large",
						"pay_type": "Postpaid", "auto_renew": "false", "vpc_id": CHECKSET,
						"period": "Month", "used_time": "1",
						"vswitch_id": CHECKSET, "security_group_id": CHECKSET, "db_type": "MySQL",
						"status": CHECKSET, "description": CHECKSET, "create_time": CHECKSET,
						"modify_time": CHECKSET, "expire_time": CHECKSET, "expired": CHECKSET,
						"latest_version": CHECKSET, "current_version": CHECKSET, "running_version": CHECKSET,
						"endpoints": CHECKSET, "security_ip_arrays": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Reuse an existing vSwitch and its VPC to avoid consuming VPC quota.
func resourcePolarDBGatewayConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_regions" "default" {
  current = true
}

data "alicloud_vswitches" "default" {
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vswitches.default.vswitches.0.vpc_id
}
`, name)
}
