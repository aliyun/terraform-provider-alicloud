---
subcategory: "esa"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_transport_layer_application"
description: |-
  Provides a Alicloud esa Transport Layer Application resource.
---

# alicloud_esa_transport_layer_application

Provides a esa Transport Layer Application resource.



For information about esa Transport Layer Application and how to use it, see [What is Transport Layer Application](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateTransportLayerApplication).

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
* `cross_border_optimization` - (Optional, Computed) CrossBorderOptimization
* `ip_access_rule` - (Optional, Computed) IP access control
* `ipv6` - (Optional, Computed) IPv6 access
* `keep_alive_protection` - (Optional, Computed, Available since v1.294.0) Keep alive protection. Valid values: `on`, `off`.
* `record_name` - (Required, ForceNew) The host record of the Layer - 4 application.
* `rules` - (Required, List) The list of forwarding rules. Rule details. For each rule, other parameters are required except comments. See [`rules`](#rules) below.
* `site_id` - (Required, ForceNew) Site ID.
* `static_ip` - (Optional, Computed, Available since v1.294.0) Static IP. Valid values: `on`, `off`.

### `rules`

The rules supports the following:
* `client_ip_pass_through_mode` - (Required) Client IP delivery
* `comment` - (Optional) Remarks
* `edge_port` - (Required) Edge Ports
* `protocol` - (Required) Agreement
* `source` - (Required) Source
* `source_port` - (Required) Source Port
* `source_type` - (Required) Source type

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<site_id>:<application_id>`.
* `application_id` - application id.
* `cname` - The CNAME domain name corresponding to the layer -4 accelerated application.
* `rules` - The list of forwarding rules.
  * `rule_id` - Rule ID.
* `rules_count` - The number of forwarding rules contained in the Layer -4 acceleration application.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 10 mins) Used when create the Transport Layer Application.
* `delete` - (Defaults to 10 mins) Used when delete the Transport Layer Application.
* `update` - (Defaults to 10 mins) Used when update the Transport Layer Application.

## Import

esa Transport Layer Application can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_transport_layer_application.example <site_id>:<application_id>
```