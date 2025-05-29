---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_ram_role_attachment"
description: |-
  Provides a Alicloud ECS Ram Role Attachment resource.
---

# alicloud_ecs_ram_role_attachment

Provides a ECS Ram Role Attachment resource.

Mount RAM role.

For information about ECS Ram Role Attachment and how to use it, see [What is Ram Role Attachment](https://next.api.alibabacloud.com/document/Ecs/2014-05-26/AttachInstanceRamRole).

-> **NOTE:** Available since v1.250.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecs_ram_role_attachment&exampleId=0ff709a5-e00b-8867-fb22-4cc62bb0ebef4f2f211e&activeTab=example&spm=docs.r.ecs_ram_role_attachment.0.0ff709a5e0&intl_lang=EN_US" target="_blank">
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

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_images" "default" {
  most_recent = true
  owners      = "system"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
  image_id          = data.alicloud_images.default.images.0.id
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_ram_role" "default" {
  name     = "${var.name}-${random_integer.default.result}"
  document = <<EOF
		{
			"Statement": [
				{
					"Action": "sts:AssumeRole",
					"Effect": "Allow",
					"Principal": {
						"Service": [
							"ecs.aliyuncs.com"
						]
					}
				}
		  	],
			"Version": "1"
		}
	  	EOF
  force    = true
}

resource "alicloud_vpc" "default" {
  vpc_name   = "${var.name}-${random_integer.default.result}"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = "${var.name}-${random_integer.default.result}"
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.192.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  image_id                   = data.alicloud_images.default.images.0.id
  instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  security_groups            = alicloud_security_group.default.*.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alicloud_vswitch.default.id
  instance_name              = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_ecs_ram_role_attachment" "default" {
  ram_role_name = alicloud_ram_role.default.id
  instance_id   = alicloud_instance.default.id
}
```

## Argument Reference

The following arguments are supported:
* `instance_id` - (Required, ForceNew) The ID of the instance.
* `policy` - (Optional) The additional policy. When you attach an instance RAM role to instances, you can specify an additional policy to further limit the permissions of the role.
* `ram_role_name` - (Required, ForceNew) The name of the instance RAM role.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<instance_id>:<ram_role_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Ram Role Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Ram Role Attachment.

## Import

ECS Ram Role Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_ram_role_attachment.example <instance_id>:<ram_role_name>
```
