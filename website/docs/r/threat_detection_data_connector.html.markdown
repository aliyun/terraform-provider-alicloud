---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_data_connector"
sidebar_current: "docs-alicloud-resource-threat-detection-data-connector"
description: |-
  Provides a Threat Detection Data Connector resource.
---

# alicloud_threat_detection_data_connector

Provides a Threat Detection Data Connector resource. Data Connector is a collector that ingests data from a source data source (such as OSS, S3, or Kafka) into a Simple Log Service (SLS) log store for threat detection analysis.

-> **NOTE:** Available since v1.236.0.

## Example Usage

Basic Usage

```terraform
variable "region" {
  default = "cn-hangzhou"
}
variable "dest_data_source_id" {}

resource "alicloud_log_project" "default" {
  name = "tf-data-connector-project"
}

resource "alicloud_log_store" "default" {
  project = alicloud_log_project.default.name
  name    = "tf-data-connector-store"
}

resource "alicloud_threat_detection_data_connector" "default" {
  data_connector_type   = "oss"
  data_connector_config = "{\"bucket\":\"my-bucket\",\"prefix\":\"logs\"}"
  src_data_type         = "OSS"
  dest_data_source_id   = var.dest_data_source_id
  log_project_name      = alicloud_log_project.default.name
  log_store_name        = alicloud_log_store.default.name
  log_region_id         = var.region
  data_connector_status = "enabled"
  auth_config_vendor    = "APACHE"
}
```

## Argument Reference

The following arguments are supported:

* `data_connector_type` - (Required, ForceNew) The type of the data connector. Valid values: `oss`, `s3`, `kafka`.
* `data_connector_config` - (Required) The configuration of the data connector, in JSON format. The content varies according to the data connector type.
* `src_data_type` - (Required, ForceNew) The source data source type, such as `OSS`, `OBS`, `COS`, `DMS(Kafka)`, `Ckafka`.
* `dest_data_source_id` - (Required, ForceNew) The destination data source id.
* `log_project_name` - (Required, ForceNew) The destination log project name.
* `log_store_name` - (Required, ForceNew) The destination log store name.
* `log_region_id` - (Required, ForceNew) The destination log store region.
* `data_connector_status` - (Optional) The status of the data connector. Valid values: `enabled`, `disabled`.
* `auth_config_id` - (Optional, ForceNew) The configuration item id of the kafka or s3 instance accessed by the collector in the multi-cloud configuration.
* `auth_config_vendor` - (Optional) The authentication vendor name. Used together with `auth_config_id` for authentication when the collector accesses the source data service. Valid values: `APACHE`, `AWS_S3`, `SALESFORCE`. Leave it empty when collecting OSS data.
* `auth_config_product` - (Optional, ForceNew) The authentication product name.
* `lang` - (Optional) The language type of the message. Valid values: `zh` (Chinese), `en` (English). Default value: `en`.
* `role_for` - (Optional, ForceNew) The user id that the administrator switches to another member's perspective.
* `region_id` - (Optional, ForceNew) The region id of the data connector.

## Attributes Reference

The following attributes are exported:

* `id` - The resource id of the data connector, equals to `data_connector_id`.
* `data_connector_id` - The id of the data connector.
* `data_connector_name` - The name of the data connector, in the format `SAS-CTDR-YYYYMMDDhhmm`.
* `sls_ingestion_job_name` - The name of the SLS data import task corresponding to the connector.
* `sls_ingestion_job_state` - The status of the SLS data import task corresponding to the connector. Valid values: `restarting`, `starting`, `running`, `stopping`, `stopped`.
* `creation_time` - The creation time of the data connector.
* `update_time` - The update time of the data connector.

## Import

Threat Detection Data Connector can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_data_connector.example <data_connector_id>
```
