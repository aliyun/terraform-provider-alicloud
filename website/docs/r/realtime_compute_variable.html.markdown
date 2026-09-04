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

-> **NOTE:** Available since v1.266.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-variable-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

resource "alicloud_vpc" "default" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = "example-tf-vpc-variable"
}

resource "alicloud_vswitch" "default" {
  is_default   = false
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-beijing-g"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = "example-tf-vSwitch-variable"
}

resource "alicloud_oss_bucket" "default" {
}

resource "alicloud_realtime_compute_vvp_instance" "default" {
  vvp_instance_name = "example-tf-vvp-variable"
  storage {
    oss {
      bucket = alicloud_oss_bucket.default.id
    }
  }
  vpc_id      = alicloud_vpc.default.id
  vswitch_ids = [alicloud_vswitch.default.id]
  resource_spec {
    cpu       = "4"
    memory_gb = "16"
  }
  payment_type = "PayAsYouGo"
  zone_id      = alicloud_vswitch.default.zone_id
}

resource "alicloud_realtime_compute_variable" "default" {
  workspace   = alicloud_realtime_compute_vvp_instance.default.resource_id
  namespace   = "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"
  name        = var.name
  kind        = "Plain"
  value       = "variable-value"
  description = "example variable for deployment"
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required, ForceNew) The ID of the workspace. The workspace groups resources together and is created by the Realtime Compute VvpInstance.
* `namespace` - (Required, ForceNew) The name of the namespace (project space). The namespace is scoped under a VvpInstance and defaults to `<vvp_instance_name>-default`.
* `name` - (Required, ForceNew) The name of the variable. It must be 1 to 64 characters in length and can contain letters, digits, underscores `_`, and hyphens `-`.
* `kind` - (Required, ForceNew) The kind of the variable. Currently only `Plain` is supported.
* `value` - (Required, Sensitive) The value of the variable.
* `description` - (Optional) The description of the variable.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of the variable, formatted as `<workspace>:<namespace>:<name>`.
* `region_id` - The region ID of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the variable.
* `update` - (Defaults to 5 mins) Used when updating the variable.
* `delete` - (Defaults to 5 mins) Used when deleting the variable.

## Import

Realtime Compute Variable can be imported using the id, e.g.

```shell
terraform import alicloud_realtime_compute_variable.example <workspace>:<namespace>:<name>
```
