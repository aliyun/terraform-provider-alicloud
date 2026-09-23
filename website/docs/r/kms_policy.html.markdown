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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_kms_policy&exampleId=01de48b0-c178-276a-16d2-cc8a7cc3dc18c8504417&activeTab=example&spm=docs.r.kms_policy.0.01de48b0c1&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_kms_policy&spm=docs.r.kms_policy.example&intl_lang=EN_US)

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

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Policy.
* `update` - (Defaults to 5 mins) Used when update the Policy.

## Import

KMS Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_kms_policy.example <id>
```