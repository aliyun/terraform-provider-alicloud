---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_ee_namespace"
sidebar_current: "docs-alicloud-resource-cr-ee-namespace"
description: |-
  Provides a Alicloud Container Registry Enterprise Edition Namespace resource.
---

# alicloud_cr_ee_namespace

Provides a Container Registry Enterprise Edition Namespace resource.

For information about Container Registry Enterprise Edition Namespace and how to use it, see [What is Namespace](https://www.alibabacloud.com/help/en/acr/developer-reference/api-cr-2018-12-01-createnamespace)

-> **NOTE:** Available since v1.86.0.

-> **NOTE:** You need to set your registry password in Container Registry Enterprise Edition console before use this resource.

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

resource "alicloud_cr_ee_instance" "default" {
  payment_type   = "Subscription"
  period         = 1
  renew_period   = 0
  renewal_status = "ManualRenewal"
  instance_type  = "Advanced"
  instance_name  = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_cr_ee_namespace" "default" {
  instance_id        = alicloud_cr_ee_instance.default.id
  name               = "${var.name}-${random_integer.default.result}"
  auto_create        = false
  default_visibility = "PUBLIC"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the Container Registry Enterprise Edition instance.
* `name` - (Required, ForceNew) The name of the Container Registry Enterprise Edition Name. It must be `2` to `120` characters in length, and can contain lowercase letters, digits, underscores (_), hyphens (-), and periods (.). It cannot start or end with a delimiter.
* `auto_create` - (Optional, Bool) Specifies whether to automatically create an image repository in the namespace. Default value: `false`. Valid values: `true`, `false`.
* `default_visibility` - (Optional) The default type of the repository that is automatically created. Valid values:
  - `PUBLIC`: A public repository.
  - `PRIVATE`: A private repository.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Namespace. It formats as `<instance_id>:<name>`.

## Import

Container Registry Enterprise Edition Namespace can be imported using the id, e.g.

```shell
$ terraform import alicloud_cr_ee_namespace.example <instance_id>:<name>
```
