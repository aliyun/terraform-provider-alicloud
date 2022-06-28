---
subcategory: "Database Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_database_gateway_gateway"
sidebar_current: "docs-alicloud-resource-database-gateway-gateway"
description: |-
  Provides a Alicloud Database Gateway Gateway resource.
---

# alicloud\_database\_gateway\_gateway

Provides a Database Gateway Gateway resource.

For information about Database Gateway Gateway and how to use it, see [What is Gateway](https://www.alibabacloud.com/help/doc-detail/123415.htm).

-> **NOTE:** Available in v1.135.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_database_gateway_gateway" "example" {
  gateway_name = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `gateway_desc` - (Optional) The description of Gateway.
* `gateway_name` - (Required) The name of the Gateway.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Gateway.
* `status` - The status of gateway. Valid values: `EXCEPTION`, `NEW`, `RUNNING`, `STOPPED`.

## Import

Database Gateway Gateway can be imported using the id, e.g.

```
$ terraform import alicloud_database_gateway_gateway.example <id>
```
