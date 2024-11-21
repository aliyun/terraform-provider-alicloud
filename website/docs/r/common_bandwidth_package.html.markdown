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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_common_bandwidth_package&exampleId=0717d382-6add-c02a-ca9f-711341f89416202c748f&activeTab=example&spm=docs.r.common_bandwidth_package.0.0717d3826a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
* `bandwidth` - (Required) The maximum bandwidth of the Internet Shared Bandwidth instance. Unit: Mbit/s. Valid values: `1` to `1000`. Default value: `1`. 
* `bandwidth_package_name` - (Optional) The name of the Internet Shared Bandwidth instance. The name must be 2 to 128 characters in length, and can contain letters, digits, underscores (\_), and hyphens (-). The name must start with a letter. 
* `deletion_protection` - (Optional) Specifies whether to enable deletion protection. Valid values:
  - `true`: yes
  - `false`: no

* `description` - (Optional) The description of the Internet Shared Bandwidth instance. The description must be 2 to 256 characters in length and start with a letter. The description cannot start with `http://` or `https://`. 
* `force` - (Optional) Specifies whether to forcefully delete the Internet Shared Bandwidth instance. Valid values:
  - `false` (default): deletes the Internet Shared Bandwidth instance only when no EIPs are associated with the Internet Shared Bandwidth instance.
  - `true`: disassociates all EIPs from the Internet Shared Bandwidth instance and deletes the Internet Shared Bandwidth instance.

* `internet_charge_type` - (Optional, ForceNew, Computed) The billing method of the Internet Shared Bandwidth instance. Set the value to `PayByTraffic`, which specifies the pay-by-data-transfer billing method. 
* `isp` - (Optional, ForceNew, Computed) The line type. Valid values:
  - `BGP` All regions support BGP (Multi-ISP).
  - `BGP_PRO` BGP (Multi-ISP) Pro lines are available in the China (Hong Kong), Singapore, Japan (Tokyo), Philippines (Manila), Malaysia (Kuala Lumpur), Indonesia (Jakarta), and Thailand (Bangkok) regions.

  If you are allowed to use single-ISP bandwidth, you can also use one of the following values:
  - `ChinaTelecom`
  - `ChinaUnicom`
  - `ChinaMobile`
  - `ChinaTelecom_L2`
  - `ChinaUnicom_L2`
  - `ChinaMobile_L2`

  If your services are deployed in China East 1 Finance, this parameter is required and you must set the value to `BGP_FinanceCloud`.

* `ratio` - (Optional, ForceNew, Computed) The percentage of the minimum bandwidth commitment. Set the parameter to `20`.

-> **NOTE:**  This parameter is available only on the Alibaba Cloud China site.

* `resource_group_id` - (Optional, Computed) The ID of the resource group to which you want to move the resource.

-> **NOTE:**   You can use resource groups to facilitate resource grouping and permission management for an Alibaba Cloud. For more information, see [What is resource management?](https://www.alibabacloud.com/help/en/doc-detail/94475.html)

* `security_protection_types` - (Optional, ForceNew) The edition of Anti-DDoS. If you do not set this parameter, Anti-DDoS Origin Basic is used. If you set the value to AntiDDoS_Enhanced, Anti-DDoS Pro(Premium) is used. It is valid when `internet_charge_type` is `PayBy95`.
* `tags` - (Optional, Map) The tag of the resource
* `zone` - (Optional) The zone of the Internet Shared Bandwidth instance. This parameter is required if you create an Internet Shared Bandwidth instance for a cloud box. 

The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated since v1.120.0). Field 'name' has been deprecated from provider version 1.120.0. New field 'bandwidth_package_name' instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `payment_type` - The billing type of the Internet Shared Bandwidth instance. Valid values: `PayAsYouGo`, `Subscription`.
* `create_time` - The creation time.
* `status` - The status of the Internet Shared Bandwidth instance. Default value: `Available`.

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