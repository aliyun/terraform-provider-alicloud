---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_network_interfaces"
sidebar_current: "docs-alicloud-datasource-ecs-network-interfaces"
description: |-
  Provides a list of Ecs Network Interfaces to the user.
---

# alicloud\_ecs\_network\_interfaces

This data source provides the Ecs Network Interfaces of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.123.1+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecs_network_interfaces" "example" {
  ids        = ["eni-abcd1234"]
  name_regex = "tf-testAcc"
}

output "first_ecs_network_interface_id" {
  value = data.alicloud_ecs_network_interfaces.example.interfaces.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Network Interface IDs.
* `instance_id` - (Optional, ForceNew) The instance id.
* `name` - (Optional, ForceNew, Deprecated in v1.123.1+) Field `name` has been deprecated from provider version 1.123.1. New field `network_interface_name` instead
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Network Interface name.
* `network_interface_name` - (Optional, ForceNew) The network interface name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `primary_ip_address` - (Optional, ForceNew) The primary private IP address of the ENI.
* `private_ip` - (Optional, ForceNew, Deprecated in v1.123.1+) Field `private_ip` has been deprecated from provider version 1.123.1. New field `primary_ip_address` instead
* `resource_group_id` - (Optional, ForceNew) The resource group id.
* `security_group_id` - (Optional, ForceNew) The security group id.
* `service_managed` - (Optional, ForceNew) Whether the user of the elastic network card is a cloud product or a virtual vendor.
* `status` - (Optional, ForceNew) The status of ENI. Valid Values: `Attaching`, `Available`, `CreateFailed`, `Creating`, `Deleting`, `Detaching`, `InUse`, `Linked`, `Linking`, `Unlinking`.
* `type` - (Optional, ForceNew) The type of ENI. Valid Values: `Primary`, `Secondary`.
* `vpc_id` - (Optional, ForceNew) The vpc id.
* `vswitch_id` - (Optional, ForceNew) The vswitch id.
* `tags` - (Optional) A map of tags assigned to ENIs.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Network Interface names.
* `interfaces` - A list of Ecs Network Interfaces. Each element contains the following attributes:
    * `creation_time` - The creation time.
    * `description` - The description of the ENI.
    * `id` - The ID of the Network Interface.
    * `instance_id` - The instance id.
    * `mac` - The MAC address of the ENI.
    * `name` - The network interface name.
    * `network_interface_id` - The network interface id.
    * `network_interface_name` - The network interface name.
    * `network_interface_traffic_mode` - The communication mode of the elastic network card.
    * `owner_id` - The ID of the account to which the ENIC belongs.
    * `primary_ip_address` - The primary private IP address of the ENI. 
    * `private_ip` - The primary private IP address of the ENI.
    * `private_ip_address` - A list of secondary private IP address that is assigned to the ENI.
    * `private_ips` - A list of secondary private IP address that is assigned to the ENI.
    * `queue_number` - Number of network card queues.
    * `resource_group_id` - The resource group id.
    * `security_group_ids` - The security group ids.
    * `security_groups` - The security groups.
    * `service_managed` - Whether the user of the elastic network card is a cloud product or a virtual vendor.
    * `service_id` - The service id.
    * `status` - The status of the ENI.
    * `tags` - The tags.
        * `tag_key` - The tagKey.
        * `tag_value` - The tagValue.
    * `type` - The type of the ENI.
    * `vpc_id` - The Vpc Id.
    * `vswitch_id` - The vswitch id.
    * `zone_id` - The zone id.
    * `associated_public_ip` - The EIP associated with the secondary private IP address of the ENI.  **NOTE:** Available in v1.163.0+.
      * `public_ip_address` - The EIP of the ENI.
