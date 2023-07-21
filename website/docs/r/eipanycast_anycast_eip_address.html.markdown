---
subcategory: "Eipanycast"
layout: "alicloud"
page_title: "Alicloud: alicloud_eipanycast_anycast_eip_address"
description: |-
  Provides a Alicloud Eipanycast Anycast Eip Address resource.
---

# alicloud_eipanycast_anycast_eip_address

Provides a Eipanycast Anycast Eip Address resource. Anycast Elastic IP Address.

For information about Eipanycast Anycast Eip Address and how to use it, see [What is Anycast Eip Address](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.208.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}


resource "alicloud_eipanycast_anycast_eip_address" "default" {
  anycast_eip_address_name = var.name
  description              = "test_1"
  bandwidth                = 200
  service_location         = "international"
  internet_charge_type     = "PayByTraffic"
  payment_type             = "PayAsYouGo"
  tags {
    tag_key   = "k1"
    tag_value = "v1"
  }
  tags {
    tag_key   = "k2"
    tag_value = "v2"
  }
  tags {
    tag_key   = "k3"
    tag_value = "v3"
  }
}
```

## Argument Reference

The following arguments are supported:
* `anycast_eip_address_name` - (Optional, Available since v1.113.0) Anycast EIP instance name.
* `bandwidth` - (Optional, Computed, Available since v1.113.0)  The peak bandwidth of the Anycast EIP instance, in Mbps.
* `description` - (Optional, Available since v1.113.0) Anycast EIP instance description.
* `internet_charge_type` - (Optional, ForceNew, Available since v1.113.0) The billing method of Anycast EIP instance. "PayByBandwidth": refers to the method of billing based on traffic.
* `payment_type` - (Optional, ForceNew, Computed, Available since v1.113.0) The payment model of Anycast EIP instance. "PayAsYouGo": Refers to the post-paid mode.
* `service_location` - (Required, ForceNew, Available since v1.113.0) Anycast EIP instance access area. "international": Refers to areas outside of Mainland China.
* `tags` - (Optional, Map) List of resource-bound tags.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` -  Anycast EIP instance creation time.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Anycast Eip Address.
* `delete` - (Defaults to 5 mins) Used when delete the Anycast Eip Address.
* `update` - (Defaults to 5 mins) Used when update the Anycast Eip Address.

## Import

Eipanycast Anycast Eip Address can be imported using the id, e.g.

```shell
$ terraform import alicloud_eipanycast_anycast_eip_address.example <id>
```