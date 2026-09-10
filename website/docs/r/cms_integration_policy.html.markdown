---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_integration_policy"
description: |-
  Provides a Alicloud Cms Integration Policy resource.
---

# alicloud_cms_integration_policy

Provides a Cms Integration Policy resource.

Policies used by the Integration Center.

For information about Cms Integration Policy and how to use it, see [What is Integration Policy](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreateIntegrationPolicy).

-> **NOTE:** Available since v1.277.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cms_integration_policy&exampleId=d312d51f-4c5a-5202-60d0-49291dffebb42d88079c&activeTab=example&spm=docs.r.cms_integration_policy.0.d312d51f4c&intl_lang=EN_US" target="_blank">
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
  policy_type             = "ECS"
  integration_policy_name = var.name
  workspace               = alicloud_cms_workspace.default.id
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cms_integration_policy&spm=docs.r.cms_integration_policy.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `entity_group` - (Optional, ForceNew, Set) The entity group used to create the policy. See [`entity_group`](#entity_group) below.
* `force` - (Optional, Bool) Specifies whether to force delete the cloud native appliance. Valid values:
  - `true`: Enable.
  - `false`: Disable.

-> **NOTE:** This parameter configures deletion behavior and is only evaluated when Terraform attempts to destroy the resource. Changes to this parameter during updates are stored but have no immediate effect.

* `integration_policy_name` - (Required) The policy name.
* `policy_type` - (Required, ForceNew) The policy type.
* `workspace` - (Required, ForceNew) The workspace.

### `entity_group`

The entity_group supports the following:
* `cluster_entity_type` - (Optional, ForceNew) The cluster entity type.
* `cluster_id` - (Optional, ForceNew) The cluster ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `region_id` - The region ID of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Integration Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Integration Policy.
* `update` - (Defaults to 5 mins) Used when update the Integration Policy.

## Import

Cms Integration Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_integration_policy.example <id>
```
