---
subcategory: "Private Link"
layout: "alicloud"
page_title: "Alicloud: alicloud_privatelink_vpc_endpoint_connection"
sidebar_current: "docs-alicloud-resource-privatelink-vpc-endpoint-connection"
description: |-
  Provides a Alicloud Private Link Vpc Endpoint Connection resource.
---

# alicloud\_privatelink\_vpc\_endpoint\_connection

Provides a Private Link Vpc Endpoint Connection resource.

For information about Private Link Vpc Endpoint Connection and how to use it, see [What is Vpc Endpoint Connection](https://help.aliyun.com/document_detail/183551.html).

-> **NOTE:** Available in v1.110.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_privatelink_vpc_endpoint_connection" "example" {
  endpoint_id = "example_value"
  service_id  = "example_value"
  bandwidth   = "1024"
}

```

## Argument Reference

The following arguments are supported:

* `bandwidth` - (Optional) The Bandwidth.
* `dry_run` - (Optional) The dry run.
* `endpoint_id` - (Required, ForceNew) The ID of the Vpc Endpoint.
* `service_id` - (Required, ForceNew) The ID of the Vpc Endpoint Service.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Vpc Endpoint Connection. The value is formatted `<service_id>:<endpoint_id>`.
* `status` - The status of Vpc Endpoint Connection.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 11 mins) Used when create the Vpc Endpoint Connection.
* `delete` - (Defaults to 6 mins) Used when delete the Vpc Endpoint Connection.
* `update` - (Defaults to 4 mins) Used when update the Vpc Endpoint Connection.

## Import

Private Link Vpc Endpoint Connection can be imported using the id, e.g.

```
$ terraform import alicloud_privatelink_vpc_endpoint_connection.example <service_id>:<endpoint_id>
```
