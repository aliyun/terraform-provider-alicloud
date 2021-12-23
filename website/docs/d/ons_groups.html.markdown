---
subcategory: "RocketMQ"
layout: "alicloud"
page_title: "Alicloud: alicloud_ons_groups"
sidebar_current: "docs-alicloud-datasource-ons-groups"
description: |-
    Provides a list of ons groups available to the user.
---

# alicloud\_ons\_groups

This data source provides a list of ONS Groups in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available in 1.53.0+

## Example Usage

```terraform 
variable "name" {
  default = "onsInstanceName"
}

variable "group_name" {
  default = "GID-onsGroupDatasourceName"
}

resource "alicloud_ons_instance" "default" {
  instance_name = var.name
  remark        = "default_ons_instance_remark"
}

resource "alicloud_ons_group" "default" {
  group_name  = var.group_name
  instance_id = alicloud_ons_instance.default.id
  remark      = "dafault_ons_group_remark"
}

data "alicloud_ons_groups" "groups_ds" {
  instance_id = alicloud_ons_group.default.instance_id
  name_regex  = var.group_id
  output_file = "groups.txt"
}

output "first_group_name" {
  value = data.alicloud_ons_groups.groups_ds.groups.0.group_name
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of the ONS Instance that owns the groups.
* `group_id_regex` - (Optional) A regex string to filter results by the group name. 
* `group_name_regex` - (Optional, Available in v1.98.0+) A regex string to filter results by the group name. 
* `group_type` - (Optional, Available in v1.98.0+) Specify the protocol applicable to the created Group ID. Valid values: `tcp`, `http`. Default to `tcp`. 
* `tags` - (Optional, Available in v1.98.0+) A map of tags assigned to the Ons instance.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of group names.
* `groups` - A list of groups. Each element contains the following attributes:
  * `id` - The name of the group.
  * `group_name` - The name of the group.
  * `group_type` - Specify the protocol applicable to the created Group ID. 
  * `owner` - The ID of the group owner, which is the Alibaba Cloud UID.
  * `independent_naming` - Indicates whether namespaces are available. Read [Fields in SubscribeInfoDo](https://www.alibabacloud.com/help/doc-detail/29619.html) for further details.
  * `remark` - Remark of the group.
  * `tags` - A map of tags assigned to the Ons group.

