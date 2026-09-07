---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_system_policys"
sidebar_current: "docs-alicloud-datasource-ram-system-policys"
description: |-
  Provides a list of Ram System Policy owned by an Alibaba Cloud account.
---

# alicloud_ram_system_policys

This data source provides Ram System Policy available to the user.[What is System Policy](https://next.api.alibabacloud.com/document/Ram/2015-05-01/GetPolicy)

-> **NOTE:** Available since v1.245.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_ram_system_policys" "default" {
  name_regex = "^AdministratorAccess$"
}

output "alicloud_ram_system_policy_example_id" {
  value = data.alicloud_ram_system_policys.default.policys.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed) A list of System Policy IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of System Policy IDs.
* `names` - A list of name of System Policys.
* `policys` - A list of System Policy Entries. Each element contains the following attributes:
  * `attachment_count` - Number of references.
  * `create_time` - Creation time.
  * `description` - The permission policy description.
  * `policy_name` - The permission policy name.
  * `policy_type` - Permission policy type.
  * `update_date` - Modification time.
  * `id` - The ID of the resource supplied above.
