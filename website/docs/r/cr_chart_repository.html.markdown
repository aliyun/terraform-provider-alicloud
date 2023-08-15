---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_chart_repository"
sidebar_current: "docs-alicloud-resource-cr-chart-repository"
description: |-
  Provides a Alicloud CR Chart Repository resource.
---

# alicloud_cr_chart_repository

Provides a CR Chart Repository resource.

For information about CR Chart Repository and how to use it, see [What is Chart Repository](https://www.alibabacloud.com/help/en/acr/developer-reference/api-cr-2018-12-01-createchartrepository).

-> **NOTE:** Available since v1.149.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
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

resource "alicloud_cr_chart_repository" "example" {
  repo_namespace_name = alicloud_cr_chart_namespace.example.namespace_name
  instance_id         = alicloud_cr_chart_namespace.example.instance_id
  repo_name           = var.name
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

```shell
$ terraform import alicloud_cr_chart_repository.example <instance_id>:<repo_namespace_name>:<repo_name>
```