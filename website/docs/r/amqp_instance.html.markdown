---
subcategory: "RabbitMQ (AMQP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_amqp_instance"
sidebar_current: "docs-alicloud-resource-amqp-instance"
description: |-
  Provides a Alicloud RabbitMQ (AMQP) Instance resource.
---

# alicloud\_amqp\_instance

Provides a RabbitMQ (AMQP) Instance resource.

For information about RabbitMQ (AMQP) Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/doc-detail/101631.htm).

-> **NOTE:** Available in v1.128.0+.

-> **NOTE:** The AMQP Instance is not support to be purchase automatically in the international site.

## Example Usage

Basic Usage

```terraform
resource "alicloud_amqp_instance" "professional" {
  instance_type  = "professional"
  max_tps        = 1000
  queue_capacity = 50
  support_eip    = true
  max_eip_tps    = 128
  payment_type   = "Subscription"
  period         = 1
}

resource "alicloud_amqp_instance" "vip" {
  instance_type  = "vip"
  max_tps        = 5000
  queue_capacity = 50
  storage_size   = 700
  support_eip    = true
  max_eip_tps    = 128
  payment_type   = "Subscription"
  period         = 1
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Optional, Available in v1.131.0+) The instance name.
* `instance_type` - (Required, ForceNew) The Instance Type. Valid values: `professional`, `enterprise`, `vip`.
* `max_eip_tps` - (Optional) The max eip tps. It is valid when `support_eip` is true. The valid value is [128, 45000] with the step size 128.
* `max_tps` - (Required) The peak TPS traffic. The smallest valid value is 1000 and the largest value is 100,000.
* `modify_type` - (Optional) The modify type. Valid values: `Downgrade`, `Upgrade`. It is required when updating other attributes.
* `payment_type` - (Required) The payment type. Valid values: `Subscription`.
* `period` - (Optional) The period. Valid values: `1`, `12`, `2`, `24`, `3`, `6`.
* `queue_capacity` - (Required) The queue capacity. The smallest value is 50 and the step size 5.
* `renewal_duration` - (Optional) RenewalDuration. Valid values: `1`, `12`, `2`, `3`, `6`.
* `renewal_duration_unit` - (Optional) Auto-Renewal Cycle Unit Values Include: Month: Month. Year: Years. Valid values: `Month`, `Year`.
* `renewal_status` - (Optional) Whether to renew an instance automatically or not. Default to "ManualRenewal".
  - `AutoRenewal`: Auto renewal.
  - `ManualRenewal`: Manual renewal.
  - `NotRenewal`: No renewal any longer. After you specify this value, Alibaba Cloud stop sending notification of instance expiry, and only gives a brief reminder on the third day before the instance expiry.
  
* `storage_size` - (Optional) The storage size. It is valid when `instance_type` is vip.
* `support_eip` - (Required) Whether to support EIP.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Instance.
* `status` - The status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 min) Used when create the Instance.

## Import

RabbitMQ (AMQP) Instance can be imported using the id, e.g.

```
$ terraform import alicloud_amqp_instance.example <id>
```
