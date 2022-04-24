---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_instance"
sidebar_current: "docs-alicloud-resource-bastionhost-instance"
description: |-
  Provides a Alicloud Bastion Host Instance Resource.
---

# alicloud_bastionhost_instance

-> **NOTE:** From the version 1.132.0, the resource has been renamed to `alicloud_bastionhost_instance`.

Cloud Bastion Host instance resource ("Yundun_bastionhost" is the short term of this product). 
For information about Resource Manager Resource Directory and how to use it, see [What is Bastionhost](https://www.alibabacloud.com/help/en/doc-detail/52922.htm).

-> **NOTE:** The endpoint of bssopenapi used only support "business.aliyuncs.com" at present.

-> **NOTE:** Available in 1.63.0+ .

-> **NOTE:** In order to destroy Cloud Bastionhost instance , users are required to apply for white list first

## Example Usage

Basic Usage

```terraform
resource "alicloud_bastionhost_instance" "default" {
  description        = "Terraform-test"
  license_code       = "bhah_ent_50_asset"
  period             = "1"
  vswitch_id         = "v-testVswitch"
  security_group_ids = ["sg-test", "sg-12345"]
}
```

```terraform
data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}
resource "alicloud_bastionhost_instance" "default" {
  description        = "Terraform-test"
  license_code       = "bhah_ent_50_asset"
  period             = 1
  security_group_ids = [alicloud_security_group.default.0.id, alicloud_security_group.default.1.id]
  vswitch_id         = "v-testVswitch"
  resource_group_id  = data.alicloud_resource_manager_resource_groups.default.ids.0
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

  lifecycle {
    ignore_changes = [ldap_auth_server.0.password, ad_auth_server.0.password]
  }
}
```

### Deleting `alicloud_bastionhost_instance` or removing it from your configuration

The `alicloud_bastionhost_instance` resource allows you to manage bastionhost instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration
will remove it from your state file and management, but will not destroy the bastionhost instance.
You can resume managing the subscription bastionhost instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:

* `license_code` - (Required)  The package type of Cloud Bastionhost instance. You can query more supported types through the [DescribePricingModule](https://help.aliyun.com/document_detail/96469.html).
* `description` - (Required) Description of the instance. This name can have a string of 1 to 63 characters.
* `period` - (Optional) Duration for initially producing the instance. Valid values: [1~9], 12, 24, 36. At present, the provider does not support modify "period".
-> **NOTE:** The attribute `period` is only used to create Subscription instance or modify the PayAsYouGo instance to Subscription. Once effect, it will not be modified that means running `terraform apply` will not effect the resource.
* `vswitch_id` - (Required, ForceNew) VSwitch ID configured to Bastionhost.
* `security_group_ids` - (Required) security group IDs configured to Bastionhost. 
  **NOTE:** There is a potential diff error because of the order of `security_group_ids` values indefinite.
  So, from version 1.160.0, `security_group_ids` type has been updated as `set` from `list`,
  and you can use [tolist](https://www.terraform.io/language/functions/tolist) to convert it to a list.
* `tags` - (Optional, Available in v1.67.0+) A mapping of tags to assign to the resource.
* `resource_group_id` - (Optional, Available in v1.87.0+) The Id of resource group which the Bastionhost Instance belongs. If not set, the resource is created in the default resource group.
* `enable_public_access` - (Optional, Available in v1.143.0+)  Whether to Enable the public internet access to a specified Bastionhost instance. The valid values: `true`, `false`.
* `ad_auth_server` - (Optional, Available from 1.169.0+) The AD auth server of the Instance. See the following `Block ad_auth_server`.
* `ldap_auth_server` - (Optional, Available from 1.169.0+) The LDAP auth server of the Instance. See the following `Block ldap_auth_server`.
  
-> **NOTE:** You can utilize the generic Terraform resource [lifecycle configuration block](https://www.terraform.io/docs/configuration/resources.html) with `ad_auth_server` or `ldap_auth_server` to configure auth server, then ignore any changes to that `password` caused externally (e.g. Application Autoscaling).
```
  # ... ignore the change about ad_auth_server.0.password and ldap_auth_server.0.password in alicloud_bastionhost_instance
  lifecycle {
    ignore_changes = [ad_auth_server.0.password,ldap_auth_server.0.password]
  }
```

#### Block ad_auth_server

The ad_auth_server supports the following:

* `account` - (Required) The username of the account that is used for the AD server.
* `base_dn` - (Required) The Base distinguished name (DN).
* `domain` - (Required) The domain on the AD server.
* `email_mapping` - (Optional) The field that is used to indicate the email address of a user on the AD server.
* `filter` - (Optional) The condition that is used to filter users.
* `instance_id` - (Required, ForceNew) The ID of the Bastion machine instance.
* `is_ssl` - (Required) Specifies whether to support SSL.
* `mobile_mapping` - (Optional) The field that is used to indicate the mobile phone number of a user on the AD server.
* `name_mapping` - (Optional) The field that is used to indicate the name of a user on the AD server.
* `password` - (Required, Sensitive) The password of the account that is used for the AD server.
* `port` - (Required) The port that is used to access the AD server.
* `server` - (Required) The address of the AD server.
* `standby_server` - (Optional) The address of the secondary AD server.

#### Block ldap_auth_server

The ldap_auth_server supports the following:

* `account` - (Required) The username of the account that is used for the LDAP server.
* `base_dn` - (Required) The Base distinguished name (DN).
* `email_mapping` - (Optional) The field that is used to indicate the email address of a user on the LDAP server.
* `filter` - (Optional) The condition that is used to filter users.
* `instance_id` - (Required, ForceNew) The ID of the Bastion machine instance.
* `is_ssl` - (Optional) Specifies whether to support SSL.
* `login_name_mapping` - (Optional) The field that is used to indicate the logon name of a user on the LDAP server.
* `mobile_mapping` - (Optional) The field that is used to indicate the mobile phone number of a user on the LDAP server.
* `name_mapping` - (Optional) The field that is used to indicate the name of a user on the LDAP server.
* `password` - (Required, Sensitive) The password of the account that is used for the LDAP server.
* `port` - (Required) The port that is used to access the LDAP server.
* `server` - (Required) The address of the LDAP server.
* `standby_server` - (Optional) The address of the secondary LDAP server.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the instance resource of Bastionhost.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 40 min) Used when create the Instance.
* `update` - (Defaults to 20 min) Used when create the Instance.

## Import

Yundun_bastionhost instance can be imported using the id, e.g.

```
$ terraform import alicloud_bastionhost_instance.example bastionhost-exampe123456
```
