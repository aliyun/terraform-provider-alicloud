---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_vpc_access"
sidebar_current: "docs-alicloud-resource-api-gateway-vpc-access"
description: |-
  Provides a Alicloud Api Gateway vpc authorization Resource.
---

# alicloud_api_gateway_vpc_access

Provides an vpc authorization resource.This authorizes the API gateway to access your VPC instances.

For information about Api Gateway vpc and how to use it, see [Set Vpc Access](https://www.alibabacloud.com/help/en/api-gateway/latest/api-cloudapi-2016-07-14-setvpcaccess)

-> **NOTE:** Available since v1.23.0.

-> **NOTE:** Terraform will auto build vpc authorization while it uses `alicloud_api_gateway_vpc_access` to build vpc.

## Example Usage

Basic Usage

```terraform
data "alicloud_zones" "example" {
  available_resource_creation = "Instance"
}

data "alicloud_instance_types" "example" {
  availability_zone = data.alicloud_zones.example.zones.0.id
  cpu_core_count    = 1
  memory_size       = 2
}

resource "alicloud_vpc" "example" {
  vpc_name   = "terraform-example"
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = "terraform-example"
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = data.alicloud_zones.example.zones.0.id
}

resource "alicloud_security_group" "example" {
  name        = "terraform-example"
  description = "New security group"
  vpc_id      = alicloud_vpc.example.id
}

data "alicloud_images" "example" {
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
  owners     = "system"
}

resource "alicloud_instance" "example" {
  availability_zone = data.alicloud_zones.example.zones.0.id
  instance_name     = "terraform-example"
  image_id          = data.alicloud_images.example.images.0.id
  instance_type     = data.alicloud_instance_types.example.instance_types.0.id
  security_groups   = [alicloud_security_group.example.id]
  vswitch_id        = alicloud_vswitch.example.id
}

resource "alicloud_api_gateway_vpc_access" "example" {
  name        = "terraform-example"
  vpc_id      = alicloud_vpc.example.id
  instance_id = alicloud_instance.example.id
  port        = 8080
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the vpc authorization. 
* `vpc_id` - (Required, ForceNew) The vpc id of the vpc authorization. 
* `instance_id` - (Required, ForceNew) ID of the instance in VPC (ECS/Server Load Balance).
* `port` - (Required, ForceNew) ID of the port corresponding to the instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the vpc authorization of api gateway.

## Import

Api gateway app can be imported using the id, e.g.

```shell
$ terraform import alicloud_api_gateway_vpc_access.example "APiGatewayVpc:vpc-aswcj19ajsz:i-ajdjfsdlf:8080"
```
