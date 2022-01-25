---
subcategory: "Data Security Center"
layout: "alicloud"
page_title: "Alicloud: alicloud_sddp_data_limits"
sidebar_current: "docs-alicloud-datasource-sddp-data-limits"
description: |-
  Provides a list of Sddp Data Limits to the user.
---

# alicloud\_sddp\_data\_limits

This data source provides the Sddp Data Limits of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.159.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_sddp_data_limits" "ids" {}
output "sddp_data_limit_id_1" {
  value = data.alicloud_sddp_data_limits.ids.limits.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Data Limit IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `parent_id` - (Optional, ForceNew) The parent asset ID of the data asset.
* `resource_type` - (Optional, ForceNew) The type of the service to which the data asset belongs. Valid values: `MaxCompute`, `OSS`, `RDS`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `limits` - A list of Sddp Data Limits. Each element contains the following attributes:
	* `audit_status` - Whether to enable the log auditing feature.
	* `check_status` - The status of the connectivity test between the data asset and SDDP.
	* `data_limit_id` - The first ID of the resource.
	* `engine_type` -The type of the database.
	* `id` - The ID of the Data Limit.
	* `local_name` - The name of the service to which the data asset belongs.
	* `log_store_day` - The retention period of raw logs after you enable the log auditing feature.
	* `parent_id` - The ID of the data asset.
	* `port` - The port that is used to connect to the database.
	* `resource_type` - The type of the service to which the data asset belongs.
	* `user_name` - The name of the user who owns the data asset.