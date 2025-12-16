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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_threat_detection_oss_scan_config&exampleId=fb929f2e-17ad-7de3-0955-641d3477a8f3eedb9658&activeTab=example&spm=docs.r.threat_detection_oss_scan_config.0.fb929f2e17&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_threat_detection_oss_scan_config&spm=docs.r.threat_detection_oss_scan_config.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `all_key_prefix` - (Optional, Computed) Indicates whether the prefixes of all objects are matched.
* `bucket_name_list` - (Required, List) The names of the buckets.
* `decompress_max_file_count` - (Optional, Int, Available since v1.255.0) The maximum number of objects that can be extracted during decompression. Valid values: 1 to 1000. If the maximum number of objects that can be extracted is reached, the decompression operation immediately ends and the detection of extracted objects is not affected.
* `decompress_max_layer` - (Optional, Int, Available since v1.255.0) The maximum number of decompression levels when multi-level packages are decompressed. Valid values: 1 to 5. If the maximum number of decompression levels is reached, the decompression operation immediately ends and the detection of extracted objects is not affected.
* `decryption_list` - (Optional, List, Available since v1.255.0) The decryption methods.
* `enable` - (Required, Int) Indicates whether the check policy is enabled. Valid values:
  - `1`: enabled.
  - `0`: disabled.
* `end_time` - (Required) The end time of the check. The time is in the HH:mm:ss format.
* `key_prefix_list` - (Optional, List) The prefixes of the objects.
* `key_suffix_list` - (Required, List) The suffixes of the objects that are checked.
* `last_modified_start_time` - (Optional, Int, Available since v1.255.0) The timestamp when the object was last modified. The time must be later than the timestamp that you specify. Unit: milliseconds.
* `oss_scan_config_name` - (Optional) The policy name.
* `scan_day_list` - (Required, List) The days when the check is performed. The value indicates the days of the week.
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
