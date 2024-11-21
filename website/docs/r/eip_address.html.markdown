---
subcategory: "Elastic IP Address (EIP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eip_address"
description: |-
  Provides a Alicloud EIP Address resource.
---

# alicloud_eip_address

Provides a EIP Address resource.

-> **NOTE:** BGP (Multi-ISP) lines are supported in all regions. BGP (Multi-ISP) Pro lines are supported only in the China (Hong Kong) region.

-> **NOTE:** The resource only supports to create `PayAsYouGo PayByTraffic`  or `Subscription PayByBandwidth` elastic IP for international account. Otherwise, you will happened error `COMMODITY.INVALID_COMPONENT`.
Your account is international if you can use it to login in [International Web Console](https://account.alibabacloud.com/login/login.htm).

For information about EIP Address and how to use it, see [What is Address](https://www.alibabacloud.com/help/en/doc-detail/36016.htm).

-> **NOTE:** Available since v1.126.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_eip_address&exampleId=81faa0cb-9221-d24b-3992-fefae457d640b16e82cb&activeTab=example&spm=docs.r.eip_address.0.81faa0cb92&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
* `activity_id` - (Optional) The promotion code. This parameter is not required. 
* `address_name` - (Optional) The EIP name.

  The name must be 1 to 128 characters in length and start with a letter, and can contain letters, digits, periods (.), underscores (\_), and hyphens (-).

-> **NOTE:**   You cannot specify this parameter if you create a subscription EIP.

* `allocation_id` - (Optional, ForceNew, Computed, Available since v1.225.1) The ID of the EIP instance.
* `auto_pay` - (Optional) Specifies whether to enable automatic payment. Valid values:
  - `false` (default): The automatic payment is disabled. If you select this option, you must go to the Order Center to complete the payment after an order is generated.
  - `true`: The automatic payment is enabled. Payments are automatically complete after an order is generated.

  If `payment_type` is set to `Subscription`, this parameter is required. If `payment_type` is set to `PayAsYouGo`, this parameter is not required.

* `bandwidth` - (Optional, Computed) The maximum bandwidth of the specified EIP. Unit: Mbit/s.
  - When `payment_type` is set to `PayAsYouGo` and `internet_charge_type` is set to `PayByBandwidth`, valid values for `bandwidth` are `1` to `500`.
  - When `payment_type` is set to `PayAsYouGo` and `internet_charge_type` is set to `PayByTraffic`, valid values for `bandwidth` are `1` to `200`.
  - When `payment_type` is set to `Subscription`, valid values for `bandwidth` are `1` to `1000`.

  Default value: `5` Mbit /s.

* `deletion_protection` - (Optional, Computed) Specifies whether to enable deletion protection. Valid values:
  - `true`: yes
  - `false`: no

* `description` - (Optional) The description of the EIP.

  The description must be 2 to 256 characters in length. The description must start with a letter but cannot start with `http://` or `https://`.

-> **NOTE:**   You cannot specify this parameter if you create a subscription EIP.

* `high_definition_monitor_log_status` - (Optional, ForceNew, Computed) The status of fine-grained monitoring. Valid values:
  - `ON`
  - `OFF`

* `internet_charge_type` - (Optional, ForceNew, Computed) The metering method of the EIP. Valid values:
  - `PayByBandwidth` (default): pay-by-bandwidth.
  - `PayByTraffic`: pay-by-data-transfer.

  When `payment_type` is set to `Subscription`, you must set `internet_charge_type` to `PayByBandwidth`.

  When `payment_type` is set to `PayAsYouGo`, set `internet_charge_type` to `PayByBandwidth` or `PayByTraffic`.

* `ip_address` - (Optional, ForceNew, Computed) The IP address of the EIP. Supports a maximum of 50 EIPs.
* `isp` - (Optional, ForceNew, Computed) The line type. Valid values:
  - `BGP` (default): BGP (Multi-ISP) line The BGP (Multi-ISP) line is supported in all regions.
  - `BGP_PRO`: BGP (Multi-ISP) Pro line The BGP (Multi-ISP) Pro line is supported in the China (Hong Kong), Singapore, Malaysia (Kuala Lumpur), Philippines (Manila), Indonesia (Jakarta), and Thailand (Bangkok) regions.

  For more information about the BGP (Multi-ISP) line and BGP (Multi-ISP) Pro line, see the "Line types" section of [What is EIP?](https://www.alibabacloud.com/help/en/doc-detail/32321.html)

  If you are allowed to use single-ISP bandwidth, you can also choose one of the following values:
  - `ChinaTelecom`
  - `ChinaUnicom`
  - `ChinaMobile`
  - `ChinaTelecom_L2`
  - `ChinaUnicom_L2`
  - `ChinaMobile_L2`

  If your services are deployed in China East 1 Finance, this parameter is required and you must set the parameter to `BGP_FinanceCloud`.

* `log_project` - (Optional) The name of the Simple Log Service (SLS) project. 
* `log_store` - (Optional) The name of the Logstore. 
* `mode` - (Optional, Computed, Available since v1.225.1) The association mode. Valid values:
  - `NAT` (default): NAT mode
  - `MULTI_BINDED`: multi-EIP-to-ENI mode
  - `BINDED`: cut-network interface controller mode

* `netmode` - (Optional, ForceNew, Computed) The network type. By default, this value is set to `public`, which specifies the public network type. 
* `payment_type` - (Optional, ForceNew, Computed) The billing method of the EIP. Valid values:
  - `Subscription`: subscription
  - `PayAsYouGo` (default): pay-as-you-go

  If `payment_type` is set to `Subscription`, set `internet_charge_type` to `PayByBandwidth`. If `payment_type` is set to `PayAsYouGo`, set `internet_charge_type` to `PayByBandwidth` or `PayByTraffic`.

* `period` - (Optional) Duration of purchase. When the value of `pricing_cycle` is `Month`, the value range of `period` is `1` to `9`. When the value of `pricing_cycle` is `Year`, the value range of `period` is `1` to `5`. If the value of the `payment_type` parameter is `Subscription`, this parameter is required. If the value of the `payment_type` parameter is `PayAsYouGo`, this parameter is left blank.
* `pricing_cycle` - (Optional) The billing cycle of the subscription EIP. Valid values:
  - `Month` (default)
  - `Year`

  If `payment_type` is set to `Subscription`, this parameter is required. If `payment_type` is set to `PayAsYouGo`, this parameter is not required.

* `public_ip_address_pool_id` - (Optional, ForceNew) The ID of the IP address pool. The EIP is allocated from the IP address pool. By default, the IP address pool feature is unavailable. To use the IP address pool, apply for the privilege in the Quota Center console. For more information, see the "Request a quota increase in the Quota Center console" section in [Manage EIP quotas](https://www.alibabacloud.com/help/en/doc-detail/108213.html). 
* `resource_group_id` - (Optional, Computed) The ID of the resource group to which you want to move the resource.

-> **NOTE:**   You can use resource groups to facilitate resource grouping and permission management for an Alibaba Cloud. For more information, see [What is resource management?](https://www.alibabacloud.com/help/en/doc-detail/94475.html)

* `security_protection_types` - (Optional, ForceNew) Security protection level.
  - When the return is empty, the basic DDoS protection is specified.
  - When `antidos_enhanced` is returned, it indicates DDoS protection (enhanced version).
* `tags` - (Optional, Map) The tag of the resource
* `zone` - (Optional, ForceNew, Computed) The zone of the EIP. When the service type of the IP address pool specified by `PublicIpAddressPoolId` is CloudBox, the default value is the zone of the IP address pool. For more information, see [ListPublicIpAddressPools](https://www.alibabacloud.com/help/en/doc-detail/429433.html). 

The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated since v1.126.0). Field 'name' has been deprecated from provider version 1.126.0. New field 'address_name' instead.
* `instance_charge_type` - (Deprecated since v1.126.0). Field 'instance_charge_type' has been deprecated from provider version 1.126.0. New field 'payment_type' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the EIP was created.
* `status` - The state of the EIP. 


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