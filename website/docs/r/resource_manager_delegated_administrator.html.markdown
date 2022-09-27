---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_delegated_administrator"
sidebar_current: "docs-alicloud-resource-resource-manager-delegated-administrator"
description: |-
  Provides a Alicloud Resource Manager Delegated Administrator resource.
---

# alicloud\_resource\_manager\_delegated\_administrator

Provides a Resource Manager Delegated Administrator resource.

For information about Resource Manager Delegated Administrator and how to use it, see [What is Delegated Administrator](https://www.alibabacloud.com/help/en/resource-management/latest/registerdelegatedadministrator#doc-api-ResourceManager-RegisterDelegatedAdministrator).

-> **NOTE:** Available in v1.181.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_accounts" "default" {
  status = "CreateSuccess"
}
resource "alicloud_resource_manager_delegated_administrator" "default" {
  account_id        = data.alicloud_resource_manager_accounts.default.accounts.0.account_id
  service_principal = "cloudfw.aliyuncs.com"
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Required, ForceNew) The ID of the member account in the resource directory.
* `service_principal` - (Required, ForceNew) The identification of the trusted service. **NOTE:** Only some trusted services support delegated administrator accounts. For more information, see [Supported trusted services](https://www.alibabacloud.com/help/en/resource-management/latest/manage-trusted-services-overview).

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Delegated Administrator. The value formats as `<account_id>:<service_principal>`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Delegated Administrator.
* `delete` - (Defaults to 1 mins) Used when deleting the Delegated Administrator.


## Import

Resource Manager Delegated Administrator can be imported using the id, e.g.

```
$ terraform import alicloud_resource_manager_delegated_administrator.example <account_id>:<service_principal>
```