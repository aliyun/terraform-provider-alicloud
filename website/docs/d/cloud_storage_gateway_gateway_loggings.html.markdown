---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_gateway_loggings"
sidebar_current: "docs-alicloud-datasource-cloud-storage-gateway-gateway-loggings"
description: |-
  Provides a list of Cloud Storage Gateway Gateway Loggings to the user.
---

# alicloud\_cloud\_storage\_gateway\_gateway\_loggings

This data source provides the Cloud Storage Gateway Gateway Loggings of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.144.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cloud_storage_gateway_gateway_loggings" "ids" {
  gateway_id = "example_value"
  status     = "Enabled"
}

output "cloud_storage_gateway_gateway_logging_id_1" {
  value = data.alicloud_cloud_storage_gateway_gateway_loggings.ids.loggings.0.id
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, ForceNew)  A list of Gateway Logging IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `loggings` - A list of Cloud Storage Gateway Gateway Loggings. Each element contains the following attributes:
    * `gateway_id` - The first ID of the resource.
    * `id` - The ID of the Gateway Logging.
    * `sls_logstore` - The ID of the Log Store.
    * `sls_project` - The ID of the Project.
    * `status` - The status of the resource. Valid values: `Enabled`, `Disable` and `None`.
      * `None` - No gateway log monitoring.
      * `Enabled` - Gateway log monitoring is enabled.
      * `Disable` - Gateway log monitoring is disabled.