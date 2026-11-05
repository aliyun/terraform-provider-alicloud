---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_dataworks_route_rule"
description: |-
  Provides a Alicloud Data Works Route Rule resource.
---

# alicloud_dataworks_route_rule

Provides a Data Works Route Rule resource.

Manages a route rule on a DataWorks network resource. A route rule routes traffic from the DataWorks network resource to a destination CIDR of type `VPC` or `IDC`. The rule is addressed by the DataWorks resource group and the network resource it is attached to; `resource_group_id` and `resource_id` select that group and network and are immutable after creation.

For information about Data Works Route Rule and how to use it, see [What is Route Rule](https://next.api.alibabacloud.com/document/dataworks-public/2018-06-01/CreateRouteRule).

-> **NOTE:** Available since v2.3.0.

-> **NOTE:** The RouteRule APIs (`CreateRouteRule`, `UpdateRouteRule`, `DeleteRouteRule`, `ListRouteRules`) are `visibility=Private`. This resource implements their full lifecycle; public availability of the underlying APIs is governed by the DataWorks product owner and does not block Terraform-side development.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-shenzhen"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

data "alicloud_vpcs" "default2" {
  name_regex = "^default-NODELETING-2$"
}

data "alicloud_vswitches" "default2" {
  vpc_id = data.alicloud_vpcs.default2.ids.0
}

resource "alicloud_data_works_dw_resource_group" "default" {
  payment_duration_unit = "Month"
  payment_type          = "PostPaid"
  specification         = "500"
  default_vswitch_id    = data.alicloud_vswitches.default.ids.0
  remark                = "tf-example resource group"
  resource_group_name   = var.name
  default_vpc_id        = data.alicloud_vpcs.default.ids.0
  auto_renew            = false
}

resource "alicloud_data_works_network" "default" {
  dw_resource_group_id = alicloud_data_works_dw_resource_group.default.id
  vpc_id               = data.alicloud_vpcs.default2.ids.0
  vswitch_id           = data.alicloud_vswitches.default2.ids.0
}

resource "alicloud_dataworks_route_rule" "default" {
  resource_group_id = alicloud_data_works_dw_resource_group.default.id
  resource_id       = alicloud_data_works_network.default.id
  destination_cidr  = "192.168.10.0/24"
  destination_type  = "VPC"
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_id` - (Required, ForceNew) The DataWorks resource group id the network resource (and its route rules) belong to.
* `resource_id` - (Required, ForceNew) The network resource id (the `alicloud_data_works_network` id) the route rule is attached to. `resource_group_id` and `resource_id` are replacement-only because `UpdateRouteRule` does not accept them; `ListRouteRules` is addressed by this pair.
* `destination_cidr` - (Required, String) The destination CIDR of the route rule, e.g. `192.168.10.0/24`.
* `destination_type` - (Required, String) The destination type of the route rule. Valid values: `VPC`, `IDC`.

## Attributes Reference

The following attributes are exported:

* `id` - (Computed, String) The composite id of the route rule: `<resource_group_id>:<resource_id>:<route_rule_id>`. `ListRouteRules` has no single Get API and is addressed by the resource group and network resource, so the composite keeps Read and Import self-contained. The middle segment (`resource_id`) may itself contain a colon; it is the part between the first and last colon.
* `route_rule_id` - (Computed, String) The route rule id returned by `CreateRouteRule` and used to address `UpdateRouteRule` and `DeleteRouteRule`.
* `destination_addressing_type` - (Computed, String) The destination addressing type returned by `ListRouteRules`. Valid values: `CIDR`, `IP`. Computed because `CreateRouteRule` and `UpdateRouteRule` do not accept it.

## Import

Data Works Route Rule can be imported using the composite id, e.g.

```shell
$ terraform import alicloud_dataworks_route_rule.example <resource_group_id>:<resource_id>:<route_rule_id>
```
