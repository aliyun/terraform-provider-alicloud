---
subcategory: "Elastic Container Instance (ECI)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eci_virtual_node"
sidebar_current: "docs-alicloud-resource-eci-virtual-node"
description: |- 
  Provides a Alicloud ECI Virtual Node resource.
---

# alicloud\_eci\_virtual\_node

Provides a ECI Virtual Node resource.

For information about ECI Virtual Node and how to use it, see [What is Virtual Node](https://www.alibabacloud.com/help/en/doc-detail/89129.html).

-> **NOTE:** Available in v1.145.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testaccvirtualnode"
}

data "alicloud_eci_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_eci_zones.default.zones.0.zone_ids.1
}

resource "alicloud_security_group" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  name   = var.name
}

resource "alicloud_eip_address" "default" {
  address_name = var.name
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_eci_virtual_node" "default" {
  security_group_id     = alicloud_security_group.default.id
  virtual_node_name     = var.name
  vswitch_id            = data.alicloud_vswitches.default.ids.1
  enable_public_network = false
  eip_instance_id       = alicloud_eip_address.default.id
  resource_group_id     = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  kube_config           = "kube config"
  tags = {
    Created = "TF"
  }
  taints {
    effect = "NoSchedule"
    key    = "Tf1"
    value  = "Test1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `eip_instance_id` - (Optional, Computed,ForceNew) The Id of eip.
* `enable_public_network` - (Optional, ForceNew) Whether to enable public network. **NOTE:** If `eip_instance_id` is not configured and `enable_public_network` is true, the system will create an elastic public network IP.
* `kube_config` - (Optional) The kube config for the k8s cluster. It needs to be connected after Base64 encoding.
* `resource_group_id` - (Optional, ForceNew) The resource group ID. 
* `security_group_id` - (Required, ForceNew) The security group ID.
* `tags` - (Optional, ForceNew) A mapping of tags to assign to the resource.
* `taints` - (Optional) The taint. See the following `Block taints`.
* `virtual_node_name` - (Optional, ForceNew) The name of the virtual node. The length of the name is limited to `2` to `128` characters. It can contain uppercase and lowercase letters, Chinese characters, numbers, half-width colon (:), underscores (_), or hyphens (-), and must start with letters.
* `vswitch_id` - (Required, ForceNew) The vswitch id.
* `zone_id` - (Optional, Computed, ForceNew) The Zone.

#### Block taints

The taints supports the following:

* `effect` - (Optional) The effect of the taint. Valid values: `NoSchedule`, `NoExecute` and `PreferNoSchedule`.
* `key` - (Optional) The key of the taint.
* `value` - (Optional) The value of the taint.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Virtual Node.
* `status` - The Status of the virtual node. Valid values: `Cleaned`, `Failed`, `Pending`, `Ready`.

### Timeouts

The `timeouts` block allows you to
specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Gateway File Share.

## Import

ECI Virtual Node can be imported using the id, e.g.

```
$ terraform import alicloud_eci_virtual_node.example <id>
```