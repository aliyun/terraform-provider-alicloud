---
subcategory: "Anycast Elastic IP Address (Eipanycast)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eipanycast_anycast_eip_address_attachment"
sidebar_current: "docs-alicloud-resource-eipanycast-anycast-eip-address-attachment"
description: |-
  Provides a Alicloud Anycast Eip Address Attachment resource.
---

# alicloud\_eipanycast\_anycast\_eip\_address\_attachment

Provides a Eipanycast Anycast Eip Address Attachment resource.

For information about Eipanycast Anycast Eip Address Attachment and how to use it, see [What is Anycast Eip Address Attachment](https://help.aliyun.com/document_detail/171857.html).

-> **NOTE:** Available in v1.113.0+.

-> **NOTE:** The following regions support currently while Slb instance support bound. 
[eu-west-1-gb33-a01,cn-hongkong-am4-c04,ap-southeast-os30-a01,us-west-ot7-a01,ap-south-in73-a01,ap-southeast-my88-a01]

## Example Usage

Basic Usage

```terraform
resource "alicloud_eipanycast_anycast_eip_address" "example" {
  service_location = "international"
}

resource "alicloud_eipanycast_anycast_eip_address_attachment" "example" {
  anycast_id              = alicloud_eipanycast_anycast_eip_address.example.id
  bind_instance_id        = "lb-j6chlcr8lffy7********"
  bind_instance_region_id = "cn-hongkong"
  bind_instance_type      = "SlbInstance"
}

```

## Argument Reference

The following arguments are supported:

* `anycast_id` - (Required, ForceNew) The ID of Anycast EIP.
* `bind_instance_id` - (Required, ForceNew) The ID of bound instance.
* `bind_instance_region_id` - (Required, ForceNew) The region ID of bound instance.
* `bind_instance_type` - (Required, ForceNew) The type of bound instance. Valid value: `SlbInstance`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Anycast Eip Address Attachment. The ID format as `anycast_id`:`bind_instance_id`:`bind_instance_region_id`:`bind_instance_type`.
* `bind_time` - The time of bound instance.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Anycast Eip Address Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Anycast Eip Address Attachment.

## Import

Eipanycast Anycast Eip Address Attachment can be imported using the id, e.g.

```
$ terraform import alicloud_eipanycast_anycast_eip_address_attachment.example `anycast_id`:`bind_instance_id`:`bind_instance_region_id`:`bind_instance_type`
```
