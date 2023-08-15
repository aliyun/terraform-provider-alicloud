---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_chart_namespace"
sidebar_current: "docs-alicloud-resource-cr-chart-namespace"
description: |-
  Provides a Alicloud CR Chart Namespace resource.
---

# alicloud_cr_chart_namespace

Provides a CR Chart Namespace resource.

For information about CR Chart Namespace and how to use it, see [What is Chart Namespace](https://www.alibabacloud.com/help/en/acr/developer-reference/api-cr-2018-12-01-createchartnamespace).

-> **NOTE:** Available since v1.149.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "example-name"
}
resource "alicloud_cr_ee_instance" "example" {
  payment_type   = "Subscription"
  period         = 1
  renew_period   = 0
  renewal_status = "ManualRenewal"
  instance_type  = "Advanced"
  instance_name  = var.name
}

resource "alicloud_cr_chart_namespace" "example" {
  instance_id    = alicloud_cr_ee_instance.example.id
  namespace_name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `auto_create_repo` - (Optional) Specifies whether to automatically create repositories in the namespace. Valid values:
  * `true` - automatically creates repositories in the namespace.
  * `false` - does not automatically create repositories in the namespace.
* `default_repo_type` - (Optional) DefaultRepoType. Valid values: `PRIVATE`, `PUBLIC`.
* `instance_id` - (Required, ForceNew) The ID of the Container Registry instance.
* `namespace_name` - (Required, ForceNew) The name of the namespace that you want to create.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Chart Namespace. The value formats as `<instance_id>:<namespace_name>`.

## Import

CR Chart Namespace can be imported using the id, e.g.

```shell
$ terraform import alicloud_cr_chart_namespace.example <instance_id>:<namespace_name>
```