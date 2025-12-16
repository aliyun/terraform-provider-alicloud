---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_instance"
sidebar_current: "docs-alicloud-resource-bastionhost-instance"
description: |-
  Provides a Alicloud Bastion Host Instance Resource.
---

# alicloud_bastionhost_instance

-> **NOTE:** Since the version 1.132.0, the resource `alicloud_yundun_bastionhost_instance` has been renamed to `alicloud_bastionhost_instance`.

Cloud Bastion Host instance resource ("Yundun_bastionhost" is the short term of this product). 
For information about Resource Manager Resource Directory and how to use it, see [What is Bastionhost](https://www.alibabacloud.com/help/en/doc-detail/52922.htm).

-> **NOTE:** The endpoint of bssopenapi used only support "business.aliyuncs.com" at present.

-> **NOTE:** Available since v1.132.0.

-> **NOTE:** In order to destroy Cloud Bastionhost instance , users are required to apply for white list first

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_bastionhost_instance&exampleId=2781e605-de83-d24d-2811-754044d92d63d7257b56&activeTab=example&spm=docs.r.bastionhost_instance.0.2781e605de&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
  cidr_block = "10.4.0.0/16"
}

data "alicloud_vswitches" "default" {
  cidr_block = "10.4.0.0/24"
  vpc_id     = data.alicloud_vpcs.default.ids.0
  zone_id    = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_bastionhost_instance" "default" {
  description        = var.name
  license_code       = "bhah_ent_50_asset"
  plan_code          = "cloudbastion"
  storage            = "5"
  bandwidth          = "5"
  period             = "1"
  vswitch_id         = data.alicloud_vswitches.default.ids[0]
  security_group_ids = [alicloud_security_group.default.id]
}
```

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_bastionhost_instance&exampleId=2cb4cc04-7299-a5cd-3df8-b702a244c58880de2a97&activeTab=example&spm=docs.r.bastionhost_instance.1.2cb4cc0472&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
  cidr_block = "10.4.0.0/16"
}

data "alicloud_vswitches" "default" {
  cidr_block = "10.4.0.0/24"
  vpc_id     = data.alicloud_vpcs.default.ids.0
  zone_id    = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_bastionhost_instance" "default" {
  description        = var.name
  license_code       = "bhah_ent_50_asset"
  plan_code          = "cloudbastion"
  storage            = "5"
  bandwidth          = "5"
  period             = 1
  security_group_ids = [alicloud_security_group.default.id]
  vswitch_id         = data.alicloud_vswitches.default.ids[0]
  ad_auth_server {
    server         = "192.168.1.1"
    standby_server = "192.168.1.3"
    port           = "80"
    domain         = "domain"
    account        = "cn=Manager,dc=test,dc=com"
    password       = "YouPassword123"
    filter         = "objectClass=person"
    name_mapping   = "nameAttr"
    email_mapping  = "emailAttr"
    mobile_mapping = "mobileAttr"
    is_ssl         = false
    base_dn        = "dc=test,dc=com"
  }
  ldap_auth_server {
    server             = "192.168.1.1"
    standby_server     = "192.168.1.3"
    port               = "80"
    login_name_mapping = "uid"
    account            = "cn=Manager,dc=test,dc=com"
    password           = "YouPassword123"
    filter             = "objectClass=person"
    name_mapping       = "nameAttr"
    email_mapping      = "emailAttr"
    mobile_mapping     = "mobileAttr"
    is_ssl             = false
    base_dn            = "dc=test,dc=com"
  }
}
```

### Deleting `alicloud_bastionhost_instance` or removing it from your configuration

