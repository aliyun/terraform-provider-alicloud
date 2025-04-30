---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_route_table_attachment"
description: |-
  Provides a Alicloud VPC Route Table Attachment resource.
---

# alicloud_route_table_attachment

Provides a VPC Route Table Attachment resource. Routing table associated resource type.

For information about VPC Route Table Attachment and how to use it, see [What is Route Table Attachment](https://www.alibabacloud.com/help/doc-detail/174112.htm).

-> **NOTE:** Available since v1.194.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_route_table_attachment&exampleId=0a510241-b12e-dc99-51e3-77d5739c450b06faab9f&activeTab=example&spm=docs.r.route_table_attachment.0.0a510241b1&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_vpc" "foo" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vswitch" "foo" {
  vpc_id       = alicloud_vpc.foo.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_route_table" "foo" {
  vpc_id           = alicloud_vpc.foo.id
  route_table_name = var.name
  description      = "route_table_attachment"
}

resource "alicloud_route_table_attachment" "foo" {
  vswitch_id     = alicloud_vswitch.foo.id
  route_table_id = alicloud_route_table.foo.id
}
```

## Argument Reference

The following arguments are supported:
* `route_table_id` - (Required, ForceNew) The ID of the route table to be bound to the switch.
* `vswitch_id` - (Required, ForceNew) The ID of the switch to bind the route table.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<route_table_id>:<vswitch_id>`.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Route Table Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Route Table Attachment.

## Import

VPC Route Table Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_route_table_attachment.example <route_table_id>:<vswitch_id>
```