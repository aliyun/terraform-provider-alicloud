---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_application"
sidebar_current: "docs-alicloud-resource-oos-application"
description: |-
  Provides a Alicloud OOS Application resource.
---

# alicloud_oos_application

Provides a OOS Application resource.

For information about OOS Application and how to use it, see [What is Application](https://www.alibabacloud.com/help/en/operation-orchestration-service/latest/api-oos-2019-06-01-createapplication).

-> **NOTE:** Available since v1.145.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_oos_application&exampleId=ad8a30d7-723a-cf73-f55a-06eee21cde1503a8fde2&activeTab=example&spm=docs.r.oos_application.0.ad8a30d772&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_oos_application" "default" {
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  application_name  = "${var.name}-${random_integer.default.result}"
  description       = var.name
  tags = {
    Created = "TF"
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_oos_application&spm=docs.r.oos_application.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `application_name` - (Required) The name of the application.
* `description` - (Optional) Application group description information.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `tags` - (Optional) The tag of the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Application. The value is formate as <application_name>.

## Import

OOS Application can be imported using the id, e.g.

```shell
$ terraform import alicloud_oos_application.example <id>
```