The `alicloud_bastionhost_instance` resource allows you to manage bastionhost instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration
will remove it from your state file and management, but will not destroy the bastionhost instance.
You can resume managing the subscription bastionhost instance via the AlibabaCloud Console.

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_bastionhost_instance&spm=docs.r.bastionhost_instance.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `license_code` - (Required)  The package type of Cloud Bastionhost instance. You can query more supported types through the [DescribePricingModule](https://help.aliyun.com/document_detail/96469.html).
* `plan_code` - (Required, ForceNew, Available since 1.193.0) The plan code of Cloud Bastionhost instance. Valid values:
  - `cloudbastion`: Basic Edition.
  - `cloudbastion_ha`: HA Edition.
* `storage` - (Required, Available since 1.193.0) The storage of Cloud Bastionhost instance. Valid values: `0` to `500`. Unit: TB. **NOTE:** From version 1.251.0, `storage` can be modified.
* `bandwidth` - (Required, Available since 1.193.0) The bandwidth of Cloud Bastionhost instance. **NOTE:** From version 1.263.0, `bandwidth` can be modified.
  If [China-Site Account](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/guides/getting-account#sign-up-for-an-alibaba-cloud-china-site-account), its valid values: 0 to 150. Unit: Mbit/s. The value must be a multiple of 5.
  If [International-Site Account](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/guides/getting-account#sign-up-for-an-alibaba-cloud-international-site-account), its valid values: 0 to 200. Unit: Mbit/s. The value must be a multiple of 10.
* `description` - (Required) Description of the instance. This name can have a string of 1 to 63 characters.
* `period` - (Optional) Duration for initially producing the instance. Valid values: [1~9], 12, 24, 36. At present, the provider does not support modify "period".
-> **NOTE:** The attribute `period` is only used to create Subscription instance or modify the PayAsYouGo instance to Subscription. Once effect, it will not be modified that means running `terraform apply` will not effect the resource.
* `vswitch_id` - (Required, ForceNew) VSwitch ID configured to Bastionhost.
* `security_group_ids` - (Required, List) security group IDs configured to Bastionhost. 
  **NOTE:** There is a potential diff error because of the order of `security_group_ids` values indefinite.
  So, from version 1.160.0, `security_group_ids` type has been updated as `set` from `list`,
  and you can use [tolist](https://www.terraform.io/language/functions/tolist) to convert it to a list.
* `slave_vswitch_id` - (Optional, ForceNew, Available since v1.263.0) Slave VSwitch ID configured to Bastionhost.
* `tags` - (Optional, Available since v1.67.0) A mapping of tags to assign to the resource.
* `resource_group_id` - (Optional, Available since v1.87.0) The Id of resource group which the Bastionhost Instance belongs. If not set, the resource is created in the default resource group.
* `enable_public_access` - (Optional, Available since v1.143.0)  Whether to Enable the public internet access to a specified Bastionhost instance. The valid values: `true`, `false`.
* `ad_auth_server` - (Optional, Available since 1.169.0) The AD auth server of the Instance. See [`ad_auth_server`](#ad_auth_server) below.
* `ldap_auth_server` - (Optional, Available since 1.169.0) The LDAP auth server of the Instance. See [`ldap_auth_server`](#ldap_auth_server) below.
* `renew_period` - (Optional, Available since 1.187.0) Automatic renewal period. Valid values: `1` to `9`, `12`, `24`, `36`. **NOTE:** The `renew_period` is required under the condition that `renewal_status` is `AutoRenewal`. From version 1.193.0, `renew_period` can be modified.
* `renewal_status` - (Optional, Computed, Available since 1.187.0) Automatic renewal status. Valid values: `AutoRenewal`, `ManualRenewal`, `NotRenewal`. From version 1.193.0, `renewal_status` can be modified.
* `renewal_period_unit` - (Optional, Computed, Available since 1.193.0) The unit of the auto-renewal period. Valid values:  **NOTE:** The `renewal_period_unit` is required under the condition that `renewal_status` is `AutoRenewal`.
  - `M`: months.
  - `Y`: years.
* `public_white_list` - (Optional, Available since 1.199.0) The public IP address that you want to add to the whitelist.

-> **NOTE:** You can utilize the generic Terraform resource [lifecycle configuration block](https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle) with `ad_auth_server` or `ldap_auth_server` to configure auth server, then ignore any changes to that `password` caused externally (e.g. Application Autoscaling).
```
  # ... ignore the change about ad_auth_server.0.password and ldap_auth_server.0.password in alicloud_bastionhost_instance
  lifecycle {
    ignore_changes = [ad_auth_server.0.password,ldap_auth_server.0.password]
  }
```

### `ad_auth_server`

The ad_auth_server supports the following:

* `account` - (Required) The username of the account that is used for the AD server.
* `base_dn` - (Required) The Base distinguished name (DN).
* `domain` - (Required) The domain on the AD server.
* `email_mapping` - (Optional) The field that is used to indicate the email address of a user on the AD server.
* `filter` - (Optional) The condition that is used to filter users.
* `is_ssl` - (Required) Specifies whether to support SSL.
* `mobile_mapping` - (Optional) The field that is used to indicate the mobile phone number of a user on the AD server.
* `name_mapping` - (Optional) The field that is used to indicate the name of a user on the AD server.
* `password` - (Optional, Sensitive) The password of the account that is used for the AD server.
* `port` - (Required) The port that is used to access the AD server.
* `server` - (Required) The address of the AD server.
* `standby_server` - (Optional) The address of the secondary AD server.

### `ldap_auth_server`

The ldap_auth_server supports the following:

* `account` - (Required) The username of the account that is used for the LDAP server.
* `base_dn` - (Required) The Base distinguished name (DN).
* `email_mapping` - (Optional) The field that is used to indicate the email address of a user on the LDAP server.
* `filter` - (Optional) The condition that is used to filter users.
* `is_ssl` - (Optional) Specifies whether to support SSL.
* `login_name_mapping` - (Optional) The field that is used to indicate the logon name of a user on the LDAP server.
* `mobile_mapping` - (Optional) The field that is used to indicate the mobile phone number of a user on the LDAP server.
* `name_mapping` - (Optional) The field that is used to indicate the name of a user on the LDAP server.
* `password` - (Optional, Sensitive) The password of the account that is used for the LDAP server.
* `port` - (Required) The port that is used to access the LDAP server.
* `server` - (Required) The address of the LDAP server.
* `standby_server` - (Optional) The address of the secondary LDAP server.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the instance resource of Bastionhost.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 40 min) Used when create the Instance.
* `update` - (Defaults to 20 min) Used when create the Instance.

## Import

Yundun_bastionhost instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_bastionhost_instance.example bastionhost-exampe123456
```
