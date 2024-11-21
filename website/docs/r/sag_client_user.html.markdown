---
subcategory: "Smart Access Gateway (Smartag)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sag_client_user"
sidebar_current: "docs-alicloud-resource-sag-client-user"
description: |-
  Provides a Sag ClientUser resource.
---

# alicloud_sag_client_user

Provides a Sag ClientUser resource. This topic describes how to manage accounts as an administrator. After you configure the network, you can create multiple accounts and distribute them to end users so that clients can access Alibaba Cloud.

For information about Sag ClientUser and how to use it, see [What is Sag ClientUser](https://www.alibabacloud.com/help/en/smart-access-gateway/latest/createsmartaccessgatewayclientuser).

-> **NOTE:** Available since v1.65.0.

-> **NOTE:** Only the following regions support. [`cn-shanghai`, `cn-shanghai-finance-1`, `cn-hongkong`, `ap-southeast-1`, `ap-southeast-3`, `ap-southeast-5`, `ap-northeast-1`, `eu-central-1`]

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_sag_client_user&exampleId=f3d00586-5ad9-704f-034a-d39935efd9bccd8508f7&activeTab=example&spm=docs.r.sag_client_user.0.f3d005865a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
variable "sag_id" {
  default = "sag-9bifkfaz4fg***"
}
provider "alicloud" {
  region = "cn-shanghai"
}

resource "alicloud_sag_client_user" "default" {
  sag_id    = var.sag_id
  bandwidth = "20"
  user_mail = "tf-example@abc.com"
  user_name = var.name
  password  = "example1234"
  client_ip = "192.1.10.0"
}
```
## Argument Reference

The following arguments are supported:

* `sag_id` - (Required, ForceNew) The ID of the SAG instance created for the SAG APP.
* `bandwidth` - (Required) The SAG APP bandwidth that the user can use. Unit: Kbit/s. Maximum value: 2000 Kbit/s.
* `user_mail` - (Required, ForceNew) The email address of the user. The administrator uses this address to send the account information for logging on to the APP to the user.
* `user_name` - (Optional, ForceNew) The user name. User names in the same SAG APP must be unique.Both the user name and the password must be specified. If you specify the user name, the password must be specified, too.
* `password` - (Optional, ForceNew) The password used to log on to the SAG APP.Both the user name and the password must be specified. If you specify the user name, the password must be specified, too.
* `client_ip` - (Optional, ForceNew) The IP address of the SAG APP. If you specify this parameter, the current account always uses the specified IP address.Note The IP address must be in the private CIDR block of the SAG client.If you do not specify this parameter, the system automatically allocates an IP address from the private CIDR block of the SAG client. In this case, each re-connection uses a different IP address.
* `kms_encrypted_password` - (Optional) The password of the KMS Encryption.
* `kms_encryption_context` - (Optional) The context of the KMS Encryption.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Sag Id and formates as `<sag_id>:<user_name>`.

## Import

The Sag ClientUser can be imported using the name, e.g.

```shell
$ terraform import alicloud_sag_client_user.example sag-abc123456:tf-username-abc123456
```

