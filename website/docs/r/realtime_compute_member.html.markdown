---
subcategory: "Realtime Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_realtime_compute_member"
description: |-
  Provides a Alicloud Realtime Compute Member resource.
---

# alicloud_realtime_compute_member

Provides a Realtime Compute Member resource.

Members are used to manage read/write permissions to Flink resources within a workspace namespace.

For information about Realtime Compute Member and how to use it, see [What is Realtime Compute Member](https://next.api.alibabacloud.com/document/ververica/2022-07-18/CreateMember).

-> **NOTE:** Available since v1.265.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

resource "alicloud_realtime_compute_vvp_instance" "default" {
  vvp_instance_name = var.name
  payment_type      = "PayAsYouGo"
  vpc_id            = alicloud_vpc.default.id
  vswitch_ids       = [alicloud_vswitch.default.id]
  zone_id           = data.alicloud_zones.default.zones.0.id
  storage {
    oss {
      bucket = alicloud_oss_bucket.default.id
    }
  }
}

resource "alicloud_ram_user" "default" {
  name = var.name
}

resource "alicloud_realtime_compute_member" "default" {
  resource_id = alicloud_realtime_compute_vvp_instance.default.resource_id
  namespace   = "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"
  member      = alicloud_ram_user.default.id
  role        = "VIEWER"
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required, ForceNew) The resource ID of the VVP workspace. The member is created under this workspace.
* `namespace` - (Required, ForceNew) The name of the namespace. The namespace is created when the VVP workspace is initialized; the default namespace is `<vvp_instance_name>-default`.
* `member` - (Required, ForceNew) The RAM user ID of the member. It must be a numeric RAM user ID; the backend prepends a `user:` prefix and validates it against the pattern `user:(WORKER_|WB)?[0-9]*`.
* `role` - (Required) The role of the member. The role controls the read/write permissions of the member to Flink resources within the namespace. It is required when creating the member, but optional when updating it.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of the Realtime Compute member, formatted as `<resource_id>:<namespace>:<member>`.
* `region_id` - The region ID of the resource.

## Import

Realtime Compute Member can be imported using the id, e.g.

```shell
$ terraform import alicloud_realtime_compute_member.example <resource_id>:<namespace>:<member>
```
