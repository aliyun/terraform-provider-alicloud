---
subcategory: "RocketMQ"
layout: "alicloud"
page_title: "Alicloud: alicloud_ons_instance"
sidebar_current: "docs-alicloud-resource-ons-instance"
description: |-
  Provides a Alicloud ONS Instance resource.
---

# alicloud\_ons\_instance

Provides an ONS instance resource.

For more information about how to use it, see [RocketMQ Instance Management API](https://www.alibabacloud.com/help/doc-detail/106354.html). 

-> **NOTE:** The number of instances in the same region cannot exceed 8. At present, the resource does not support region "mq-internet-access" and "ap-southeast-5". 

-> **NOTE:** Available in 1.51.0+

## Example Usage

Basic Usage

```
resource "alicloud_ons_instance" "example" {
  name   = "tf-example-ons-instance"
  remark = "tf-example-ons-instance-remark"
}
```

## Argument Reference

The following arguments are supported:


* `name` - (Required)Two instances on a single account in the same region cannot have the same name. The length must be 3 to 64 characters. Chinese characters, English letters digits and hyphen are allowed.
* `remark` - (Optional)This attribute is a concise description of instance. The length cannot exceed 128.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above.
* `instance_type` - The edition of instance. 1 represents the postPaid edition, and 2 represents the platinum edition.
* `instance_status` - The status of instance. 1 represents the platinum edition instance is in deployment. 2 represents the postpaid edition instance are overdue. 5 represents the postpaid or platinum edition instance is in service. 7 represents the platinum version instance is in upgrade and the service is available.
* `release_time` - Platinum edition instance expiration time.

## Import

ONS INSTANCE can be imported using the id, e.g.

```
$ terraform import alicloud_ons_instance.instance MQ_INST_1234567890_Baso1234567
```
