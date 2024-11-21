---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_router_interface"
sidebar_current: "docs-alicloud-resource-express-connect--router-interface"
description: |-
  Provides a Alicloud Express Connect Router Interface resource.
---

# alicloud_express_connect_router_interface

Provides a Express Connect Router Interface resource.

For information about Express Connect Router Interface and how to use it, see [What is Router Interface](https://www.terraform.io/docs/providers/alicloud/r/router_interface_connection).

-> **NOTE:** Available since v1.199.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_express_connect_router_interface&exampleId=a1db9877-f2a0-0901-f367-7d0882d2d8c34857b390&activeTab=example&spm=docs.r.express_connect_router_interface.0.a1db9877f2&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_regions" "default" {
  current = true
}
resource "alicloud_express_connect_router_interface" "default" {
  description           = var.name
  opposite_region_id    = data.alicloud_regions.default.regions.0.id
  router_id             = data.alicloud_vpcs.default.vpcs.0.router_id
  role                  = "InitiatingSide"
  router_type           = "VRouter"
  payment_type          = "PayAsYouGo"
  router_interface_name = var.name
  spec                  = "Mini.2"
}
```

## Argument Reference

The following arguments are supported:
* `access_point_id` - (Optional, ForceNew) The access point ID to which the VBR belongs.
* `auto_pay` - (Optional) Whether to pay automatically, value:-**false** (default): automatic payment is not enabled. After generating an order, you need to complete the payment at the order center.-**true**: Enable automatic payment to automatically pay for orders.> **InstanceChargeType** is required when the value of the parameter is **PrePaid.
* `delete_health_check_ip` - (Optional) Whether to delete the health check IP address configured on the router interface. Value:-**true**: deletes the health check IP address.-**false** (default): does not delete the health check IP address.
* `description` - (Optional) The description of the router interface. The description must be 2 to 256 characters in length and cannot start with http:// or https://.
* `hc_rate` - (Optional) The health check rate. Unit: seconds. The recommended value is 2. This indicates the interval between successive probe messages sent during the specified health check.
* `hc_threshold` - (Optional) The health check thresholds. Unit: pcs. The recommended value is 8. This indicates the number of probe messages to be sent during the specified health check.
* `health_check_source_ip` - (Optional) The health check source IP address, must be an unused IP within the local VPC.
* `health_check_target_ip` - (Optional) The IP address for health screening purposes.
* `opposite_access_point_id` - (Optional, ForceNew) The Access point ID to which the other end belongs.
* `opposite_interface_id` - (Optional) The Interface ID of the router at the other end.
* `opposite_interface_owner_id` - (Optional) The AliCloud account ID of the owner of the router interface on the other end.
* `opposite_region_id` - (Required, ForceNew) The geographical ID of the location of the receiving end of the connection.
* `opposite_router_id` - (Optional) The id of the router at the other end.
* `opposite_router_type` - (Optional) The opposite router type of the router on the other side. Valid Values: `VRouter`, `VBR`.
* `payment_type` - (Optional, ForceNew) The payment methods for router interfaces. Valid Values: `PayAsYouGo`, `Subscription`.
* `period` - (Optional) Purchase duration, value:-When you choose to pay on a monthly basis, the value range is **1 to 9 * *.-When you choose to pay per year, the value range is **1 to 3 * *.> **InstanceChargeType** is required when the value of the parameter is **PrePaid.
* `pricing_cycle` - (Optional) The billing cycle of the prepaid fee. Valid values:-**Month** (default): monthly payment.-**Year**: Pay per Year.> **InstanceChargeType** is required when the value of the parameter is **PrePaid.
* `role` - (Required, ForceNew) The role of the router interface. Valid Values: `InitiatingSide`, `AcceptingSide`.
* `router_id` - (Required, ForceNew) The router id associated with the router interface.
* `router_interface_name` - (Optional) The name of the resource.
* `router_type` - (Required, ForceNew) The type of router associated with the router interface. Valid Values: `VRouter`, `VBR`.
* `spec` - (Required) The specification of the router interface. Valid Values: `Mini.2`, `Mini.5`, `Mini.5`, `Small.2`, `Small.5`, `Middle.1`, `Middle.2`, `Middle.5`, `Large.1`, `Large.2`, `Large.5`, `XLarge.1`, `Negative`.
* `status` - (Optional) The status of the resource. Valid Values: `Idle`, `AcceptingConnecting`, `Connecting`, `Activating`, `Active`, `Modifying`, `Deactivating`, `Inactive`, `Deleting`.
* `router_interface_id` - (Optional) The first ID of the resource.


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `bandwidth` - The bandwidth of the resource.
* `business_status` - The businessStatus of the resource. Valid Values: `Normal`, `FinancialLocked`, `SecurityLocked`.
* `connected_time` - The connected time of the resource.
* `create_time` - The creation time of the resource.
* `cross_border` - The cross border of the resource.
* `end_time` - The end time of the resource.
* `has_reservation_data` - The has reservation data of the resource.
* `opposite_bandwidth` - The opposite bandwidth of the router on the other side.
* `opposite_interface_business_status` - The opposite interface business status of the router on the other side. Valid Values: `Normal`, `FinancialLocked`, `SecurityLocked`.
* `opposite_interface_spec` - The opposite interface spec of the router on the other side. Valid Values: `Mini.2`, `Mini.5`, `Mini.5`, `Small.2`, `Small.5`, `Middle.1`, `Middle.2`, `Middle.5`, `Large.1`, `Large.2`, `Large.5`, `XLarge.1`, `Negative`.
* `opposite_interface_status` - The opposite interface status of the router on the other side. Valid Values: `Idle`, `AcceptingConnecting`, `Connecting`, `Activating`, `Active`, `Modifying`, `Deactivating`, `Inactive`, `Deleting`.
* `opposite_vpc_instance_id` - The opposite vpc instance id of the router on the other side.
* `reservation_active_time` - The reservation active time of the resource.
* `reservation_bandwidth` - The reservation bandwidth of the resource.
* `reservation_internet_charge_type` - The reservation internet charge type of the resource.
* `reservation_order_type` - The reservation order type of the resource.
* `vpc_instance_id` - The vpc instance id of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Router Interface.
* `delete` - (Defaults to 5 mins) Used when delete the Router Interface.
* `update` - (Defaults to 5 mins) Used when update the Router Interface.

## Import

Express Connect Router Interface can be imported using the id, e.g.

```shell
$ terraform import alicloud_expressconnect_router_interface.example <id>
```