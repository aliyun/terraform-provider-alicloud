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
resource "alicloud_nas_access_group" "foo" {
  access_group_name = "tf-NasConfigName"
  access_group_type = "Vpc"
  description       = "tf-testAccNasConfig"
}

resource "alicloud_nas_access_rule" "foo" {
  access_group_name = alicloud_nas_access_group.foo.access_group_name
  source_cidr_ip    = "168.1.1.0/16"
  rw_access_type    = "RDWR"
  user_access_type  = "no_squash"
  priority          = 2
}


```

## Argument Reference

The following arguments are supported:

* `access_group_name` - (Required, ForceNew) Permission group name.
* `source_cidr_ip` - (Required) Address or address segment.
* `rw_access_type` - (Optional) Read-write permission type: `RDWR` (default), `RDONLY`.
* `user_access_type` - (Optional) User permission type: `no_squash` (default), `root_squash`, `all_squash`.
* `priority` - (Optional) Priority level. Range: 1-100. Default value: `1`.

## Attributes Reference

The following attributes are exported:

* `id` - This ID of this resource. The value is formate as `<access_group_name>:<access rule id>`.
* `access_rule_id` - The nas access rule ID.

## Import

Nas Access Rule can be imported using the id, e.g.

```
$ terraform import alicloud_nas_access_rule.foo tf-testAccNasConfigName:1
```

