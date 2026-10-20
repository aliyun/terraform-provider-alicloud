---
subcategory: "PaiFeatureStore"
layout: "alicloud"
page_title: "Alicloud: alicloud_paifeaturestore_feature_view"
description: |-
  Provides an Alicloud PaiFeatureStore Feature View resource.
---

# alicloud_paifeaturestore_feature_view

Provides an Alicloud PaiFeatureStore Feature View resource. A feature view belongs to a FeatureStore instance and project, and maps an offline feature table to an online one.

For information about PaiFeatureStore Feature View and how to use it, see [PaiFeatureStore](https://www.alibabacloud.com/help/en/pai/feature-store/).

-> **NOTE:** The FeatureStore instance and project that a feature view belongs to are not yet onboarded as Terraform resources. Provide `instance_id` and `project_id` of an existing instance and project created in the PAI Feature Store console or via the OpenAPI.

-> **NOTE:** Every argument of this resource forces replacement when changed, including the nested `fields`, `transform` and `input` blocks. The `UpdateFeatureView` API accepts only a subset of the field sub-properties; rather than ship a silently lossy partial update, the resource recreates the feature view on any change until a later revision implements the mapped update subset precisely.

-> **NOTE:** Available since v2.1.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_paifeaturestore_feature_view" "default" {
  instance_id       = "fsi-xxxxxx" # an existing FeatureStore instance
  feature_view_name = var.name
  project_id        = "fp-xxxxxx" # an existing project under the instance
  type              = "Batch"
  write_method      = "ByReadyMadeTable"
  sync_online_table = true
  ttl               = 86400

  fields {
    name       = "user_id"
    type       = "STRING"
    attributes = ["PrimaryKey"]
  }

  fields {
    name = "age"
    type = "INT64"
  }
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the FeatureStore instance. It is the parent of the feature view. Changing this creates a new resource.
* `feature_view_name` - (Required, ForceNew) The name of the feature view. It must start with a letter or a number and can contain only letters, numbers, and underscores, up to 128 characters in length.
* `project_id` - (Required, ForceNew) The ID of the project that the feature view belongs to.
* `type` - (Required, ForceNew) The type of the feature view. Valid values: `Batch`, `Stream`, `Sequence`.
* `write_method` - (Required, ForceNew) The method that is used to write data to the online table. Valid values: `ByReadyMadeTable`, `Custom`.
* `sync_online_table` - (Required, ForceNew) Specifies whether to synchronize the online feature table.
* `ttl` - (Required, ForceNew) The data lifecycle. Unit: seconds.
* `fields` - (Optional) The list of feature view fields. The `fields` block is documented below.
* `feature_entity_id` - (Optional, ForceNew, Computed) The ID of the feature entity.
* `write_to_feature_db` - (Optional, ForceNew, Computed) Specifies whether data is written to FeatureDB. Defaults to `false`.
* `register_table` - (Optional, ForceNew, Computed) The name of the offline registry table.
* `register_datasource_id` - (Optional, ForceNew, Computed) The ID of the data source where the registry table resides.
* `config` - (Optional, ForceNew, Computed) The advanced configuration, in JSON format.
* `tags` - (Optional, ForceNew, Computed) The tag list.

The `fields` block supports:

* `name` - (Required, ForceNew) The name of the feature field.
* `type` - (Required, ForceNew) The type of the feature field. Valid values: `INT32`, `INT64`, `FLOAT`, `DOUBLE`, `STRING`, `BOOLEAN`, `TIMESTAMP`, and the complex types `ARRAY` and `MAP`.
* `attributes` - (Optional, ForceNew, Computed) The attributes of the field. Valid values: `Partition`, `PrimaryKey`, `EventTime`.
* `transform` - (Optional) The feature production configuration. The `transform` block is documented below.
* `dimension` - (Optional, ForceNew, Computed) The dimension.

The `transform` block supports:

* `type` - (Optional, ForceNew, Computed) The feature production type. Valid value: `LLMEmbedding`.
* `llm_config_id` - (Optional, ForceNew, Computed) The LLM configuration ID.
* `input` - (Optional) The feature production inputs. The `input` block is documented below.

The `input` block supports:

* `name` - (Optional, ForceNew, Computed) The input field name.
* `type` - (Optional, ForceNew, Computed) The input field type.
* `modality` - (Optional, ForceNew, Computed) The input field modality. Valid values: `NONE`, `TEXT`, `IMAGE`, `VIDEO`, `MULTI_IMAGES`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. It is the composite identity in the format `<InstanceId>:<FeatureViewId>`.
* `feature_view_id` - The ID of the feature view, assigned by the API at create time.
* `region_id` - The region ID of the resource.
* `create_time` - The creation time of the resource.

## Import

PaiFeatureStore Feature View can be imported using the composite identity `<InstanceId>:<FeatureViewId>`, e.g.

```shell
terraform import alicloud_paifeaturestore_feature_view.default fsi-xxxxxx:fv-yyyyyy
```
