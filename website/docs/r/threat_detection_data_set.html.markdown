---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_data_set"
sidebar_current: "docs-alicloud-resource-threat-detection-data-set"
description: |-
  Provides a Alicloud Threat Detection Data Set resource.
---

# alicloud_threat_detection_data_set

Provides a Threat Detection Data Set resource.

For information about Threat Detection Data Set and how to use it, see [What is Data Set](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-cloud-siem-2024-12-12-createdataset).

-> **NOTE:** Available since v1.245.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

resource "alicloud_threat_detection_data_set" "default" {
  data_set_name           = var.name
  data_set_field_key_name = "ip"
  data_set_file_name      = "example.csv"
  data_set_description    = "description for the data set"
  data_set_type           = "userDefined"
  data_set_status         = 1
  ip_whitelist_recognizers {
    auto_recognize_status        = "on"
    recognize_scope              = "all"
    ip_whitelist_recognizer_type = "ip"
  }
}
```

## Argument Reference

The following arguments are supported:
* `data_set_name` - (Required) The name of the data set.
* `data_set_field_key_name` - (Required, ForceNew) The primary key name of the data set.
* `data_set_file_name` - (Required) The uploaded file name of the data set.
* `data_set_description` - (Optional) The description of the data set.
* `data_set_type` - (Optional, ForceNew) The type of the data set.
* `data_set_status` - (Optional) The status of the data set.
* `role_for` - (Optional) The user ID of another role to switch to.
* `ip_whitelist_recognizers` - (Optional) IP whitelist recognizer configuration. See [`ip_whitelist_recognizers`](#ip_whitelist_recognizers) below for details.
* `lang` - (Optional) The language type of the request.
* `region_id` - (Optional) The region of the user.

### `ip_whitelist_recognizers`

The ip_whitelist_recognizers supports the following:

* `auto_recognize_status` - (Optional) Automatic identification status.
* `recognize_scope` - (Optional) Identification range.
* `ip_whitelist_recognizer_type` - (Optional) The IP type recognized by the identifier.

## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above. The value is the same as `data_set_id`.
* `data_set_id` - The ID of the data set.
* `create_time` - The creation time of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Data Set.
* `delete` - (Defaults to 5 mins) Used when delete the Data Set.
* `update` - (Defaults to 5 mins) Used when update the Data Set.

## Import

Threat Detection Data Set can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_data_set.example <data_set_id>
```
