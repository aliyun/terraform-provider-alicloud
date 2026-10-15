---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_addon_release"
description: |-
  Provides a Alicloud Cms Addon Release resource.
---

# alicloud_cms_addon_release

Provides a Cms Addon Release resource.

Release package of observability addon.

For information about Cms Addon Release and how to use it, see [What is Addon Release](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreateAddonRelease).

-> **NOTE:** Available since v1.280.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cms_addon_release&exampleId=518edb85-cdbe-eb3c-96ca-12f2747a3da6f7212960&activeTab=example&spm=docs.r.cms_addon_release.0.518edb85cd&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_log_project" "default" {
  project_name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = var.name
  sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_integration_policy" "default" {
  policy_type             = "Cloud"
  integration_policy_name = var.name
  workspace               = alicloud_cms_workspace.default.id
}

resource "alicloud_cms_addon_release" "default" {
  integration_policy_id = alicloud_cms_integration_policy.default.id
  addon_name            = "cloud-acs-ecs"
  addon_version         = "2.0.7"
  workspace             = alicloud_cms_integration_policy.default.workspace
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cms_addon_release&spm=docs.r.cms_addon_release.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:

* `addon_name` - (Required, ForceNew) The name of the add-on to integrate.
* `addon_release_name` - (Optional, ForceNew) The name of the release after the integration.
* `addon_version` - (Required) The version of the add-on to integrate.
* `aliyun_lang` - (Optional, ForceNew) The language of the add-on. Valid values: `zh`, `en`.
* `config` - (Optional) The metadata.
* `env_type` - (Optional, ForceNew) The environment type. Valid values: `CS`, `ECS`, `Cloud`.
* `integration_policy_id` - (Required, ForceNew) The ID of the environment policy.
* `workspace` - (Optional, ForceNew) The name of the workspace where the add-on is installed.
* `entity_rules` - (Optional, Set, Available since v2.1.0) The entity discovery rule model. See [`entity_rules`](#entity_rules) below.
* `dry_run` - (Optional, Bool) Specifies whether to perform a dry run. Valid values: `true`, `false`.

### `entity_rules`

The entity_rules supports the following:

* `resource_group_id` - (Optional) The resource group ID.
* `instance_ids` - (Optional, List) A collection of instance IDs.
* `entity_types` - (Optional, List) A collection of entity types.
* `region_ids` - (Optional, List) A collection of region IDs.
* `ip_match_rule` - (Optional, Set) A collection of IP address matching rules. See [`ip_match_rule`](#entity_rules-ip_match_rule) below.
* `tags` - (Optional, Set) The tag-based service discovery rules. See [`tags`](#entity_rules-tags) below.
* `labels` - (Optional, Set) A collection of label matching rules. See [`labels`](#entity_rules-labels) below.
* `annotations` - (Optional, Set) A collection of annotation matching rules. See [`annotations`](#entity_rules-annotations) below.
* `field_rules` - (Optional, Set) The property matching rules. See [`field_rules`](#entity_rules-field_rules) below.

### `entity_rules-ip_match_rule`

The entity_rules-ip_match_rule supports the following:

* `ip_field_key` - (Optional) The key of the IP address field.
* `ip_cidr` - (Optional) The CIDR block to match against.

### `entity_rules-tags`

The entity_rules-tags supports the following:

* `op` - (Optional) The operation.
* `tag_key` - (Optional) The key of the tag.
* `tag_values` - (Optional, List) A collection of tag values.

### `entity_rules-labels`

The entity_rules-labels supports the following:

* `op` - (Optional) The operation.
* `tag_key` - (Optional) The key of the label.
* `tag_values` - (Optional, List) A collection of label values.


### `entity_rules-annotations`

The entity_rules-annotations supports the following:

* `op` - (Optional) The operation.
* `tag_key` - (Optional) The key of the annotation.
* `tag_values` - (Optional, List) A collection of annotation values.

### `entity_rules-field_rules`

The entity_rules-field_rules supports the following:

* `op` - (Optional) The operation.
* `field_key` - (Optional) The key of the property.
* `field_values` - (Optional, List) A collection of property values.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Addon Release. It formats as `<integration_policy_id>:<addon_release_name>`.
* `create_time` - The time when the add-on was accessed.
* `region_id` - The ID of the region.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Addon Release.
* `update` - (Defaults to 5 mins) Used when update the Addon Release.
* `delete` - (Defaults to 5 mins) Used when delete the Addon Release.

## Import

Cms Addon Release can be imported using the id, which consists of integration_policy_id and addon_release_name, e.g.

```shell
$ terraform import alicloud_cms_addon_release.example <integration_policy_id>:<addon_release_name>
```
