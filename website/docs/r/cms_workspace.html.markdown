---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_workspace"
description: |-
  Provides a Alicloud Cms Workspace resource.
---

# alicloud_cms_workspace

Provides a Cms Workspace resource.



For information about Cms Workspace and how to use it, see [What is Workspace](https://next.api.alibabacloud.com/document/Cms/2024-03-30/PutWorkspace).

-> **NOTE:** Available since v1.276.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cms_workspace&exampleId=d2d6c6fe-09eb-d84e-e913-5db9a10d8651bc42bb91&activeTab=example&spm=docs.r.cms_workspace.0.d2d6c6fe09&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = var.name
  sls_project    = alicloud_log_project.default.project_name
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cms_workspace&spm=docs.r.cms_workspace.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `description` - (Optional) The description of the workspace.
* `display_name` - (Optional) The dispalyName of the workspace.
* `sls_project` - (Required, ForceNew) The project bind to workspace.
* `workspace_name` - (Required, ForceNew) The name of the workspace.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `create_time` - The creation time of the workspace.
* `region_id` - The region of the workspace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Workspace.
* `delete` - (Defaults to 5 mins) Used when delete the Workspace.
* `update` - (Defaults to 5 mins) Used when update the Workspace.

## Import

Cms Workspace can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_workspace.example <workspace_name>
```
