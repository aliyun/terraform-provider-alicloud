---
subcategory: "Simple Application Server"
layout: "alicloud"
page_title: "Alicloud: alicloud_simple_application_server_instance"
sidebar_current: "docs-alicloud-resource-simple-application-server-instance"
description: |-
  Provides a Alicloud Simple Application Server Instance resource.
---

# alicloud_simple_application_server_instance

Provides a Simple Application Server Instance resource.

For information about Simple Application Server Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/doc-detail/190440.htm).

-> **NOTE:** Available since v1.135.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_simple_application_server_instance&exampleId=e8653110-c705-9b88-ae31-9f257ba9092308ccdc5a&activeTab=example&spm=docs.r.simple_application_server_instance.0.e8653110c7&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf_example"
}

data "alicloud_simple_application_server_images" "default" {
  platform = "Linux"
}
data "alicloud_simple_application_server_plans" "default" {
  platform = "Linux"
}

resource "alicloud_simple_application_server_instance" "default" {
  payment_type   = "Subscription"
  plan_id        = data.alicloud_simple_application_server_plans.default.plans.0.id
  instance_name  = var.name
  image_id       = data.alicloud_simple_application_server_images.default.images.0.id
  period         = 1
  data_disk_size = 100
}
```

### Deleting `alicloud_simple_application_server_instance` or removing it from your configuration

The `alicloud_simple_application_server_instance` resource allows you to manage `payment_type = "Subscription"` instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the resource Instance.
You can resume managing the subscription instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:

* `auto_renew` - (Optional) Specifies whether to enable auto-renewal. Unit: months. Valid values: `true` and `false`.
* `auto_renew_period` - (Optional) The auto renew period. Valid values: `1`,`3`, `6`, `12`, `24`, `36`. **NOTE:** The attribute `auto_renew` is valid when the attribute is `true`.
* `data_disk_size` - (Optional) The size of the data disk. Unit: GB. Valid values: `0` to `16380`.
* `image_id` - (Required) The ID of the image.  You can use the `alicloud_simple_application_server_images` to query the available images in the specified region. The value must be an integral multiple of 20.
* `instance_name` - (Optional) The name of the simple application server.
* `password` - (Optional) The password of the simple application server. The password must be 8 to 30 characters in length. It must contain at least three of the following character types: uppercase letters, lowercase letters, digits, and special characters. Special characters include: `( ) ~ ! @ # $ % ^ & * - + = | { } [ ] : ; < > , . ? /`.
* `payment_type` - (Optional, ForceNew) The paymen type of the resource. Valid values: `Subscription`.
* `period` - (Required) The period. Unit: months. Valid values: `1`,`3`, `6`, `12`, `24`, `36`.
* `plan_id` - (Required) The ID of the plan. You can use the `alicloud_simple_application_server_plans`  to query all the plans provided by Simple Application Server in the specified region.
* `status` - (Optional) The status of the simple application server. Valid values: `Resetting`, `Running`, `Stopped`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Instance.
* `update` - (Defaults to 5 mins) Used when update the Instance.

## Import

Simple Application Server Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_simple_application_server_instance.example <id>
```
