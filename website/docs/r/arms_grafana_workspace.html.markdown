---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_grafana_workspace"
description: |-
  Provides a Alicloud ARMS Grafana Workspace resource.
---

# alicloud_arms_grafana_workspace

Provides a ARMS Grafana Workspace resource. 

For information about ARMS Grafana Workspace and how to use it, see [What is Grafana Workspace](https://next.api.alibabacloud.com/document/ARMS/2019-08-08/ListGrafanaWorkspace).

-> **NOTE:** Available since v1.215.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_arms_grafana_workspace&exampleId=a1305a5b-7e08-ac51-2ea2-69366274c76a2e4bcd03&activeTab=example&spm=docs.r.arms_grafana_workspace.0.a1305a5b7e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

resource "alicloud_arms_grafana_workspace" "default" {
  grafana_version           = "9.0.x"
  description               = var.name
  resource_group_id         = data.alicloud_resource_manager_resource_groups.default.ids.0
  grafana_workspace_edition = "standard"
  grafana_workspace_name    = var.name
  tags = {
    Created = "tf"
    For     = "example"
  }
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) Description.
* `grafana_version` - (Optional) The version of the grafana.
* `grafana_workspace_edition` - (Optional, ForceNew) The edition of the grafana.
* `grafana_workspace_name` - (Optional) The name of the resource.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `tags` - (Optional, Map) The tag of the resource.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Grafana Workspace.
* `delete` - (Defaults to 5 mins) Used when delete the Grafana Workspace.
* `update` - (Defaults to 5 mins) Used when update the Grafana Workspace.

## Import

ARMS Grafana Workspace can be imported using the id, e.g.

```shell
$ terraform import alicloud_arms_grafana_workspace.example <id>
```