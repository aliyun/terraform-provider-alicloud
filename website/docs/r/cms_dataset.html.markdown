---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_dataset"
description: |-
  Provides a Alicloud Cms Dataset resource.
---

# alicloud_cms_dataset

Provides a Cms Dataset resource.



For information about Cms Dataset and how to use it, see [What is Dataset](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreateDataset).

-> **NOTE:** Available since v1.292.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cms_dataset&exampleId=a369a658-e27e-2615-dce3-a47f747dc0d99cc86233&activeTab=example&spm=docs.r.cms_dataset.0.a369a658e2&intl_lang=EN_US" target="_blank">
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

resource "alicloud_cms_dataset" "default" {
  workspace    = alicloud_cms_workspace.default.workspace_name
  dataset_name = var.name
  description  = "terraform-example"
  schema       = jsonencode({ type = "record", name = "example", fields = [{ name = "metric", type = "string" }] })
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cms_dataset&spm=docs.r.cms_dataset.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `dataset_name` - (Required, ForceNew) The name of the resource.
* `description` - (Optional) The description of the dataset.
* `schema` - (Required, ForceNew) The schema definition of the dataset, in JSON format.
* `workspace` - (Required, ForceNew) The name of the workspace to which the dataset belongs.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. It is formatted as `<workspace>:<dataset_name>`.
* `create_time` - The creation time of the resource.
* `region_id` - The region ID of the resource.
* `update_time` - The last modified time of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Dataset.
* `delete` - (Defaults to 5 mins) Used when delete the Dataset.
* `update` - (Defaults to 5 mins) Used when update the Dataset.

## Import

Cms Dataset can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_dataset.example <workspace>:<dataset_name>
```
