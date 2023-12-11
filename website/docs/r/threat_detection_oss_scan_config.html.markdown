---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_oss_scan_config"
description: |-
  Provides a Alicloud Threat Detection Oss Scan Config resource.
---

# alicloud_threat_detection_oss_scan_config

Provides a Threat Detection Oss Scan Config resource. Oss detection configuration.

For information about Threat Detection Oss Scan Config and how to use it, see [What is Oss Scan Config](https://www.alibabacloud.com/help/zh/security-center/developer-reference/api-sas-2018-12-03-createossscanconfig/).

-> **NOTE:** Available since v1.214.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_oss_bucket" "default8j4t1R" {
  bucket_name = var.name

  storage_class = "Standard"
}

resource "alicloud_oss_bucket" "default9HMqfT" {
  bucket_name = var.name

  storage_class = "Standard"
}

resource "alicloud_oss_bucket" "defaultxBXqFQ" {
  bucket_name = var.name

  storage_class = "Standard"
}

resource "alicloud_oss_bucket" "defaulthZvCmR" {
  bucket_name = var.name

  storage_class = "Standard"
}


resource "alicloud_threat_detection_oss_scan_config" "default" {
  bucket_name_list     = ["gcx-test-oss-71", "gcx-test-oss-72", "gcx-test-oss-73"]
  key_prefix_list      = ["/root", "/usr", "/123"]
  oss_scan_config_name = var.name

  end_time        = "00:00:01"
  start_time      = "00:00:00"
  enable          = "0"
  key_suffix_list = [".html", ".php", ".k"]
  scan_day_list   = ["1", "2", "4", "3"]
}
```

## Argument Reference

The following arguments are supported:
* `all_key_prefix` - (Optional) Match all prefixes.
* `bucket_name_list` - (Required) Bucket List.
* `enable` - (Required) Enable configuration.
* `end_time` - (Required) End time, hours, minutes and seconds.
* `key_prefix_list` - (Optional) File prefix list.
* `key_suffix_list` - (Required) File Suffix List.
* `oss_scan_config_name` - (Optional) Configuration Name.
* `scan_day_list` - (Required) Scan cycle.
* `start_time` - (Required) Start time, hours, minutes and seconds.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Oss Scan Config.
* `delete` - (Defaults to 5 mins) Used when delete the Oss Scan Config.
* `update` - (Defaults to 5 mins) Used when update the Oss Scan Config.

## Import

Threat Detection Oss Scan Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_oss_scan_config.example <id>
```