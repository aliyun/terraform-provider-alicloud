---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_prometheus_view"
description: |-
  Provides a Alicloud Cms Prometheus View resource.
---

# alicloud_cms_prometheus_view

Provides a Cms Prometheus View resource.



For information about Cms Prometheus View and how to use it, see [What is Prometheus View](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreatePrometheusView).

-> **NOTE:** Available since v1.278.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cms_prometheus_view&exampleId=b6f10b94-17df-62a2-9d4a-519bdef167cc2dc8c68b&activeTab=example&spm=docs.r.cms_prometheus_view.0.b6f10b9417&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_account" "default" {
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_log_project" "default" {
  project_name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = "${var.name}-${random_integer.default.result}"
  sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_prometheus_instance" "default" {
  prometheus_instance_name = "${var.name}-${random_integer.default.result}"
  workspace                = alicloud_cms_workspace.default.id
}

resource "alicloud_cms_prometheus_view" "default" {
  prometheus_view_name = "${var.name}-${random_integer.default.result}"
  version              = "V2"
  prometheus_instances {
    prometheus_instance_id = alicloud_cms_prometheus_instance.default.id
    region_id              = alicloud_cms_prometheus_instance.default.region_id
    user_id                = data.alicloud_account.default.id
  }
  workspace = alicloud_cms_prometheus_instance.default.workspace
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cms_prometheus_view&spm=docs.r.cms_prometheus_view.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `auth_free_read_policy` - (Optional) Read password-free address whitelist policy.
* `enable_auth_free_read` - (Optional, Bool) Specifies whether to enable password-free read access. Valid values: `true`, `false`.
* `prometheus_instances` - (Required, Set) The list of Prometheus instances. See [`prometheus_instances`](#prometheus_instances) below.
* `prometheus_view_name` - (Required) The name of the Prometheus view.
* `version` - (Required, ForceNew) The version. Valid values: `V1`, `V2`.
* `workspace` - (Required, ForceNew) The workspace to which the environment belongs.

### `prometheus_instances`

The prometheus_instances supports the following:
* `prometheus_instance_id` - (Required) The ID of the prometheus instance.
* `region_id` - (Required) The region ID of the prometheus instance.
* `user_id` - (Required) The user ID of the prometheus instance.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the instance was created.
* `region_id` - The region ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Prometheus View.
* `delete` - (Defaults to 5 mins) Used when delete the Prometheus View.
* `update` - (Defaults to 5 mins) Used when update the Prometheus View.

## Import

Cms Prometheus View can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_prometheus_view.example <id>
```
