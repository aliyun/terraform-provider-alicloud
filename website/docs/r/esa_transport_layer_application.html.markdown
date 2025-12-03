---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_transport_layer_application"
description: |-
  Provides a Alicloud ESA Transport Layer Application resource.
---

# alicloud_esa_transport_layer_application

Provides a ESA Transport Layer Application resource.

Transport Layer Acceleration Application.

For information about ESA Transport Layer Application and how to use it, see [What is Transport Layer Application](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateTransportLayerApplication).

-> **NOTE:** Available since v1.260.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
  site_name           = "gositecdn.cn"
}

resource "alicloud_esa_transport_layer_application" "default" {
  record_name               = "resource2.gositecdn.cn"
  site_id                   = data.alicloud_esa_sites.default.sites.0.site_id
  ip_access_rule            = "off"
  ipv6                      = "off"
  cross_border_optimization = "off"
  rules {
    source                      = "1.2.3.4"
    comment                     = "transportLayerApplication"
    edge_port                   = "80"
    source_type                 = "ip"
    protocol                    = "TCP"
    source_port                 = "8080"
    client_ip_pass_through_mode = "off"
  }
}
```

### Deleting `alicloud_esa_transport_layer_application` or removing it from your configuration

The `alicloud_esa_transport_layer_application` resource allows you to manage  `status = "active"`  instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Instance.
You can resume managing the subscription instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:
* `cross_border_optimization` - (Optional) Whether to enable China mainland network access optimization, default is disabled. Value range:
  - `on`: Enabled.
  - `off`: Disabled.
* `ip_access_rule` - (Optional) IP access rule switch. When enabled, the WAF's IP access rules apply to the transport layer application.
  - `on`: Enabled.
  - `off`: Disabled.
* `ipv6` - (Optional) IPv6 switch.
* `record_name` - (Required, ForceNew) Domain name of the transport layer application
* `rules` - (Required, List) The list of forwarding rules. Rule details. For each rule, other parameters are required except comments. See [`rules`](#rules) below.
* `site_id` - (Required, ForceNew) Site ID.

### `rules`

The rules supports the following:
* `client_ip_pass_through_mode` - (Required) Client IP pass-through protocol, supporting:
  - `off`: No pass-through.
  - `PPv1`: PROXY Protocol v1, supports client IP pass-through for TCP protocol.
  - `PPv2`: PROXY Protocol v2, supports client IP pass-through for TCP and UDP protocols.
  - `SPP`: Simple Proxy Protocol, supports client IP pass-through for UDP protocol.
* `comment` - (Optional) Comment information for the rule (optional).
* `edge_port` - (Required) Edge port. Supports:
  - A single port, such as 80.
  - Port range, such as 81-85, representing ports 81, 82, 83, 84, and 85.
  - Combination of ports and port ranges, separated by commas, such as 80,81-85,90, representing ports 80, 81, 82, 83, 84, 85, and 90.

Edge ports within a single rule and between multiple rules must not overlap.

* `protocol` - (Required) Forwarding rule protocol, with values:
  - `TCP`: TCP protocol.
  - `UDP`: UDP protocol.
* `source` - (Required) Specific value of the origin, which needs to match the origin type.
* `source_port` - (Required) Source Port
* `source_type` - (Required) Origin type, supporting:
  - `ip`: IP address.
  - `domain`: Domain name.
  - `OP`: Origin pool.
  - `LB`: Load balancer.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<application_id>`.
* `application_id` - Layer 4 application ID.
* `rules` - The list of forwarding rules. Rule details. For each rule, other parameters are required except comments.
  * `rule_id` - Rule ID
* `status` - Status of the transport layer application, modification and deletion are not allowed.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 25 mins) Used when create the Transport Layer Application.
* `delete` - (Defaults to 5 mins) Used when delete the Transport Layer Application.
* `update` - (Defaults to 17 mins) Used when update the Transport Layer Application.

## Import

ESA Transport Layer Application can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_transport_layer_application.example <site_id>:<application_id>
```