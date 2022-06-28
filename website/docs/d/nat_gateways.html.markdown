---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_nat_gateways"
sidebar_current: "docs-alicloud-datasource-nat-gateways"
description: |-
    Provides a list of Nat Gateways owned by an Alibaba Cloud account.
---

# alicloud\_nat\_gateways

This data source provides a list of Nat Gateways owned by an Alibaba Cloud account.

-> **NOTE:** Available in 1.37.0+.

## Example Usage

```
variable "name" {
  default = "natGatewaysDatasource"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "foo" {
  vpc_name   = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_nat_gateway" "foo" {
  vpc_id        = "${alicloud_vpc.foo.id}"
  specification = "Small"
  nat_gate_name = "${var.name}"
}

data "alicloud_nat_gateways" "foo" {
  vpc_id     = "${alicloud_vpc.foo.id}"
  name_regex = "${alicloud_nat_gateway.foo.name}"
  ids        = ["${alicloud_nat_gateway.foo.id}"]
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of NAT gateways IDs.
* `name_regex` - (Optional) A regex string to filter nat gateways by name.
* `vpc_id` - (Optional) The ID of the VPC.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `dry_run` - (Optional, ForceNew, Available in 1.121.0+) Specifies whether to only precheck the request.
* `nat_gateway_name` - (Optional, ForceNew, Available in 1.121.0+) The name of NAT gateway.
* `nat_type` - (Optional, ForceNew, Available in 1.121.0+) The nat type of NAT gateway. Valid values `Enhanced` and `Normal`.
* `payment_type` - (Optional, ForceNew, Available in 1.121.0+) The payment type of NAT gateway. Valid values `PayAsYouGo` and `Subscription`.
* `resource_group_id` - (Optional, ForceNew, Available in 1.121.0+) The resource group id of NAT gateway.
* `specification` - (Optional, ForceNew, Available in 1.121.0+) The specification of NAT gateway. Valid values `Middle`, `Large`, `Small` and `XLarge.1`. Default value is `Small`.
* `status` - (Optional, ForceNew, Available in 1.121.0+) The status of NAT gateway. Valid values `Available`, `Converting`, `Creating`, `Deleting` and `Modifying`.
* `tags` - (Optional, ForceNew, Available in 1.121.0+) The tags of NAT gateway.
* `enable_details` - (Optional, Available in 1.121.0+) Default to `false`. Set it to `true` can output more details about resource attributes.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - (Optional) A list of Nat gateways IDs.
* `names` - A list of Nat gateways names.
* `gateways` - A list of Nat gateways. Each element contains the following attributes:
  * `id` - The ID of the NAT gateway.
  * `name` - Name of the NAT gateway.
  * `description` - The description of the NAT gateway.
  * `creation_time` - (Deprecated form v1.121.0) Time of creation.
  * `spec` - The specification of the NAT gateway.
  * `status` - The status of the NAT gateway.
  * `snat_table_id` - Deprecated from v1.121.0, replace by snat_table_ids.
  * `snat_table_ids` - The ID of the SNAT table that is associated with the NAT gateway.
  * `forward_table_id` - Deprecated from v1.121.0, replace by forward_table_ids.
  * `forward_table_ids` - The ID of the DNAT table.
  * `vpc_id` - The ID of the VPC.
  * `ip_lists` - The ip address of the bind eip.
  * `business_status` - The state of the NAT gateway.
  * `deletion_protection` - Indicates whether deletion protection is enabled.
  * `ecs_metric_enabled` - Indicates whether the traffic monitoring feature is enabled.
  * `expired_time` - The time when the NAT gateway expires.
  * `internet_charge_type` - The metering method of the NAT gateway.  
  * `nat_gateway_id` - The ID of the NAT gateway.
  * `nat_gateway_name` - The name of the NAT gateway.
  * `network_type` - (Available in 1.137.0+) Indicates the type of the created NAT gateway. Valid values `internet` and `intranet`.
  * `nat_type` - The type of the NAT gateway. 
  * `payment_type` - The billing method of the NAT gateway. 
  * `resource_group_id` - The ID of the resource group.
  * `specification` - The specification of the NAT gateway.
  * `vswitch_id` - The ID of the vSwitch to which the NAT gateway belongs.
  * `tags` - The tags of NAT gateway.

