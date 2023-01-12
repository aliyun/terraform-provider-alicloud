---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_custom_routing_endpoint"
sidebar_current: "docs-alicloud-resource-ga-custom-routing-endpoint"
description: |-
  Provides a Alicloud Global Accelerator (GA) Custom Routing Endpoint resource.
---

# alicloud\_ga\_custom\_routing\_endpoint

Provides a Global Accelerator (GA) Custom Routing Endpoint resource.

For information about Global Accelerator (GA) Custom Routing Endpoint and how to use it, see [What is Custom Routing Endpoint](https://www.alibabacloud.com/help/en/global-accelerator/latest/createcustomroutingendpoints).

-> **NOTE:** Available in v1.197.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ga_custom_routing_endpoint" "default" {
  endpoint_group_id          = "your_custom_routing_endpoint_group_id"
  endpoint                   = "your_vswitch_id"
  type                       = "PrivateSubNet"
  traffic_to_endpoint_policy = "DenyAll"
}
```

## Argument Reference

The following arguments are supported:

* `endpoint_group_id` - (Required, ForceNew) The ID of the endpoint group in which to create endpoints.
* `endpoint` - (Required, ForceNew) The ID of the endpoint (vSwitch).
* `type` - (Required, ForceNew) The backend service type of the endpoint. Valid values: `PrivateSubNet`.
* `traffic_to_endpoint_policy` - (Optional, Computed) The access policy of traffic to the endpoint. Default value: `DenyAll`. Valid values:
  - `DenyAll`: denies all traffic to the endpoint.
  - `AllowAll`: allows all traffic to the endpoint.
  - `AllowCustom`: allows traffic only to specified destinations in the endpoint.
  
## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Custom Routing Endpoint. It formats as `<endpoint_group_id>:<custom_routing_endpoint_id>`.
* `accelerator_id` - The ID of the GA instance with which the endpoint is associated.
* `listener_id` - The ID of the listener with which the endpoint is associated.
* `custom_routing_endpoint_id` - The ID of the Custom Routing Endpoint.
* `status` - The status of the Custom Routing Endpoint.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Custom Routing Endpoint.
* `update` - (Defaults to 5 mins) Used when update the Custom Routing Endpoint.
* `delete` - (Defaults to 5 mins) Used when delete the Custom Routing Endpoint.

## Import

Global Accelerator (GA) Custom Routing Endpoint can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_custom_routing_endpoint.example <endpoint_group_id>:<custom_routing_endpoint_id>
```
