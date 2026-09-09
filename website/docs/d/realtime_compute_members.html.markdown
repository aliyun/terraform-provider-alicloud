---
subcategory: "Realtime Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_realtime_compute_members"
description: |-
  Provides a list of Realtime Compute Members to the user.
---

# alicloud_realtime_compute_members

This data source provides the Realtime Compute Member of the current Alibaba Cloud user.

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

data "alicloud_realtime_compute_members" "ids" {
  ids         = [alicloud_realtime_compute_member.default.id]
  namespace   = alicloud_realtime_compute_member.default.namespace
  resource_id = alicloud_realtime_compute_member.default.resource_id
}

output "realtime_compute_members_id_0" {
  value = data.alicloud_realtime_compute_members.ids.members.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, List) A list of Member IDs. It formats as `<resource_id>:<namespace>:<member>`.
* `namespace` - (Required) The name of the namespace.
* `resource_id` - (Required) The workspace ID.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `members` - A list of Members. Each element contains the following attributes:
  * `id` - The ID of the Member.
  * `resource_id` - The workspace ID.
  * `namespace` - The name of the namespace.
  * `member` - The member UID.
  * `role` - The member role.
