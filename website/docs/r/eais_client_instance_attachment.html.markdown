---
subcategory: "Elastic Accelerated Computing Instances (EAIS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eais_client_instance_attachment"
description: |-
  Provides a Alicloud EAIS Client Instance Attachment resource.
---

# alicloud_eais_client_instance_attachment

Provides a EAIS Client Instance Attachment resource.

Bind an ECS or ECI instance.

For information about EAIS Client Instance Attachment and how to use it, see [What is Client Instance Attachment](https://www.alibabacloud.com/help/en/resource-orchestration-service/latest/aliyun-eais-clientinstanceattachment).

-> **NOTE:** Available since v1.246.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_eais_client_instance_attachment&exampleId=2848594c-4550-b39b-1cd9-a3b990cc5cf9f0066c29&activeTab=example&spm=docs.r.eais_client_instance_attachment.0.2848594c45&intl_lang=EN_US" target="_blank">
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

variable "zone" {
  default = "cn-hangzhou-i"
}

variable "ecs_image" {
  default = "ubuntu_20_04_x64_20G_alibase_20230316.vhd"
}

variable "ecs_type" {
  default = "ecs.g7.large"
}

variable "region" {
  default = "cn-hangzhou"
}

variable "category" {
  default = "ei"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "example" {
  availability_zone = "cn-hangzhou-i"
  cpu_core_count    = 1
  memory_size       = 2
}

data "alicloud_images" "example" {
  name_regex = "^ubuntu_18.*64"
  owners     = "system"
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = "cn-hangzhou-i"
}

resource "alicloud_security_group" "example" {
  security_group_name = var.name
  description         = var.name
  vpc_id              = alicloud_vpc.example.id
}

resource "alicloud_instance" "example" {
  availability_zone          = "cn-hangzhou-i"
  vswitch_id                 = alicloud_vswitch.example.id
  image_id                   = data.alicloud_images.example.images.0.id
  instance_type              = data.alicloud_instance_types.example.instance_types.0.id
  system_disk_category       = "cloud_efficiency"
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = 5
  security_groups            = [alicloud_security_group.example.id]
  instance_name              = var.name
  user_data                  = "echo 'net.ipv4.ip_forward=1'>> /etc/sysctl.conf"
}

resource "alicloud_eais_instance" "eais" {
  instance_name     = var.name
  vswitch_id        = alicloud_vswitch.example.id
  security_group_id = alicloud_security_group.example.id
  instance_type     = "eais.ei-a6.2xlarge"
  category          = "ei"
}


resource "alicloud_eais_client_instance_attachment" "default" {
  instance_id        = alicloud_eais_instance.eais.id
  client_instance_id = alicloud_instance.example.id
  category           = "ei"
  status             = "Bound"
  ei_instance_type   = "eais.ei-a6.2xlarge"
}
```

### Deleting `alicloud_eais_client_instance_attachment` or removing it from your configuration

The `alicloud_eais_client_instance_attachment` resource allows you to manage  `category = "eais"`  instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Instance.
You can resume managing the subscription instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:
* `category` - (Optional) EAIS instance category, valid values: `eais`, `ei`, default is `eais`.
* `client_instance_id` - (Required, ForceNew) The ID of the ECS or ECI instance bound to the EAIS instance.
* `ei_instance_type` - (Optional, ForceNew) The Ei instance specification, which is used to filter matching specifications for updating.
* `instance_id` - (Required, ForceNew) The EAIS instance ID.
* `status` - (Optional, Computed) The status of the resource

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<instance_id>:<client_instance_id>`.
* `create_time` - The creation time of the resource
* `region_id` - The region ID of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Client Instance Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Client Instance Attachment.
* `update` - (Defaults to 5 mins) Used when update the Client Instance Attachment.

## Import

EAIS Client Instance Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_eais_client_instance_attachment.example <instance_id>:<client_instance_id>
```