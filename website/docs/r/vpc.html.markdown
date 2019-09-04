---
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc"
sidebar_current: "docs-alicloud-resource-vpc"
description: |-
  Provides a Alicloud VPC resource.
---

# alicloud\_vpc

Provides a VPC resource.

-> **NOTE:** Terraform will auto build a router and a route table while it uses `alicloud_vpc` to build a vpc resource.

## Example Usage

Basic Usage

```
resource "alicloud_vpc" "vpc" {
  name       = "tf_test_foo"
  cidr_block = "172.16.0.0/12"
}
```
## Argument Reference

The following arguments are supported:

* `cidr_block` - (Required, ForceNew) The CIDR block for the VPC.
* `name` - (Optional) The name of the VPC. Defaults to null.
* `description` - (Optional) The VPC description. Defaults to null.
* `resource_group_id` - (Optional, Available in 1.40.0+) The Id of resource group which the VPC belongs.
* `tags` - (Optional, Available in v1.55.3+) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the VPC.
* `cidr_block` - The CIDR block for the VPC.
* `name` - The name of the VPC.
* `description` - The description of the VPC.
* `router_id` - The ID of the router created by default on VPC creation.
* `route_table_id` - The route table ID of the router created by default on VPC creation.

## Import

VPC can be imported using the id, e.g.

```
$ terraform import alicloud_vpc.example vpc-abc123456
```

