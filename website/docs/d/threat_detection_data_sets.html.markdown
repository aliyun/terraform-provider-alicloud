---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_data_sets"
sidebar_current: "docs-alicloud-datasource-threat-detection-data-sets"
description: |-
  Provides a list of Threat Detection Data Set items to the user.
---

# alicloud_threat_detection_data_sets

This data source provides the Threat Detection Data Set items of a current Alibaba Cloud account.

-> **NOTE:** Available since v1.245.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_threat_detection_data_sets" "default" {
  name_regex = "tf-example"
}

output "first_data_set_id" {
  value = "${data.alicloud_threat_detection_data_sets.default.data_sets.0.data_set_id}"
}
```

## Argument Reference

The following arguments are supported:
* `data_set_name` - (Optional, ForceNew) The name of the data set used to filter results.
* `data_set_status` - (Optional, ForceNew) The status of the data set used to filter results.
* `data_set_type` - (Optional, ForceNew) The type of the data set used to filter results.
* `ids` - (Optional) A list of Data Set IDs.
* `name_regex` - (Optional) A regex string to filter results by Data Set name.
* `output_file` - (Optional) File path where data results should be saved after running `terraform plan`.

## Attributes Reference

The following attributes are exported:
* `ids` - A list of Data Set IDs.
* `names` - A list of Data Set names.
* `data_sets` - A list of Threat Detection Data Set items. Each element contains the following attributes:
  * `id` - The ID of the data set.
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
