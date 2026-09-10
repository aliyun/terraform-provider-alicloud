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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_threat_detection_monitor_account&exampleId=054ecd72-ec07-62a7-9040-60e1c58095adb7162aca&activeTab=example&spm=docs.r.threat_detection_monitor_account.0.054ecd72ec&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_threat_detection_monitor_account&spm=docs.r.threat_detection_monitor_account.example&intl_lang=EN_US)


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
