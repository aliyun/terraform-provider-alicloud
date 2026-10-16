---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_application_instances"
sidebar_current: "docs-alicloud-datasource-sae-application-instances"
description: |-
  Provides a list of SAE application instances to the user.
---

# alicloud\_sae\_application\_instances

This data source provides the instances of an SAE application, including the ENI IP (`instance_container_ip`) of each instance.

-> **NOTE:** Available in v1.289.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_sae_application_instances" "default" {
  application_id = "47a14017-2587-4db1-80f1-3285453f7459"
}

output "sae_instance_ips" {
  value = data.alicloud_sae_application_instances.default.instances.*.instance_container_ip
}
```

Filter by instance ids:

```terraform
data "alicloud_sae_application_instances" "default" {
  application_id = "47a14017-2587-4db1-80f1-3285453f7459"
  ids            = ["ai-gateway-litellm-test-47a14017-2587-4db1-3285-qc4rh"]
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Required) The ID of the SAE application.
* `group_id` - (Optional) The ID of the instance group. If not set, instances of all groups of the application are returned.
* `ids` - (Optional) A list of instance IDs to filter results.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of instance IDs.
* `instances` - A list of instances. Each element contains the following attributes:
    * `id` - The ID of the instance.
    * `instance_id` - The ID of the instance.
    * `group_id` - The ID of the group the instance belongs to.
    * `instance_container_ip` - The ENI IP address of the instance container.
    * `instance_container_status` - The status of the instance container.
    * `instance_health_status` - The health status of the instance.
    * `main_container_status` - The status of the main container.
    * `image_url` - The image address the instance is running.
    * `package_version` - The package version the instance is running.
    * `vswitch_id` - The vSwitch ID of the ENI of the instance.
