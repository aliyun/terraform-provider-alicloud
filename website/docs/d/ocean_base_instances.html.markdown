---
subcategory: "Ocean Base"
layout: "alicloud"
page_title: "Alicloud: alicloud_ocean_base_instances"
sidebar_current: "docs-alicloud-datasource-ocean-base-instances"
description: |-
  Provides a list of Ocean Base Instances to the user.
---

# alicloud\_ocean\_base\_instances

This data source provides the Ocean Base Instances of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.203.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ocean_base_instances" "ids" {}
output "ocean_base_instance_id_1" {
  value = data.alicloud_ocean_base_instances.ids.instances.0.id
}

data "alicloud_ocean_base_instances" "nameRegex" {
  name_regex = "^my-Instance"
}
output "ocean_base_instance_id_2" {
  value = data.alicloud_ocean_base_instances.nameRegex.instances.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Instance IDs.
* `instance_id` - (ForceNew,Optional) OceanBase cluster ID.
* `resource_group_id` - (ForceNew,Optional) The ID of the enterprise resource group to which the instance resides.
* `search_key` - (ForceNew,Optional) The filter keyword for the query list.
* `instance_name` - (Optional, ForceNew) OceanBase cluster name. The length is 1 to 20 English or Chinese characters. If this parameter is not specified, the default value is the InstanceId of the cluster.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Instance name.
* `status` - (Optional, ForceNew) The status of the Instance. Value range: `PENDING_CREATE`, `ONLINE`, `TENANT_CREATING`, `TENANT_SPEC_MODIFYING`, `EXPANDING`, `REDUCING`, `SPEC_UPGRADING`, `DISK_UPGRADING`, `WHITE_LIST_MODIFYING`, `PARAMETER_MODIFYING`, `SSL_MODIFYING`, `PREPAID_EXPIRE_CLOSED`, `ARREARS_CLOSED`, `PENDING_DELETE`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Instance names.
* `instances` - A list of Ocean Base Instances. Each element contains the following attributes:
	* `commodity_code` - The product code of the OceanBase cluster.
	* `cpu` - The number of CPU cores of the cluster.
	* `node_num` - The number of nodes in the cluster.
	* `create_time` - The creation time of the resource.
	* `disk_size` - The size of the storage space, in GB.
	* `id` - The ID of the Instance.
	* `instance_class` - Cluster specification information.
	* `instance_id` - OceanBase cluster ID.
	* `instance_name` - OceanBase cluster name.
	* `payment_type` - The payment method of the instance.
	* `resource_group_id` - The ID of the enterprise resource group to which the instance resides.
	* `series` - Series of OceanBase clusters.
	* `status` - The status of the resource.
	* `zones` - Information about the zone where the cluster is deployed.