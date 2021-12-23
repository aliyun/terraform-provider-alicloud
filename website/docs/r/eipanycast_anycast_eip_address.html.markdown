---
subcategory: "Anycast Elastic IP Address (Eipanycast)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eipanycast_anycast_eip_address"
sidebar_current: "docs-alicloud-resource-eipanycast-anycast-eip-address"
description: |-
  Provides a Alicloud Anycast Eip Address resource.
---

# alicloud\_eipanycast\_anycast\_eip\_address

Provides a Eipanycast Anycast Eip Address resource.

For information about Eipanycast Anycast Eip Address and how to use it, see [What is Anycast Eip Address](https://help.aliyun.com/document_detail/169284.html).

-> **NOTE:** Available in v1.113.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_eipanycast_anycast_eip_address" "example" {
  service_location = "international"
}

```

## Argument Reference

The following arguments are supported:

* `anycast_eip_address_name` - (Optional) Anycast EIP instance name.
* `bandwidth` - (Optional)  The peak bandwidth of the Anycast EIP instance, in Mbps. It can not be changed when the internet_charge_type is `PayByBandwidth` and the default value is 200.
* `description` - (Optional) Anycast EIP instance description.
* `internet_charge_type` - (Optional, ForceNew) The billing method of Anycast EIP instance. `PayByBandwidth`: refers to the method of billing based on traffic. Valid value: `PayByBandwidth`.
* `payment_type` - (Optional, ForceNew) The payment model of Anycast EIP instance. `PayAsYouGo`: Refers to the post-paid mode. Valid value: `PayAsYouGo`. Default value is `PayAsYouGo`.
* `service_location` - (Required, ForceNew)  Anycast EIP instance access area. `international`: Refers to areas outside of Mainland China.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Anycast Eip Address.
* `status` - The IP status.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 11 mins) Used when create the Anycast Eip Address.

## Import

Eipanycast Anycast Eip Address can be imported using the id, e.g.

```
$ terraform import alicloud_eipanycast_anycast_eip_address.example <id>
```
