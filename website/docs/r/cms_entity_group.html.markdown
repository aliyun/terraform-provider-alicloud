---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_entity_group"
sidebar_current: "docs-alicloud-resource-cms-entity-group"
description: |-
  Provides a Alicloud Cloud Monitor Service Entity Group resource.
---

# alicloud_cms_entity_group

Provides a Cloud Monitor Service Entity Group resource.

For information about Cloud Monitor Service Entity Group and how to use it, see [What is Entity Group](https://www.alibabacloud.com/help/en/cloudmonitor/developer-reference/api-cms-2024-03-30-createentitygroup).

-> **NOTE:** Available since v1.254.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_log_project" "default" {
  project_name = "tf-entity-group-example"
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = "tf-entity-group-example"
  sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_entity_group" "default" {
  entity_group_name = "tf-entity-group-example"
  description       = "terraform entity group example"
  workspace         = alicloud_cms_workspace.default.id

  entity_rules {
    entity_types = ["ECS"]
    instance_ids = ["i-bp1example0001"]

    tags {
      op         = "add"
      tag_key    = "Env"
      tag_values = ["test"]
    }

    ip_match_rule {
      ip_field_key = "ip"
      ip_cidr      = "192.168.0.0/16"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `entity_group_name` - (Optional) The name of the entity group.
* `description` - (Optional) The description of the entity group.
* `resource_group_id` - (Optional) The ID of the resource group. Maps to `entityRules.resourceGroupId` in the API request.
* `workspace` - (Required, ForceNew) The workspace to which the entity group belongs.
* `entity_rules` - (Optional) The rules of the entity group. A list with at most one element. See [`entity_rules`](#entity_rules) below.

### entity_rules

The `entity_rules` block supports:

* `resource_group_id` - (Optional) The resource group id within the entity rules.
* `tags` - (Optional) A list of tag rules. See [`tags`](#tags) below.
* `labels` - (Optional) A list of label rules. See [`labels`](#labels) below.
* `annotations` - (Optional) A list of annotation rules. See [`annotations`](#annotations) below.
* `ip_match_rule` - (Optional) The IP match rule. A list with at most one element. See [`ip_match_rule`](#ip_match_rule) below.
* `instance_ids` - (Optional) A list of instance IDs.
* `field_rules` - (Optional) A list of field rules. See [`field_rules`](#field_rules) below.
* `entity_types` - (Optional) A list of entity types. Valid values: `Cloud`, `ECS`, `CS`.

### tags

The `tags` block supports:

* `op` - (Optional) The operation type.
* `tag_key` - (Optional) The tag key.
* `tag_values` - (Optional) A list of tag values.

### labels

The `labels` block supports:

* `op` - (Optional) The operation type.
* `tag_key` - (Optional) The label key.
* `tag_values` - (Optional) A list of label values.

### annotations

The `annotations` block supports:

* `op` - (Optional) The operation type.
* `tag_key` - (Optional) The annotation key.
* `tag_values` - (Optional) A list of annotation values.

### ip_match_rule

The `ip_match_rule` block supports:

* `ip_field_key` - (Optional) The IP field key.
* `ip_cidr` - (Optional) The IP CIDR block.

### field_rules

The `field_rules` block supports:

* `field_key` - (Optional) The field key.
* `op` - (Optional) The operation type.
* `field_values` - (Optional) A list of field values.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the entity group. It is the value of `entity_group_id`.
* `entity_group_id` - The ID of the entity group.

## Import

Cloud Monitor Service Entity Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_entity_group.example <entity_group_id>
```
