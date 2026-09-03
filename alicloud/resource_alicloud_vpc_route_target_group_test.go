// Hand-written ACC test for alicloud_vpc_route_target_group.
// The generator stub set Computed-only nested fields (enable_status /
// health_check_status) in the HCL config, which is invalid; this file
// rewrites the test body to a valid, fully-covered hand-written test.
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccAliCloudVpcRouteTargetGroup_basic covers the full lifecycle of the
// RouteTargetGroup resource: create (Active-Standby, 2 GWLB endpoint members
// with weight 100/0), update (name/description), resource group
// move, tags add/update/remove, and import.
//
// Coverage (every Optional/Required field appears in at least one Step Config):
//   - config_mode (Required, ForceNew)
//   - vpc_id (Required, ForceNew)
//   - route_target_member_list (Required) + nested member_id/member_type/weight
//   - route_target_group_name (Optional, updatable)
//   - route_target_group_description (Optional, updatable)
//   - resource_group_id (Optional+Computed, via MoveResourceGroup)
//   - tags (Optional, via SetResourceTags)
//
// Computed fields asserted: status, resource_group_id.
// member_list count asserted (TypeSet order is non-deterministic).
func TestAccAliCloudVpcRouteTargetGroup_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_route_target_group.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcRouteTargetGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcRouteTargetGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc-rtg%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcRouteTargetGroupBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-wulanchabu"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			// Step 1: create with every Optional/Required field set.
			// Active-Standby requires exactly 2 members (weight 100 active + 0
			// standby) in different AZs.
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":                     []string{"alicloud_privatelink_vpc_endpoint_zone.getVpcEndpointZoneA", "alicloud_privatelink_vpc_endpoint_zone.getVpcEndpointZoneB"},
					"config_mode":                    "Active-Standby",
					"vpc_id":                         "${alicloud_vpc.defaultVpc.id}",
					"route_target_group_name":        name,
					"route_target_group_description": "tf-test-route-target-group-desc",
					"resource_group_id":              "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"route_target_member_list": []map[string]interface{}{
						{
							"member_id":   "${alicloud_privatelink_vpc_endpoint.getVpcEndpointA.id}",
							"member_type": "GatewayLoadBalancerEndpoint",
							"weight":      100,
						},
						{
							"member_id":   "${alicloud_privatelink_vpc_endpoint.getVpcEndpointB.id}",
							"member_type": "GatewayLoadBalancerEndpoint",
							"weight":      0,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"config_mode":                    "Active-Standby",
						"vpc_id":                         CHECKSET,
						"route_target_group_name":        name,
						"route_target_group_description": "tf-test-route-target-group-desc",
						"resource_group_id":              CHECKSET,
						"route_target_member_list.#":     "2",
						"status":                         CHECKSET,
					}),
				),
			},
			// Step 2: update name + description (the supported UpdateRouteTargetGroup
			// path). Weight is immutable and the API forbids updating an enabled
			// member, so the member list is inherited unchanged from Step 1 and the
			// provider only sends RouteTargetMemberList when it actually changed.
			// Switching active/standby (a weight swap) is a separate
			// SwitchActiveRouteTarget operation, not covered by this update path.
			{
				Config: testAccConfig(map[string]interface{}{
					"route_target_group_name":        name + "-update",
					"route_target_group_description": "tf-test-desc-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_target_group_name":        name + "-update",
						"route_target_group_description": "tf-test-desc-update",
						"route_target_member_list.#":     "2",
					}),
				),
			},
			// Step 3: move resource group (MoveResourceGroup).
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			// Step 4: move resource group back.
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			// Step 5: tags add (SetResourceTags).
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			// Step 6: tags update.
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			// Step 7: tags remove.
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			// Step 8: import by route_target_group_id.
			// All Computed attributes (enable_status, health_check_status,
			// status, create_time) are verified on import; none are ignored.
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// AlicloudVpcRouteTargetGroupMap asserts the top-level Computed fields that
// are expected to be populated on every step once the resource exists.
//
// create_time is documented in the GetRouteTargetGroup response but the backend
// does not currently return it (consistent across all Gets in the ACC run); the
// schema keeps it Computed and the Read maps CreateTime per the API contract, so
// it will populate once the backend honors its documented response. It is not
// asserted here because asserting an unreturned backend field would test a
// backend gap rather than provider behavior.
var AlicloudVpcRouteTargetGroupMap = map[string]string{
	"status": CHECKSET,
}

// AlicloudVpcRouteTargetGroupBasicDependence builds the precondition fixture for
// Active-Standby: one VPC with two vswitches in two different AZs, two GWLB
// load balancers (one per AZ), each backing a GWLB-type PrivateLink endpoint
// service (service_resource_type=gwlb + an endpoint_service_resource that
// attaches the GWLB), and two GatewayLoadBalancer endpoints — one per chain —
// used as the two route-target members. The two members land in two different
// availability zones, satisfying Active-Standby's two-different-AZ rule.
// Each endpoint additionally carries an alicloud_privatelink_vpc_endpoint_zone
// attaching it to its AZ's vswitch: the RouteTargetGroup backend looks up the
// member endpoint by zone, so the endpoint must have a non-empty zone or
// CreateRouteTargetGroup returns ResourceNotFound.GatewayLoadBalancerEndpoint.
// The route_target_group depends_on both endpoint zones so they exist before
// Create is called.
//
// Region is cn-wulanchabu (a GWLB-supported region); eu-central-1 has no GWLB
// and produced Mismatch.EndpointType because the endpoint services were not
// GWLB-backed. Each endpoint depends_on its endpoint_service_resource so the
// GWLB is attached to the service before the endpoint is created.
func AlicloudVpcRouteTargetGroupBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

