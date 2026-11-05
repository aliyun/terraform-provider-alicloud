---
subcategory: "Elastic IP Address (EIP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eip"
sidebar_current: "docs-alicloud-resource-eip"
description: |-
  Provides a ECS EIP resource.
---

# alicloud\_eip

Provides an elastic IP resource.

-> **DEPRECATED:**  This resource  has been deprecated from version `1.126.0`. Please use new resource [alicloud_eip_address](https://www.terraform.io/docs/providers/alicloud/r/eip_address).

-> **NOTE:** The resource only supports to create `PostPaid PayByTraffic`  or `PrePaid PayByBandwidth` elastic IP for international account. Otherwise, you will happened error `COMMODITY.INVALID_COMPONENT`.
Your account is international if you can use it to login in [International Web Console](https://account.alibabacloud.com/login/login.htm).

-> **NOTE:** From version 1.10.1, this resource supports creating "PrePaid" EIP. In addition, it supports setting EIP name and description.

## Example Usage

```terraform
# Create a new EIP.
resource "alicloud_eip" "example" {
  bandwidth            = "10"
  internet_charge_type = "PayByBandwidth"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_eip&spm=docs.r.eip.example&intl_lang=EN_US)

## Module Support

You can use the existing [eip module](https://registry.terraform.io/modules/terraform-alicloud-modules/eip/alicloud) 
to create several EIP instances and associate them with other resources one-click, like ECS instances, SLB, Nat Gateway and so on.

## Argument Reference

The following arguments are supported:
* `activity_id` - (Optional) The promotion code. This parameter is not required.
* `address_name` - (Optional) The EIP name.

  The name must be 1 to 128 characters in length and start with a letter, and can contain letters, digits, periods (.), underscores (\_), and hyphens (-).

* `allocation_id` - (Optional, ForceNew, Computed, Available since v1.225.1) The ID of the EIP instance.
* `auto_pay` - (Optional) Specifies whether to enable automatic payment. Valid values:
  - `false` (default): The automatic payment is disabled. If you select this option, you must go to the Order Center to complete the payment after an order is generated.
  - `true`: The automatic payment is enabled. Payments are automatically complete after an order is generated.

* `bandwidth` - (Optional, Computed) The maximum bandwidth of the specified EIP. Unit: Mbit/s.
* `deletion_protection` - (Optional, Computed) Specifies whether to enable deletion protection. Valid values:
  - `true`: yes
  - `false`: no

* `description` - (Optional) The description of the EIP.

* `high_definition_monitor_log_status` - (Optional, ForceNew, Computed) The status of fine-grained monitoring. Valid values:
  - `ON`
  - `OFF`

* `internet_charge_type` - (Optional, ForceNew, Computed) The metering method of the EIP. Valid values:
  - `PayByBandwidth` (default): pay-by-bandwidth.
  - `PayByTraffic`: pay-by-data-transfer.

* `ip_address` - (Optional, ForceNew, Computed) The IP address of the EIP. Supports a maximum of 50 EIPs.
* `isp` - (Optional, ForceNew, Computed) The line type. Valid values:
  - `BGP` (default): BGP (Multi-ISP) line
  - `BGP_PRO`: BGP (Multi-ISP) Pro line

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

* `period` - (Optional) Duration of purchase. When the value of `pricing_cycle` is `Month`, the value range of `period` is `1` to `9`. When the value of `pricing_cycle` is `Year`, the value range of `period` is `1` to `5`. If the value of the `payment_type` parameter is `Subscription`, this parameter is required. If the value of the `payment_type` parameter is `PayAsYouGo`, this parameter is left blank.
* `pricing_cycle` - (Optional) The billing cycle of the subscription EIP. Valid values:
  - `Month` (default)
  - `Year`

* `public_ip_address_pool_id` - (Optional, ForceNew) The ID of the IP address pool. The EIP is allocated from the IP address pool.
* `resource_group_id` - (Optional, Computed) The ID of the resource group to which you want to move the resource.
* `security_protection_types` - (Optional, ForceNew) Security protection level.
  - When the return is empty, the basic DDoS protection is specified.
  - When `antidos_enhanced` is returned, it indicates DDoS protection (enhanced version).
* `tags` - (Optional, Map) The tag of the resource
* `zone` - (Optional, ForceNew, Computed) The zone of the EIP.

The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated since v1.126.0). Field 'name' has been deprecated from provider version 1.126.0. New field `address_name` instead.
* `instance_charge_type` - (Deprecated since v1.126.0). Field 'instance_charge_type' has been deprecated from provider version 1.126.0. New field `payment_type` instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the EIP was created.
* `status` - The state of the EIP.

## Import

Elastic IP address can be imported using the id, e.g.

```shell
$ terraform import alicloud_eip_address.example eip-abc12345678
```
