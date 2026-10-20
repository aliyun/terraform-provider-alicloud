---
subcategory: "Pai"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_feature_store_label_table"
sidebar_current: "docs-alicloud-resource-pai-feature-store-label-table"
description: |-
  Provides a Pai Feature Store Label Table resource.
---

# alicloud_pai_feature_store_label_table

Provides a Pai Feature Store Label Table resource.

A label table describes the label fields registered against a FeatureStore project and datasource.

For information about Pai Feature Store Label Table and how to use it, see
[CreateLabelTable](https://www.alibabacloud.com/help/en/pai/developer-reference/api-paifeaturestore-2023-06-21-createlabeltable).

-> **NOTE:** Available since v2.1.0.

## Example Usage

```terraform
resource "alicloud_pai_feature_store_label_table" "example" {
  instance_id      = "fs-cn-xxxxx"
  project_id       = "xxxxx"
  label_table_name = "terraform-example"
  datasource_id    = "xxxxx"
  fields = [
    {
      name                = "user_id"
      type                = "STRING"
      attributes          = ["Partition"]
      aligned_entity_name = "user_id"
    },
    {
      name                = "label_value"
      type                = "DOUBLE"
      attributes          = ["LabelField"]
      aligned_entity_name = "user_id"
    },
  ]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew, Available since v2.1.0) The ID of the FeatureStore instance, which can be obtained through the ListInstances interface.
* `project_id` - (Required, ForceNew, Available since v2.1.0) The ID of the FeatureStore project, which can be obtained through the ListProjects interface.
* `label_table_name` - (Required, Available since v2.1.0) The label table name. A name can contain only letters, digits and underscores and can be up to 64 characters long.
* `datasource_id` - (Required, Available since v2.1.0) The ID of the datasource the label table belongs to.
* `fields` - (Required, List, Available since v2.1.0) A list of fields of the label table. See [`fields`](#fields) below. The API can only change the `attributes` and `aligned_entity_name` of existing fields, so changing the length or the order of the list, or the `name` or `type` of any field, forces a new resource.

  ~> **NOTE:** Within each `fields` object all four attributes must be set because the object is modeled as a single typed list element.

### `fields`

The `fields` object supports the following:

* `name` - (Required, Available since v2.1.0) The field name.
* `type` - (Required, Available since v2.1.0) The field data type. Valid values: `DOUBLE`, `INT32`, `INT64`, `FLOAT`, `STRING`, `BOOLEAN`, `TIMESTAMP`.
* `attributes` - (Required, Set, Available since v2.1.0) A set of attributes of the field. Valid values: `Partition`, `FeatureField`, `FeatureGenerationReserveField`, `EventTime`, `LabelField`, `EventTimeIsUnixTime`. Use an empty list when the field has no attributes.
* `aligned_entity_name` - (Required, Available since v2.1.0) The name of the entity the field is aligned to. Use an empty string when the field is not aligned to an entity.

## Attributes Reference

The following attributes are exported:

* `id` - (Available since v2.1.0) The ID of the resource. Format: `<instance_id>:<label_table_id>`.
* `label_table_id` - (Available since v2.1.0) The first ID of the resource.
* `create_time` - (Available since v2.1.0) The creation time of the resource.
* `region_id` - (Available since v2.1.0) The region ID of the resource.

## Import

Pai Feature Store Label Table can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_feature_store_label_table.example <instance_id>:<label_table_id>
```
