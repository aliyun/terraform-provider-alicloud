---
subcategory: "File Storage (NAS)"
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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nas_access_rule&exampleId=a173801b-1e04-f167-4152-54456b0837be7cead684&activeTab=example&spm=docs.r.nas_access_rule.0.a173801b1e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_nas_access_group" "default" {
  access_group_type = "Vpc"
  description       = "ExtremeAccessGroup"
  access_group_name = "terraform-example-${random_integer.default.result}"
  file_system_type  = "extreme"
}

resource "alicloud_nas_access_rule" "default" {
  access_group_name   = alicloud_nas_access_group.default.access_group_name
  rw_access_type      = "RDONLY"
  ipv6_source_cidr_ip = "::1"
  user_access_type    = "no_squash"
  priority            = "1"
  file_system_type    = "extreme"
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
* `region_id` - (Available since v1.256.0) The region ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Access Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Access Rule.
* `update` - (Defaults to 5 mins) Used when update the Access Rule.

## Import

NAS Access Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_nas_access_rule.example <access_group_name>:<file_system_type>:<access_rule_id>
```
