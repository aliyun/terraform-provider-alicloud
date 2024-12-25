---
subcategory: "Schedulerx"
layout: "alicloud"
page_title: "Alicloud: alicloud_schedulerx_namespace"
description: |-
  Provides a Alicloud Schedulerx Namespace resource.
---

# alicloud_schedulerx_namespace

Provides a Schedulerx Namespace resource.



For information about Schedulerx Namespace and how to use it, see [What is Namespace](https://www.alibabacloud.com/help/en/schedulerx/schedulerx-serverless/developer-reference/api-schedulerx2-2019-04-30-listnamespaces).

-> **NOTE:** Available since v1.173.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_schedulerx_namespace" "default" {
  namespace_name = var.name
  description    = var.name
}
```

### Deleting `alicloud_schedulerx_namespace` or removing it from your configuration

Terraform cannot destroy resource `alicloud_schedulerx_namespace`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `description` - (Optional, ForceNew) Namespace description.
* `namespace_name` - (Required, ForceNew) Namespace name.
* `namespace_uid` - (Optional, ForceNew, Computed,  Available since v1.240.0) Namespace uid.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Namespace.

## Import

Schedulerx Namespace can be imported using the id, e.g.

```shell
$ terraform import alicloud_schedulerx_namespace.example <id>
```