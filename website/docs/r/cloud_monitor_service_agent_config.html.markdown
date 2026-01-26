---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_monitor_service_agent_config"
description: |-
  Provides a Alicloud Cloud Monitor Service Agent Config resource.
---

# alicloud_cloud_monitor_service_agent_config

Provides a Cloud Monitor Service Agent Config resource.

Cloud monitoring plug-in global configuration.

For information about Cloud Monitor Service Agent Config and how to use it, see [What is Agent Config](https://next.api.alibabacloud.com/document/Cms/2019-01-01/PutMonitoringConfig).

-> **NOTE:** Available since v1.270.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_cloud_monitor_service_agent_config" "default" {
  enable_install_agent_new_ecs = false
}
```

### Deleting `alicloud_cloud_monitor_service_agent_config` or removing it from your configuration

Terraform cannot destroy resource `alicloud_cloud_monitor_service_agent_config`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `enable_install_agent_new_ecs` - (Optional) Whether the cloud monitoring plug-in is automatically installed on the newly purchased ECS host. Value:
  - true (default): The cloud monitoring plug-in is automatically installed on the newly purchased ECS host.
  - false: The cloud monitoring plug-in is not automatically installed on the newly purchased ECS host.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<Alibaba Cloud Account ID>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Agent Config.
* `update` - (Defaults to 5 mins) Used when update the Agent Config.

## Import

Cloud Monitor Service Agent Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_monitor_service_agent_config.example <Alibaba Cloud Account ID>
```