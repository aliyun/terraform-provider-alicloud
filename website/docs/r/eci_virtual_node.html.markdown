---
subcategory: "Elastic Container Instance (ECI)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eci_virtual_node"
sidebar_current: "docs-alicloud-resource-eci-virtual-node"
description: |-
  Provides a Alicloud ECI Virtual Node resource.
---

# alicloud_eci_virtual_node

Provides a ECI Virtual Node resource.

For information about ECI Virtual Node and how to use it, see [What is Virtual Node](https://www.alibabacloud.com/help/en/doc-detail/89129.html).

-> **NOTE:** Available since v1.145.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_eci_virtual_node&exampleId=a3803358-0917-e69e-2bcc-e7580c32d772d71ffc4d&activeTab=example&spm=docs.r.eci_virtual_node.0.a380335809&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_eci_zones" "default" {}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.1.0.0/16"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_eci_zones.default.zones.0.zone_ids.0
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_eip_address" "default" {
  isp                       = "BGP"
  address_name              = var.name
  netmode                   = "public"
  bandwidth                 = "1"
  security_protection_types = ["AntiDDoS_Enhanced"]
  payment_type              = "PayAsYouGo"
}
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_eci_virtual_node" "default" {
  security_group_id     = alicloud_security_group.default.id
  virtual_node_name     = var.name
  vswitch_id            = alicloud_vswitch.default.id
  enable_public_network = false
  eip_instance_id       = alicloud_eip_address.default.id
  resource_group_id     = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  kube_config           = "kube_config"
  tags = {
    Created = "TF"
  }
  taints {
    effect = "NoSchedule"
    key    = "TF"
    value  = "example"
  }
}
```

## Argument Reference

The following arguments are supported:

* `eip_instance_id` - (Optional, ForceNew) The Id of eip.
* `enable_public_network` - (Optional, ForceNew) Whether to enable public network. **NOTE:** If `eip_instance_id` is not configured and `enable_public_network` is true, the system will create an elastic public network IP.
* `kube_config` - (Required) The kube config for the k8s cluster. It needs to be connected after Base64 encoding.
* `resource_group_id` - (Optional, ForceNew) The resource group ID. 
* `security_group_id` - (Required, ForceNew) The security group ID.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.
* `taints` - (Optional) The taint. See [`taints`](#taints) below.
* `virtual_node_name` - (Optional, ForceNew) The name of the virtual node. The length of the name is limited to `2` to `128` characters. It can contain uppercase and lowercase letters, Chinese characters, numbers, half-width colon (:), underscores (_), or hyphens (-), and must start with letters.
* `vswitch_id` - (Required, ForceNew) The vswitch id.
* `zone_id` - (Optional, ForceNew) The Zone.

### `taints`

The taints supports the following:

* `effect` - (Optional) The effect of the taint. Valid values: `NoSchedule`, `NoExecute` and `PreferNoSchedule`.
* `key` - (Optional) The key of the taint.
* `value` - (Optional) The value of the taint.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Virtual Node.
* `status` - The Status of the virtual node. Valid values: `Cleaned`, `Failed`, `Pending`, `Ready`.

## Timeouts

The `timeouts` block allows you to
specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Gateway File Share.

## Import

ECI Virtual Node can be imported using the id, e.g.

```shell
$ terraform import alicloud_eci_virtual_node.example <id>
```