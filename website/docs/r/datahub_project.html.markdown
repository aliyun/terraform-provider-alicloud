---
subcategory: "Datahub Service (DataHub)"
layout: "alicloud"
page_title: "Alicloud: alicloud_datahub_project"
sidebar_current: "docs-alicloud-resource-datahub-project"
description: |-
  Provides a Alicloud datahub project resource.
---

# alicloud_datahub_project

The project is the basic unit of resource management in Datahub Service and is used to isolate and control resources. It contains a set of Topics. You can manage the datahub sources of an application by using projects. [Refer to details](https://www.alibabacloud.com/help/en/datahub/latest/nerbcz).

-> **NOTE:** Available since v1.19.0.

-> **NOTE:** Currently Datahub service only can be supported in the regions: cn-beijing, cn-hangzhou, cn-shanghai, cn-shenzhen,  ap-southeast-1.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_datahub_project&exampleId=92f7ddde-ea3d-510f-e2bf-175880dee423c5120e27&activeTab=example&spm=docs.r.datahub_project.0.92f7dddeea&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}
resource "alicloud_datahub_project" "example" {
  name    = var.name
  comment = "created by terraform"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the datahub project. Its length is limited to 3-32 and only characters such as letters, digits and '_' are allowed. It is case-insensitive.
* `comment` - (Optional) Comment of the datahub project. It cannot be longer than 255 characters.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the datahub project. It is the same as its name.
* `create_time` - Create time of the datahub project. It is a human-readable string rather than 64-bits UTC.
* `last_modify_time` - Last modify time of the datahub project. It is the same as *create_time* at the beginning. It is also a human-readable string rather than 64-bits UTC.

## Import

Datahub project can be imported using the *name* or ID, e.g.

```shell
$ terraform import alicloud_datahub_project.example tf_datahub_project
```
