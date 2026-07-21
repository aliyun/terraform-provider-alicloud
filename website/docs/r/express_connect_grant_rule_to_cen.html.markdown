---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_grant_rule_to_cen"
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
  <a href="https://api.aliyun.com/terraform?resource=alicloud_express_connect_grant_rule_to_cen&exampleId=7805dce3-9982-8444-26d8-ccf9207a27ab4ed215ce&activeTab=example&spm=docs.r.express_connect_grant_rule_to_cen.0.7805dce399&intl_lang=EN_US" target="_blank">
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

data "alicloud_express_connect_physical_connections" "default" {
  name_regex = "^preserved-NODELETING"
}

resource "random_integer" "default" {
  max = 2999
  min = 1
}

resource "alicloud_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alicloud_express_connect_physical_connections.default.connections.0.id
  virtual_border_router_name = var.name
  vlan_id                    = random_integer.default.id
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
}

resource "alicloud_express_connect_grant_rule_to_cen" "default" {
  cen_id       = alicloud_cen_instance.default.id
  cen_owner_id = data.alicloud_account.default.id
  instance_id  = alicloud_express_connect_virtual_border_router.default.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_express_connect_grant_rule_to_cen&spm=docs.r.express_connect_grant_rule_to_cen.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `cen_id` - (Required, ForceNew) The ID of the CEN instance to which you want to grant permissions.
* `cen_owner_id` - (Required, ForceNew) The user ID (UID) of the Alibaba Cloud account to which the CEN instance belongs.
* `instance_id` - (Required, ForceNew) The ID of the VBR.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<cen_id>:<cen_owner_id>:<instance_id>`.
* `create_time` - (Available since v1.263.0) The time when the instance was created.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Grant Rule To Cen.
* `delete` - (Defaults to 5 mins) Used when delete the Grant Rule To Cen.

## Import

Express Connect Grant Rule To Cen can be imported using the id, e.g.

```shell
$ terraform import alicloud_express_connect_grant_rule_to_cen.example <cen_id>:<cen_owner_id>:<instance_id>
```
