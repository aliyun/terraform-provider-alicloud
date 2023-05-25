---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_gateway_route_table_attachment"
sidebar_current: "docs-alicloud-resource-vpc-gateway-route-table-attachment"
description: |-
  Provides a Alicloud VPC Gateway Route Table Attachment resource.
---

# alicloud\_vpc\_gateway\_route\_table\_attachment

Provides a VPC Gateway Route Table Attachment resource.

For information about VPC Gateway Route Table Attachment and how to use it, see [What is Gateway Route Table Attachment](https://www.alibabacloud.com/help/doc-detail/174112.htm).

-> **NOTE:** Available in v1.194.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_vpc" "example" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "terraform-example"
}

resource "alicloud_route_table" "example" {
  vpc_id           = alicloud_vpc.example.id
  route_table_name = "terraform-example"
  description      = "terraform-example"
  associate_type   = "Gateway"
}

resource "alicloud_vpc_ipv4_gateway" "example" {
  ipv4_gateway_name = "terraform-example"
  vpc_id            = alicloud_vpc.example.id
  enabled           = "true"
}

resource "alicloud_vpc_gateway_route_table_attachment" "example" {
  ipv4_gateway_id = alicloud_vpc_ipv4_gateway.example.id
  route_table_id  = alicloud_route_table.example.id
}

```

## Argument Reference

The following arguments are supported:

* `dry_run` - (Optional) Specifies whether to only precheck this request. Default value: `false`.
* `ipv4_gateway_id` - (Required, ForceNew) The ID of the IPv4 Gateway instance.
* `route_table_id` - (Required, ForceNew) The ID of the Gateway route table to be bound.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Gateway Route Table Attachment. The value formats as `<route_table_id>:<ipv4_gateway_id>`.
* `status` - The status of the IPv4 Gateway instance. Value:

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Gateway Route Table Attachment.
* `delete` - (Defaults to 2 mins) Used when delete the Gateway Route Table Attachment.

## Import

VPC Gateway Route Table Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_gateway_route_table_attachment.example <route_table_id>:<ipv4_gateway_id>
```