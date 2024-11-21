---
subcategory: "Classic Load Balancer (SLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_server_group_server_attachment"
sidebar_current: "docs-alicloud-resource-slb-server-group-server-attachment"
description: |-
  Provides a Load Banlancer Virtual Backend Server Group Server Attachment resource.
---

# alicloud\_slb\_server\_group\_server\_attachment

-> **NOTE:** Available in v1.163.0+.

For information about server group server attachment and how to use it, see [Configure a server group server attachment](https://www.alibabacloud.com/help/en/doc-detail/35218.html).

-> **NOTE:** Applying this resource may conflict with applying `alicloud_slb_listener`, 
and the `alicloud_slb_listener` block should use `depends_on = [alicloud_slb_server_group_server_attachment.xxx]` to avoid it.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_slb_server_group_server_attachment&exampleId=9da2adb6-0785-e634-d06d-b1838eaac79a1252ba4e&activeTab=example&spm=docs.r.slb_server_group_server_attachment.0.9da2adb607&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "slb_server_group_server_attachment" {
  default = "terraform-example"
}

variable "slb_server_group_server_attachment_count" {
  default = 5
}

data "alicloud_zones" "server_attachment" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "server_attachment" {
  availability_zone = data.alicloud_zones.server_attachment.zones[0].id
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "server_attachment" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}


resource "alicloud_vpc" "server_attachment" {
  vpc_name   = var.slb_server_group_server_attachment
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "server_attachment" {
  vswitch_name = var.slb_server_group_server_attachment
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.server_attachment.id
  zone_id      = data.alicloud_zones.server_attachment.zones.0.id
}

resource "alicloud_security_group" "server_attachment" {
  name   = var.slb_server_group_server_attachment
  vpc_id = alicloud_vpc.server_attachment.id
}

resource "alicloud_instance" "server_attachment" {
  count                      = var.slb_server_group_server_attachment_count
  image_id                   = data.alicloud_images.server_attachment.images[0].id
  instance_type              = data.alicloud_instance_types.server_attachment.instance_types[0].id
  instance_name              = var.slb_server_group_server_attachment
  security_groups            = alicloud_security_group.server_attachment.*.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_zones.server_attachment.zones[0].id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alicloud_vswitch.server_attachment.id
}

resource "alicloud_slb_load_balancer" "server_attachment" {
  load_balancer_name = var.slb_server_group_server_attachment
  vswitch_id         = alicloud_vswitch.server_attachment.id
  load_balancer_spec = "slb.s2.small"
  address_type       = "intranet"
}

resource "alicloud_slb_server_group" "server_attachment" {
  load_balancer_id = alicloud_slb_load_balancer.server_attachment.id
  name             = var.slb_server_group_server_attachment
}

resource "alicloud_slb_server_group_server_attachment" "server_attachment" {
  count           = var.slb_server_group_server_attachment_count
  server_group_id = alicloud_slb_server_group.server_attachment.id
  server_id       = alicloud_instance.server_attachment[count.index].id
  port            = 8080
  weight          = 0
}

```

## Argument Reference

The following arguments are supported:

* `server_group_id` - (Required, ForceNew) The ID of the server group.
* `server_id` - (Required, ForceNew) The ID of the backend server. You can specify the ID of an Elastic Compute Service (ECS) instance or an elastic network interface (ENI).
* `port` - (Required, ForceNew) The port that is used by the backend server. Valid values: `1` to `65535`.
* `weight` - (Optional, ForceNew, Computed) The weight of the backend server. Valid values: `0` to `100`. Default value: `100`. If the value is set to `0`, no requests are forwarded to the backend server.
* `type` - (Optional, ForceNew, Computed) The type of backend server. Valid values: `ecs`, `eni`.
* `description` - (Optional, ForceNew, Computed) The description of the backend server.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the virtual server group server attachment. The value formats as `<server_group_id>:<server_id>:<port>`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the resource.
* `delete` - (Defaults to 5 mins) Used when delete the resource.


## Import

Load balancer backend server group server attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_slb_server_group_server_attachment.example <server_group_id>:<server_id>:<port>
```
