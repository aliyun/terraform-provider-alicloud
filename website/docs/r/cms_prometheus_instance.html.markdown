---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_prometheus_instance"
description: |-
  Provides a Alicloud Cms Prometheus Instance resource.
---

# alicloud_cms_prometheus_instance

Provides a Cms Prometheus Instance resource.



For information about Cms Prometheus Instance and how to use it, see [What is Prometheus Instance](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreatePrometheusInstance).

-> **NOTE:** Available since v1.277.0.

## Example Usage

Basic Usage

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

resource "alicloud_cms_prometheus_instance" "default" {
  prometheus_instance_name = var.name
  workspace                = alicloud_cms_workspace.default.id
}
```

## Argument Reference

The following arguments are supported:
* `archive_duration` - (Optional, Int) The number of days that data is automatically archived after the storage duration expires. Valid values: `60` to `3650`.
* `auth_free_read_policy` - (Optional) The policy for password-free read access.
* `auth_free_write_policy` - (Optional) The policy for password-free write access.
* `enable_auth_free_read` - (Optional, Bool) Specifies whether to enable password-free read access. Valid values: `true`, `false`.
* `enable_auth_free_write` - (Optional, Bool) Specifies whether to enable password-free write access. Valid values: `true`, `false`.
* `prometheus_instance_name` - (Required) The name of the instance.
* `storage_duration` - (Optional, Int) The storage duration of the instance in days.
* `workspace` - (Required, ForceNew) The workspace to which the instance belongs.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `create_time` - Instance creation time, using UTC +0 time, in the format of yyyy-MM-ddTHH:mmZ.
* `payment_type` - Payment Type.
* `region_id` - The region ID of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Prometheus Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Prometheus Instance.
* `update` - (Defaults to 5 mins) Used when update the Prometheus Instance.

## Import

Cms Prometheus Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_prometheus_instance.example <id>
```
