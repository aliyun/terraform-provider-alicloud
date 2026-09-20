---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_data_set"
sidebar_current: "docs-alicloud-datasource-threat-detection-data-set"
description: |-
  Provides a Threat Detection Data Set data source.
---

# alicloud_threat_detection_data_set

This data source provides the Threat Detection Data Set resource of a current Alibaba Cloud account.

-> **NOTE:** Available since v1.245.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_threat_detection_data_set" "default" {
  id = "example-data-set-id"
}

output "data_set_name" {
  value = "${data.alicloud_threat_detection_data_set.default.data_set_name}"
}
```

## Argument Reference

The following arguments are supported:
* `id` - (Required) The ID of the Data Set.
* `output_file` - (Optional) File path where data results should be saved after running `terraform plan`.

## Attributes Reference

The following attributes are exported:
* `data_set_id` - The ID of the data set.
* `data_set_name` - The name of the data set.
* `data_set_field_key_name` - The primary key name of the data set.
* `data_set_file_name` - The uploaded file name of the data set.
* `data_set_description` - The description of the data set.
* `data_set_type` - The type of the data set.
* `data_set_status` - The status of the data set.
* `role_for` - The user ID of another role.
* `lang` - The language type.
* `region_id` - The region of the user.
* `create_time` - The creation time of the resource.
* `ip_whitelist_recognizers` - IP whitelist recognizer configuration.
  * `auto_recognize_status` - Automatic identification status.
  * `recognize_scope` - Identification range.
  * `ip_whitelist_recognizer_type` - The IP type recognized by the identifier.
