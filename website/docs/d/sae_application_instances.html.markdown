---
subcategory: "Serverless App Engine (SAE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sae_application_instances"
sidebar_current: "docs-alicloud-datasource-sae-application-instances"
description: |-
  Provides a list of application instances of the specified Serverless App Engine (SAE) application.
---

# alicloud_sae_application_instances

This data source provides the application instances of the specified Serverless App Engine (SAE) application, including the ENI (Elastic Network Interface) IP addresses of the instances.

-> **NOTE:** Available since v2.1.0.

## Example Usage

```terraform
data "alicloud_sae_application_instances" "default" {
  application_id = "d700e680-aa4d-4ec1-afc2-6566b5ff****"
}

output "instance_container_ips" {
  value = [for i in data.alicloud_sae_application_instances.default.instances : i.instance_container_ip]
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Required) The ID of the SAE application.
* `group_id` - (Optional) The ID of the application group to which the instances belong. If omitted, the data source queries the instances of every group of the application and aggregates the results.
* `ids` - (Optional) A list of instance IDs, used as a filter. It is also exported after the query.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the data source, which is the same as `application_id`.
* `ids` - A list of SAE application instance IDs.
* `instances` - A list of SAE application instances. Each element contains the following attributes:
  * `instance_id` - The ID of the application instance.
  * `group_id` - The ID of the application group to which the instance belongs.
  * `instance_container_ip` - The intranet IP (ENI IP) of the application instance.
  * `instance_container_status` - The running status of the application instance.
  * `instance_health_status` - The health check status of the application instance.
  * `main_container_status` - The status of the main container of the application instance.
  * `image_url` - The image repository address of the application instance.
  * `package_version` - The version of the deployment package of the application instance.
  * `vswitch_id` - The ID of the VSwitch where the application instance resides.
