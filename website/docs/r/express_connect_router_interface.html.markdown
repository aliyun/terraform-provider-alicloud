---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_router_interface"
description: |-
  Provides a Alicloud Express Connect Router Interface resource.
---

# alicloud_express_connect_router_interface

Provides a Express Connect Router Interface resource.



For information about Express Connect Router Interface and how to use it, see [What is Router Interface](https://next.api.alibabacloud.com/document/Vpc/2016-04-28/CreateRouterInterface).

-> **NOTE:** Available since v1.263.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}
variable "name" {
  default = "tfexample"
}
data "alicloud_resource_manager_resource_groups" "default" {
}
data "alicloud_account" "this" {
}
data "alicloud_regions" "default" {
  current = true
}
data "alicloud_express_connect_physical_connections" "nameRegex" {
  name_regex = "^preserved-NODELETING-JG"
}
data "alicloud_alb_zones" "default" {
}
resource "alicloud_vpc" "default" {
  vpc_name    = var.name
  cidr_block  = "172.16.0.0/16"
  enable_ipv6 = "true"
}
resource "alicloud_vswitch" "zone_a" {
  vswitch_name         = var.name
  vpc_id               = alicloud_vpc.default.id
  cidr_block           = "172.16.0.0/24"
  zone_id              = data.alicloud_alb_zones.default.zones.0.id
  ipv6_cidr_block_mask = "6"
}
resource "alicloud_express_connect_virtual_border_router" "default" {
  physical_connection_id = data.alicloud_express_connect_physical_connections.nameRegex.connections.0.id
  vlan_id                = "1001"
  peer_gateway_ip        = "192.168.254.2"
  peering_subnet_mask    = "255.255.255.0"
  local_gateway_ip       = "192.168.254.1"
}
resource "alicloud_express_connect_router_interface" "default" {
  auto_renew                  = "true"
  spec                        = "Mini.2"
  opposite_router_type        = "VRouter"
  router_id                   = alicloud_express_connect_virtual_border_router.default.id
  description                 = "terraform-example"
  access_point_id             = "ap-cn-hangzhou-jg-B"
  resource_group_id           = data.alicloud_resource_manager_resource_groups.default.ids.0
  period                      = "1"
  opposite_router_id          = alicloud_vpc.default.router_id
  role                        = "InitiatingSide"
  payment_type                = "PayAsYouGo"
  auto_pay                    = "true"
  opposite_interface_owner_id = data.alicloud_account.this.id
  router_interface_name       = var.name
  fast_link_mode              = "true"
  opposite_region_id          = "cn-hangzhou"
  router_type                 = "VBR"
}
```

## Argument Reference

The following arguments are supported:
* `access_point_id` - (Optional, ForceNew) Access point ID
* `auto_renew` - (Optional) Whether to enable automatic renewal. Value:
  - `false` (default): disabled.
  - `true`: enabled.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `delete_health_check_ip` - (Optional) Whether to delete the health check IP address configured on the router interface. Value:
  - `true`: deletes the health check IP address.
  - `false` (default): does not delete the health check IP address.

-> **NOTE:** This parameter only applies during resource update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `description` - (Optional) The router interface description. It must be 2 to 256 characters in length and must start with a letter or a Chinese character, but cannot start with http:// or https.
* `fast_link_mode` - (Optional, ForceNew) Whether the VBR router interface is created by using the fast connection mode. The fast connection mode can automatically complete the connection after the VBR and the router interfaces at both ends of the VPC are created. Value:
  - `true` : Yes.
  - `false`(default) : No.


-> **NOTE:** - This parameter takes effect only when the value of `RouterType` is `VBR` and the value of `OppositeRouterType` is **VRouter.
  - When the `FastLinkMode` parameter is set to `true`, the `Role` parameter must be set to `InitiatingSide`,`AccessPointId`,`OppositeRouterType`,`opppiterouterid`, and `OppositeInterfaceOwnerId` are required.
* `hc_rate` - (Optional, Int) Health check rate. Unit: milliseconds. The recommend value is 2000. Indicates the time interval for sending continuous detection packets during a specified health check.
* `hc_threshold` - (Optional) Health check threshold. Unit: One. The recommend value is 8. Indicates the number of detection packets sent during the specified health check.
* `health_check_source_ip` - (Optional) Health check source IP address
* `health_check_target_ip` - (Optional) Health check destination IP address
* `opposite_access_point_id` - (Optional, ForceNew) Peer access point ID
* `opposite_interface_owner_id` - (Optional, ForceNew) Account ID of the peer router interface
* `opposite_region_id` - (Required, ForceNew) Region of the connection peer
* `opposite_router_id` - (Optional, ForceNew) The ID of the router to which the opposite router interface belongs.
* `opposite_router_type` - (Optional, ForceNew, Computed) The router type associated with the peer router interface. Valid values:
  - VRouter: VPC router.
  - VBR: Virtual Border Router.
* `payment_type` - (Optional, ForceNew, Computed) The payment method of the router interface. Valid values:
  - Subscription : PrePaid.
  - PayAsYouGo : PostPaid.
* `period` - (Optional, Int) Purchase duration, value:
  - When you choose to pay on a monthly basis, the value range is **1 to 9**.
  - When you choose to pay per year, the value range is **1 to 3**.

-> **NOTE:**  `period` is required when the value of the parameter `payment_type` is `Subscription`.


-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `pricing_cycle` - (Optional) The billing cycle of the prepaid fee. Valid values:
  - `Month` (default): monthly payment.
  - `Year`: Pay per Year.


-> **NOTE:**  `period` is required when the value of the parameter `payment_type` is `Subscription`.


-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `resource_group_id` - (Optional, Computed) The ID of the resource group
* `role` - (Required, ForceNew) The role of the router interface. Valid values:
  - InitiatingSide : the initiator of the connection.
  - AcceptingSide : Connect to the receiving end.
* `router_id` - (Required, ForceNew) The ID of the router where the route entry is located.
* `router_interface_name` - (Optional) Resource attribute field representing the resource name. It must be 2 to 128 characters in length and must start with a letter or a Chinese character, but cannot start with http:// or https.
* `router_type` - (Required, ForceNew) The type of the router where the routing table resides. Valid values:
  - VRouter:VPC router
  - VBR: Border Router
* `spec` - (Required) The specification of the router interface. The available specifications and corresponding bandwidth values are as follows:
  - Mini.2: 2 Mbps
  - Mini.5: 5 Mbps
  - Small.1: 10 Mbps
  - Small.2: 20 Mbps
  - Small.5: 50 Mbps
  - Middle.1: 100 Mbps
  - Middle.2: 200 Mbps
  - Middle.5: 500 Mbps
  - Large.1: 1000 Mbps
  - Large.2: 2000 Mbps
  - Large.5: 5000 Mbps
  - Xlarge.1: 10000 Mbps

When the Role is AcceptingSide (connecting to the receiving end), the Spec value is Negative, which means that the specification is not involved in creating the receiving end router interface.
* `status` - (Optional, Computed) Resource attribute fields that represent the status of the resource. Value range:
  - Idle : Initialize.
  - Connecting : the initiator is in the process of Connecting.
  - AcceptingConnecting : the receiving end is being connected.
  - Activating : Restoring.
  - Active : Normal.
  - Modifying : Modifying.
  - Deactivating : Freezing.
  - Inactive : Frozen.
  - Deleting : Deleting.
  - Deleted : Deleted.
* `tags` - (Optional, Map) The tag of the resource

The following arguments will be discarded. Please use new fields as soon as possible:
* `auto_pay` - (Deprecated since v1.263.0). Field 'name' has been deprecated from provider version 1.263.0.
* `opposite_interface_id` - (Deprecated since v1.263.1). Field 'router_table_id' has been deprecated from provider version 1.263.0.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `bandwidth` - The bandwidth of the router interface
* `business_status` - The service status of the router interface. 
* `connected_time` - Time the connection was established
* `create_time` - The creation time of the resource
* `cross_border` - CrossBorder
* `end_time` - End Time of Prepaid
* `has_reservation_data` - Whether there is renewal data
* `opposite_bandwidth` - opposite bandwidth
* `opposite_interface_business_status` - The service status of the router interface on the opposite end of the connection. 
* `opposite_interface_spec` - Specifications of the interface of the peer router.
* `opposite_interface_status` - The status of the router interface on the peer of the connection.
* `opposite_vpc_instance_id` - The peer VPC ID
* `reservation_active_time` - ReservationActiveTime
* `reservation_bandwidth` - Renew Bandwidth
* `reservation_internet_charge_type` - Payment Type for Renewal
* `reservation_order_type` - Renewal Order Type
* `router_interface_id` - The first ID of the resource
* `vpc_instance_id` - ID of the local VPC in the peering connection

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Router Interface.
* `delete` - (Defaults to 5 mins) Used when delete the Router Interface.
* `update` - (Defaults to 5 mins) Used when update the Router Interface.

## Import

Express Connect Router Interface can be imported using the id, e.g.

```shell
$ terraform import alicloud_express_connect_router_interface.example <id>
```