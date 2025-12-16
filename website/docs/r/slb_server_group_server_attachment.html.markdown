---
subcategory: "Classic Load Balancer (SLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_server_group_server_attachment"
description: |-
  Provides a Alicloud Load Balancer Virtual Backend Server Group Server Attachment resource.
---

# alicloud_slb_server_group_server_attachment

Provides a Load Balancer Virtual Backend Server Group Server Attachment resource.

-> **NOTE:** Available since v1.163.0.

For information about Load Balancer Virtual Backend Server Group Server Attachment and how to use it, see [What is Virtual Backend Server Group Server Attachment](https://www.alibabacloud.com/help/en/slb/classic-load-balancer/developer-reference/api-slb-2014-05-15-addvservergroupbackendservers).

-> **NOTE:** Applying this resource may conflict with applying `alicloud_slb_listener`, 
and the `alicloud_slb_listener` block should use `depends_on = [alicloud_slb_server_group_server_attachment.xxx]` to avoid it.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_slb_server_group_server_attachment&exampleId=39bab162-8a3c-1f1f-a90b-4a7f993c1a94a72d5800&activeTab=example&spm=docs.r.slb_server_group_server_attachment.0.39bab1628a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_slb_zones" "default" {
  available_slb_address_type = "vpc"
}

data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_slb_zones.default.zones.0.id
  instance_type_family = "ecs.sn1ne"
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.192.0/24"
  zone_id      = data.alicloud_slb_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = var.name
  vswitch_id         = alicloud_vswitch.default.id
  load_balancer_spec = "slb.s2.small"
  address_type       = "intranet"
}

resource "alicloud_slb_server_group" "default" {
  load_balancer_id = alicloud_slb_load_balancer.default.id
  name             = var.name
}

resource "alicloud_instance" "default" {
  image_id                   = data.alicloud_images.default.images[0].id
  instance_type              = data.alicloud_instance_types.default.instance_types[0].id
  instance_name              = var.name
  security_groups            = alicloud_security_group.default.*.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_slb_zones.default.zones.0.id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alicloud_vswitch.default.id
}

resource "alicloud_slb_server_group_server_attachment" "server_attachment" {
  server_group_id = alicloud_slb_server_group.default.id
  server_id       = alicloud_instance.default.id
  port            = 8080
  type            = "ecs"
  weight          = 0
  description     = var.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_slb_server_group_server_attachment&spm=docs.r.slb_server_group_server_attachment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `server_group_id` - (Required, ForceNew) The ID of the server group.
* `server_id` - (Required, ForceNew) The ID of the backend server. You can specify the ID of an Elastic Compute Service (ECS) instance or an elastic network interface (ENI).
* `port` - (Required, ForceNew) The port that is used by the backend server. Valid values: `1` to `65535`.
* `type` - (Optional, ForceNew) The type of backend server. Valid values: `ecs`, `eni`, `eci`. **NOTE:** From version 1.246.0, `type` can be set to `eci`.
* `weight` - (Optional, ForceNew) The weight of the backend server. Valid values: `0` to `100`. Default value: `100`. If the value is set to `0`, no requests are forwarded to the backend server.
* `description` - (Optional, ForceNew) The description of the backend server.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Virtual Backend Server Group Server Attachment. It formats as `<server_group_id>:<server_id>:<port>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Virtual Backend Server Group Server Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Virtual Backend Server Group Server Attachment.

## Import

Load Balancer Virtual Backend Server Group Server Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_slb_server_group_server_attachment.example <server_group_id>:<server_id>:<port>
```
