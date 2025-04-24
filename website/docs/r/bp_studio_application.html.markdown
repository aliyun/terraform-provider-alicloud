---
subcategory: "Cloud Architect Design Tools (BPStudio)"
layout: "alicloud"
page_title: "Alicloud: alicloud_bp_studio_application"
sidebar_current: "docs-alicloud-resource-bp-studio-application"
description: |-
  Provides a Alicloud Cloud Architect Design Tools (BPStudio) Application resource.
---

# alicloud\_bp\_studio\_application

Provides a Cloud Architect Design Tools Application resource.

For information about Cloud Architect Design Tools Application and how to use it, see [What is Application](https://help.aliyun.com/document_detail/428263.html).

-> **NOTE:** Available in v1.192.0+.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_bp_studio_application&exampleId=8c892a8d-2442-8ce4-72d0-7fde902dbace5557fe35&activeTab=example&spm=docs.r.bp_studio_application.0.8c892a8d24&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

data "alicloud_instances" "default" {
  status = "Running"
}

resource "alicloud_bp_studio_application" "default" {
  application_name  = var.name
  template_id       = "YAUUQIYRSV1CMFGX"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  area_id           = "cn-hangzhou"
  instances {
    id        = "data.alicloud_instances.default.instances.0.id"
    node_name = "data.alicloud_instances.default.instances.0.name"
    node_type = "ecs"
  }
  configuration = {
    enableMonitor = "1"
  }
  variables = {
    test = "1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `application_name` - (Required, ForceNew) The name of the application.
* `template_id` - (Required, ForceNew) The id of the template.
* `resource_group_id` - (Optional, ForceNew, Computed) The id of the resource group.
* `area_id` - (Optional, ForceNew) The id of the area.
* `instances` - (Optional, ForceNew) The instance list. Support the creation of instances in the existing vpc under the application. See the following `Block instances`.
* `configuration` - (Optional, ForceNew) The configuration of the application.
* `variables` - (Optional, ForceNew) The variables of the application.

#### Block instances

The instances supports the following:

* `id` - (Optional, ForceNew) The id of the instance.
* `node_name` - (Optional, ForceNew) The name of the instance.
* `node_type` - (Optional, ForceNew) The type of the instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Application.
* `status` - The status of the Application.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 120 mins) Used when create the Application.
* `delete` - (Defaults to 120 mins) Used when delete the Application.

## Import

Cloud Architect Design Tools Application can be imported using the id, e.g.

```shell
$ terraform import alicloud_bp_studio_application.example <id>
```
