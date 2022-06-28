---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_deployment_sets"
sidebar_current: "docs-alicloud-datasource-ecs-deployment-sets"
description: |-
  Provides a list of Ecs Deployment Sets to the user.
---

# alicloud\_ecs\_deployment\_sets

This data source provides the Ecs Deployment Sets of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.140.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecs_deployment_sets" "ids" {
  ids = ["example_id"]
}
output "ecs_deployment_set_id_1" {
  value = data.alicloud_ecs_deployment_sets.ids.sets.0.id
}

data "alicloud_ecs_deployment_sets" "nameRegex" {
  name_regex = "^my-DeploymentSet"
}
output "ecs_deployment_set_id_2" {
  value = data.alicloud_ecs_deployment_sets.nameRegex.sets.0.id
}

```

## Argument Reference

The following arguments are supported:

* `deployment_set_name` - (Optional, ForceNew) The name of the deployment set.
* `ids` - (Optional, ForceNew, Computed)  A list of Deployment Set IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Deployment Set name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `strategy` - (Optional, ForceNew) The deployment strategy. Valid values: `Availability`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Deployment Set names.
* `sets` - A list of Ecs Deployment Sets. Each element contains the following attributes:
	* `create_time` - The time when the deployment set was created.
	* `deployment_set_id` - The ID of the Deployment Set.
	* `deployment_set_name` - The name of the deployment set.
	* `description` - The description of the deployment set.
	* `domain` - The deployment domain.
	* `granularity` - The deployment granularity.
	* `id` - The ID of the Deployment Set.
	* `instance_amount` - The number of instances in the deployment set.
	* `instance_ids` - The IDs of the instances in the deployment set.
	* `strategy` - The deployment strategy.
