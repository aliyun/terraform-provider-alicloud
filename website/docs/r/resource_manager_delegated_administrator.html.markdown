---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_delegated_administrator"
sidebar_current: "docs-alicloud-resource-resource-manager-delegated-administrator"
description: |-
  Provides a Alicloud Resource Manager Delegated Administrator resource.
---

# alicloud_resource_manager_delegated_administrator

Provides a Resource Manager Delegated Administrator resource.

For information about Resource Manager Delegated Administrator and how to use it, see [What is Delegated Administrator](https://www.alibabacloud.com/help/en/resource-management/latest/registerdelegatedadministrator#doc-api-ResourceManager-RegisterDelegatedAdministrator).

-> **NOTE:** Available since v1.181.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_resource_manager_delegated_administrator&exampleId=18ddb1ce-0a79-b7d3-fd29-cb3683ea5cc56ee77a65&activeTab=example&spm=docs.r.resource_manager_delegated_administrator.0.18ddb1ce0a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
variable "display_name" {
  default = "EAccount"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

data "alicloud_resource_manager_folders" "example" {

}

resource "alicloud_resource_manager_account" "example" {
  display_name = "${var.display_name}-${random_integer.default.result}"
  folder_id    = data.alicloud_resource_manager_folders.example.ids.0
}

resource "alicloud_resource_manager_delegated_administrator" "example" {
  account_id        = alicloud_resource_manager_account.example.id
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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Delegated Administrator.
* `delete` - (Defaults to 1 mins) Used when deleting the Delegated Administrator.


## Import

Resource Manager Delegated Administrator can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_delegated_administrator.example <account_id>:<service_principal>
```