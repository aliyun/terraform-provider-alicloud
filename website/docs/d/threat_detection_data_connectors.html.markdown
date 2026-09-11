---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_data_connectors"
sidebar_current: "docs-alicloud-datasource-threat-detection-data-connectors"
description: |-
  Provides a list of Threat Detection Data Connectors to the user.
---

# alicloud_threat_detection_data_connectors

This data source provides a list of Threat Detection Data Connectors of current Alibaba Cloud user.

-> **NOTE:** Available since v1.236.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_threat_detection_data_connectors" "default" {
  ids            = ["<data_connector_id>"]
  enable_details = true
}

output "first_connector_id" {
  value = data.alicloud_threat_detection_data_connectors.default.data_connectors.0.data_connector_id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of data connector ids used to filter the results.
* `name_regex` - (Optional) A regex string to filter the results by the data connector name.
* `lang` - (Optional) The language type of the message. Valid values: `zh` (Chinese), `en` (English).
* `region_id` - (Optional) The region id of the data connectors.
* `enable_details` - (Optional) Whether to get the details of each data connector. Default value: `false`.
* `output_file` - (Optional) The name of output file that saves the filter results.

## Attributes Reference

The following attributes are exported in addition to the `ids` and `name_regex`:

* `ids` - A list of data connector ids.
* `data_connectors` - A list of data connectors. Each element contains the following attributes:
  * `id` - The resource id of the data connector, equals to `data_connector_id`.
  * `data_connector_id` - The id of the data connector.
  * `data_connector_name` - The name of the data connector.
  * `data_connector_type` - The type of the data connector.
  * `data_connector_status` - The status of the data connector.
  * `data_connector_config` - The configuration of the data connector.
  * `src_data_type` - The source data source type.
  * `dest_data_source_id` - The destination data source id.
  * `log_project_name` - The destination log project name.
  * `log_store_name` - The destination log store name.
  * `log_region_id` - The destination log store region.
  * `auth_config_id` - The configuration item id of the kafka or s3 instance.
  * `auth_config_vendor` - The authentication vendor name.
  * `auth_config_product` - The authentication product name.
  * `sls_ingestion_job_name` - The name of the SLS data import task.
  * `sls_ingestion_job_state` - The status of the SLS data import task.
  * `creation_time` - The creation time of the data connector.
  * `update_time` - The update time of the data connector.
