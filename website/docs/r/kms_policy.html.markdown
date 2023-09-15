---
subcategory: "KMS"
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_policy"
description: |-
  Provides a Alicloud KMS Policy resource.
---

# alicloud_kms_policy

Provides a KMS Policy resource. Permission policies which can be bound to the Application Access Points.

For information about KMS Policy and how to use it, see [What is Policy](https://www.alibabacloud.com/help/zh/key-management-service/latest/api-createpolicy).

-> **NOTE:** Available since v1.210.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_kms_network_rule" "networkRule1" {
  description       = "dummy"
  source_private_ip = ["10.10.10.10"]
  network_rule_name = format("%s1", var.name)
}

resource "alicloud_kms_network_rule" "networkRule2" {
  description       = "dummy"
  source_private_ip = ["10.10.10.10"]
  network_rule_name = format("%s2", var.name)
}

resource "alicloud_kms_network_rule" "networkRule3" {
  description       = "dummy"
  source_private_ip = ["10.10.10.10"]
  network_rule_name = format("%s3", var.name)
}


resource "alicloud_kms_policy" "default" {
  description          = "terraformpolicy"
  permissions          = ["RbacPermission/Template/CryptoServiceKeyUser", "RbacPermission/Template/CryptoServiceSecretUser"]
  resources            = ["secret/*", "key/*"]
  policy_name          = var.name
  kms_instance_id      = "shared"
  access_control_rules = <<EOF
  {
      "NetworkRules":[
          "alicloud_kms_network_rule.networkRule1.network_rule_name"
      ]
  }
  EOF
}
```

## Argument Reference

The following arguments are supported:
* `access_control_rules` - (Required) Network Rules in JSON struct.
* `description` - (Optional) Description.
* `kms_instance_id` - (Required, ForceNew) KMS instance .
* `permissions` - (Required) Allowed permissions (RBAC)Optional values:"RbacPermission/Template/CryptoServiceKeyUser" and "RbacPermission/Template/CryptoServiceSecretUser".
* `policy_name` - (Required, ForceNew) Policy Name.
* `resources` - (Required) The resources that the permission policy allows to access.Use "key/${KeyId}" or "key/*"  to specify a key or all keys.Use "secret/${SecretName}" or "secret/*" to specify a secret or all secrets.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Policy.
* `update` - (Defaults to 5 mins) Used when update the Policy.

## Import

KMS Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_kms_policy.example <id>
```