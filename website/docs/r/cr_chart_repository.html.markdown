---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_chart_repository"
sidebar_current: "docs-alicloud-resource-cr-chart-repository"
description: |-
  Provides a Alicloud CR Chart Repository resource.
---

# alicloud\_cr\_chart\_repository

Provides a CR Chart Repository resource.

For information about CR Chart Repository and how to use it, see [What is Chart Repository](https://www.alibabacloud.com/help/doc-detail/145318.htm).

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

resource "alicloud_cr_chart_repository" "default" {
  repo_namespace_name = alicloud_cr_chart_namespace.default.namespace_name
  instance_id         = local.instance
  repo_name           = "repo_name"
}

```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the Container Registry instance.
* `repo_name` - (Required, ForceNew) The name of the repository that you want to create.
* `repo_namespace_name` - (Required, ForceNew) The namespace to which the repository belongs.
* `repo_type` - (Optional) The default repository type. Valid values: `PUBLIC`,`PRIVATE`.
* `summary` - (Optional) The summary about the repository.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Chart Repository. The value formats as `<instance_id>:<repo_namespace_name>:<repo_name>`.

## Import

CR Chart Repository can be imported using the id, e.g.

```
$ terraform import alicloud_cr_chart_repository.example <instance_id>:<repo_namespace_name>:<repo_name>
```