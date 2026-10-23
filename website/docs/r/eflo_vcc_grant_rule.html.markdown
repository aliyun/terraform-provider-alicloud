---
subcategory: "unreleased"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_vcc_grant_rule"
description: |-
  Provides a Alicloud Eflo Vcc Grant Rule resource.
---

# alicloud_eflo_vcc_grant_rule

Provides a Eflo Vcc Grant Rule resource.

Lingjun Connect Cross Account Authorization.

For information about Eflo Vcc Grant Rule and how to use it, see [What is Vcc Grant Rule](https://next.api.alibabacloud.com/document/eflo/2022-05-30/CreateVccGrantRule).

-> **NOTE:** Available since v2.1.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-wulanchabu"
}

variable "vcc_id" {
  default = "vcc-xxxxxxxx"
}

data "alicloud_account" "default" {
}

resource "alicloud_eflo_er" "default" {
  er_name        = "tf-example-er"
  master_zone_id = "cn-wulanchabu-b"
}

resource "alicloud_eflo_vcc_grant_rule" "default" {
  grant_tenant_id = data.alicloud_account.default.id
  er_id           = alicloud_eflo_er.default.id
  instance_id     = var.vcc_id
}
```

## Argument Reference

The following arguments are supported:
* `er_id` - (Required, ForceNew) The ID of the ER instance
* `grant_tenant_id` - (Required, ForceNew) Authorized tenant ID
* `instance_id` - (Required, ForceNew) Network instance ID

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<instance_id>:<grant_rule_id>`.
* `create_time` - Creation time.
* `grant_rule_id` - The ID of the authorized resource instance.
* `instance_name` - Network instance name.
* `product` - Network product CODE: VPD, VCC.
* `tenant_id` - Tenant ID.
* `used` - Whether the current authorization information has been used; The default is false.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vcc Grant Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Vcc Grant Rule.

## Import

Eflo Vcc Grant Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_vcc_grant_rule.example <instance_id>:<grant_rule_id>
```