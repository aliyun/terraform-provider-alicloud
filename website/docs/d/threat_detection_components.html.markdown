---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_components"
sidebar_current: "docs-alicloud-datasource-threat-detection-components"
description: |-
  Provides a list of Threat Detection Components owned by an Alibaba Cloud account.
---

# alicloud_threat_detection_components

This data source provides the Threat Detection Components available to the user. A component is a reusable action unit that can be used to build security orchestration playbooks.

-> **NOTE:** Available since v1.245.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_threat_detection_components" "default" {
  name_regex = "^my"
  ids        = ["my-component"]
}

output "first_component_name" {
  value = data.alicloud_threat_detection_components.default.components.0.component_name
}
```

## Argument Reference

The following arguments are supported:

* `component_name` - (Optional, ForceNew) The name of the component. The data source only returns the component that exactly matches this name.
* `ids` - (Optional, ForceNew) A list of component names. The data source returns only the components whose names match the list.
* `lang` - (Optional, ForceNew) The language type of the request and received messages. Valid values: `zh`, `en`.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by component name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `role_for` - (Optional, ForceNew) The account ID of the Resource Directory member.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of component names.
* `components` - A list of Component Entries. Each element contains the following attributes:
  * `id` - The name of the component. The value is the same as `component_name`.
  * `component_name` - The name of the component.
  * `component_alias` - The alias of the component, used for display.
  * `component_description` - The description of the component.
  * `component_logo` - The logo URL of the component.
  * `component_extension` - The extension information of the component.
  * `create_time` - The creation time of the component.
  * `update_time` - The update time of the component.
  * `component_actions` - The list of actions supported by the component. Each element contains the following attributes:
    * `component_action_name` - The name of the component action.
    * `component_action_description` - The description of the component action.
    * `input_configs` - The input configurations of the component action. Each element contains the following attributes:
      * `field_name` - The field name.
      * `field_type` - The field type.
      * `field_description` - The field description.
      * `field_display_config` - The field display configuration.
      * `default_value` - The default value of the field.
      * `required` - Whether the field is required.
    * `output_configs` - The output configurations of the component action. Each element contains the following attributes:
      * `field_name` - The field name.
      * `field_type` - The field type.
  * `component_asset_configs` - The asset configuration list of the component. Each element contains the following attributes:
    * `field_name` - The field name.
    * `field_type` - The field type.
    * `field_description` - The field description.
    * `required` - Whether the field is required.
    * `encrypted` - Whether the field is stored encrypted.
    * `default_value` - The default value of the field.
