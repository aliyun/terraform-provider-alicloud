---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_entity_groups"
sidebar_current: "docs-alicloud-datasource-cms-entity-groups"
description: |-
  Provides a list of Cloud Monitor Service Entity Groups to the user.
---

# alicloud_cms_entity_groups

This data source provides the Cloud Monitor Service Entity Groups of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.254.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_entity_groups" "default" {
  workspace         = alicloud_cms_workspace.default.id
  entity_group_type = "ECS"
  ids               = ["<entity_group_id>"]
}

output "first_entity_group_name" {
  value = data.alicloud_cms_entity_groups.default.groups.0.entity_group_name
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of entity group IDs.
* `name` - (Optional) The name of the entity group used to filter results.
* `entity_group_type` - (Optional) The type of the entity group. Valid values: `Cloud`, `ECS`, `CS`.
* `workspace` - (Required) The workspace to which the entity groups belong.
* `output_file` - (Optional) The name of the file that saves the data source results.

## Attributes Reference

The following attributes are exported in addition to the `ids` and `name` arguments:

* `groups` - A list of entity groups. Each element contains the following attributes:
  * `entity_group_id` - The ID of the entity group.
  * `entity_group_name` - The name of the entity group.
  * `description` - The description of the entity group.
  * `workspace` - The workspace of the entity group.
  * `resource_group_id` - The resource group id of the entity group.
  * `entity_rules` - The rules of the entity group. See [`entity_rules`](#entity_rules) below.

### entity_rules

The `entity_rules` block contains:

* `resource_group_id` - The resource group id within the entity rules.
* `tags` - A list of tag rules.
* `labels` - A list of label rules.
* `annotations` - A list of annotation rules.
* `ip_match_rule` - The IP match rule.
* `instance_ids` - A list of instance IDs.
* `field_rules` - A list of field rules.
* `entity_types` - A list of entity types.
