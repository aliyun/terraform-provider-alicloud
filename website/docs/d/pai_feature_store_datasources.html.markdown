---
subcategory: "PAI Feature Store"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_feature_store_datasources"
sidebar_current: "docs-alicloud-datasource-pai-feature-store-datasources"
description: |-
  Provides a list of PAI Feature Store Datasources available to the user.
---

# alicloud_pai_feature_store_datasources

This data source provides a list of PAI Feature Store Datasources owned by an Alibaba Cloud account.[What is Datasource](https://www.alibabacloud.com/help/en/pai/developer-reference/api-paifeaturestore-2023-06-21-listdatasources)

-> **NOTE:** Available since v1.278.0.

## Example Usage

```terraform
data "alicloud_pai_feature_store_datasources" "default" {
  instance_id = "fs-xxxx"
  name        = "terraform-example"
}

output "first_datasource_id" {
  value = data.alicloud_pai_feature_store_datasources.default.datasources.0.datasource_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The ID of the FeatureStore instance.
* `workspace_id` - (Optional) Filter the results by the PAI workspace ID.
* `type` - (Optional) Filter the results by the data source type. Valid values: `Hologres`, `GraphCompute`, `FeatureDB`, `MaxCompute`.
* `name` - (Optional) Filter the results by the datasource name.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the data source, same as `instance_id`.
* `ids` - A list of datasource IDs.
* `total_count` - The total number of matching datasources.
* `datasources` - A list of FeatureStore Datasources. Each element contains the following attributes:
  * `datasource_id` - The ID of the datasource.
  * `instance_id` - The ID of the FeatureStore instance.
  * `type` - The data source type.
  * `datasource_name` - The datasource name.
  * `uri` - The resource address of the datasource.
  * `config` - The resource configuration information as a JSON string.
  * `create_time` - The creation time of the resource.
