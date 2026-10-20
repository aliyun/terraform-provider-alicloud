---
subcategory: "Pai"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_feature_store_feature_entity"
sidebar_current: "docs-alicloud-resource-pai-feature-store-feature-entity"
description: |-
  Provides a Pai Feature Store Feature Entity resource.
---

# alicloud_pai_feature_store_feature_entity

Provides a Pai Feature Store Feature Entity resource.

A feature entity is the name of a set of feature tables, such as users, items, etc.

For information about Pai Feature Store Feature Entity and how to use it, see
[CreateFeatureEntity](https://www.alibabacloud.com/help/en/pai/developer-reference/api-paifeaturestore-2023-06-21-createfeatureentity).

-> **NOTE:** Available since v1.295.0.

## Example Usage

```terraform
resource "alicloud_pai_feature_store_feature_entity" "default" {
  instance_id         = "fs-cn-xxxxx"
  project_id          = "xxxxx"
  feature_entity_name = "tf-example"
  join_id             = "user_id"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the FeatureStore instance, which can be obtained through the ListInstances interface.
* `project_id` - (Required, ForceNew) The ID of the FeatureStore project, which can be obtained through the ListProjects interface.
* `feature_entity_name` - (Required, ForceNew) The feature entity name. Start with a letter or number and contain only letters, numbers, underscores, and decimal points. The maximum length cannot exceed 63.
* `join_id` - (Required, ForceNew) The ID that associates multiple sets of feature sets. Generally, enter the primary key field in the user/item table, such as `user_id` or `item_id`.
* `parent_feature_entity_id` - (Optional, ForceNew) The parent feature entity ID. Empty or `0` indicates a root feature entity.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. Format: `<instance_id>:<feature_entity_id>`.
* `feature_entity_id` - The first ID of the resource.
* `create_time` - The creation time of the resource.
* `region_id` - The region ID of the resource.

## Import

Pai Feature Store Feature Entity can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_feature_store_feature_entity.default <instance_id>:<feature_entity_id>
```
