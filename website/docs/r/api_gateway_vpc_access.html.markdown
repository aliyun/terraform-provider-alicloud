---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_vpc_access"
description: |-
  Provides a Alicloud Api Gateway Vpc Access resource.
---

# alicloud_api_gateway_vpc_access

Provides an Api Gateway Vpc Access resource.

For information about Api Gateway Vpc Access and how to use it, see [What is Vpc Access](https://www.alibabacloud.com/help/en/api-gateway/latest/api-cloudapi-2016-07-14-setvpcaccess)

-> **NOTE:** Available since v1.23.0.

-> **NOTE:** Terraform will auto build vpc authorization while it uses `alicloud_api_gateway_vpc_access` to build vpc.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_api_gateway_vpc_access&exampleId=b5341684-a6a0-1482-f91e-797559cda57c9e1cb314&activeTab=example&spm=docs.r.api_gateway_vpc_access.0.b5341684a6&intl_lang=EN_US" target="_blank">
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
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  image_id             = data.alicloud_images.default.images.0.id
  system_disk_category = "cloud_efficiency"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.192.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name   = var.name
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
  instance_name              = var.name
  description                = var.name
}

resource "alicloud_api_gateway_vpc_access" "default" {
  name        = var.name
  vpc_id      = alicloud_vpc.default.id
  instance_id = alicloud_instance.default.id
  port        = 8080
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_api_gateway_vpc_access&spm=docs.r.api_gateway_vpc_access.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the authorization. The name must be unique.
* `vpc_id` - (Required, ForceNew) The ID of the VPC. The VPC must be an available one that belongs to the same account as the API.
* `instance_id` - (Required, ForceNew) The ID of an ECS or SLB instance in the VPC.
* `port` - (Required, ForceNew) The port number that corresponds to the instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Vpc Access. It formats as `<name>:<vpc_id>:<instance_id>:<port>`.

## Import

Api Gateway Vpc Access can be imported using the id, e.g.

```shell
$ terraform import alicloud_api_gateway_vpc_access.example <name>:<vpc_id>:<instance_id>:<port>
```
