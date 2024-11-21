---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_project"
description: |-
  Provides a Alicloud Data Works Project resource.
---

# alicloud_data_works_project

Provides a Data Works Project resource.



For information about Data Works Project and how to use it, see [What is Project](https://www.alibabacloud.com/help/en/dataworks/developer-reference/api-dataworks-public-2020-05-18-createproject).

-> **NOTE:** Available since v1.229.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_data_works_project&exampleId=1477296e-1cb3-70c6-612d-94ff543c341c267b7202&activeTab=example&spm=docs.r.data_works_project.0.1477296e1c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_data_works_project" "default" {
  project_name = "${var.name}_${random_integer.default.result}"
  project_mode = "2"
  description  = "${var.name}_${random_integer.default.result}"
  display_name = "${var.name}_${random_integer.default.result}"
  status       = "0"
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Required) Description of the workspace
* `display_name` - (Required) The display name of the workspace.
* `project_mode` - (Optional, ForceNew) The mode of the workspace, with the following values:
  - 2, indicates the simple workspace mode.
  - 3, indicating the standard workspace mode.
* `project_name` - (Required, ForceNew) Immutable Name of the workspace.
* `status` - (Optional, Computed) The status of the resource

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Project.
* `delete` - (Defaults to 5 mins) Used when delete the Project.
* `update` - (Defaults to 5 mins) Used when update the Project.

## Import

Data Works Project can be imported using the id, e.g.

```shell
$ terraform import alicloud_data_works_project.example <id>
```