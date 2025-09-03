---
subcategory: "Anti-DDoS Pro (DdosBgp)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddosbgp_ip"
description: |-
  Provides a Alicloud Anti-DDoS Pro (DdosBgp) Ip resource.
---

# alicloud_ddosbgp_ip

Provides a Anti-DDoS Pro (DdosBgp) Ip resource.



For information about Anti-DDoS Pro (DdosBgp) Ip and how to use it, see [What is Ip](https://www.alibabacloud.com/help/en/ddos-protection/latest/addip).

-> **NOTE:** Available since v1.180.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ddosbgp_ip&exampleId=c68b456c-68bc-362a-1369-f2de392cfc0d379785d5&activeTab=example&spm=docs.r.ddosbgp_ip.0.c68b456c68&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "terraform-example"
}

data "alicloud_account" "default" {
}

resource "alicloud_ddosbgp_instance" "default" {
  name             = var.name
  base_bandwidth   = 20
  bandwidth        = -1
  ip_count         = 100
  ip_type          = "IPv4"
  normal_bandwidth = 100
  type             = "Enterprise"
}

resource "alicloud_eip_address" "default" {
  address_name = var.name
}

resource "alicloud_ddosbgp_ip" "default" {
  instance_id = alicloud_ddosbgp_instance.default.id
  ip          = alicloud_eip_address.default.ip_address
  member_uid  = data.alicloud_account.default.id
}
```

## Argument Reference

The following arguments are supported:
* `instance_id` - (Required, ForceNew) The ID of the Anti-DDoS Origin instance.
* `ip` - (Required, ForceNew) The IP address that you want to add.
* `member_uid` - (Optional, ForceNew, Available since v1.225.1) The member to which the asset belongs.
* `resource_group_id` - (Deprecated since v1.259.0) Field `resource_group_id` has been deprecated from provider version 1.259.0.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<instance_id>:<ip>`.
* `status` - The status of the IP address.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Ip.
* `delete` - (Defaults to 5 mins) Used when delete the Ip.

## Import

Anti-DDoS Pro (DdosBgp) Ip can be imported using the id, e.g.

```shell
$ terraform import alicloud_ddosbgp_ip.example <instance_id>:<ip>
```
