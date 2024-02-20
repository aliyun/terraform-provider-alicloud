---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_load_balancer"
description: |-
  Provides a Alicloud ALB Load Balancer resource.
---

# alicloud_alb_load_balancer

Provides a ALB Load Balancer resource.

For information about ALB Load Balancer and how to use it, see [What is Load Balancer](https://www.alibabacloud.com/help/en/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-createloadbalancer).

-> **NOTE:** Available since v1.132.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
}
data "alicloud_alb_zones" "default" {}
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  count        = 2
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = format("10.4.%d.0/24", count.index + 1)
  zone_id      = data.alicloud_alb_zones.default.zones[count.index].id
  vswitch_name = format("${var.name}_%d", count.index + 1)
}

resource "alicloud_alb_load_balancer" "default" {
  vpc_id                 = alicloud_vpc.default.id
  address_type           = "Internet"
  address_allocated_mode = "Fixed"
  load_balancer_name     = var.name
  load_balancer_edition  = "Basic"
  resource_group_id      = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  load_balancer_billing_config {
    pay_type = "PayAsYouGo"
  }
  tags = {
    Created = "TF"
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default.0.id
    zone_id    = data.alicloud_alb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default.1.id
    zone_id    = data.alicloud_alb_zones.default.zones.1.id
  }
  modification_protection_config {
    status = "NonProtection"
  }
}
```

## Argument Reference

The following arguments are supported:
* `access_log_config` - (Optional, ForceNew, Computed) The Access Logging Configuration Structure. See [`access_log_config`](#access_log_config) below.
* `address_allocated_mode` - (Optional, ForceNew) The method in which IP addresses are assigned. Valid values:  Fixed: The ALB instance uses a fixed IP address. Dynamic (default): An IP address is dynamically assigned to each zone of the ALB instance.
* `address_ip_version` - (Optional, ForceNew, Computed) The protocol version. Value:
  - **IPv4**:IPv4 type.
  - **DualStack**: the dual-stack type.
* `address_type` - (Required) The type of IP address that the SLB instance uses to provide services.
* `bandwidth_package_id` - (Optional, ForceNew, Available since v1.211.2) The ID of the EIP bandwidth plan which is associated with an ALB instance that uses a public IP address.
* `deletion_protection_enabled` - (Optional) Whether to enable delete protection. Value:true: on.false (default): Off.
* `dry_run` - (Optional) Whether to PreCheck only this request, value:true: sends a check request and does not create a resource. Check items include whether required parameters are filled in, request format, and business restrictions. If the check fails, the corresponding error is returned. If the check passes, the error code DryRunOperation is returned.false (default): Sends a normal request, returns the HTTP_2xx status code after the check, and directly performs the operation.
* `ipv6_address_type` - (Optional, Available since v1.211.2) The address type of Ipv6.
* `load_balancer_billing_config` - (Required, ForceNew) The configuration of the billing method. See [`load_balancer_billing_config`](#load_balancer_billing_config) below.
* `load_balancer_edition` - (Required) The edition of the ALB instance.
* `load_balancer_name` - (Optional) The name of the resource.
* `modification_protection_config` - (Optional, Computed) Modify the Protection Configuration. See [`modification_protection_config`](#modification_protection_config) below.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `tags` - (Optional, Computed, Map) The tag of the resource.
* `vpc_id` - (Required, ForceNew) The ID of the virtual private cloud (VPC) where the SLB instance is deployed.
* `zone_mappings` - (Required) The zones and vSwitches. You must specify at least two zones. See [`zone_mappings`](#zone_mappings) below.

### `access_log_config`

The access_log_config supports the following:
* `log_project` - (Required) This Log Storage Project.
* `log_store` - (Required) This Log Storage Method Is Increased.

### `load_balancer_billing_config`

The load_balancer_billing_config supports the following:
* `pay_type` - (Required, ForceNew) Pay Type. Valid values: `PayAsYouGo`. **Note:** provider changes the payment type to `PayAsYouGo`, while the actual parameter on api is `PostPay`.

### `modification_protection_config`

The modification_protection_config supports the following:
* `reason` - (Optional, Computed) Managed Instance.
* `status` - (Optional, Computed) Load Balancing Modify the Protection Status.

### `zone_mappings`

The zone_mappings supports the following:
* `vswitch_id` - (Required) The ID of the vSwitch that corresponds to the zone. Each zone can use only one vSwitch and subnet.
* `zone_id` - (Required) The ID of the zone to which the SLB instance belongs.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `dns_name` - The domain name of the ALB instance. **NOTE:** Available in v1.158.0+.
* `create_time` - The creation time of the resource.
* `status` - Server Load Balancer Instance Status:, indicating that the instance listener will no longer forward traffic.(default).
* `zone_mappings` - The zones and vSwitches. You must specify at least two zones.
  * `load_balancer_addresses` - The SLB Instance Address.
    * `address` - IP Address. The Public IP Address, and Private IP Address from the Address Type.
    * `allocation_id` - The ID of the EIP instance.
    * `eip_type` - The type of the EIP instance.
    * `ipv6_address` - Ipv6 address.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Load Balancer.
* `delete` - (Defaults to 5 mins) Used when delete the Load Balancer.
* `update` - (Defaults to 5 mins) Used when update the Load Balancer.

## Import

Alb Load Balancer can be imported using the id, e.g.

```shell
$ terraform import alicloud_alb_load_balancer.example <id>
```