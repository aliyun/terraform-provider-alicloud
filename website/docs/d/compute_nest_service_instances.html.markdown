---
subcategory: "Compute Nest"
layout: "alicloud"
page_title: "Alicloud: alicloud_compute_nest_service_instances"
sidebar_current: "docs-alicloud-datasource-compute-nest-service-instances"
description: |-
  Provides a list of Compute Nest Service Instances to the user.
---

# alicloud_compute_nest_service_instances

This data source provides the Compute Nest Service Instances of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.205.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_compute_nest_service_instances" "ids" {
  ids = ["example_id"]
}

output "arms_prometheis_id_1" {
  value = data.alicloud_compute_nest_service_instances.ids.service_instances.0.id
}

data "alicloud_compute_nest_service_instances" "nameRegex" {
  name_regex = "tf-example"
}

output "arms_prometheis_id_2" {
  value = data.alicloud_compute_nest_service_instances.nameRegex.service_instances.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, Computed) A list of Service Instance IDs.
* `name_regex` - (Optional) A regex string to filter results by Service Instance name.
* `status` - (Optional) The status of the Service Instance. Valid Values: `Created`, `Deploying`, `DeployedFailed`, `Deployed`, `Upgrading`, `Deleting`, `Deleted`, `DeletedFailed`.
* `filter` - (Optional) The conditions that are used to filter. See [`filter`](#filter) below.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

---

### `filter`

The filter supports the following:

* `name` - (Optional) The name of the filter. Valid Values: `Name`, `ServiceInstanceName`, `ServiceInstanceId`, `ServiceId`, `Version`, `Status`, `DeployType`, `ServiceType`, `OperationStartTimeBefore`, `OperationStartTimeAfter`, `OperationEndTimeBefore`, `OperationEndTimeAfter`, `OperatedServiceInstanceId`, `OperationServiceInstanceId`, `EnableInstanceOps`.
* `value` - (Optional) Set of values that are accepted for the given field.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Service Instance names.
* `service_instances` - A list of Service Instances.

---

### `service_instances`

The service_instances supports the following attributes:

* `id` - The ID of the Service Instance.
* `service_instance_id` - The ID of the Service Instance.
* `service_instance_name` - The name of the Service Instance.
* `parameters` - The parameters entered by the deployment service instance.
* `enable_instance_ops` - Whether the service instance has the O&M function.
* `template_name` - The name of the template.
* `operation_start_time` - The start time of O&M.
* `operation_end_time` - The end time of O&M.
* `resources` - The list of imported resources.
* `operated_service_instance_id` - The ID of the imported service instance.
* `source` - The source of the Service Instance.
* `service` - Service details.
* `tags` - The tag of the Service Instance.
* `status` - The status of the Service Instance.

---

### `service_instances.service`

The service supports the following attributes:

* `service_id` - The id of the service.
* `service_type` - The type of the service.
* `deploy_type` - The type of the deployment.
* `supplier_name` - The name of the supplier.
* `supplier_url` - The url of the supplier.
* `publish_time` - The time of publish.
* `version` - The version of the service.
* `version_name` - The version name of the service.
* `service_infos` - Service information.
* `status` - The status of the service.

---

### `service_instances.service.service_infos`

The service_infos supports the following attributes:

* `name` - The name of the service.
* `short_description` - The short description of the service.
* `image` - The image of the service.
* `locale` - The locale of the service.
