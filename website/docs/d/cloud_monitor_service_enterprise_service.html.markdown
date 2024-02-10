---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_monitor_service_enterprise_service"
description: |-
  Provides a Alicloud Cloud Monitor Service Enterprise Public resource.
---

# alicloud_cloud_monitor_service_enterprise_service

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

data "alicloud_cloud_monitor_service_enterprise_service" "current" {
  enable = "On"
}
```

## Argument Reference

The following arguments are supported:
* `enable` - (Required) Setting the value to `On` to enable the service. Valid values: `On` or `Off`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.