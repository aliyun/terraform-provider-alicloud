---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_network_route_table"
description: |-
  Provides a Alicloud ENS Network Route Table resource.
---

# alicloud_ens_network_route_table

Provides a ENS Network Route Table resource.

Routing table of ENS network.

For information about ENS Network Route Table and how to use it, see [What is Network Route Table](https://next.api.alibabacloud.com/document/Ens/2017-11-10/CreateEnsRouteTable).

-> **NOTE:** Available since v1.291.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ens_network_route_table&exampleId=cac3e4e3-e786-5999-bed7-15c1a06189459654f882&activeTab=example&spm=docs.r.ens_network_route_table.0.cac3e4e3e7&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "ens_region_id" {
  default = "cn-hangzhou-31"
}

resource "alicloud_ens_network" "default" {
  network_name  = var.name
  description   = var.name
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_network_route_table" "default" {
  network_id       = alicloud_ens_network.default.id
  associate_type   = "Gateway"
  route_table_name = var.name
  description      = var.name
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ens_network_route_table&spm=docs.r.ens_network_route_table.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `network_id` - (Required, ForceNew) The network ID.
* `associate_type` - (Required, ForceNew) The binding type of the created routing table. Value: `Gateway`.
* `route_table_name` - (Optional) The name of the routing table. The length is 1 to 128 characters, but it cannot start with `http://` or `https://`.
* `description` - (Optional) The description of the routing table. The length is 1 to 256 characters, but it cannot start with `http://` or `https://`.
* `is_default_gateway_route_table` - (Computed) Whether it is the default gateway route table.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. Same as `route_table_id`.
* `route_table_id` - The ID of the routing table.
* `route_table_type` - The type of the routing table. Value: `Custom`, `System`.
* `status` - The status of the routing table.
* `create_time` - The creation time of the routing table.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Network Route Table.
* `delete` - (Defaults to 5 mins) Used when delete the Network Route Table.
* `update` - (Defaults to 5 mins) Used when update the Network Route Table.

## Import

ENS Network Route Table can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_network_route_table.example <route_table_id>
```
