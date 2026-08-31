---
subcategory: "Cloud SSO"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_sso_delegate_account"
description: |-
  Provides a Alicloud Cloud SSO Delegate Account resource.
---

# alicloud_cloud_sso_delegate_account

Provides a Cloud SSO Delegate Account resource.

Delegated Administrator Account.

For information about Cloud SSO Delegate Account and how to use it, see [What is Delegate Account](https://next.api.alibabacloud.com/document/cloudsso/2021-05-15/EnableDelegateAccount).

-> **NOTE:** Available since v1.259.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_sso_delegate_account&exampleId=3d14572f-c899-f5a7-3a39-b443c460236ac3a32505&activeTab=example&spm=docs.r.cloud_sso_delegate_account.0.3d14572fc8&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-shanghai"
}

data "alicloud_resource_manager_accounts" "default" {
  status = "CreateSuccess"
}
resource "alicloud_resource_manager_delegated_administrator" "default" {
  account_id        = data.alicloud_resource_manager_accounts.default.accounts.0.account_id
  service_principal = "cloudsso.aliyuncs.com"
}

resource "alicloud_cloud_sso_delegate_account" "default" {
  account_id = alicloud_resource_manager_delegated_administrator.default.account_id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cloud_sso_delegate_account&spm=docs.r.cloud_sso_delegate_account.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `account_id` - (Required, ForceNew) Delegate administrator account Id

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Delegate Account.
* `delete` - (Defaults to 5 mins) Used when delete the Delegate Account.

## Import

Cloud SSO Delegate Account can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_sso_delegate_account.example <id>
```