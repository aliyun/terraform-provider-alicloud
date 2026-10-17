---
subcategory: "PAI Feature Store"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_feature_store_datasource"
sidebar_current: "docs-alicloud-resource-pai-feature-store-datasource"
description: |-
  Provides a PAI Feature Store Datasource resource.
---

# alicloud_pai_feature_store_datasource

Provides a PAI Feature Store Datasource resource.

For information about PAI Feature Store Datasource and how to use it, see [What is Datasource](https://www.alibabacloud.com/help/en/pai/developer-reference/api-paifeaturestore-2023-06-21-createdatasource).

-> **NOTE:** Available since v1.278.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_pai_feature_store_datasource" "default" {
  instance_id     = "fs-xxxx"
  workspace_id    = "ws-xxxx"
  type            = "FeatureDB"
  datasource_name = "terraform-example"
  uri             = "featuredb://example"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the FeatureStore instance, which can be obtained through the ListInstances interface.
* `type` - (Required, ForceNew) The data source type. Valid values: `Hologres`, `GraphCompute`, `FeatureDB`, `MaxCompute`.
* `datasource_name` - (Required) The datasource name. Start with a letter or number and contain only letters, numbers, and underscores. The maximum length cannot exceed 63.
* `workspace_id` - (Required, ForceNew) The PAI workspace ID, which can be obtained through the ListWorkspaces interface.
* `uri` - (Required) The resource address of the datasource. The maximum length cannot exceed 64.
* `config` - (Optional) The resource configuration information. The value is a JSON string, and the expected contents depend on `type`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, in the format `<InstanceId>:<DatasourceId>`.
* `datasource_id` - The ID of the datasource, allocated by the server.
* `create_time` - The creation time of the resource.

## Import

PAI Feature Store Datasource can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_feature_store_datasource.example <instance_id>:<datasource_id>
```
