---
subcategory: "Smart Access Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_sag_client_user"
sidebar_current: "docs-alicloud-resource-sag-client-user"
description: |-
  Provides a Sag ClientUser resource.
---

# alicloud\_sag\_client_user

Provides a Sag ClientUser resource. This topic describes how to manage accounts as an administrator. After you configure the network, you can create multiple accounts and distribute them to end users so that clients can access Alibaba Cloud.

For information about Sag ClientUser and how to use it, see [What is Sag ClientUser](https://www.alibabacloud.com/help/doc-detail/108326.htm).

-> **NOTE:** Available in 1.65.0+

-> **NOTE:** Only the following regions support. [`cn-shanghai`, `cn-shanghai-finance-1`, `cn-hongkong`, `ap-southeast-1`, `ap-southeast-2`, `ap-southeast-3`, `ap-southeast-5`, `ap-northeast-1`, `eu-central-1`]

## Example Usage

Basic Usage

```
resource "alicloud_sag_client_user" "default" {
  sag_id    = "sag-xxxxx"
  bandwidth = "20"
  user_mail = "tftest-xxxxx@test.com"
  user_name = "th-username-xxxxx"
  password  = "xxxxxxx"
  client_ip = "192.1.10.0"
}
```
## Argument Reference

The following arguments are supported:

* `sag_id` - (Required,ForceNew) The ID of the SAG instance created for the SAG APP.
* `bandwidth` - (Required) The SAG APP bandwidth that the user can use. Unit: Kbit/s. Maximum value: 2000 Kbit/s.
* `user_mail` - (Required,ForceNew) The email address of the user. The administrator uses this address to send the account information for logging on to the APP to the user.
* `user_name` - (Optional,ForceNew) The user name. User names in the same SAG APP must be unique.Both the user name and the password must be specified. If you specify the user name, the password must be specified, too.
* `password` - (Optional,ForceNew) The password used to log on to the SAG APP.Both the user name and the password must be specified. If you specify the user name, the password must be specified, too.
* `client_ip` - (Optional,ForceNew) The IP address of the SAG APP. If you specify this parameter, the current account always uses the specified IP address.Note The IP address must be in the private CIDR block of the SAG client.If you do not specify this parameter, the system automatically allocates an IP address from the private CIDR block of the SAG client. In this case, each re-connection uses a different IP address.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Sag Id and formates as `<sag_id>:<user_name>`.

## Import

The Sag ClientUser can be imported using the name, e.g.

```
$ terraform import alicloud_sag_client_user.example sag-abc123456:tf-username-abc123456
```

