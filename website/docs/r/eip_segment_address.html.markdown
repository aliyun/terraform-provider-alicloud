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

## Argument Reference

The following arguments are supported:
* `bandwidth` - (Optional) The peak bandwidth of the EIP. Unit: Mbps. When the value of instancargetype is PostPaid and the value of InternetChargeType is PayByBandwidth, the range of Bandwidth is 1 to 500. If the value of instancargetype is PostPaid and the value of InternetChargeType is PayByTraffic, the range of Bandwidth is 1 to 200. When instancargetype is set to PrePaid, the range of Bandwidth is 1 to 1000. The default value is 5 Mbps.
* `eip_mask` - (Required) Mask of consecutive EIPs. Value:28: For a single call, the system will allocate 16 consecutive EIPs.27: For a single call, the system will allocate 32 consecutive EIPs.26: For a single call, the system will allocate 64 consecutive EIPs.25: For a single call, the system will allocate 128 consecutive EIPs.24: For a single call, the system will allocate 256 consecutive EIPs.
* `internet_charge_type` - (Optional) Continuous EIP billing method, valid values:
  - **PayByBandwidth** (default): Billing based on fixed bandwidth.
  - **PayByTraffic**: Billing by usage flow.
* `isp` - (Optional) Line type. Valid values:
  - **BGP** (default):BGP (multi-line) line. BGP (multi-line) EIP is supported in all regions.
  - **BGP_PRO** :BGP (multi-line)_boutique line. Currently, only Hong Kong, Singapore, Japan (Tokyo), Malaysia (Kuala Lumpur), the Philippines (Manila), Indonesia (Jakarta), and Thailand (Bangkok) regions support BGP (multi-line)_boutique route EIP.
For more information about BGP (multi-line) lines and BGP (multi-line) premium lines, see EIP line types.
If you are a whitelist user with single-line bandwidth, you can also select the following types:
  - **ChinaTelecom** : China Telecom
  - **ChinaUnicom** : China Unicom
  - **ChinaMobile** : China Mobile
  - **ChinaTelecom_L2** : China Telecom L2
  - **ChinaUnicom_L2** : China Unicom L2
  - **ChinaMobile_L2** : China Mobile L2
If you are a user of Hangzhou Financial Cloud, this field is required. The value is `BGP_FinanceCloud`.
* `netmode` - (Optional) The network type. Set the value to **public**.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the contiguous Elastic IP address group was created. The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time is displayed in UTC.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Segment Address.
* `delete` - (Defaults to 5 mins) Used when delete the Segment Address.

## Import

EIP Segment Address can be imported using the id, e.g.

```shell
$ terraform import alicloud_eip_segment_address.example <id>
```