---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_monitor_account"
description: |-
  Provides a Alicloud Threat Detection Monitor Account resource.
---

# alicloud_threat_detection_monitor_account

Provides a Threat Detection Monitor Account resource.

Multi-account management account. The member accounts in the resource directory are added to the monitor account list of Threat Detection (Security Center), so that the security administrator account can manage the security of these member accounts in a unified manner.

For information about Threat Detection Monitor Account and how to use it, see [What is Monitor Account](https://next.api.alibabacloud.com/document/Sas/2018-12-03/CreateMonitorAccount).

-> **NOTE:** Available since v1.292.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "member_account_ids" {
  description = "The IDs of the member accounts in the resource directory. Multiple IDs must be separated by commas (,)."
  type        = string
}

resource "alicloud_threat_detection_monitor_account" "default" {
  account_ids = var.member_account_ids
}
```

### Deleting `alicloud_threat_detection_monitor_account` or removing it from your configuration

Terraform cannot destroy resource `alicloud_threat_detection_monitor_account`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `account_ids` - (Optional) The list of member account IDs in the resource directory. You can call [ListAccountsInResourceDirectory](https://next.api.alibabacloud.com/document/Sas/2018-12-03/ListAccountsInResourceDirectory) to obtain the member account IDs. Multiple member account IDs must be separated by commas (,). The monitor account list is fully replaced by the incoming list. If this parameter is not specified, the existing monitor account list is cleared.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<Alibaba Cloud Account ID>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Monitor Account.
* `update` - (Defaults to 5 mins) Used when update the Monitor Account.
* `delete` - (Defaults to 5 mins) Used when delete the Monitor Account.

## Import

Threat Detection Monitor Account can be imported using the id, e.g.

```shell
$ terraform import alicloud_threat_detection_monitor_account.example <Alibaba Cloud Account ID>
```
