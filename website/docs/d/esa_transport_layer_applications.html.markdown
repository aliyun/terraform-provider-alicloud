---
subcategory: "esa"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_transport_layer_applications"
sidebar_current: "docs-alicloud-datasource-esa-transport-layer-applications"
description: |-
  Provides a list of Esa Transport Layer Application owned by an Alibaba Cloud account.
---

# alicloud_esa_transport_layer_applications

This data source provides Esa Transport Layer Application available to the user.[What is Transport Layer Application](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateTransportLayerApplication)

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
  site_name           = "gositecdn.cn"
}

data "alicloud_esa_transport_layer_applications" "default" {
  site_id        = data.alicloud_esa_sites.default.sites.0.site_id
  match_type     = "exact"
  record_name    = "resource2.gositecdn.cn"
  enable_details = true
}

output "first_application_id" {
  value = data.alicloud_esa_transport_layer_applications.default.applications.0.application_id
}
```

## Argument Reference

The following arguments are supported:
* `match_type` - (ForceNew, Optional) The following four query types are supported for four-layer application host records. The default is exact query.
  -`fuzzy`: fuzzy query.
  -`exact`: exact query.
  -`prefix`: prefix matching query.
  -`suffix`: suffix matching query.
* `record_name` - (ForceNew, Optional) The host record of the Layer - 4 application.
* `site_id` - (Required, ForceNew) Site ID.
* `ids` - (Optional, Computed) A list of Transport Layer Application IDs. The value is formulated as `<site_id>:<application_id>`.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Transport Layer Application IDs.
* `applications` - A list of Transport Layer Application Entries. Each element contains the following attributes:
  * `application_id` - application id.
  * `cname` - The CNAME domain name corresponding to the layer -4 accelerated application.
  * `cross_border_optimization` - CrossBorderOptimization.
  * `ip_access_rule` - IP access control.
  * `ipv6` - IPv6 access.
  * `keep_alive_protection` - Keep alive protection.
  * `record_name` - The host record of the Layer -4 application.
  * `rules` - The list of forwarding rules.
    * `client_ip_pass_through_mode` - Client IP delivery.
    * `comment` - Remarks.
    * `edge_port` - Edge Ports.
    * `protocol` - Agreement.
    * `rule_id` - Rule ID.
    * `source` - Source.
    * `source_port` - Source Port.
    * `source_type` - Source type.
  * `rules_count` - The number of forwarding rules contained in the Layer -4 acceleration application.
  * `site_id` - Site ID.
  * `static_ip` - Static IP.
  * `status` - **NOTE:** This field is only available when `enable_details` is `true`. The status of the resource.
  * `id` - The ID of the resource supplied above.
