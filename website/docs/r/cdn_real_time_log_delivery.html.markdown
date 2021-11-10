---
subcategory: "CDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_cdn_real_time_log_delivery"
sidebar_current: "docs-alicloud-resource-cdn-real-time-log-delivery"
description: |-
  Provides a Alicloud CDN Real Time Log Delivery resource.
---

# alicloud\_cdn\_real\_time\_log\_delivery

Provides a CDN Real Time Log Delivery resource.

For information about CDN Real Time Log Delivery and how to use it, see [What is Real Time Log Delivery](https://www.alibabacloud.com/help/doc-detail/100456.htm).

-> **NOTE:** Available in v1.134.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cdn_real_time_log_delivery" "example" {
  domain     = "example_value"
  logstore   = "example_value"
  project    = "example_value"
  sls_region = "cn-hanghzou"
}

```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, ForceNew) The accelerated domain name for which you want to configure real-time log delivery. You can specify multiple domain names and separate them with commas (,).
* `logstore` - (Required, ForceNew) The name of the Logstore that collects log data from Alibaba Cloud Content Delivery Network (CDN) in real time.
* `project` - (Required, ForceNew) The name of the Log Service project that is used for real-time log delivery.
* `sls_region` - (Required, ForceNew) The region where the Log Service project is deployed.

-> **NOTE:** If your Project and Logstore services already exist, if you continue to create existing content, the created content will overwrite your existing indexes and custom reports. Please be careful to create your existing services to avoid affecting your online services after coverage.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Real Time Log Delivery. Its value is same as `domain`.
* `status` - The status of the real-time log delivery feature. Valid Values: `online` and `offline`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Real Time Log Delivery.

## Import

CDN Real Time Log Delivery can be imported using the id, e.g.

```
$ terraform import alicloud_cdn_real_time_log_delivery.example <domain>
```
