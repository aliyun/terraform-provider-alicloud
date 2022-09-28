---
subcategory: "Network Attached Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_access_rule"
sidebar_current: "docs-alicloud-resource-nas-access-rule"
description: |-
  Provides a Alicloud Nas Access Rule resource.
---

# alicloud\_all_nas_access_rule

Provides a Nas Access Rule resource.

When NAS is activated, the Default VPC Permission Group is automatically generated. It allows all IP addresses in a VPC to access the mount point with full permissions. Full permissions include Read/Write permission with no restriction on root users.

-> **NOTE:** Available in v1.191.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_nas_access_group" "foo" {
  access_group_name = "tf-NasConfigName"
  access_group_type = "Vpc"
  description       = "tf-testAccNasConfig"
  file_system_type  = "extreme"
}

resource "alicloud_all_nas_access_rule" "foo" {
  access_group_name = alicloud_nas_access_group.foo.access_group_name
  source_cidr_ip    = "168.1.1.0/16"
  rw_access_type    = "RDWR"
  user_access_type  = "no_squash"
  priority          = 2
  file_system_type  = "extreme"
}

resource "alicloud_all_nas_access_rule" "foo1" {
  access_group_name   = alicloud_nas_access_group.foo.access_group_name
  ipv6_source_cidr_ip = "::1"
  rw_access_type      = "RDWR"
  user_access_type    = "no_squash"
  priority            = 2
  file_system_type    = "extreme"
}

```

## Argument Reference

The following arguments are supported:

* `access_group_name` - (Required, ForceNew) Permission group name.
* `source_cidr_ip` - (Optional) Address or address segment.
* `rw_access_type` - (Optional) Read-write permission type: `RDWR` , `RDONLY`.
* `user_access_type` - (Optional) User permission type: `no_squash` , `root_squash`, `all_squash`.
* `priority` - (Optional) Priority level. Range: 1-100. Default value: `1`.
* `file_system_type` - (Required, ForceNew) The type of the file system: `standard `, `extreme`.
* `ipv6_source_cidr_ip` - (Optional) The IPv6 address or CIDR block of the authorized object.
                                     You must set one of "ipv6_source_cidr_ip" and "source_cidr_ip".

## Attributes Reference

The following attributes are exported:

* `id` - This ID of this resource. The value is format as `<access_group_name>:<access_rule_id>:<file_system_type>`.
* `access_rule_id` - The nas access rule ID.

## Import

Nas Access Rule can be imported using the id, e.g.

```
$ terraform import alicloud_nas_access_rule.foo tf-testAccNasConfigName:1:standard
```

