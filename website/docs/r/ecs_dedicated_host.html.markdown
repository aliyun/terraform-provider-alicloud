---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_dedicated_host"
sidebar_current: "docs-alicloud-resource-ecs-dedicated-host"
description: |-
  Provides a Alibaba Cloud ecs dedicated host resource.
---

# alicloud\_ecs\_dedicated\_host

This resouce used to create a dedicated host and store its initial version. For information about Aliecs Dedicated Host and how to use it, see [What is Resource Aliecs Dedicated Host](https://www.alibabacloud.com/help/doc-detail/134238.htm).

-> **NOTE:** Available in 1.91.0+.

## Example Usage
Basic Usage

```
resource "alicloud_ecs_dedicated_host" "default" {
  dedicated_host_type = "ddh.g5"
  tags = {
    Create = "Terraform",
    For    = "DDH",
  }
  description         = "From_Terraform"
  dedicated_host_name = "dedicated_host_name"
}
```

Create Prepaid DDH

```
resource "alicloud_ecs_dedicated_host" "default" {
  dedicated_host_type = "ddh.g5"
  tags = {
    Create = "Terraform",
    For    = "DDH",
  }
  description = "From_Terraform"
  dedicated_host_name = "dedicated_host_name"
  payment_type = "PrePaid"
  expired_time = 1
  sale_cycle = "Month"
}
```
### Deleting alicloud_ecs_dedicated_host or removing it from your configuration

The alicloud_ecs_dedicated_host resource allows you to manage payment_type = "PrePaid" dedicated host, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration
will remove it from your state file and management, but will not destroy the Dedicated Host.
You can resume managing the subscription dedicated host via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:

* `action_on_maintenance` - (Optional) The policy used to migrate the instances from the dedicated host when the dedicated host fails or needs to be repaired online. Valid values: `Migrate`, `Stop`.
* `auto_placement` - (Optional, Computed) Specifies whether to add the dedicated host to the resource pool for automatic deployment. If you do not specify the DedicatedHostId parameter when you create an instance on a dedicated host, Alibaba Cloud automatically selects a dedicated host from the resource pool to host the instance. Valid values: `on`, `off`. Default: `on`.
* `auto_release_time` - (Optional, Computed) The automatic release time of the dedicated host. Specify the time in the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time must be in UTC+0.
* `auto_renew` - (Optional) Specifies whether to automatically renew the subscription dedicated host.
* `auto_renew_period` - (Optional) The auto-renewal period of the dedicated host. Unit: months. Valid values: `1`, `2`, `3`, `6`, and `12`. takes effect and is required only when the AutoRenew parameter is set to true.
* `dedicated_host_name` - (Optional, Computed) The name of the dedicated host. The name must be 2 to 128 characters in length. It must start with a letter but cannot start with http:// or https://. It can contain letters, digits, colons (:), underscores (_), and hyphens (-).
* `dedicated_host_type` - (Required, ForceNew, Computed) The type of the dedicated host. You can call the [DescribeDedicatedHostTypes](https://www.alibabacloud.com/help/doc-detail/134240.htm) operation to obtain the most recent list of dedicated host types.
* `description` - (Optional, Computed) The description of the dedicated host. The description must be 2 to 256 characters in length and cannot start with http:// or https://.
* `detail_fee` - (Optional) Specifies whether to return the billing details of the order when the billing method is changed from subscription to pay-as-you-go. Default: `false`.
* `dry_run` - (Optional) Specifies whether to only validate the request. Default: `false`.
* `expired_time` - (Optional, Computed) The subscription period of the dedicated host. The Period parameter takes effect and is required only when the ChargeType parameter is set to PrePaid.
* `network_attributes` - (Optional) dedicated host network parameters. contains the following attributes:
  * `slb_udp_timeout` - The timeout period for a UDP session between Server Load Balancer (SLB) and the dedicated host. Unit: seconds. Valid values: 15 to 310.
  * `udp_timeout` - The timeout period for a UDP session between a user and an Alibaba Cloud service on the dedicated host. Unit: seconds. Valid values: 15 to 310.
* `payment_type` - (Optional, Computed) The billing method of the dedicated host. Valid values: `PrePaid`, `PostPaid`. Default: `PostPaid`.
* `resource_group_id` - (Optional, Computed) The ID of the resource group to which the dedicated host belongs.
* `sale_cycle` - (Optional, Computed) The unit of the subscription period of the dedicated host.
* `zone_id` - (Optional, ForceNew, Computed) The zone ID of the dedicated host. This parameter is empty by default. If you do not specify this parameter, the system automatically selects a zone.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `cpu_over_commit_ratio` - (Optional, Available in 1.123.1+) CPU oversold ratio. Only custom specifications g6s, c6s, r6s support setting the CPU oversold ratio.
* `dedicated_host_cluster_id` - (Optional, Available in 1.123.1+) The dedicated host cluster ID to which the dedicated host belongs.
* `min_quantity` - (Optional, Available in 1.123.1+) Specify the minimum purchase quantity of a dedicated host.

## Attributes Reference

* `id` - The ID of the dedicated host.
* `status` - The status of the dedicated host.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 11 mins) Used when create the dedicated host.
* `delete` - (Defaults to 1 mins) Used when delete the dedicated host.
* `update` - (Defaults to 11 mins) Used when update the dedicated host.

## Import

Ecs dedicated host can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_dedicated_host.default dh-2zedmxxxx
```
