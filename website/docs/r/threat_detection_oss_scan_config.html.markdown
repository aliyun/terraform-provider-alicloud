---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_oss_scan_config"
description: |-
  Provides a Alicloud Threat Detection Oss Scan Config resource.
---

# alicloud_threat_detection_oss_scan_config

Provides a Threat Detection Oss Scan Config resource.

Oss detection configuration.

For information about Threat Detection Oss Scan Config and how to use it, see [What is Oss Scan Config](https://www.alibabacloud.com/help/zh/security-center/developer-reference/api-sas-2018-12-03-createossscanconfig/).

-> **NOTE:** Available since v1.214.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

locals {
  bucket_random = random_integer.default.result
}

resource "alicloud_oss_bucket" "default8j4t1R" {
  bucket        = "${var.name}-1-${local.bucket_random}"
  storage_class = "Standard"
}

resource "alicloud_oss_bucket" "default9HMqfT" {
  bucket        = "${var.name}-2-${local.bucket_random}"
  storage_class = "Standard"
}

resource "alicloud_oss_bucket" "defaultxBXqFQ" {
  bucket        = "${var.name}-3-${local.bucket_random}"
  storage_class = "Standard"
}

resource "alicloud_oss_bucket" "defaulthZvCmR" {
  bucket        = "${var.name}-4-${local.bucket_random}"
  storage_class = "Standard"
}

resource "alicloud_threat_detection_oss_scan_config" "default" {
  key_suffix_list = [
    ".jsp",
    ".php",
    ".k"
  ]
  scan_day_list = [
    "2",
    "5",
    "4",
    "3"
  ]
  oss_scan_config_name = var.name
  end_time             = "00:00:02"
  start_time           = "00:00:01"
  enable               = "1"
  all_key_prefix       = "false"
  bucket_name_list = [
    alicloud_oss_bucket.default8j4t1R.bucket,
    alicloud_oss_bucket.default9HMqfT.bucket,
    alicloud_oss_bucket.defaultxBXqFQ.bucket
  ]
  key_prefix_list = [
    "/root",
    "/usr",
    "/123"
  ]
}
```

## Argument Reference

The following arguments are supported:
* `all_key_prefix` - (Optional, Computed) Specifies whether to match the prefixes of all objects.
* `bucket_name_list` - (Required, Set) The names of the buckets.
* `decompress_max_file_count` - (Optional, Int, Available since v1.255.0) The maximum number of objects that can be extracted during decompression. Valid values: 1 to 1000. If the maximum number of objects that can be extracted is reached, the decompression operation immediately ends and the detection of extracted objects is not affected.
* `decompress_max_layer` - (Optional, Int, Available since v1.255.0) The maximum number of decompression levels when multi-level packages are decompressed. Valid values: 1 to 5. If the maximum number of decompression levels is reached, the decompression operation immediately ends and the detection of extracted objects is not affected.
* `decryption_list` - (Optional, List, Available since v1.255.0) The decryption methods.
* `enable` - (Required, Int) Indicates whether the check policy is enabled. Valid values:

  - `1`: enabled.
  - `0`: disabled.
* `end_time` - (Required) The end time of the check. The time is in the HH:mm:ss format.
* `key_prefix_list` - (Optional, Set) The prefixes of the objects.
* `key_suffix_list` - (Required, Set) The suffixes of the objects that are checked.
* `last_modified_start_time` - (Optional, Int, Available since v1.255.0) The timestamp when the object was last modified. The time must be later than the timestamp that you specify. Unit: milliseconds.
* `oss_scan_config_name` - (Optional) The policy name.
* `scan_day_list` - (Required, Set) The days when the check is performed. The value indicates the days of the week.
* `start_time` - (Required) The start time of the check. The time is in the HH:mm:ss format.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Oss Scan Config.
* `delete` - (Defaults to 5 mins) Used when delete the Oss Scan Config.
* `update` - (Defaults to 5 mins) Used when update the Oss Scan Config.

## Import

Threat Detection Oss Scan Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_oss_scan_config.example <id>
```