---
subcategory: "Elastic IP Address (EIP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eip_address"
description: |-
  Provides a Alicloud EIP Address resource.
---

# alicloud_eip_address

Provides a EIP Address resource. -> **NOTE:** BGP (Multi-ISP) lines are supported in all regions. BGP (Multi-ISP) Pro lines are supported only in the China (Hong Kong) region.

-> **NOTE:** The resource only supports to create `PayAsYouGo PayByTraffic`  or `Subscription PayByBandwidth` elastic IP for international account. Otherwise, you will happened error `COMMODITY.INVALID_COMPONENT`.
Your account is international if you can use it to login in [International Web Console](https://account.alibabacloud.com/login/login.htm).

For information about EIP Address and how to use it, see [What is Address](https://www.alibabacloud.com/help/en/doc-detail/36016.htm).

-> **NOTE:** Available since v1.126.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}


resource "alicloud_eip_address" "default" {
  description               = var.name
  isp                       = "BGP"
  address_name              = var.name
  netmode                   = "public"
  bandwidth                 = "1"
  security_protection_types = ["AntiDDoS_Enhanced"]
  payment_type              = "PayAsYouGo"
}
```

### Deleting `alicloud_eip_address` or removing it from your configuration

The `alicloud_eip_address` resource allows you to manage  `payment_type = "Subscription"`  instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Instance.
You can resume managing the subscription instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:
* `activity_id` - (Optional) Special activity ID. This parameter is not required.
* `address_name` - (Optional) The name of the EIP instance. This name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin or end with a hyphen, and must not begin with http:// or https://.
* `allocation_id` - (Optional, Available since v1.224.0) The ID of the EIP instance. If you specify the instance ID of An EIP that has already been applied for, the IpAddress of that instance will be reused. Only one of the IpAddress and InstanceId parameters needs to be specified. If neither parameter is specified, the system will randomly apply for an EIP.
* `auto_pay` - (Optional) Whether to pay automatically. Valid values: `true` and `false`. Default value: `true`. When `auto_pay` is `true`, The order will be automatically paid. When `auto_pay` is `false`, The order needs to go to the order center to complete the payment. **NOTE:** When `payment_type` is `Subscription`, this parameter is valid.
* `bandwidth` - (Optional, Computed) The maximum bandwidth of the EIP. Valid values: `1` to `200`. Unit: Mbit/s. Default value: `5`.
* `deletion_protection` - (Optional, Computed, Available since v1.207.0) Whether the delete protection function is turned on.
  - **true**: enabled.
  - **false**: not enabled.
* `description` - (Optional) The description of the EIP.
* `high_definition_monitor_log_status` - (Optional, ForceNew, Computed) Whether the second-level monitoring is enabled for the EIP.
  - **OFF**: not enabled.
  - **ON**: enabled.
* `internet_charge_type` - (Optional, ForceNew, Computed) Renewal Payment type.
  - **PayByBandwidth**: billed by fixed bandwidth.
  - **PayByTraffic**: Billing by traffic.
* `ip_address` - (Optional, ForceNew, Computed) The IP address of the EIP.
* `isp` - (Optional, ForceNew, Computed) The line type. You can set this parameter only when you create a `PayAsYouGo` EIP. Valid values: 
  - `BGP`: BGP (Multi-ISP) lines.Up to 89 high-quality BGP lines are available worldwide. Direct connections with multiple Internet Service Providers (ISPs), including Telecom, Unicom, Mobile, Railcom, Netcom, CERNET, China Broadcast Network, Dr. Peng, and Founder, can be established in all regions in mainland China.
  - `BGP_PRO`: BGP (Multi-ISP) Pro lines optimize data transmission to mainland China and improve connection quality for international services. Compared with BGP (Multi-ISP), when BGP (Multi-ISP) Pro provides services to clients in mainland China (excluding data centers), cross-border connections are established without using international ISP services. This reduces network latency.
  - `ChinaTelecom`: China Telecom.
  - `ChinaUnicom`: China Unicom.
  - `ChinaMobile`: China Mobile.
  - `ChinaTelecom_L2`: China Telecom L2.
  - `ChinaUnicom_L2`: China Unicom L2.
  - `ChinaMobile_L2`: China Mobile L2.
  - `BGP_FinanceCloud`: If your services are deployed in China East 1 Finance, this parameter is required and you must set the value to `BGP_FinanceCloud`.
  - `BGP_International`: BGP_International.
-> **NOTE:** From version 1.203.0, `isp` can be set to `ChinaTelecom`, `ChinaUnicom`, `ChinaMobile`, `ChinaTelecom_L2`, `ChinaUnicom_L2`, `ChinaMobile_L2`, `BGP_FinanceCloud`, `BGP_International`.
* `log_project` - (Optional) The Name of the logging service LogProject. Current parameter is required when configuring high precision second-by-second monitoring for EIP.
* `log_store` - (Optional) The Name of the logging service LogStore. Current parameter is required when configuring high precision second-by-second monitoring for EIP.
* `mode` - (Optional, Computed, Available since v1.224.0) Binding mode, value:
  - **NAT** (default):NAT mode (normal mode).
  - **MULTI_BINDED**: indicates the multi-EIP NIC visible mode.
  - **BINDED**: indicates the mode in which the EIP NIC is visible.
* `netmode` - (Optional, ForceNew, Computed) The type of the network. Valid value is `public` (Internet).
* `payment_type` - (Optional, ForceNew, Computed) The billing method of the EIP. Valid values:  `Subscription`, `PayAsYouGo`.
* `period` - (Optional) When the PricingCycle is set to Month, the Period value ranges from 1 to 9.  When the PricingCycle is set to Year, the Period range is 1 to 5.  If the value of the InstanceChargeType parameter is PrePaid, this parameter is required. If the value of the InstanceChargeType parameter is PostPaid, this parameter is not filled in.
* `pricing_cycle` - (Optional, Available since v1.207.0) Value: Month (default): Pay monthly. Year: Pay per Year. This parameter is required when the value of the InstanceChargeType parameter is Subscription(PrePaid). This parameter is optional when the value of the InstanceChargeType parameter is PayAsYouGo(PostPaid).
* `public_ip_address_pool_id` - (Optional, ForceNew) The ID of the IP address pool to which the EIP belongs.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `security_protection_types` - (Optional, ForceNew) Security protection level.
  - When the return is empty, the basic DDoS protection is specified.
  - When **antidos_enhanced** is returned, it indicates DDoS protection (enhanced version).
* `tags` - (Optional, Map) The tag of the resource.
* `zone` - (Optional, ForceNew, Computed, Available since v1.207.0) The zone of the EIP.  This parameter is returned only for whitelist users that are visible to the zone.

The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated since v1.126.0). Field 'name' has been deprecated from provider version 1.126.0. New field 'address_name' instead.
* `instance_charge_type` - (Deprecated since v1.126.0). Field 'instance_charge_type' has been deprecated from provider version 1.126.0. New field 'payment_type' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the EIP was created.
* `status` - The status of the EIP.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 9 mins) Used when create the Address.
* `delete` - (Defaults to 9 mins) Used when delete the Address.
* `update` - (Defaults to 9 mins) Used when update the Address.

## Import

EIP Address can be imported using the id, e.g.

```shell
$ terraform import alicloud_eip_address.example <id>
```