---
subcategory: "RabbitMQ (AMQP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_amqp_instance"
description: |-
  Provides a Alicloud Amqp Instance resource.
---

# alicloud_amqp_instance

Provides a Amqp Instance resource. The instance of Amqp.

For information about RabbitMQ (AMQP) Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/en/message-queue-for-rabbitmq/latest/createinstance).

-> **NOTE:** Available since v1.128.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_amqp_instance&exampleId=45cec220-2e60-4a5f-0eaa-39abcd8e5e46e0861292&activeTab=example&spm=docs.r.amqp_instance.0.45cec2202e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-shanghai"
}


resource "alicloud_amqp_instance" "default" {
  instance_name  = var.name
  instance_type  = "professional"
  max_tps        = "1000"
  queue_capacity = "50"
  period_cycle   = "Year"
  support_eip    = "false"
  period         = "1"
  auto_renew     = "true"
  payment_type   = "Subscription"
}
```

### Deleting `alicloud_amqp_instance` or removing it from your configuration

The `alicloud_amqp_instance` resource allows you to manage  `payment_type = "PayAsYouGo"`  instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Instance.
You can resume managing the subscription instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:
* `auto_renew` - (Optional, Available since v1.129.0) Renewal method. Automatic renewal: true; Manual renewal: false. When RenewalStatus has a value, the value of RenewalStatus shall prevail.
* `instance_name` - (Optional, Computed) The instance name.
* `instance_type` - (Optional, Computed) Instance type. Valid values are as follows:  professional: professional Edition enterprise: enterprise Edition vip: Platinum Edition.
* `max_connections` - (Optional, Computed, Available since v1.129.0) The maximum number of connections, according to the value given on the purchase page of the cloud message queue RabbitMQ version console.
* `max_eip_tps` - (Optional, Computed) Peak TPS traffic of the public network, which must be a multiple of 128, unit: times per second.
* `max_tps` - (Optional, Computed) Configure the private network TPS traffic peak, please set the value according to the cloud message queue RabbitMQ version of the console purchase page given.
* `modify_type` - (Optional) This parameter must be provided while you change the instance specification. Type of instance lifting and lowering:
  - Upgrade: Upgrade
  - Downgrade: Downgrading.
* `payment_type` - (Required, ForceNew) The Payment type. Valid value: Subscription: prepaid. PayAsYouGo: Post-paid.
* `period` - (Optional) Prepayment cycle, unit: periodCycle.  This parameter is valid when PaymentType is set to Subscription.
* `period_cycle` - (Optional, Available since v1.129.0) Prepaid cycle units. Value: Month. Year: Year.
* `queue_capacity` - (Optional, Computed) Configure the maximum number of queues. The value range is as follows:  Professional version:[50,1000], minimum modification step size is 5  Enterprise Edition:[200,6000], minimum modification step size is 100  Platinum version:[10000,80000], minimum modification step size is 100.
* `renewal_duration` - (Optional, Computed) The number of automatic renewal cycles.
* `renewal_duration_unit` - (Optional, Computed) Auto-Renewal Cycle Unit Values Include: Month: Month. Year: Years.
* `renewal_status` - (Optional, Computed) The renewal status. Value: AutoRenewal: automatic renewal. ManualRenewal: manual renewal. NotRenewal: no renewal.
* `serverless_charge_type` - (Optional, Available since v1.129.0) The billing type of the serverless instance. Value: onDemand.
* `storage_size` - (Optional, Computed) Configure the message storage space. Unit: GB. The value is as follows:  Professional Edition and Enterprise Edition: Fixed to 0. Description A value of 0 indicates that the Professional Edition and Enterprise Edition instances do not charge storage fees, but do not have storage space. Platinum version example: m × 100, where the value range of m is [7,28].
* `support_eip` - (Optional) Whether to support public network.
* `support_tracing` - (Optional, Computed) Whether to activate the message trace function. The values are as follows:  true: Enable message trace function false: message trace function is not enabled Description The Platinum Edition instance provides the 15-day message trace function free of charge. The trace function can only be enabled and the trace storage duration can only be set to 15 days. For instances of other specifications, you can enable or disable the trace function.
* `tracing_storage_time` - (Optional, Computed) Configure the storage duration of message traces. Unit: Days. The value is as follows:  3:3 days 7:7 days 15:15 days This parameter is valid when SupportTracing is true.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - OrderCreateTime.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Instance.
* `update` - (Defaults to 5 mins) Used when update the Instance.

## Import

Amqp Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_amqp_instance.example <id>
```