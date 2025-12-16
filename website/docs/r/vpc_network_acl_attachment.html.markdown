---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_network_acl_attachment"
description: |-
  Provides a Alicloud VPC Network Acl Attachment resource.
---

# alicloud_vpc_network_acl_attachment

Provides a VPC Network Acl Attachment resource.

Resources associated with network Acl.

For information about VPC Network Acl Attachment and how to use it, see [What is Network Acl Attachment](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/associatenetworkacl).

-> **NOTE:** Available since v1.193.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 2)
  zone_id    = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_network_acl" "default" {
  vpc_id = alicloud_vswitch.default.vpc_id
}

resource "alicloud_vpc_network_acl_attachment" "default" {
  network_acl_id = alicloud_network_acl.default.id
  resource_id    = alicloud_vswitch.default.id
  resource_type  = "VSwitch"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_vpc_network_acl_attachment&spm=docs.r.vpc_network_acl_attachment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `network_acl_id` - (Required, ForceNew) The ID of the network ACL instance.
* `resource_id` - (Required, ForceNew) The ID of the associated  vswitch instance.
* `resource_type` - (Required, ForceNew) The type of the associated resource. Valid values: `VSwitch`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<network_acl_id>:<resource_id>`.
* `status` - The status of the Network Acl Attachment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Network Acl Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Network Acl Attachment.

## Import

VPC Network Acl Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_network_acl_attachment.example <network_acl_id>:<resource_id>
```