---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_basic_endpoint_group"
sidebar_current: "docs-alicloud-resource-ga-basic-endpoint-group"
description: |-
  Provides a Alicloud Global Accelerator (GA) Basic Endpoint Group resource.
---

# alicloud_ga_basic_endpoint_group

Provides a Global Accelerator (GA) Basic Endpoint Group resource.

For information about Global Accelerator (GA) Basic Endpoint Group and how to use it, see [What is Basic Endpoint Group](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-createbasicendpointgroup).

-> **NOTE:** Available since v1.194.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ga_basic_endpoint_group&exampleId=1a0bd89d-569a-7fc8-0a7e-2ea8935e46ab52cdc787&activeTab=example&spm=docs.r.ga_basic_endpoint_group.0.1a0bd89d56&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "region" {
  default = "cn-hangzhou"
}

variable "endpoint_group_region" {
  default = "cn-beijing"
}

provider "alicloud" {
  region  = var.region
  profile = "default"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = "terraform-example"
  vswitch_id         = alicloud_vswitch.default.id
  load_balancer_spec = "slb.s2.small"
  address_type       = "intranet"
}

resource "alicloud_ga_basic_accelerator" "default" {
  duration               = 1
  basic_accelerator_name = "terraform-example"
  description            = "terraform-example"
  bandwidth_billing_type = "CDT"
  auto_use_coupon        = "true"
  auto_pay               = true
}

resource "alicloud_ga_basic_endpoint_group" "default" {
  accelerator_id            = alicloud_ga_basic_accelerator.default.id
  endpoint_group_region     = var.endpoint_group_region
  endpoint_type             = "SLB"
  endpoint_address          = alicloud_slb_load_balancer.default.id
  endpoint_sub_address      = "192.168.0.1"
  basic_endpoint_group_name = "terraform-example"
  description               = "terraform-example"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ga_basic_endpoint_group&spm=docs.r.ga_basic_endpoint_group.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the basic GA instance.
* `endpoint_group_region` - (Required, ForceNew) The ID of the region where you want to create the endpoint group.
* `endpoint_type` - (Optional, ForceNew) The type of the endpoint. Valid values: `ENI`, `SLB` and `ECS`.
* `endpoint_address` - (Optional, ForceNew) The address of the endpoint.
* `endpoint_sub_address` - (Optional, ForceNew) The sub address of the endpoint.
* `basic_endpoint_group_name` - (Optional) The name of the endpoint group. The `basic_endpoint_group_name` must be 2 to 128 characters in length, and can contain letters, digits, underscores (_), and hyphens (-). The name must start with a letter.
* `description` - (Optional) The description of the endpoint group. The `description` cannot exceed 256 characters in length and cannot contain http:// or https://.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Basic Endpoint Group.
* `status` - The status of the Basic Endpoint Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Basic Endpoint Group.
* `update` - (Defaults to 3 mins) Used when update the Basic Endpoint Group.
* `delete` - (Defaults to 3 mins) Used when delete the Basic Endpoint Group.

## Import

Global Accelerator (GA) Basic Endpoint Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_basic_endpoint_group.example <id>
```
