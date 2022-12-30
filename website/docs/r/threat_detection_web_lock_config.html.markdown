---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_web_lock_config"
sidebar_current: "docs-alicloud-resource-threat_detection-web-lock-config"
description: |-
  Provides a Alicloud Threat Detection Web Lock Config resource.
---

# alicloud_threat_detection_web_lock_config

Provides a Threat Detection Web Lock Config resource.

For information about Threat Detection Web Lock Config and how to use it, see [What is Web Lock Config](https://www.alibabacloud.com/help/en/security-center/latest/api-doc-sas-2018-12-03-api-doc-modifyweblockstart).

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_threat_detection_assets" "default" {
  machine_types = "ecs"
}
resource "alicloud_threat_detection_web_lock_config" "default" {
  inclusive_file_type = "php;jsp;asp;aspx;js;cgi;html;htm;xml;shtml;shtm;jpg"
  uuid                = data.alicloud_threat_detection_assets.default.ids.0
  mode                = "whitelist"
  local_backup_dir    = "/usr/local/aegis/bak"
  dir                 = "/tmp/"
  defence_mode        = "audit"
}
```

## Argument Reference

The following arguments are supported:
* `defence_mode` - (Required,ForceNew) Protection mode. Value:-**block**: Intercept-**audit**: Alarm
* `dir` - (Required,ForceNew) Specify the protection directory.
* `exclusive_dir` - (ForceNew,Optional) Specify a directory address that does not require Web tamper protection (I. E. Excluded directories).> The protection Mode **Mode** is set to **blacklist**, you need to configure this parameter.
* `exclusive_file` - (ForceNew,Optional) Specify files that do not need to enable tamper protection for web pages (that is, exclude files).> The protection Mode **Mode** is set to **blacklist**, you need to configure this parameter.
* `exclusive_file_type` - (ForceNew,Optional) Specify the type of file that does not require Web tamper protection (that is, the type of excluded file). When there are multiple file types, use semicolons (;) separation. Value:-php-jsp-asp-aspx-js-cgi-html-htm-xml-shtml-shtm-jpg-gif-png > The protection Mode **Mode** is set to **blacklist**, you need to configure this parameter.
* `inclusive_file_type` - (ForceNew,Optional) Specify the type of file that requires tamper protection. When there are multiple file types, use semicolons (;) separation. Value:-php-jsp-asp-aspx-js-cgi-html-htm-xml-shtml-shtm-jpg-gif-png> The protection Mode **Mode** is set to **whitelist**, you need to configure this parameter.
* `local_backup_dir` - (Required,ForceNew) The local backup path is used to protect the safe backup of the Directory.
* `mode` - (Required,ForceNew) Specify the protected directory mode. Value:-**whitelist**: whitelist mode, which protects the added protected directories and file types.-**blacklist**: blacklist mode, which protects all unexcluded subdirectories, file types, and specified files under the added protection directory.
* `uuid` - (Required,ForceNew) Specify the UUID of the server to which you want to add a protection directory.> You can call the [DescribeCloudCenterInstances](~~ 141932 ~~) interface to obtain the UUID of the server.



## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Web Lock Config.
* `delete` - (Defaults to 5 mins) Used when delete the Web Lock Config.

## Import

Threat Detection Web Lock Config can be imported using the id, e.g.

```shell
$terraform import alicloud_threat_detection_web_lock_config.example <id>
```