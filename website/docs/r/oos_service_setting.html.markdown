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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oos_service_setting&exampleId=38b88eb0-0982-c5c7-9f9c-2756d361bf6e390ea938&activeTab=example&spm=docs.r.oos_service_setting.0.38b88eb009&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-testaccoossetting"
}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
}

resource "alicloud_oss_bucket_acl" "default" {
  bucket = alicloud_oss_bucket.default.bucket
  acl    = "public-read-write"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_oos_service_setting" "default" {
  delivery_oss_enabled      = true
  delivery_oss_key_prefix   = "path1/"
  delivery_oss_bucket_name  = alicloud_oss_bucket.default.bucket
  delivery_sls_enabled      = true
  delivery_sls_project_name = alicloud_log_project.default.project_name
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

```shell
$ terraform import alicloud_oos_service_setting.example <id>
```