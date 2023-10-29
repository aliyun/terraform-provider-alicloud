---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint_service"
description: |-
  Provides a Alicloud Private Link Vpc Endpoint Service resource.
---

# alicloud_privatelink_vpc_endpoint_service

Provides a Private Link Vpc Endpoint Service resource. 

For information about Private Link Vpc Endpoint Service and how to use it, see [What is Vpc Endpoint Service](https://www.alibabacloud.com/help/en/privatelink/latest/api-privatelink-2020-04-15-createvpcendpointservice).

-> **NOTE:** Available since v1.109.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}

resource "alicloud_privatelink_vpc_endpoint_service" "example" {
  service_description    = var.name
  connect_bandwidth      = 103
  auto_accept_connection = false
}
```

## Argument Reference

The following arguments are supported:
* `auto_accept_connection` - (Optional) Whether to automatically accept terminal node connections.
* `connect_bandwidth` - (Optional) The connection bandwidth.
* `dry_run` - (Optional) Whether to pre-check this request only. Default to: `false`
* `payer` - (Optional, ForceNew) The payer type. Valid Value: `EndpointService`, `Endpoint`. Default to: `Endpoint`.
* `resource_group_id` - (Optional, Computed, Available since v1.212.0) The ID of the resource group.
* `service_description` - (Optional) The description of the terminal node service.
* `service_resource_type` - (Optional, ForceNew, Available since v1.212.0) The service type of resource. Valid values:
  - **slb**: indicates that the service resource type is classic load balancer (clb).
  - **alb**: indicates that the service resource type is application load balancer (alb).
  - **nlb**: indicates that the service resource type is network load balancer (nlb).
* `service_support_ipv6` - (Optional, Available since v1.212.0) Whether the endpoint supports IPv6. Value:
  - **true**: Yes.
  - **false** (default): No.
* `tags` - (Optional, Map, Available since v1.212.0) The tag of the resource.
* `zone_affinity_enabled` - (Optional, Available since v1.212.0) Whether to support the nearby resolution of the available area. Valid values:
  - **true**: Yes.
  - **false**: No.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.
* `service_business_status` - The business status of Vpc Endpoint Service.
* `service_domain` - Service Domain.
* `status` - The status of the resource.
* `vpc_endpoint_service_name` - VpcEndpointServiceName.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vpc Endpoint Service.
* `delete` - (Defaults to 5 mins) Used when delete the Vpc Endpoint Service.
* `update` - (Defaults to 5 mins) Used when update the Vpc Endpoint Service.

## Import

Private Link Vpc Endpoint Service can be imported using the id, e.g.

```shell
$ terraform import alicloud_privatelink_vpc_endpoint_service.example <id>
```