---
subcategory: "EIP Bandwidth Plan (CBWP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_common_bandwidth_package"
description: |-
  Provides a Alicloud CBWP Common Bandwidth Package resource.
---

# alicloud_common_bandwidth_package

Provides a CBWP Common Bandwidth Package resource. -> **NOTE:** Terraform will auto build common bandwidth package instance while it uses `alicloud_common_bandwidth_package` to build a common bandwidth package resource.

For information about common bandwidth package billing methods, see [Common Bandwidth Package Billing Methods](https://www.alibabacloud.com/help/doc-detail/67459.html).

For information about CBWP Common Bandwidth Package and how to use it, see [What is Common Bandwidth Package](https://www.alibabacloud.com/help/en/eip-bandwidth-plan).

-> **NOTE:** Available since v1.23.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_common_bandwidth_package" "default" {
  bandwidth_package_name    = var.name
  description               = var.name
  isp                       = "BGP"
  bandwidth                 = "1000"
  ratio                     = 100
  internet_charge_type      = "PayByBandwidth"
  resource_group_id         = data.alicloud_resource_manager_resource_groups.default.ids.0
  security_protection_types = ["AntiDDoS_Enhanced"]
}
```

### Deleting `alicloud_common_bandwidth_package` or removing it from your configuration

The `alicloud_common_bandwidth_package` resource allows you to manage  `internet_charge_type = "PayBy95"`  instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Instance.
You can resume managing the subscription instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:
* `bandwidth` - (Required) The peak bandwidth of the shared bandwidth. Unit: Mbps. 
  Valid values: [2, 20000] for China-Site account; [1, 5000] for International-Site account. See [Account Guide](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/guides/getting-account) details.
* `bandwidth_package_name` - (Optional, Available since v1.120.0) The name of the Internet Shared Bandwidth instance.
* `deletion_protection` - (Optional, Available since v1.124.4) Whether enable the deletion protection or not. Default value: false.
  - **true**: Enable deletion protection.
  - **false**: Disable deletion protection.
* `description` - (Optional) The description of the shared bandwidth.
* `force` - (Optional) Whether to forcibly delete an Internet Shared Bandwidth instance. Value:
  - **false** (default): only the internet shared bandwidth that does not contain the EIP is deleted.
  - **true**: removes all EIPs from the internet shared bandwidth instance and deletes the internet shared bandwidth.
* `internet_charge_type` - (Optional, ForceNew) The billing method of the common bandwidth package. Valid values are `PayByBandwidth` and `PayBy95` and `PayByTraffic`, `PayByDominantTraffic`. `PayBy95` is pay by classic 95th percentile pricing. International-Site Account doesn't support `PayByBandwidth` and `PayBy95`. Default to `PayByTraffic`. **NOTE:** From 1.176.0+, `PayByDominantTraffic` is available.
* `isp` - (Optional, ForceNew, Computed, Available since v1.90.1) The type of the Internet Service Provider. Valid values: `BGP`, `BGP_PRO`, `ChinaTelecom`, `ChinaUnicom`, `ChinaMobile`, `ChinaTelecom_L2`, `ChinaUnicom_L2`, `ChinaMobile_L2` and `BGP_FinanceCloud`. Default to `BGP`. **NOTE:** From version 1.203.0, isp can be set to `ChinaTelecom`, `ChinaUnicom`, `ChinaMobile`, `ChinaTelecom_L2`, `ChinaUnicom_L2`, `ChinaMobile_L2`, `BGP_FinanceCloud`, `BGP_International`.
* `ratio` - (Optional, ForceNew, Computed, Available since v1.55.3) Ratio of the common bandwidth package. It is valid when `internet_charge_type` is `PayBy95`. Default to 100. Valid values: [10-100].
* `resource_group_id` - (Optional, Computed, Available since v1.115.0) The Id of resource group which the common bandwidth package belongs.
* `security_protection_types` - (Optional, ForceNew, Available since v1.184.0) The edition of Anti-DDoS. If you do not set this parameter, Anti-DDoS Origin Basic is used. If you set the value to AntiDDoS_Enhanced, Anti-DDoS Pro(Premium) is used. It is valid when `internet_charge_type` is `PayBy95`.
* `tags` - (Optional, Map, Available since v1.207.0) The tag of the resource.
* `zone` - (Optional, Available since v1.120.0) The available area of the shared bandwidth.

The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated since v1.120.0). Field 'name' has been deprecated from provider version 1.120.0. New field 'bandwidth_package_name' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The create time.
* `payment_type` - The billing type of the Internet Shared Bandwidth instance. Valid values: `PayAsYouGo`, `Subscription`.
* `status` - The status of the Internet Shared Bandwidth instance. Default value: **Available**.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 10 mins) Used when create the Common Bandwidth Package.
* `delete` - (Defaults to 10 mins) Used when delete the Common Bandwidth Package.
* `update` - (Defaults to 10 mins) Used when update the Common Bandwidth Package.

## Import

CBWP Common Bandwidth Package can be imported using the id, e.g.

```shell
$ terraform import alicloud_common_bandwidth_package.example <id>
```