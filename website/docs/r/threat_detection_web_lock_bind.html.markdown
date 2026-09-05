---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_web_lock_bind"
sidebar_current: "docs-alicloud-resource-threat-detection-web-lock-bind"
description: |-
  Provides a Alicloud Threat Detection Web Lock Bind resource.
---

# alicloud_threat_detection_web_lock_bind

Provides a Threat Detection Web Lock Bind resource.

Bind a server to the web tamper proofing feature, configure the protected directories and file types, and manage the protection status.

For information about Threat Detection Web Lock Bind and how to use it, see [ModifyWebLockStart](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-modifyweblockstart).

-> **NOTE:** Available since v1.248.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_threat_detection_assets" "default" {
  machine_types = "ecs"
}

resource "alicloud_threat_detection_web_lock_bind" "default" {
  uuid                = data.alicloud_threat_detection_assets.default.ids.0
  dir                 = "/tmp/"
  mode                = "whitelist"
  local_backup_dir    = "/usr/local/aegis/bak"
  defence_mode        = "audit"
  inclusive_file_type = "php;jsp;asp;aspx;js;cgi;html;htm;xml;shtml;shtm;jpg"
  status              = "on"
  lang                = "zh"
}
```

## Argument Reference

The following arguments are supported:

* `uuid` - (Required, ForceNew) The UUID of the server to which the web tamper proofing feature is bound. You can call the [DescribeWebLockBindList](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-describeweblockbindlist) operation to query the bound servers.
* `dir` - (Required, ForceNew) The protected directory. Separate multiple directories with commas (,).
* `local_backup_dir` - (Required, ForceNew) The local backup path used to safely back up the protected directory.
* `mode` - (Required, ForceNew) The protected directory mode. Valid values: `whitelist`, `blacklist`.
* `defence_mode` - (Required, ForceNew) The protection mode. Valid values: `block`, `audit`.
* `exclusive_dir` - (Optional, ForceNew) The excluded directory address that does not require web tamper protection. Required when `mode` is `blacklist`.
* `exclusive_file_type` - (Optional, ForceNew) The excluded file types that do not require web tamper protection. Separate multiple file types with semicolons (;). Required when `mode` is `blacklist`.
* `inclusive_file_type` - (Optional, ForceNew) The file types that require tamper protection. Separate multiple file types with semicolons (;). Required when `mode` is `whitelist`.
* `exclusive_file` - (Optional, ForceNew) The excluded files that do not require web tamper protection. Required when `mode` is `blacklist`.
* `status` - (Optional) The protection status of the server. Valid values: `on`, `off`.
* `lang` - (Optional) The language type of the request and response. Valid values: `zh`, `en`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, which is the `uuid` of the bound server.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the Web Lock Bind.
* `delete` - (Defaults to 5 mins) Used when deleting the Web Lock Bind.

## Import

Threat Detection Web Lock Bind can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_web_lock_bind.example <uuid>
```
