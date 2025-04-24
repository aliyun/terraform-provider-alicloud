---
subcategory: "Anti-DDoS Pro (DdosBgp)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddosbgp_ip"
sidebar_current: "docs-alicloud-resource-ddos-bgp-ip"
description: |-
  Provides a Alicloud Ddos Bgp Ip resource.
---

# alicloud_ddosbgp_ip

Provides a Ddos Bgp Ip resource.

For information about Ddos Bgp Ip and how to use it, see [What is Ip](https://www.alibabacloud.com/help/en/ddos-protection/latest/addip).

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
  region = "cn-beijing"
}

variable "name" {
  default = "tf-example"
}
data "alicloud_resource_manager_resource_groups" "default" {}
resource "alicloud_ddosbgp_instance" "instance" {
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
  instance_id       = alicloud_ddosbgp_instance.instance.id
  ip                = alicloud_eip_address.default.ip_address
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the native protection enterprise instance to be operated.
* `ip` - (Required, ForceNew) The IP address.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `member_uid` - (Optional, ForceNew, Available since v1.225.1) The member account id of the IP address.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Ip. The value formats as `<instance_id>:<ip>`.
* `status` - The current state of the IP address. Valid Value: `normal`, `hole_begin`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Ddos Bgp Ip.
* `delete` - (Defaults to 1 mins) Used when deleting the Ddos Bgp Ip.

## Import

Ddos Bgp Ip can be imported using the id, e.g.

```shell
$ terraform import alicloud_ddosbgp_ip.example <instance_id>:<ip>
```
