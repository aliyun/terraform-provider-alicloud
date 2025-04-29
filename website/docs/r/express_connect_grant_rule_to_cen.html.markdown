---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_grant_rule_to_cen"
sidebar_current: "docs-alicloud-resource-express-connect-grant-rule-to-cen"
description: |-
  Provides a Alicloud Express Connect Grant Rule To Cen resource.
---

# alicloud_express_connect_grant_rule_to_cen

Provides a Express Connect Grant Rule To Cen resource.

For information about Express Connect Grant Rule To Cen and how to use it, see [What is Grant Rule To Cen](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/grantinstancetocen).

-> **NOTE:** Available since v1.196.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_express_connect_grant_rule_to_cen&exampleId=a51ebeb2-c935-de2b-801c-bed796c4ed59886c348f&activeTab=example&spm=docs.r.express_connect_grant_rule_to_cen.0.a51ebeb2c9&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}
variable "name" {
  default = "tf-example"
}
data "alicloud_express_connect_physical_connections" "example" {
  name_regex = "^preserved-NODELETING"
}
resource "random_integer" "vlan_id" {
  max = 2999
  min = 1
}
resource "alicloud_express_connect_virtual_border_router" "example" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.example.connections.0.id
  virtual_border_router_name = var.name
  vlan_id                    = random_integer.vlan_id.id
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

resource "alicloud_cen_instance" "example" {
  cen_instance_name = var.name
}
data "alicloud_account" "default" {}

resource "alicloud_express_connect_grant_rule_to_cen" "example" {
  cen_id       = alicloud_cen_instance.example.id
  cen_owner_id = data.alicloud_account.default.id
  instance_id  = alicloud_express_connect_virtual_border_router.example.id
}
```

## Argument Reference

The following arguments are supported:

* `cen_id` - (Required, ForceNew) The ID of the CEN instance to which you want to grant permissions.
* `cen_owner_id` - (Required, ForceNew) The user ID (UID) of the Alibaba Cloud account to which the CEN instance belongs.
* `instance_id` - (Required, ForceNew) The ID of the VBR.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Grant Rule To Cen. It formats as `<cen_id>:<cen_owner_id>:<instance_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Grant Rule To Cen.
* `delete` - (Defaults to 3 mins) Used when delete the Grant Rule To Cen.

## Import

Express Connect Grant Rule To Cen can be imported using the id, e.g.

```shell
$ terraform import alicloud_express_connect_grant_rule_to_cen.example <cen_id>:<cen_owner_id>:<instance_id>
```
