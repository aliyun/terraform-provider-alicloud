---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_state_configurations"
sidebar_current: "docs-alicloud-datasource-oos-state-configurations"
description: |-
  Provides a list of Oos State Configurations to the user.
---

# alicloud\_oos\_state\_configurations

This data source provides the Oos State Configurations of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.147.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_oos_state_configurations" "ids" {}
output "oos_state_configuration_id_1" {
  value = data.alicloud_oos_state_configurations.ids.configurations.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of State Configuration IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `tags` - (Optional, Computed) The tag of the resource.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `configurations` - A list of Oos State Configurations. Each element contains the following attributes:
    * `configure_mode` - The configuration mode.
    * `create_time` - The creation time.
    * `update_time` - The time when the configuration is updated.
    * `description` - The description.
    * `id` - The ID of the State Configuration.
    * `parameters` - The parameters.
    * `resource_group_id` - The ID of the resource group.
    * `schedule_expression` - The schedule expression.
    * `schedule_type` - The schedule type.
    * `state_configuration_id` - The ID of the final state configuration.
    * `tags` - The tag of the resource.
    * `targets` - The target resource.
    * `template_id` - The ID of the template.
    * `template_name` - The name of the template.
    * `template_version` - The version of the template.