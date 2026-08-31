---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_log_shipper"
sidebar_current: "docs-alicloud-datasource-threat-detection-log-shipper"
description: |-
  Provides a datasource to open the Threat Detection Log Shipper automatically.
---

# alicloud\_threat\_detection\_log\_shipper

Using this data source can open Threat Detection Log Shipper automatically. If the service has been enabled, it will return `Opened`.

For information about Threat Detection Log Shipper and how to use it, see [What is Log Shipper](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-modifyopenlogshipper).

-> **NOTE:** Available in v1.195.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_threat_detection_log_shipper" "open" {
  enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to enable the service. Valid values: `On` or `Off`. Default to `Off`.

-> **NOTE:** Setting `enable = "On"` to open the Threat Detection Log Shipper that means you have read and agreed the [Threat Detection Log Shipper Terms of Service](https://help.aliyun.com/document_detail/170157.html). The service can not closed once it is opened.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status.
* `open_status` - Log analysis shipping activation status.
* `auth_status` - Log Analysis Service authorization status.
* `buy_status` - Cloud Security Center purchase status.
* `sls_project_status` - Log analysis project status.
* `sls_service_status` - Log Analysis Service is activated.
