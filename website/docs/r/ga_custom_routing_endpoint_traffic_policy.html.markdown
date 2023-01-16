---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_custom_routing_endpoint_traffic_policy"
sidebar_current: "docs-alicloud-resource-ga-custom-routing-endpoint-traffic-policy"
description: |-
  Provides a Alicloud Global Accelerator (GA) Custom Routing Endpoint Traffic Policy resource.
---

# alicloud\_ga\_custom\_routing\_endpoint\_traffic\_policy

Provides a Global Accelerator (GA) Custom Routing Endpoint Traffic Policy resource.

For information about Global Accelerator (GA) Custom Routing Endpoint Traffic Policy and how to use it, see [What is Custom Routing Endpoint Traffic Policy](https://www.alibabacloud.com/help/en/global-accelerator/latest/createcustomroutingendpointtrafficpolicies).

-> **NOTE:** Available in v1.197.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ga_custom_routing_endpoint_traffic_policy" "default" {
  endpoint_id = "your_custom_routing_endpoint_id"
  address     = "192.168.192.2"
  port_ranges {
    from_port = 10001
    to_port   = 10002
  }
}
```

## Argument Reference

The following arguments are supported:

* `endpoint_id` - (Required, ForceNew) The ID of the Custom Routing Endpoint.
* `address` - (Required) The IP address of the destination to which traffic is allowed.
* `port_ranges` - (Optional) Port rangeSee the following. See the following `Block port_ranges`.

#### Block port_ranges

The port_ranges supports the following:

* `from_port` - (Optional) The start port of the port range of the traffic destination. The specified port must fall within the port range of the specified endpoint group.
* `to_port` - (Optional) The end port of the port range of the traffic destination. The specified port must fall within the port range of the specified endpoint group.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Custom Routing Endpoint Traffic Policy. It formats as `<endpoint_id>:<custom_routing_endpoint_traffic_policy_id>`.
* `accelerator_id` - The ID of the GA instance.
* `listener_id` - The ID of the listener.
* `endpoint_group_id` - The ID of the endpoint group.
* `custom_routing_endpoint_traffic_policy_id` - The ID of the Custom Routing Endpoint Traffic Policy.
* `status` - The status of the Custom Routing Endpoint Traffic Policy.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Custom Routing Endpoint Traffic Policy.
* `update` - (Defaults to 5 mins) Used when update the Custom Routing Endpoint Traffic Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Custom Routing Endpoint Traffic Policy.

## Import

Global Accelerator (GA) Custom Routing Endpoint Traffic Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_custom_routing_endpoint_traffic_policy.example <endpoint_id>:<custom_routing_endpoint_traffic_policy_id>
```
