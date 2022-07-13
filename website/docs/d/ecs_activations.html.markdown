---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_activations"
sidebar_current: "docs-alicloud-datasource-ecs-activations"
description: |-
  Provides a list of Ecs Activations to the user.
---

# alicloud\_ecs\_activations

This data source provides the Ecs Activations of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.177.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecs_activations" "ids" {}
output "ecs_activation_id_1" {
  value = data.alicloud_ecs_activations.ids.activations.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Activation IDs.
* `instance_name` - (Optional, ForceNew) The default prefix of the instance name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `activations` - A list of Ecs Activations. Each element contains the following attributes:
	* `activation_id` - The ID of the activation code.
	* `create_time` - The time when the activation code was created.
	* `deregistered_count` - The number of instances that have been logged out.
	* `description` - Description of the corresponding activation code.
	* `disabled` - Indicates whether the activation code is disabled.
	* `id` - The ID of the Activation.
	* `instance_count` - The maximum number of times the activation code is used to register a managed instance.
	* `instance_name` - The default prefix of the instance name.
	* `ip_address_range` - The IP address of the host that allows the activation code to be used.
	* `registered_count` - The number of instances that were registered.
	* `time_to_live_in_hours` - The validity period of the activation code. Unit: hours.