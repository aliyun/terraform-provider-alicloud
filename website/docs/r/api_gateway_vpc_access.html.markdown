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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_api_gateway_vpc_access&exampleId=b5341684-a6a0-1482-f91e-797559cda57c9e1cb314&activeTab=example&spm=docs.r.api_gateway_vpc_access.0.b5341684a6&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
  name_regex = "^ubuntu_18.*64"
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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_api_gateway_vpc_access&spm=docs.r.api_gateway_vpc_access.example&intl_lang=EN_US)
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
