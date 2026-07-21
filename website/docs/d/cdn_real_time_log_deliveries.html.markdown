---
subcategory: "CDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_cdn_real_time_log_deliveries"
sidebar_current: "docs-alicloud-datasource-cdn-real-time-log-deliveries"
description: |-
  Provides a list of Cdn Real Time Log Deliveries to the user.
---

# alicloud\_cdn\_real\_time\_log\_deliveries

This data source provides the Cdn Real Time Log Deliveries of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.134.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cdn_real_time_log_deliveries" "example" {
  domain = "example_value"
}
output "cdn_real_time_log_delivery_1" {
  value = data.alicloud_cdn_real_time_log_deliveries.example.deliveries.0.id
}

```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the real-time log delivery feature. Valid Values: `online` and `offline`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `deliveries` - A list of Cdn Real Time Log Deliveries. Each element contains the following attributes:
	* `domain` - Real-Time Log Service Domain.
	* `id` - The ID of the Real Time Log Delivery.
	* `logstore` - The name of the Logstore that collects log data from Alibaba Cloud Content Delivery Network (CDN) in real time.
	* `project` - The name of the Log Service project that is used for real-time log delivery.
	* `sls_region` - The region where the Log Service project is deployed.
	* `status` -The status of the real-time log delivery feature. Valid Values: `online` and `offline`.
