---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_monitor_service_enterprise_public"
description: |-
  Provides a Alicloud Cloud Monitor Service Enterprise Public resource.
---

# alicloud_cloud_monitor_service_enterprise_public

Provides a Cloud Monitor Service Enterprise Public resource. Hybrid Cloud Monitoring.

For information about Cloud Monitor Service Enterprise Public and how to use it, see [What is Enterprise Public](https://www.alibabacloud.com/help/en/cms/user-guide/overview-3).

-> **NOTE:** Available since v1.215.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_cloud_monitor_service_enterprise_public" "default" {
}
```

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Enterprise Public.
* `delete` - (Defaults to 5 mins) Used when delete the Enterprise Public.

## Import

Cloud Monitor Service Enterprise Public can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_monitor_service_enterprise_public.example <id>
```