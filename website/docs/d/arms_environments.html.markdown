---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_environments"
description: |-
  Provides a list of ARMS Environments to the user.
---

# alicloud_arms_environments

This data source provides the ARMS Environments of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.258.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_arms_environment" "default" {
  bind_resource_id     = data.alicloud_vpcs.default.ids.0
  environment_sub_type = "ECS"
  environment_type     = "ECS"
  environment_name     = "${var.name}-${random_integer.default.result}"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids.1
  tags = {
    Created = "TF"
    For     = "Environment"
  }
}

data "alicloud_arms_environments" "ids" {
  ids = [alicloud_arms_environment.default.id]
}

output "arms_environments_id_0" {
  value = data.alicloud_arms_environments.ids.environments.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of ARMS Environment IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by ARMS Environment name.
* `environment_type` - (Optional, ForceNew) The environment type. Valid values: `CS`, `ECS`, `Cloud`.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of ARMS Environment names.
* `environments` - A list of ARMS Environments. Each element contains the following attributes:
  * `id` - The ID of the environment instance.
  * `bind_resource_id` - The ID of the resource bound to the environment instance.
  * `bind_resource_type` - The resource type.
  * `bind_vpc_cidr` - The CIDR block that is bound to the VPC.
  * `environment_id` - The ID of the environment instance.
  * `environment_name` - The name of the environment instance.
  * `environment_type` - The type of the environment instance.
  * `grafana_datasource_uid` - The unique ID of the Grafana data source.
  * `grafana_folder_uid` - The unique ID of the Grafana directory.
  * `managed_type` - Indicates whether agents or exporters are managed.
  * `prometheus_instance_id` - The ID of the Prometheus instance.
  * `region_id` - The region ID.
  * `resource_group_id` - The ID of the resource group.
  * `tags` - The tags of the environment resource.
  * `user_id` - The user ID.
