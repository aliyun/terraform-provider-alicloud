---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_route_target_groups"
sidebar_current: "docs-alicloud-datasource-vpc-route-target-groups"
description: |-
  Provides a list of VPC Route Target Group owned by an Alibaba Cloud account.
---

# alicloud_vpc_route_target_groups

This data source provides VPC Route Target Group available to the user.[What is Route Target Group](https://next.api.alibabacloud.com/document/Vpc/2016-04-28/CreateRouteTargetGroup)

-> **NOTE:** Available since v1.287.0.

## Example Usage

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
# service-resource attachment + GatewayLoadBalancer endpoint.
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

# Standby member (zone B): identical chain in a different zone.
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

data "alicloud_vpc_route_target_groups" "default" {
  ids        = [alicloud_vpc_route_target_group.default.id]
  name_regex = alicloud_vpc_route_target_group.default.route_target_group_name
  vpc_id     = alicloud_vpc.default.id
}

output "alicloud_vpc_route_target_group_example_id" {
  value = data.alicloud_vpc_route_target_groups.default.groups[0].id
}
```

## Argument Reference

The following arguments are supported:
* `resource_group_id` - (Optional) The ID of the resource group to which the route target group belongs.
* `route_target_group_id` - (Optional) The ID of the route target group.
A maximum of 50 instance IDs can be specified in a single query.
* `route_target_member_list` - (Optional) The member list of the route target group.
In active/standby mode, the following restrictions apply to route target group members:
1. The route target group must contain exactly two members.
2. The route target group members must belong to different zones. See [`route_target_member_list`](#route_target_member_list) below.
* `tags` - (Optional) The tags of the route target group.
* `vpc_id` - (Optional) The ID of the VPC to which the route target group belongs.
* `ids` - (Optional, Computed) A list of Route Target Group IDs.
* `name_regex` - (Optional) A regex string to filter Route Target Groups by name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

### `route_target_member_list`

The route_target_member_list supports the following:
* `member_id` - (Required) The instance ID of the route target member. Used to filter route target groups that contain the specified member.
* `member_type` - (Required) The instance type of the route target configuration. The following type is currently supported: GatewayLoadBalancerEndpoint.
* `enable_status` - (Computed) Indicates the enable status of the current route target configuration. Valid values: `Enable`, `Disable`.
* `health_check_status` - (Computed) The health check status of the current route target configuration.
* `weight` - (Required) Sets the weight attribute for the current route target configuration.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Route Target Group IDs.
* `names` - A list of name of Route Target Groups.
* `groups` - A list of Route Target Group Entries. Each element contains the following attributes:
  * `config_mode` - The configuration mode of the route target group.
  * `create_time` - The time when the route target group was created.
  * `region_id` - The region ID of the VPC to which the route target group belongs.
  * `resource_group_id` - The ID of the resource group to which the route target group belongs.
  * `route_target_group_description` - The description of the route target group.
  * `route_target_group_id` - The ID of the route target group.
  * `route_target_group_name` - The name of the route target group.
  * `route_target_member_list` - The member list of the route target group.
    * `enable_status` - Indicates the enable status of the current route target configuration. Valid values: `Enable`, `Disable`.
    * `health_check_status` - The health check status of the current route target configuration.
    * `member_id` - The instance ID of the route target member.
    * `member_type` - The instance type of the route target configuration.
    * `weight` - Sets the weight attribute for the current route target configuration.
  * `status` - The status of the route target group. Valid values: `Pending`, `Available`.
  * `tags` - The tags of the route target group.
  * `vpc_id` - The ID of the VPC to which the route target group belongs.
  * `id` - The ID of the resource supplied above.
