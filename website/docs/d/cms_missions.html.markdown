---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_missions"
sidebar_current: "docs-alicloud-datasource-cms-missions"
description: |-
  Provides a list of Cms Missions to the user.
---

# alicloud\_cms\_missions

This data source provides the Cms Missions of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.293.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_cms_missions" "default" {
  digital_employee_name = "my-digital-employee"
  name_regex            = var.name
  ids                   = ["terraform-example"]
}

output "cms_mission_id_1" {
  value = data.alicloud_cms_missions.default.missions.0.name
}
```

## Argument Reference

The following arguments are supported:

* `digital_employee_name` - (Optional, ForceNew) Filter results by the associated digital employee name.
* `ids` - (Optional, ForceNew, Computed) A list of Mission IDs (mission names).
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Mission name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `missions` - A list of Cms Missions. Each element contains the following attributes:
  * `name` - The name of the mission.
  * `display_name` - The display name of the mission.
  * `digital_employee_name` - The name of the associated digital employee.
  * `enabled` - Whether the mission is enabled.
  * `create_time` - The creation time of the mission.
  * `update_time` - The last update time of the mission.
  * `variables` - The variables of the mission.
  * `notification_policy` - The notification policy.
    * `region` - The notification policy region.
    * `workspace` - The notification policy workspace name.
