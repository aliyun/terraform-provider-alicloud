---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_custom_routing_endpoint_group_destination"
sidebar_current: "docs-alicloud-resource-ga-custom-routing-endpoint-group-destination"
description: |-
  Provides a Alicloud Global Accelerator (GA) Custom Routing Endpoint Group Destination resource.
---

# alicloud\_ga\_custom\_routing\_endpoint\_group\_destination

Provides a Global Accelerator (GA) Custom Routing Endpoint Group Destination resource.

For information about Global Accelerator (GA) Custom Routing Endpoint Group Destination and how to use it, see [What is Custom Routing Endpoint Group Destination](https://www.alibabacloud.com/help/en/global-accelerator/latest/createcustomroutingendpointgroupdestinations).

-> **NOTE:** Available in v1.197.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ga_custom_routing_endpoint_group_destination" "default" {
  endpoint_group_id = "your_custom_routing_endpoint_group_id"
  protocols         = ["tcp", "udp"]
  from_port         = 1
  to_port           = 2
}
```

## Argument Reference

The following arguments are supported:

* `endpoint_group_id` - (Required, ForceNew) The ID of the endpoint group.
* `protocols` - (Required) The backend service protocol of the endpoint group. Valid values: `tcp`, `udp`, `tcp, udp`.
* `from_port` - (Required) The start port of the backend service port range of the endpoint group. The `from_port` value must be smaller than or equal to the `to_port` value. Valid values: `1` to `65499`.
* `to_port` - (Required) The end port of the backend service port range of the endpoint group. The `from_port` value must be smaller than or equal to the `to_port` value. Valid values: `1` to `65499`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Custom Routing Endpoint Group Destination. It formats as `<endpoint_group_id>:<custom_routing_endpoint_group_destination_id>`.
* `accelerator_id` - The ID of the GA instance.
* `listener_id` - The ID of the listener.
* `custom_routing_endpoint_group_destination_id` - The ID of the Custom Routing Endpoint Group Destination.
* `status` - The status of the Custom Routing Endpoint Group Destination.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Custom Routing Endpoint Group Destination.
* `update` - (Defaults to 5 mins) Used when update the Custom Routing Endpoint Group Destination.
* `delete` - (Defaults to 5 mins) Used when delete the Custom Routing Endpoint Group Destination.

## Import

Global Accelerator (GA) Custom Routing Endpoint Group Destination can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_custom_routing_endpoint_group_destination.example <endpoint_group_id>:<custom_routing_endpoint_group_destination_id>
```
