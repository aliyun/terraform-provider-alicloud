---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_route_target_group"
description: |-
  Provides a Alicloud VPC Route Target Group resource.
---

# alicloud_vpc_route_target_group

Provides a VPC Route Target Group resource.

Route target group.

For information about VPC Route Target Group and how to use it, see [What is Route Target Group](https://next.api.alibabacloud.com/document/Vpc/2016-04-28/CreateRouteTargetGroup).

-> **NOTE:** Available since v1.287.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

variable "region" {
  default = "cn-wulanchabu"
}

variable "zone_id_1" {
  default = "cn-wulanchabu-b"
}

variable "zone_id_2" {
  default = "cn-wulanchabu-c"
}

provider "alicloud" {
  region = var.region
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "zone_a" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = var.zone_id_1
  cidr_block = "192.168.0.0/24"
}

resource "alicloud_vswitch" "zone_b" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = var.zone_id_2
  cidr_block = "192.168.1.0/24"
}

# Active member (zone A): GWLB load balancer + GWLB-type endpoint service +
# service-resource attachment + GatewayLoadBalancer endpoint. The endpoint
# depends_on the service-resource so the GWLB is attached to the service
# before the endpoint is created.
resource "alicloud_gwlb_load_balancer" "active" {
  load_balancer_name = "${var.name}-gwlb-active"
  address_ip_version = "Ipv4"
  vpc_id             = alicloud_vpc.default.id
  zone_mappings {
    vswitch_id = alicloud_vswitch.zone_a.id
    zone_id    = var.zone_id_1
  }
}

resource "alicloud_privatelink_vpc_endpoint_service" "active" {
  auto_accept_connection = true
  service_description    = "${var.name}-eps-active"
  service_resource_type  = "gwlb"
}

resource "alicloud_privatelink_vpc_endpoint_service_resource" "active" {
  resource_id   = alicloud_gwlb_load_balancer.active.id
  resource_type = "gwlb"
  service_id    = alicloud_privatelink_vpc_endpoint_service.active.id
  zone_id       = var.zone_id_1
  dry_run       = "false"
}

resource "alicloud_privatelink_vpc_endpoint" "active" {
  service_id        = alicloud_privatelink_vpc_endpoint_service.active.id
  vpc_endpoint_name = "${var.name}-ep-active"
  vpc_id            = alicloud_vpc.default.id
  service_name      = alicloud_privatelink_vpc_endpoint_service.active.vpc_endpoint_service_name
  endpoint_type     = "GatewayLoadBalancer"
}

# Attach zone A to the GWLB endpoint. The route target group backend looks up
# the member endpoint by zone, so the endpoint must carry a non-empty zone.
resource "alicloud_privatelink_vpc_endpoint_zone" "active" {
  endpoint_id = alicloud_privatelink_vpc_endpoint.active.id
  vswitch_id  = alicloud_vswitch.zone_a.id
}

# Standby member (zone B): identical chain in a different zone so the two
# members satisfy active-standby's two-different-zone rule.
resource "alicloud_gwlb_load_balancer" "standby" {
  load_balancer_name = "${var.name}-gwlb-standby"
  address_ip_version = "Ipv4"
  vpc_id             = alicloud_vpc.default.id
  zone_mappings {
    vswitch_id = alicloud_vswitch.zone_b.id
    zone_id    = var.zone_id_2
  }
}

resource "alicloud_privatelink_vpc_endpoint_service" "standby" {
  auto_accept_connection = true
  service_description    = "${var.name}-eps-standby"
  service_resource_type  = "gwlb"
}

resource "alicloud_privatelink_vpc_endpoint_service_resource" "standby" {
  resource_id   = alicloud_gwlb_load_balancer.standby.id
  resource_type = "gwlb"
  service_id    = alicloud_privatelink_vpc_endpoint_service.standby.id
  zone_id       = var.zone_id_2
  dry_run       = "false"
}

resource "alicloud_privatelink_vpc_endpoint" "standby" {
  service_id        = alicloud_privatelink_vpc_endpoint_service.standby.id
  vpc_endpoint_name = "${var.name}-ep-standby"
  vpc_id            = alicloud_vpc.default.id
  service_name      = alicloud_privatelink_vpc_endpoint_service.standby.vpc_endpoint_service_name
  endpoint_type     = "GatewayLoadBalancer"
}

# Attach zone B to the standby GWLB endpoint (different zone from active).
resource "alicloud_privatelink_vpc_endpoint_zone" "standby" {
  endpoint_id = alicloud_privatelink_vpc_endpoint.standby.id
  vswitch_id  = alicloud_vswitch.zone_b.id
}

# The route target group depends_on both endpoint zones: the backend looks up
# each member endpoint by zone, so the zones must exist before Create is called.
resource "alicloud_vpc_route_target_group" "default" {
  route_target_group_name        = var.name
  route_target_group_description = var.name
  vpc_id                         = alicloud_vpc.default.id
  config_mode                    = "Active-Standby"
  route_target_member_list {
    member_id   = alicloud_privatelink_vpc_endpoint.active.id
    member_type = "GatewayLoadBalancerEndpoint"
    weight      = 100
  }
  route_target_member_list {
    member_id   = alicloud_privatelink_vpc_endpoint.standby.id
    member_type = "GatewayLoadBalancerEndpoint"
    weight      = 0
  }
}
```

## Argument Reference

The following arguments are supported:
* `config_mode` - (Required, ForceNew) The configuration mode of the route target group. Supported modes include:
  - **Active-Standby**: active-standby mode.
* `resource_group_id` - (Optional, Computed) The ID of the resource group to which the route target group belongs.
* `route_target_group_description` - (Optional) The description of the route target group.
The description must be 1 to 256 characters in length and cannot start with http:// or https://.
* `route_target_group_name` - (Optional) The name of the route target group.
The name must be 1 to 128 characters in length and cannot start with http:// or https://.
* `route_target_member_list` - (Required, List) The member list of the route target group.
**Note: The parameter is immutable after resource creation. In active-standby mode, member weight and type cannot be changed via UpdateRouteTargetGroup; switching active/standby uses a separate SwitchActiveRouteTarget operation.
In active/standby mode, the following restrictions apply to route target group members:
1. The route target group must contain exactly two members.
2. The route target group members must belong to different zones. See [`route_target_member_list`](#route_target_member_list) below.
* `tags` - (Optional, Map) The tags of the route target group.
* `vpc_id` - (Required, ForceNew) The ID of the VPC to which the route target group belongs.

### `route_target_member_list`

The route_target_member_list supports the following:
* `member_id` - (Required) The instance ID of the route target member.
* `member_type` - (Required) The instance type of the route target configuration. The following type is currently supported:
  - GatewayLoadBalancerEndpoint.
* `weight` - (Required, Int) Sets the weight attribute for the current route target configuration.

In active-standby mode, the weight can only be set to 0 or 100:
- Only one route target configuration can be set to 100, serving as the active instance.
- Only one route target configuration can be set to 0, serving as the standby instance.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the route target group was created.
* `route_target_member_list` - The member list of the route target group.
  * `enable_status` - Indicates the enable status of the current route target configuration. Valid values: `Enable`, `Disable`.
  * `health_check_status` - The health check status of the current route target configuration.
* `status` - The status of the route target group. Valid values: `Pending`, `Available`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Route Target Group.
* `delete` - (Defaults to 5 mins) Used when delete the Route Target Group.
* `update` - (Defaults to 5 mins) Used when update the Route Target Group.

## Import

VPC Route Target Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_route_target_group.example <route_target_group_id>
```