variable "zone_id_1" {
  default = "cn-wulanchabu-b"
}

variable "zone_id_2" {
  default = "cn-wulanchabu-c"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultVpc" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "vswitchA" {
  vpc_id     = alicloud_vpc.defaultVpc.id
  zone_id    = var.zone_id_1
  cidr_block = "192.168.0.0/24"
}

resource "alicloud_vswitch" "vswitchB" {
  vpc_id     = alicloud_vpc.defaultVpc.id
  zone_id    = var.zone_id_2
  cidr_block = "192.168.1.0/24"
}

# Chain A (zone 1): GWLB load balancer + GWLB-type endpoint service +
# service-resource attachment + GatewayLoadBalancer endpoint.
resource "alicloud_gwlb_load_balancer" "gwlbA" {
  load_balancer_name  = "${var.name}-gwlb-a"
  address_ip_version  = "Ipv4"
  vpc_id              = alicloud_vpc.defaultVpc.id
  zone_mappings {
    vswitch_id = alicloud_vswitch.vswitchA.id
    zone_id    = var.zone_id_1
  }
}

resource "alicloud_privatelink_vpc_endpoint_service" "getVpcEndpointServiceA" {
  auto_accept_connection = true
  service_description     = "${var.name}-eps-a"
  service_resource_type   = "gwlb"
}

resource "alicloud_privatelink_vpc_endpoint_service_resource" "getVpcEndpointServiceResourceA" {
  resource_id   = alicloud_gwlb_load_balancer.gwlbA.id
  resource_type = "gwlb"
  service_id    = alicloud_privatelink_vpc_endpoint_service.getVpcEndpointServiceA.id
  zone_id       = var.zone_id_1
  dry_run       = "false"
}

resource "alicloud_privatelink_vpc_endpoint" "getVpcEndpointA" {
  service_id        = alicloud_privatelink_vpc_endpoint_service.getVpcEndpointServiceA.id
  vpc_endpoint_name = "${var.name}-ep-a"
  vpc_id            = alicloud_vpc.defaultVpc.id
  service_name      = alicloud_privatelink_vpc_endpoint_service.getVpcEndpointServiceA.vpc_endpoint_service_name
  endpoint_type     = "GatewayLoadBalancer"

  depends_on = [alicloud_privatelink_vpc_endpoint_service_resource.getVpcEndpointServiceResourceA]
}

# Zone A: attach zone 1 to the GWLB endpoint. The RouteTargetGroup backend
# looks up the member endpoint by zone, so the endpoint must carry a
# non-empty zone (added via VpcEndpointZone) or CreateRouteTargetGroup
# returns ResourceNotFound.GatewayLoadBalancerEndpoint.
resource "alicloud_privatelink_vpc_endpoint_zone" "getVpcEndpointZoneA" {
  endpoint_id = alicloud_privatelink_vpc_endpoint.getVpcEndpointA.id
  vswitch_id  = alicloud_vswitch.vswitchA.id
}

# Chain B (zone 2): identical structure in a different AZ so the two
# route-target members satisfy Active-Standby's two-different-AZ rule.
resource "alicloud_gwlb_load_balancer" "gwlbB" {
  load_balancer_name  = "${var.name}-gwlb-b"
  address_ip_version  = "Ipv4"
  vpc_id              = alicloud_vpc.defaultVpc.id
  zone_mappings {
    vswitch_id = alicloud_vswitch.vswitchB.id
    zone_id    = var.zone_id_2
  }
}

resource "alicloud_privatelink_vpc_endpoint_service" "getVpcEndpointServiceB" {
  auto_accept_connection = true
  service_description     = "${var.name}-eps-b"
  service_resource_type   = "gwlb"
}

resource "alicloud_privatelink_vpc_endpoint_service_resource" "getVpcEndpointServiceResourceB" {
  resource_id   = alicloud_gwlb_load_balancer.gwlbB.id
  resource_type = "gwlb"
  service_id    = alicloud_privatelink_vpc_endpoint_service.getVpcEndpointServiceB.id
  zone_id       = var.zone_id_2
  dry_run       = "false"
}

resource "alicloud_privatelink_vpc_endpoint" "getVpcEndpointB" {
  service_id        = alicloud_privatelink_vpc_endpoint_service.getVpcEndpointServiceB.id
  vpc_endpoint_name = "${var.name}-ep-b"
  vpc_id            = alicloud_vpc.defaultVpc.id
  service_name      = alicloud_privatelink_vpc_endpoint_service.getVpcEndpointServiceB.vpc_endpoint_service_name
  endpoint_type     = "GatewayLoadBalancer"

  depends_on = [alicloud_privatelink_vpc_endpoint_service_resource.getVpcEndpointServiceResourceB]
}

# Zone B: attach zone 2 to the second GWLB endpoint (different zone from
# member A, satisfying active-standby's two-different-AZ rule at the
# endpoint-zone level too).
resource "alicloud_privatelink_vpc_endpoint_zone" "getVpcEndpointZoneB" {
  endpoint_id = alicloud_privatelink_vpc_endpoint.getVpcEndpointB.id
  vswitch_id  = alicloud_vswitch.vswitchB.id
}
`, name)
}
