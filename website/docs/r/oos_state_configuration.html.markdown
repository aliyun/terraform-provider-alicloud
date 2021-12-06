---
subcategory: "Operation Orchestration Service (OOS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_oos_state_configuration"
sidebar_current: "docs-alicloud-resource-oos-state-configuration"
description: |-
  Provides a Alicloud OOS State Configuration resource.
---

# alicloud\_oos\_state\_configuration

Provides a OOS State Configuration resource.

For information about OOS State Configuration and how to use it, see [What is State Configuration](https://www.alibabacloud.com/help/en/doc-detail/208728.html).

-> **NOTE:** Available in v1.147.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_oos_state_configuration" "default" {
  template_name       = "ACS-ECS-InventoryDataCollection"
  configure_mode      = "ApplyOnly"
  description         = var.name
  schedule_type       = "rate"
  schedule_expression = "1 hour"
  resource_group_id   = data.alicloud_resource_manager_resource_groups.default.ids.0
  targets             = "{\"Filters\": [{\"Type\": \"All\", \"Parameters\": {\"InstanceChargeType\": \"PrePaid\"}}], \"ResourceType\": \"ALIYUN::ECS::Instance\"}"
  parameters          = "{\"policy\": {\"ACS:Application\": {\"Collection\": \"Enabled\"}}}"
  tags = {
    Created = "TF"
    For     = "Test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `configure_mode` - (Optional, Computed) Configuration mode. Valid values: `ApplyAndAutoCorrect`, `ApplyAndMonitor`, `ApplyOnly`.
* `description` - (Optional) The description of the resource.
* `parameters` - (Optional) The parameter of the Template. This field is in the format of JSON strings. For detailed definition instructions, please refer to [Metadata types that are supported by a configuration list](https://www.alibabacloud.com/help/en/doc-detail/208276.html).
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `schedule_expression` - (Required) Timing expression.
* `schedule_type` - (Required) Timing type. Valid values: `rate`.
* `tags` - (Optional) The tag of the resource.
* `targets` - (Required) The Target resources.  This field is in the format of JSON strings. For detailed definition instructions, please refer to [Parameter](https://www.alibabacloud.com/help/en/doc-detail/120674.html).
* `template_name` - (Required, ForceNew) The name of the template.
* `template_version` - (Optional, Computed, ForceNew) The version number. If you do not specify this parameter, the system uses the latest version.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of State Configuration. The value is formate as <state_configuration_id>.

## Import

OOS State Configuration can be imported using the id, e.g.

```
$ terraform import alicloud_oos_state_configuration.example <id>
```