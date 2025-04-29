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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_kms_network_rule&exampleId=a306da17-a94a-9576-3d5e-2ffb37a61a24c57d24aa&activeTab=example&spm=docs.r.kms_network_rule.0.a306da17a9&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Network Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Network Rule.
* `update` - (Defaults to 5 mins) Used when update the Network Rule.

## Import

KMS Network Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_kms_network_rule.example <id>
```