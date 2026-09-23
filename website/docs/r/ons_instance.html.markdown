---
subcategory: "RocketMQ (Ons)"
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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ons_instance&exampleId=cae9cec5-f3cd-e263-257b-9a3b9ff8461b081583fd&activeTab=example&spm=docs.r.ons_instance.0.cae9cec5f3&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_ons_instance" "example" {
  instance_name = "${var.name}-${random_integer.default.result}"
  remark        = var.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ons_instance&spm=docs.r.ons_instance.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:


* `name` - (Optional, Deprecated from v1.97.0+) Replaced by `instance_name` after version 1.97.0.
* `instance_name` - (Optional, Available in v1.97.0+) Two instances on a single account in the same region cannot have the same name. The length must be 3 to 64 characters. Chinese characters, English letters digits and hyphen are allowed.
* `remark` - (Optional) This attribute is a concise description of instance. The length cannot exceed 128.
* `tags` - (Optional, Available in v1.97.0+) A mapping of tags to assign to the resource.
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above.
* `instance_type` - The edition of instance. 1 represents the postPaid edition, and 2 represents the platinum edition.
* `instance_status` - The status of instance. 1 represents the platinum edition instance is in deployment. 2 represents the postpaid edition instance are overdue. 5 represents the postpaid or platinum edition instance is in service. 7 represents the platinum version instance is in upgrade and the service is available.
* `release_time` - Platinum edition instance expiration time.
* `status` - The status of instance. 1 represents the platinum edition instance is in deployment. 2 represents the postpaid edition instance are overdue. 5 represents the postpaid or platinum edition instance is in service. 7 represents the platinum version instance is in upgrade and the service is available.

## Import

ONS INSTANCE can be imported using the id, e.g.

```shell
$ terraform import alicloud_ons_instance.instance MQ_INST_1234567890_Baso1234567
```
