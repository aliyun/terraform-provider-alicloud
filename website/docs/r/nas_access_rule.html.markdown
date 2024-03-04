---
subcategory: "Network Attached Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_access_rule"
description: |-
  Provides a Alicloud NAS Access Rule resource.
---

# alicloud_nas_access_rule

Provides a NAS Access Rule resource. 

For information about NAS Access Rule and how to use it, see [What is Access Rule](https://www.alibabacloud.com/help/en/nas/developer-reference/api-nas-2017-06-26-createaccessrule).

-> **NOTE:** Available since v1.34.0.

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
* `access_group_name` - (Required, ForceNew) AccessGroupName.
* `file_system_type` - (Optional, ForceNew) filesystem type. include standard, extreme.
* `ipv6_source_cidr_ip` - (Optional, Available since v1.218.0) Ipv6SourceCidrIp.
* `priority` - (Optional, Computed) Priority.
* `rw_access_type` - (Optional, Computed) RWAccess.
* `source_cidr_ip` - (Optional) SourceCidrIp.
* `user_access_type` - (Optional, Computed) UserAccess.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<access_group_name>:<file_system_type>:<access_rule_id>`.
* `access_rule_id` - The first ID of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Access Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Access Rule.
* `update` - (Defaults to 5 mins) Used when update the Access Rule.

## Import

NAS Access Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_nas_access_rule.example <access_group_name>:<file_system_type>:<access_rule_id>
```