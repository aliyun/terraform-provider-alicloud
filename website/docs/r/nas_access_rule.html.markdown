---
subcategory: "Network Attached Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_access_rule"
sidebar_current: "docs-alicloud-resource-nas-access-rule"
description: |-
  Provides a Alicloud Nas Access Rule resource.
---

# alicloud\_nas_access_rule

Provides a Nas Access Rule resource.

When NAS is activated, the Default VPC Permission Group is automatically generated. It allows all IP addresses in a VPC to access the mount point with full permissions. Full permissions include Read/Write permission with no restriction on root users.

-> **NOTE:** Available in v1.34.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_nas_access_group" "standard" {
  access_group_name = "tf-NasConfigName"
  access_group_type = "Vpc"
  description       = "tf-testAccNasConfig"
}

resource "alicloud_nas_access_rule" "ipv4" {
  access_group_name = alicloud_nas_access_group.standard.access_group_name
  source_cidr_ip    = "168.1.1.0/16"
  rw_access_type    = "RDWR"
  user_access_type  = "no_squash"
  priority          = 2
}

resource "alicloud_nas_access_group" "extreme" {
  access_group_name = "tf-NasConfigName_extreme"
  access_group_type = "Vpc"
  description       = "tf-testAccNasConfig"
  file_system_type  = "extreme"
}

resource "alicloud_nas_access_rule" "ipv6" {
  access_group_name   = alicloud_nas_access_group.extreme.access_group_name
  ipv6_source_cidr_ip = "0:0:0:0:0:0:0:0/0"
  rw_access_type      = "RDWR"
  user_access_type    = "no_squash"
  priority            = 2
}

```

## Argument Reference

The following arguments are supported:

* `access_group_name` - (Required, ForceNew) Permission group name.
* `source_cidr_ip` - (Optional) Address or address segment.
* `rw_access_type` - (Optional) Read-write permission type: `RDWR` (default), `RDONLY`.
* `user_access_type` - (Optional) User permission type: `no_squash` (default), `root_squash`, `all_squash`.
* `priority` - (Optional) Priority level. Range: 1-100. Default value: `1`.
* `file_system_type` - (Optional, Available in v1.152.0+) The type of the file system. Valid values: `standard` or `extreme`. Default value: `standard`.
* `ipv6_source_cidr_ip` - (Optional, Available in v1.152.0+ and when the `file_system_type` is `extreme`) The IPv6 address or IPv6 CIDR block of the authorized object. **NOTE:** Only Extreme NAS file systems that reside in the China (Hohhot) region support IPv6. Only permission groups that reside in VPC support IPv6. This parameter is unavailable if you specify the `source_cidr_ip` parameter.

## Attributes Reference

The following attributes are exported:

* `id` - This ID of this resource. The value is formate as `<access_group_name>:<access_rule_id>:<file_system_type>`. **NOTE:** Before v1.152.0 ,The value is formats as `<access_group_name>:<access_rule_id>`
* `access_rule_id` - The nas access rule ID.

## Import

Nas Access Rule can be imported using the id, e.g.

```
$ terraform import alicloud_nas_access_rule.example <access_group_name>:<access_rule_id>:<file_system_type>
```

**NOTE:** Before v1.152.0, Nas Access Rule can be imported using the id, e.g.

```
$ terraform import alicloud_nas_access_rule.example <access_group_name>:<access_rule_id>
```

