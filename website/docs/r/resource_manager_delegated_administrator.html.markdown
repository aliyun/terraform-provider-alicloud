---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_delegated_administrator"
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
  <a href="https://api.aliyun.com/terraform?resource=alicloud_resource_manager_delegated_administrator&exampleId=549a0973-ff56-05e0-b333-4321e882f95aeeb26b1a&activeTab=example&spm=docs.r.resource_manager_delegated_administrator.0.549a0973ff&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_resource_manager_accounts" "default" {
  status = "CreateSuccess"
}

resource "alicloud_resource_manager_delegated_administrator" "default" {
  account_id        = data.alicloud_resource_manager_accounts.default.accounts.0.account_id
  service_principal = "cloudfw.aliyuncs.com"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_resource_manager_delegated_administrator&spm=docs.r.resource_manager_delegated_administrator.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `account_id` - (Required, ForceNew) The Alibaba Cloud account ID of the member in the resource directory.
* `service_principal` - (Required, ForceNew) The identifier of the trusted service.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<account_id>:<service_principal>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Delegated Administrator.
* `delete` - (Defaults to 5 mins) Used when delete the Delegated Administrator.


## Import

Resource Manager Delegated Administrator can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_delegated_administrator.example <account_id>:<service_principal>
```
