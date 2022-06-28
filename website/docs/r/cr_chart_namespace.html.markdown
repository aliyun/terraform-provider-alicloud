---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_chart_namespace"
sidebar_current: "docs-alicloud-resource-cr-chart-namespace"
description: |-
  Provides a Alicloud CR Chart Namespace resource.
---

# alicloud\_cr\_chart\_namespace

Provides a CR Chart Namespace resource.

For information about CR Chart Namespace and how to use it, see [What is Chart Namespace](https://www.alibabacloud.com/help/doc-detail/145313.htm).

-> **NOTE:** Available in v1.149.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cr_ee_instance" "default" {
  payment_type  = "Subscription"
  period        = 1
  instance_type = "Advanced"
  instance_name = "name"
}

resource "alicloud_cr_chart_namespace" "default" {
  instance_id    = alicloud_cr_ee_instance.default.id
  namespace_name = "name"
}
```

## Argument Reference

The following arguments are supported:

* `auto_create_repo` - (Optional) Specifies whether to automatically create repositories in the namespace. Valid values:
  * `true` - automatically creates repositories in the namespace.
  * `false` - does not automatically create repositories in the namespace.
* `default_repo_type` - (Optional, Computed) DefaultRepoType. Valid values: `PRIVATE`, `PUBLIC`.
* `instance_id` - (Required, ForceNew) The ID of the Container Registry instance.
* `namespace_name` - (Required, ForceNew) The name of the namespace that you want to create.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Chart Namespace. The value formats as `<instance_id>:<namespace_name>`.

## Import

CR Chart Namespace can be imported using the id, e.g.

```
$ terraform import alicloud_cr_chart_namespace.example <instance_id>:<namespace_name>
```