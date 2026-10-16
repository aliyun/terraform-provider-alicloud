---
subcategory: "PAI"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_resource_log"
description: |-
  Provides a PAI-EAS Resource SLS log delivery resource.
---

# alicloud_pai_resource_log

Provides a PAI-EAS Resource SLS log delivery resource. It enables delivering the
EAS resource group logs to a Log Service (SLS) project and log store. Destroying
the resource disables log delivery.

The log delivery configuration is a singleton per EAS resource group, located by
the parent resource group's `cluster_id` and `resource_id`. There is no update
API; changing `cluster_id`, `resource_id`, `project_name` or `log_store` forces
replacement (delete then recreate, i.e. disable then re-enable log delivery).

For information about PAI-EAS Resource log delivery and how to use it, see
[What is ResourceLog](https://www.alibabacloud.com/help/en/pai/developer-reference/api-eas-2021-07-01-createresourcelog).

-> **NOTE:** Available since v1.295.0.

-> **NOTE:** The SLS project referenced by `project_name` and the `log_store`
   must already exist in Log Service; this resource does not create them.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_pai_resource_log" "default" {
  cluster_id   = "eas-gateway-xxxxx"
  resource_id  = "eas-r-xxxxx"
  project_name = "eas-r-xxxxx-sls"
  log_store    = "access_log"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) The cluster id of the parent EAS resource group.
* `resource_id` - (Required, ForceNew) The resource id of the parent EAS resource group.
* `project_name` - (Required, ForceNew) The SLS project name that the resource group logs are delivered to. The project must already exist in Log Service; this resource does not create it.
* `log_store` - (Required, ForceNew) The SLS log store that the resource group logs are delivered to.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, formatted as `<cluster_id>:<resource_id>`.
* `status` - The status of the resource group.
* `message` - The SLS log delivery message returned by DescribeResourceLog. Empty when there is no extra message.

## Import

PAI-EAS Resource log delivery can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_resource_log.example <cluster_id>:<resource_id>
```
