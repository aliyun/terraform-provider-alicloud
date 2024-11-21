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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alb_load_balancer&exampleId=69375c6d-bff4-e697-2696-1baac7f74f3a4c97779e&activeTab=example&spm=docs.r.alb_load_balancer.0.69375c6dbf&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

data "alicloud_alb_zones" "default" {
}

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
  load_balancer_edition  = "Basic"
  address_type           = "Internet"
  vpc_id                 = alicloud_vpc.default.id
  address_allocated_mode = "Fixed"
  resource_group_id      = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  load_balancer_name     = var.name
  load_balancer_billing_config {
    pay_type = "PayAsYouGo"
  }
  modification_protection_config {
    status = "NonProtection"
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default.0.id
    zone_id    = data.alicloud_alb_zones.default.zones.0.id
  }
  zone_mappings {
    vswitch_id = alicloud_vswitch.default.1.id
    zone_id    = data.alicloud_alb_zones.default.zones.1.id
  }
  tags = {
    Created = "TF"
  }
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_edition` - (Required) The edition of the ALB instance. The features and billing rules vary based on the edition of the ALB instance. Valid values: `Basic`, `Standard`, `StandardWithWaf`.
* `address_type` - (Required) The type of the address of the ALB instance. Valid values: `Internet`, `Intranet`.
* `vpc_id` - (Required, ForceNew) The ID of the VPC.
* `address_allocated_mode` - (Optional, ForceNew) The mode in which IP addresses are allocated. Valid values: `Fixed`, `Dynamic`.
* `address_ip_version` - (Optional, ForceNew) The protocol version. Valid values: `IPv4`, `DualStack`.
* `ipv6_address_type` - (Optional, Available since v1.211.2) The address type of the Ipv6 address. Valid values: `Internet`, `Intranet`.
* `bandwidth_package_id` - (Optional, ForceNew, Available since v1.211.2) The ID of the Internet Shared Bandwidth instance that is associated with the Internet-facing ALB instance.
* `resource_group_id` - (Optional) The ID of the resource group.
* `load_balancer_name` - (Optional) The name of the ALB instance.
* `deletion_protection_enabled` - (Optional, Bool) Specifies whether to enable deletion protection. Default value: `false`. Valid values:
  - `true`: Enables deletion protection.
  - `false`: Disables deletion protection.
* `load_balancer_billing_config` - (Required, ForceNew, Set) The billing method of the ALB instance. See [`load_balancer_billing_config`](#load_balancer_billing_config) below.
* `modification_protection_config` - (Optional, Set) The configuration of the read-only mode. See [`modification_protection_config`](#modification_protection_config) below.
* `access_log_config` - (Optional, Set) The configuration of the access log. See [`access_log_config`](#access_log_config) below.
* `zone_mappings` - (Required, Set) The list of zones and vSwitch mappings. You must specify at least two zones. See [`zone_mappings`](#zone_mappings) below.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `dry_run` - (Optional, Bool) Specifies whether to perform a dry run. Default value: `false`. Valid values: `true`, `false`.

### `load_balancer_billing_config`

The load_balancer_billing_config supports the following:

* `pay_type` - (Required, ForceNew) The billing method of the ALB instance. Valid values: `PayAsYouGo`.

### `modification_protection_config`

The modification_protection_config supports the following:

* `status` - (Optional) Specifies whether to enable the configuration read-only mode. Valid values: `ConsoleProtection`, `NonProtection`.
* `reason` - (Optional) The reason for enabling the configuration read-only mode. **NOTE:** `reason` takes effect only if `status` is set to `ConsoleProtection`.

### `access_log_config`

The access_log_config supports the following:

* `log_project` - (Required) The project to which the access log is shipped.
* `log_store` - (Required) The Logstore to which the access log is shipped.

### `zone_mappings`

The zone_mappings supports the following:

* `vswitch_id` - (Required) The ID of the VSwitch.
* `zone_id` - (Required) The zone ID of the ALB instance.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Load Balancer.
* `dns_name` - (Available since v1.158.0) The domain name of the ALB instance.
* `status` - The status of the Load Balancer.
* `create_time` - The time when the resource was created.
* `zone_mappings` - The list of zones and vSwitch mappings.
  * `load_balancer_addresses` - The IP address of the ALB instance.
    * `allocation_id` - The ID of the EIP.
    * `eip_type` - The type of the EIP.
    * `address` - IP address. The Public IP Address, and Private IP Address from the Address Type.
    * `ipv6_address` - Ipv6 address.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Load Balancer.
* `update` - (Defaults to 5 mins) Used when update the Load Balancer.
* `delete` - (Defaults to 5 mins) Used when delete the Load Balancer.


## Import

Alb Load Balancer can be imported using the id, e.g.

```shell
$ terraform import alicloud_alb_load_balancer.example <id>
```
