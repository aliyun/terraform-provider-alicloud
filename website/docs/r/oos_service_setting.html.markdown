---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_service_setting"
sidebar_current: "docs-alicloud-resource-oos-service-setting"
description: |-
  Provides a Alicloud OOS Service Setting resource.
---

# alicloud\_oos\_service\_setting

Provides a OOS Service Setting resource.

For information about OOS Service Setting and how to use it, see [What is Service Setting](https://www.alibabacloud.com/help/en/doc-detail/268700.html).

-> **NOTE:** Available in v1.147.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testaccoossetting"
}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
  acl    = "public-read-write"
}

resource "alicloud_log_project" "default" {
  name = var.name
}

resource "alicloud_oos_service_setting" "default" {
  delivery_oss_enabled      = true
  delivery_oss_key_prefix   = "path1/"
  delivery_oss_bucket_name  = alicloud_oss_bucket.default.bucket
  delivery_sls_enabled      = true
  delivery_sls_project_name = alicloud_log_project.default.name
}
```

## Argument Reference

The following arguments are supported:

* `delivery_oss_bucket_name` - (Optional) The name of the OSS bucket. **NOTE:** When the `delivery_oss_enabled` is `true`, The `delivery_oss_bucket_name` is valid.
* `delivery_oss_enabled` - (Optional) Is the recording function for the OSS delivery template enabled.  
* `delivery_oss_key_prefix` - (Optional) The Directory of the OSS bucket. **NOTE:** When the `delivery_oss_enabled` is `true`, The `delivery_oss_bucket_name` is valid.
* `delivery_sls_enabled` - (Optional) Is the execution record function to SLS delivery Template turned on.
* `delivery_sls_project_name` - (Optional) The name of SLS  Project. **NOTE:** When the `delivery_sls_enabled` is `true`, The `delivery_sls_project_name` is valid.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Service Setting.

## Import

OOS Service Setting can be imported using the id, e.g.

```
$ terraform import alicloud_oos_service_setting.example <id>
```