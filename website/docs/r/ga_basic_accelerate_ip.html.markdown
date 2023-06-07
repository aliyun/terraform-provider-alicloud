---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_basic_accelerate_ip"
sidebar_current: "docs-alicloud-resource-ga-basic-accelerate-ip"
description: |-
  Provides a Alicloud Global Accelerator (GA) Basic Accelerate IP resource.
---

# alicloud_ga_basic_accelerate_ip

Provides a Global Accelerator (GA) Basic Accelerate IP resource.

For information about Global Accelerator (GA) Basic Accelerate IP and how to use it, see [What is Basic Accelerate IP](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-doc-ga-2019-11-20-api-doc-createbasicaccelerateip).

-> **NOTE:** Available since v1.194.0.

## Example Usage

Basic Usage

```terraform
variable "region" {
  default = "cn-hangzhou"
}

provider "alicloud" {
  region  = var.region
  profile = "default"
}

resource "alicloud_ga_basic_accelerator" "default" {
  duration               = 1
  basic_accelerator_name = "terraform-example"
  description            = "terraform-example"
  bandwidth_billing_type = "CDT"
  auto_use_coupon        = "true"
  auto_pay               = true
}

resource "alicloud_ga_basic_ip_set" "default" {
  accelerator_id       = alicloud_ga_basic_accelerator.default.id
  accelerate_region_id = var.region
  isp_type             = "BGP"
  bandwidth            = "5"
}

resource "alicloud_ga_basic_accelerate_ip" "default" {
  accelerator_id = alicloud_ga_basic_accelerator.default.id
  ip_set_id      = alicloud_ga_basic_ip_set.default.id
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the Basic GA instance.
* `ip_set_id` - (Required, ForceNew) The ID of the Basic Ip Set.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Basic Accelerate IP.
* `accelerate_ip_address` - The address of the Basic Accelerate IP.
* `status` - The status of the Basic Accelerate IP instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Basic Accelerate IP.
* `delete` - (Defaults to 3 mins) Used when delete the Basic Accelerate IP.

## Import

Global Accelerator (GA) Basic Accelerate IP can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_basic_accelerate_ip.example <id>
```
