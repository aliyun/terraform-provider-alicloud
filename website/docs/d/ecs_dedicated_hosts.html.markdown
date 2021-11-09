---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_dedicated_hosts"
sidebar_current: "docs-alicloud-datasource-ecs-dedicated-hosts"
description: |-
    Provides a list of available ECS Dedicated Hosts.
---

# alicloud\_ecs\_dedicated\_hosts

This data source provides a list of ECS Dedicated Hosts in an Alibaba Cloud account according to the specified filters.
 
-> **NOTE:** Available in v1.91.0+.

## Example Usage

```
# Declare the data source
data "alicloud_ecs_dedicated_hosts" "dedicated_hosts_ds" {
  name_regex = "tf-testAcc"
  dedicated_host_type = "ddh.g5"
  status = "Available"
}

output "first_dedicated_hosts_id" {
  value = "${data.alicloud_ecs_dedicated_hosts.dedicated_hosts_ds.hosts.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of ECS Dedicated Host ids.
* `name_regex` - (Optional) A regex string to filter results by the ECS Dedicated Host name.
* `dedicated_host_id` - (Optional) The ID of ECS Dedicated Host.
* `dedicated_host_name` - (Optional) The name of ECS Dedicated Host.
* `dedicated_host_type` - (Optional) The type of the dedicated host.
* `operation_locks` - (Optional, Available in 1.123.1+) The reason why the dedicated host resource is locked.
* `resource_group_id` - (Optional) The ID of the resource group to which the ECS Dedicated Host belongs.
* `status` - (Optional) The status of the ECS Dedicated Host. validate value: `Available`, `Creating`, `PermanentFailure`, `Released`, `UnderAssessment`.
* `zone_id` - (Optional) The zone ID of the ECS Dedicated Host.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `output_file` - (Optional) Save the result to the file.

#### Block operation_locks

The operation_locks supports the following: 

* `lock_reason` - (Optional, ForceNew) The reason why the dedicated host resource is locked.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` -  A list of ECS Dedicated Host ids. 
* `names` -  A list of ECS Dedicated Host names.
* `hosts` - A list of ECS Dedicated Hosts. Each element contains the following attributes:
  * `id` - ID of the ECS Dedicated Host.
  * `action_on_maintenance` - The policy used to migrate the instances from the dedicated host when the dedicated host fails or needs to be repaired online.
  * `auto_placement` - Specifies whether to add the dedicated host to the resource pool for automatic deployment.
  * `auto_release_time` - The automatic release time of the dedicated host.
  * `cores` - A mapping of tags to assign to the resource.
  * `dedicated_host_id` - ID of the ECS Dedicated Host.
  * `dedicated_host_name` - The name of the dedicated host.
  * `dedicated_host_type` - The type of the dedicated host.
  * `description` - The description of the dedicated host.
  * `expired_time` - The expiration time of the subscription dedicated host.
  * `gpu_spec` - The GPU model.
  * `machine_id` - The machine code of the dedicated host.
  * `payment_type` - The billing method of the dedicated host.
  * `physical_gpus` - The number of physical GPUs.
  * `resource_group_id` - The ID of the resource group to which the dedicated host belongs.
  * `sale_cycle` - The unit of the subscription billing method.
  * `sockets` - The number of physical CPUs.
  * `status` - The service status of the dedicated host.
  * `supported_instance_types_list` - The list of ECS instance
  * `capacity` - (Available in 1.123.1+) A collection of proprietary host performance indicators.
    * `available_local_storage` - The remaining local disk capacity. Unit: GiB.
    * `available_memory` - The remaining memory capacity, unit: GiB.
    * `available_vcpus` - The number of remaining vCPU cores.
    * `available_vgpus` - The number of available virtual GPUs.
    * `local_storage_category` - Local disk type.
    * `total_local_storage` - The total capacity of the local disk, in GiB.
    * `total_memory` - The total memory capacity, unit: GiB.
    * `total_vcpus` - The total number of vCPU cores.
    * `total_vgpus` - The total number of virtual GPUs.
  * `cpu_over_commit_ratio` - (Available in 1.123.1+) CPU oversold ratio.
  * `network_attributes` - dedicated host network parameters. contains the following attributes:
    * `slb_udp_timeout` - The timeout period for a UDP session between Server Load Balancer (SLB) and the dedicated host. Unit: seconds.
    * `udp_timeout` - (Available in 1.123.1+) The timeout period for a UDP session between a user and an Alibaba Cloud service on the dedicated host. Unit: seconds.
  * `operation_locks` - (Available in 1.123.1+) The operation_locks. contains the following attribute:
    * `lock_reason` - The reason why the dedicated host resource is locked.
  * `supported_instance_type_families` - (Available in 1.123.1+) ECS instance type family supported by the dedicated host.
  * `supported_custom_instance_type_families` - (Available in 1.123.1+) A custom instance type family supported by a dedicated host.
  
   
    
