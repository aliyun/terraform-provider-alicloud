---
subcategory: "Elastic IP Address (EIP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eip_segment_address"
description: |-
  Provides a Alicloud EIP Segment Address resource.
---

# alicloud_eip_segment_address

Provides a EIP Segment Address resource.

For information about EIP Segment Address and how to use it, see [What is Segment Address](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/allocateeipsegmentaddress).

-> **NOTE:** Available since v1.207.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_eip_segment_address&exampleId=baf965c9-33d1-daec-0185-00b89374d20325467b6e&activeTab=example&spm=docs.r.eip_segment_address.0.baf965c933&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "terraform-example"
}


resource "alicloud_eip_segment_address" "default" {
  eip_mask             = "28"
  bandwidth            = "5"
  isp                  = "BGP"
  internet_charge_type = "PayByBandwidth"
  netmode              = "public"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_eip_segment_address&spm=docs.r.eip_segment_address.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `bandwidth` - (Optional, Available since v1.207.0) The maximum bandwidth of the contiguous EIP group. Unit: Mbit/s.
  - Valid values when `InstanceChargeType` is set to `PostPaid` and `InternetChargeType` is set to `PayByBandwidth`: `1` to `500`.****
  - Valid values when `InstanceChargeType` is set to `PostPaid` and `InternetChargeType` is set to `PayByTraffic`: `1` to `200`.****
  - Valid values when `InstanceChargeType` is set to `PrePaid`: `1` to `1000`.****

  Default value: `5`. Unit: Mbit/s.

* `eip_mask` - (Required, Available since v1.207.0) The subnet mask of the contiguous EIP group. Valid values:
  - `28`: applies for 16 contiguous EIPs in each call.
  - `27`: applies for 32 contiguous EIPs in each call.
  - `26`: applies for 64 contiguous EIPs in each call.
  - `25`: applies for 128 contiguous EIPs in each call.
  - `24`: applies for 256 contiguous EIPs in each call.

-> **NOTE:**   Some IP address are reserved for specific purposes. Therefore, the actual number of the contiguous EIPs may be one, three, or four less than the expected number.

* `internet_charge_type` - (Optional, Available since v1.207.0) The metering method of the contiguous EIP group. Valid values:
  - `PayByBandwidth` (default)
  - `PayByTraffic`

* `isp` - (Optional, Available since v1.207.0) The line type. Valid values:
  - `BGP` (default): BGP (Multi-ISP) line The BGP (Multi-ISP) line is supported in all regions.
  - `BGP_PRO`: BGP (Multi-ISP) Pro line BGP (Multi-ISP) Pro line is supported only in the China (Hong Kong), Singapore, Japan (Tokyo), Malaysia (Kuala Lumpur), Philippines (Manila), Indonesia (Jakarta), and Thailand (Bangkok) regions.

  For more information about the BGP (Multi-ISP) line and BGP (Multi-ISP) Pro line, see [EIP line types](https://www.alibabacloud.com/help/en/doc-detail/32321.html).

  If you are allowed to use single-ISP bandwidth, you can also use one of the following values:
  - `ChinaTelecom`
  - `ChinaUnicom`
  - `ChinaMobile`
  - `ChinaTelecom_L2`
  - `ChinaUnicom_L2`
  - `ChinaMobile_L2`

  If your services are deployed in China East 1 Finance, this parameter is required and you must set the parameter to `BGP_FinanceCloud`.

* `netmode` - (Optional, Available since v1.207.0) The network type. Set the value to `public`, which specifies the public network type. 
* `resource_group_id` - (Optional) The resource group ID. 
* `zone` - (Optional, ForceNew, Computed) The zone of the contiguous EIP group. 

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the contiguous Elastic IP address group was created. The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time is displayed in UTC.
* `segment_address_name` - The name of the contiguous Elastic IP address group.
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Segment Address.
* `delete` - (Defaults to 5 mins) Used when delete the Segment Address.

## Import

EIP Segment Address can be imported using the id, e.g.

```shell
$ terraform import alicloud_eip_segment_address.example <id>
```