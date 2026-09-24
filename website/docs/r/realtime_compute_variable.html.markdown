---
subcategory: "Realtime Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_realtime_compute_variable"
description: |-
  Provides a Alicloud Realtime Compute Variable resource.
---

# alicloud_realtime_compute_variable

Provides a Realtime Compute Variable resource.

Variable of Realtime Compute for Apache Flink, used to reference reusable values in deployments.

For information about Realtime Compute Variable and how to use it, see [What is Variable](https://next.api.alibabacloud.com/document/ververica/2022-07-18/CreateVariable).

-> **NOTE:** Available since v1.294.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_oss_buckets" "default" {
}

resource "alicloud_vpc" "default" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default" {
  is_default   = false
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-hangzhou-i"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
}

resource "alicloud_realtime_compute_vvp_instance" "default" {
  vvp_instance_name = var.name
  storage {
    oss {
      bucket = data.alicloud_oss_buckets.default.buckets.0.name
    }
  }
  vpc_id      = alicloud_vpc.default.id
  vswitch_ids = [alicloud_vswitch.default.id]
  resource_spec {
    cpu       = "8"
    memory_gb = "32"
  }
  payment_type = "PayAsYouGo"
  zone_id      = alicloud_vswitch.default.zone_id
}

resource "alicloud_realtime_compute_variable" "default" {
  name      = var.name
  namespace = "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"
  workspace = alicloud_realtime_compute_vvp_instance.default.resource_id
  kind      = "Clear"
  value     = "YourPassword123!"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of the variable.
* `kind` - (Required, ForceNew) The kind of the variable. Valid values: `Clear`, `Encrypted`.
* `name` - (Required, ForceNew) The name of the variable.
* `namespace` - (Required, ForceNew) The name of the namespace.
* `value` - (Required) The value of the variable.
* `workspace` - (Required, ForceNew) The ID of the workspace.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Variable. It formats as `<workspace>:<namespace>:<name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Variable.
* `update` - (Defaults to 5 mins) Used when update the Variable.
* `delete` - (Defaults to 5 mins) Used when delete the Variable.

## Import

Realtime Compute Variable can be imported using the id, e.g.

```shell
$ terraform import alicloud_realtime_compute_variable.example <workspace>:<namespace>:<name>
```
