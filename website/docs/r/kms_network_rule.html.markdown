---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_network_rule"
description: |-
  Provides a Alicloud KMS Network Rule resource.
---

# alicloud_kms_network_rule

Provides a KMS Network Rule resource. Network rules that can be bound by Application Access Point's policies.

For information about KMS Network Rule and how to use it, see [What is Network Rule](https://www.alibabacloud.com/help/zh/key-management-service/latest/api-createnetworkrule).

-> **NOTE:** Available since v1.210.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}


resource "alicloud_kms_network_rule" "default" {
  description       = "example-description"
  source_private_ip = ["10.10.10.10/24", "192.168.17.13", "100.177.24.254"]
  network_rule_name = var.name
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) Description.
* `network_rule_name` - (Optional, ForceNew, Computed) Network Rule Name.
* `source_private_ip` - (Required) Allowed private network addresses.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Network Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Network Rule.
* `update` - (Defaults to 5 mins) Used when update the Network Rule.

## Import

KMS Network Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_kms_network_rule.example <id>
```