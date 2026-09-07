---
subcategory: "Realtime Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_realtime_compute_member"
description: |-
  Provides a Alicloud Realtime Compute Member resource.
---

# alicloud_realtime_compute_member

Provides a Realtime Compute Member resource.



For information about Realtime Compute Member and how to use it, see [What is Member](https://next.api.alibabacloud.com/document/ververica/2022-07-18/CreateMember).

-> **NOTE:** Available since v1.292.0.

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

resource "alicloud_ram_user" "default" {
  name         = var.name
  display_name = "displayname"
  mobile       = "86-18888888888"
  email        = "hello.uuu@aaa.com"
  comments     = "yoyoyo"
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

resource "alicloud_realtime_compute_member" "default" {
  member      = alicloud_ram_user.default.id
  namespace   = "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"
  resource_id = alicloud_realtime_compute_vvp_instance.default.resource_id
  role        = "viewer"
}
```

## Argument Reference

The following arguments are supported:

* `member` - (Required, ForceNew) The name of the member.
* `namespace` - (Required, ForceNew) The name of the namespace.
* `resource_id` - (Required, ForceNew) The workspace ID.
* `role` - (Required) The role of the member. Valid values: `viewer`, `owner`, `editor`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Member. It formats as `<resource_id>:<namespace>:<member>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Member.
* `update` - (Defaults to 5 mins) Used when update the Member.
* `delete` - (Defaults to 5 mins) Used when delete the Member.

## Import

Realtime Compute Member can be imported using the id, e.g.

```shell
$ terraform import alicloud_realtime_compute_member.example <resource_id>:<namespace>:<member>
```
