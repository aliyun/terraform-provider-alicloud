---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_application_instances"
description: |-
  Provides a list of Sae Application Instances to the user.
---

# alicloud_sae_application_instances

This data source provides the Sae Application Instances of the current Alibaba Cloud user.

-> **NOTE:** Available since v2.1.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_sae_applications" "default" {
  app_name = "default-NODELETING-image"
}

data "alicloud_sae_application_groups" "default" {
  app_id = data.alicloud_sae_applications.default.applications.0.id
}

data "alicloud_sae_application_instances" "default" {
  app_id   = data.alicloud_sae_application_groups.default.app_id
  group_id = data.alicloud_sae_application_groups.default.groups.0.id
}

output "sae_application_instances_id_0" {
  value = data.alicloud_sae_application_instances.default.instances.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, List) A list of Instance IDs.
* `app_id` - (Required) The application ID.
* `group_id` - (Required) The ID of the application group. If you want to get the `group_id`, please use the data source [alicloud_sae_application_groups](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/sae_application_groups).
* `instance_id` - (Optional) The application instance ID.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `instances` - A list of Instances. Each element contains the following attributes:
  * `id` - The ID of the Instance.
  * `instance_id` - The instance ID.
  * `group_id` - The ID of the group to which the instance belongs.
  * `instance_container_ip` - The internal IP address of the instance.
  * `instance_container_restarts` - The number of times the instance container has restarted.
  * `instance_container_status` - The status of the instance container.
  * `instance_health_status` - The health check status of the instance.
  * `vswitch_id` - The ID of the vSwitch where the instance is deployed.
  * `package_version` - The version of the deployment package.
  * `image_url` - The URL of the container image.
  * `debug_status` - Whether the instance is being debugged.
  * `eip` - The elastic IP address (EIP).
  * `unhealthy_message` - The error message for a failed health check on the application instance. This field is empty if the check is successful.
  * `traffic_status` - The traffic status of the instance.
  * `main_container_status` - The status of the main container.
  * `timestamp` - The timestamp.
  * `create_time_stamp` - The time when instance creation began.
  * `finish_time_stamp` - The time when instance creation was completed.
  * `sidecar_containers_status` - The status of the sidecar containers.
    * `container_id` - The ID of the sidecar container.
    * `image_url` - The URL of the container image.
    * `container_status` - The status of the container.
