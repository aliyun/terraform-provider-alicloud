---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_route_table"
sidebar_current: "docs-alicloud-resource-route-table"
description: |-
  Provides a Alicloud Route Table resource.
---

# alicloud\_route_table

Provides a route table resource to add customized route tables.

-> **NOTE:** Terraform will auto build route table instance while it uses `alicloud_route_table` to build a route table resource.

Currently, customized route tables are available in most regions apart from China (Beijing), China (Hangzhou), and China (Shenzhen) regions.
For information about route table and how to use it, see [What is Route Table](https://www.alibabacloud.com/help/doc-detail/87057.htm).

## Example Usage

Basic Usage

```terraform
resource "alicloud_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "vpc-example-name"
}

resource "alicloud_route_table" "foo" {
  vpc_id           = alicloud_vpc.foo.id
  route_table_name = "route-table-example-name"
  description      = "route-table-example-description"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) The vpc_id of the route table, the field can't be changed.
* `name` - (Optional) Field `name` has been deprecated from provider version 1.119.1. New field `route_table_name` instead.
* `route_table_name` - (Optional, Available in v1.119.1+) The name of the route table.
* `description` - (Optional) The description of the route table instance.
* `tags` - (Optional, Available in v1.55.3+) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the route table instance id.
* `status` - (Available in v1.119.1+) The status of the route table.

### Timeouts

-> **NOTE:** Available in 1.119.1+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the route table (until it reaches the initial `Available` status). 

## Import

The route table can be imported using the id, e.g.

```
$ terraform import alicloud_route_table.foo vtb-abc123456
```


