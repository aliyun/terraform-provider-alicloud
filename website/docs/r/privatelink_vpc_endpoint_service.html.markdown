---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint_service"
sidebar_current: "docs-alicloud-resource-privatelink-vpc-endpoint-service"
description: |-
  Provides a Alicloud Private Link Vpc Endpoint Service resource.
---

# alicloud\_privatelink\_vpc\_endpoint\_service

Provides a Private Link Vpc Endpoint Service resource.

For information about Private Link Vpc Endpoint Service and how to use it, see [What is Vpc Endpoint Service](https://help.aliyun.com/document_detail/183540.html).

-> **NOTE:** Available in v1.109.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_privatelink_vpc_endpoint_service" "example" {
  service_description    = "tftest"
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
* `service_description` - (Optional) The description of the terminal node service.

-> **NOTE:** The `resources` only support load balancing instance with private network type and PrivateLink function.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Vpc Endpoint Service. Value as `service_id`.
* `service_business_status` - The business status of Vpc Endpoint Service.
* `service_domain` - Service Domain.
* `status` - The status of Vpc Endpoint Service.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 6 mins) Used when create the Vpc Endpoint Service.
* `update` - (Defaults to 4 mins) Used when update the Vpc Endpoint Service.

## Import

Private Link Vpc Endpoint Service can be imported using the id, e.g.

```
$ terraform import alicloud_privatelink_vpc_endpoint_service.example <service_id>
```
