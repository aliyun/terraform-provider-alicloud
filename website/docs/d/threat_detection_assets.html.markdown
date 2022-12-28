---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_assets"
sidebar_current: "docs-alicloud-datasource-threat_detection-assets"
description: |-
  Provides a list of Threat Detection Asset owned by an Alibaba Cloud account.
---

# alicloud_threat_detection_assets

This data source provides Threat Detection Asset available to the user.[What is Asset](https://www.alibabacloud.com/help/en/security-center/latest/api-doc-sas-2018-12-03-api-doc-describecloudcenterinstances)

-> **NOTE:** Available in 1.195.0+

## Example Usage

```
data "alicloud_threat_detection_assets" "default" {}

output "alicloud_threat_detection_asset_example_id" {
  value = data.alicloud_threat_detection_assets.default.assets.0.id
}
```

## Argument Reference

The following arguments are supported:
* `criteria` - (ForceNew,Optional) Set the conditions for searching assets. This parameter is in JSON format. Note the case when you enter the parameter. **NOTE:** You can search for assets by using conditions such as the instance ID, instance name, VPC ID, region, and public IP address of the asset.
* `importance` - (ForceNew,Optional) Set asset importance. Value:
  - **2**: Significant assets
  - **1**: General assets
  - **0**: Test asset
* `logical_exp` - (ForceNew,Optional) Set the logical relationship between multiple search conditions. The default value is **OR**. Valid values:
  - **OR**: indicates that the relationship between multiple search conditions is **OR**.
  - **AND**: indicates that the relationship between multiple search conditions is **AND**.
* `machine_types` - (ForceNew,Optional) The type of asset to query. Value:
  - **ecs**: server.
  - **cloud_product**: Cloud product.
* `no_group_trace` - (Optional, ForceNew) Specifies whether to internationalize the name of the default group. Default value: false
* `ids` - (Optional, ForceNew, Computed) A list of Asset IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Asset IDs.
* `assets` - A list of Asset Entries. Each element contains the following attributes:
  * `create_time` - The creation time of the resource
  * `uuid` - The UUID of the instance.
  * `id` - The ID of the instance.
