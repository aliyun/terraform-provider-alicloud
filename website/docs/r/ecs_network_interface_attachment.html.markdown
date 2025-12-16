---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_network_interface_attachment"
sidebar_current: "docs-alicloud-resource-ecs-network-interface-attachment"
description: |-
  Provides a Alicloud ECS Network Interface Attachment resource.
---

# alicloud_ecs_network_interface_attachment

Provides a ECS Network Interface Attachment resource.

For information about ECS Network Interface Attachment and how to use it, see [What is Network Interface Attachment](https://www.alibabacloud.com/help/en/doc-detail/58515.htm).

-> **NOTE:** Available since v1.123.1.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecs_network_interface_attachment&exampleId=5dfac723-b704-382e-9cd0-0b6ef0f827304d79bf95&activeTab=example&spm=docs.r.ecs_network_interface_attachment.0.5dfac723b7&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "Instance"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
  eni_amount        = 3
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/24"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "192.168.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vpc_id       = alicloud_vpc.default.id
}

resource "alicloud_security_group" "default" {
  name        = var.name
  description = "New security group"
  vpc_id      = alicloud_vpc.default.id
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alicloud_instance" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
  instance_name     = var.name
  host_name         = "tf-example"
  image_id          = data.alicloud_images.default.images.0.id
  instance_type     = data.alicloud_instance_types.default.instance_types.0.id
  security_groups   = [alicloud_security_group.default.id]
  vswitch_id        = alicloud_vswitch.default.id
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_ecs_network_interface" "default" {
  network_interface_name = var.name
  vswitch_id             = alicloud_vswitch.default.id
  security_group_ids     = [alicloud_security_group.default.id]
  description            = "Basic example"
  primary_ip_address     = "192.168.0.2"
  tags = {
    Created = "TF",
    For     = "example",
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}

resource "alicloud_ecs_network_interface_attachment" "default" {
  network_interface_id = alicloud_ecs_network_interface.default.id
  instance_id          = alicloud_instance.default.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ecs_network_interface_attachment&spm=docs.r.ecs_network_interface_attachment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `network_interface_id` - (Required, ForceNew)  The ID of the network interface.
* `instance_id` - (Required, ForceNew) The ID of the ECS instance.
* `trunk_network_instance_id` - (Optional, ForceNew) The ID of the trunk network instance.
* `network_card_index` - (Optional, ForceNew, Int, Available since v1.223.1) The index of the network card.
* `wait_for_network_configuration_ready` - (Optional, Bool) The wait for network configuration ready.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Network Interface Attachment. It formats as `<network_interface_id>:<instance_id>`.

## Timeouts

-> **NOTE:** Available since v1.223.1.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Network Interface Attachment.
* `delete` - (Defaults to 1 mins) Used when delete the Network Interface Attachment.

## Import

ECS Network Interface Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_network_interface_attachment.example <network_interface_id>:<instance_id>
```
