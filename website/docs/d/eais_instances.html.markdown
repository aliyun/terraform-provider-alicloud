---
subcategory: "Elastic Accelerated Computing Instances (EAIS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eais_instances"
sidebar_current: "docs-alicloud-datasource-eais-instances"
description: |-
  Provides a list of Eais Instances to the user.
---

# alicloud\_eais\_instances

This data source provides the Eais Instances of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.137.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_eais_instances" "ids" {
  id = ["example_id"]
}
output "eais_instance_id_1" {
  value = data.alicloud_eais_instances.ids.instances.0.id
}

data "alicloud_eais_instances" "nameRegex" {
  name_regex = "^my-Instance"
}
output "eais_instance_id_2" {
  value = data.alicloud_eais_instances.nameRegex.instances.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Instance IDs.
* `instance_name` - (Optional, ForceNew) The Name of the instance.
* `instance_type` - (Optional, ForceNew) EAIS instance type. Valid values: `eais.ei-a6.4xlarge`, `eais.ei-a6.2xlarge`, `eais.ei-a6.xlarge`, `eais.ei-a6.large`, `eais.ei-a6.medium`.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Instance name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `Attaching`, `Available`, `Detaching`, `InUse`, `Starting`, `Unavailable`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Instance names.
* `instances` - A list of EAIS Instances. Each element contains the following attributes:
	* `client_instance_id` - The ID of the ECS instance to be bound.
	* `client_instance_name` - The name of the ECS instance bound to the EAIS instance.
	* `client_instance_type` - The type of the ECS instance bound to the EAIS instance.
	* `id` - The ID of the Instance.
	* `instance_id` - The ID of the resource.
	* `instance_name` - The name of the resource.
	* `instance_type` - The type of the resource. Valid values: `eais.ei-a6.4xlarge`, `eais.ei-a6.2xlarge`, `eais.ei-a6.xlarge`, `eais.ei-a6.large`, `eais.ei-a6.medium`.
	* `status` - The status of the resource. Valid values: `Attaching`, `Available`, `Detaching`, `InUse`, `Starting`, `Unavailable`.
	* `zone_id` - The ID of the region to which the EAIS instance belongs.
