---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_network_acl_attachment"
description: |-
  Provides a Alicloud VPC Network Acl Attachment resource.
---

# alicloud_vpc_network_acl_attachment

Provides a VPC Network Acl Attachment resource. Resources associated with network Acl.

For information about VPC Network Acl Attachment and how to use it, see [What is Network Acl Attachment](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/associatenetworkacl).

-> **NOTE:** Available since v1.193.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_network_acl_attachment&exampleId=a640762e-86ca-6d0a-df56-d5a6e3b393fe0e6d9482&activeTab=example&spm=docs.r.vpc_network_acl_attachment.0.a640762e86&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

## Argument Reference

The following arguments are supported:

* `network_acl_id` - (Required, ForceNew) The ID of the network ACL.
* `resource_id` - (Required, ForceNew) The ID of the associated resource.
* `resource_type` - (Required, ForceNew) The type of the associated resource. Valid values: `VSwitch`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Network Acl Attachment. The value formats as `<network_acl_id>:<resource_id>`.
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