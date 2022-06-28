---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_load_balancer"
sidebar_current: "docs-alicloud-resource-alb-load-balancer"
description: |-
  Provides a Alicloud ALB Load Balancer resource.
---

# alicloud\_alb\_load\_balancer

Provides a ALB Load Balancer resource.

For information about ALB Load Balancer and how to use it, see [What is Load Balancer](https://www.alibabacloud.com/help/doc-detail/197341.htm).

-> **NOTE:** Available in v1.132.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "example-name"
}

data "alicloud_alb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default_1" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_alb_zones.default.zones.0.id
}
resource "alicloud_vswitch" "vswitch_1" {
  count        = length(data.alicloud_vswitches.default_1.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  zone_id      = data.alicloud_alb_zones.default.zones.0.id
  vswitch_name = var.name
}

data "alicloud_vswitches" "default_2" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_alb_zones.default.zones.1.id
}
resource "alicloud_vswitch" "vswitch_2" {
  count        = length(data.alicloud_vswitches.default_2.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 4)
  zone_id      = data.alicloud_alb_zones.default.zones.1.id
  vswitch_name = var.name
}
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_alb_load_balancer" "default" {
  vpc_id                 = data.alicloud_vpcs.default.ids.0
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
    vswitch_id = length(data.alicloud_vswitches.default_1.ids) > 0 ? data.alicloud_vswitches.default_1.ids[0] : concat(alicloud_vswitch.vswitch_1.*.id, [""])[0]
    zone_id    = data.alicloud_alb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = length(data.alicloud_vswitches.default_2.ids) > 0 ? data.alicloud_vswitches.default_2.ids[0] : concat(alicloud_vswitch.vswitch_2.*.id, [""])[0]
    zone_id    = data.alicloud_alb_zones.default.zones.1.id
  }
  modification_protection_config {
    status = "NonProtection"
  }
}

```

## Argument Reference

The following arguments are supported:

* `access_log_config` - (Optional, Optional) The Access Logging Configuration Structure.
* `address_allocated_mode` - (Optional, ForceNew) The method in which IP addresses are assigned. Valid values: `Fixed` and `Dynamic`. Default value: `Dynamic`.
  *`Fixed`: The ALB instance uses a fixed IP address. 
  *`Dynamic`: An IP address is dynamically assigned to each zone of the ALB instance.
* `address_type` - (Required, ForceNew) The type of IP address that the ALB instance uses to provide services. Valid values: `Intranet`, `Internet`.
* `deletion_protection_enabled` - (Optional) The deletion protection enabled. Valid values: `true` and `false`. Default value: `false`.
* `dry_run` - (Optional) Specifies whether to precheck the API request. Valid values: `true` and `false`.
* `load_balancer_billing_config` - (Required, ForceNew) The configuration of the billing method.
* `load_balancer_edition` - (Required) The edition of the ALB instance. Different editions have different limits and billing methods.  Valid values: `Basic` and `Standard`.
* `load_balancer_name` - (Required) The name of the resource.
* `modification_protection_config` - (Optional, Computed) Modify the Protection Configuration.
* `resource_group_id` - (Optional) The ID of the resource group.
* `vpc_id` - (Required, ForceNew) The ID of the virtual private cloud (VPC) where the ALB instance is deployed.
* `zone_mappings` - (Required, ForceNew) The zones and vSwitches. You must specify at least two zones.
* `tags` - (Optional) A mapping of tags to assign to the resource.

#### Block load_balancer_billing_config

The load_balancer_billing_config supports the following:

* `pay_type` - (Required) The billing method of the ALB instance. Valid value: `PayAsYouGo`.

#### Block zone_mappings

The zone_mappings supports the following: 

* `vswitch_id` - (Required) The ID of the vSwitch that corresponds to the zone. Each zone can use only one vSwitch and subnet.
* `zone_id` - (Required) The ID of the zone to which the ALB instance belongs.

#### Block modification_protection_config

The modification_protection_config supports the following: 

* `status` - (Optional, Computed, Available in v1.132.0+) Specifies whether to enable the configuration read-only mode for the ALB instance. Valid values: `NonProtection` and `ConsoleProtection`.
  * `NonProtection` - disables the configuration read-only mode. After you disable the configuration read-only mode, you cannot set the ModificationProtectionReason parameter. If the parameter is set, the value is cleared.
  * `ConsoleProtection` - enables the configuration read-only mode. After you enable the configuration read-only mode, you can set the ModificationProtectionReason parameter.
* `reason` - (Optional, Computed,Available in v1.132.0+) The reason for modification protection. This parameter must be 2 to 128 characters in length, and can contain letters, digits, periods, underscores, and hyphens. The reason must start with a letter. This parameter is required only if `ModificationProtectionStatus` is set to `ConsoleProtection`.

#### Block access_log_config

The access_log_config supports the following: 

* `log_project` - (Optional) The log service that access logs are shipped to.
* `log_store` - (Optional) The log service that access logs are shipped to.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Load Balancer.
* `status` - The load balancer status. Valid values: `Active`, `Configuring`, `CreateFailed`
* `dns_name` - The domain name of the ALB instance. **NOTE:** Available in v1.158.0+.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Load Balancer.
* `delete` - (Defaults to 2 mins) Used when delete the Load Balancer.
* `update` - (Defaults to 2 mins) Used when update the Load Balancer.

## Import

ALB Load Balancer can be imported using the id, e.g.

```
$ terraform import alicloud_alb_load_balancer.example <id>
```
