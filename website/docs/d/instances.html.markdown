---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_instances"
sidebar_current: "docs-alicloud-datasource-instances"
description: |-
  Provides a list of ECS instances to the user.
---

# alicloud_instances

The Instances data source list ECS instance resources according to their ID, name regex, image id, status and other fields.

-> **NOTE:** Available since v1.7.0.

## Example Usage

```terraform
data "alicloud_instances" "instances_ds" {
  name_regex = "web_server"
  status     = "Running"
}

output "first_instance_id" {
  value = "${data.alicloud_instances.instances_ds.instances.0.id}"
}

output "instance_ids" {
  value = "${data.alicloud_instances.instances_ds.ids}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of ECS instance IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by instance name.
* `image_id` - (Optional, ForceNew) The image ID of some ECS instance used.
* `status` - (Optional, ForceNew) Instance status. Valid values: "Creating", "Starting", "Running", "Stopping" and "Stopped". If undefined, all statuses are considered.
* `vpc_id` - (Optional, ForceNew) ID of the VPC linked to the instances.
* `vswitch_id` - (Optional, ForceNew) ID of the vSwitch linked to the instances.
* `availability_zone` - (Optional, ForceNew) Availability zone where instances are located.
* `resource_group_id` - (Optional, ForceNew, Available since v1.57.0) The ID of resource group which the instance belongs.
* `ram_role_name` - (Optional, ForceNew, Available since v1.69.0) The RAM role name which the instance attaches.
* `instance_name` - (Optional, ForceNew, Available since v1.203.0) The name of the instance. Fuzzy search with the asterisk (*) wildcard characters is supported.
* `tags` - (Optional, ForceNew) A map of tags assigned to the ECS instances. It must be in the format:
  ```
  data "alicloud_instances" "taggedInstances" {
    tags = {
      tagKey1 = "tagValue1",
      tagKey2 = "tagValue2"
    }
  }
  ```
* `enable_details` - (Optional, Available since v1.204.0) Default to `true`. If false, the attributes `ram_role_name` and `disk_device_mappings` will not be fetched and output.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of instances names. 
* `instances` - A list of instances. Each element contains the following attributes:
  * `id` - ID of the instance.
  * `region_id` - Region ID the instance belongs to.
  * `availability_zone` - Availability zone the instance belongs to.
  * `status` - Instance current status.
  * `name` - Instance name.
  * `description` - Instance description.
  * `instance_type` - Instance type.
  * `vpc_id` - ID of the VPC the instance belongs to.
  * `vswitch_id` - ID of the vSwitch the instance belongs to.
  * `image_id` - Image ID the instance is using.
  * `private_ip` - Instance private IP address.
  * `public_ip` - Instance public IP address.
  * `eip` - EIP address the VPC instance is using.
  * `security_groups` - List of security group IDs the instance belongs to.
  * `key_name` - Key pair the instance is using.
  * `creation_time` - Instance creation time.
  * `instance_charge_type` - Instance charge type.
  * `internet_charge_type` - Instance network charge type.
  * `internet_max_bandwidth_out` - Max output bandwidth for internet.
  * `spot_strategy` - Spot strategy the instance is using.
  * `disk_device_mappings` - Description of the attached disks.
    * `device` - Device information of the created disk: such as /dev/xvdb.
    * `size` - Size of the created disk.
    * `category` - Cloud disk category.
    * `type` - Cloud disk type: system disk or data disk.
    * `disk_id` - The ID of the Disk.
    * `disk_name` - The name of the Disk.
  * `tags` - A map of tags assigned to the ECS instance.
  * `resource_group_id` - The Id of resource group.
  * `ram_role_name` - The Ram role name.
