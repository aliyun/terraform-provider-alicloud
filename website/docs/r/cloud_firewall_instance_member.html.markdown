---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_instance_member"
sidebar_current: "docs-alicloud-resource-cloud_firewall-instance-member"
description: |-
  Provides a Alicloud Cloud Firewall Instance Member resource.
---

# alicloud_cloud_firewall_instance_member

Provides a Cloud Firewall Instance Member resource.

For information about Cloud Firewall Instance Member and how to use it, see [What is Instance Member](https://www.alibabacloud.com/help/en/cloud-firewall/cloudfirewall/developer-reference/api-cloudfw-2017-12-07-addinstancemembers).

-> **NOTE:** Available since v1.194.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_firewall_instance_member&exampleId=3548ac26-5b3c-ab03-0909-25f42801a7ce7bd210cb&activeTab=example&spm=docs.r.cloud_firewall_instance_member.0.3548ac265b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "AliyunTerraform"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_resource_manager_account" "default" {
  display_name = "${var.name}-${random_integer.default.result}"
  timeouts {
    delete = "5m"
  }
}

resource "alicloud_cloud_firewall_instance_member" "default" {
  member_desc = "${var.name}-${random_integer.default.result}"
  member_uid  = alicloud_resource_manager_account.default.id
}
```

## Argument Reference

The following arguments are supported:
* `member_desc` - (Optional) Remarks of cloud firewall member accounts.
* `member_uid` - (Required, ForceNew) The UID of the cloud firewall member account.

## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above. Its value same as `member_uid`.
* `create_time` - When the cloud firewall member account was added.> use second-level timestamp format.
* `member_display_name` - The name of the cloud firewall member account.
* `modify_time` - The last modification time of the cloud firewall member account.> use second-level timestamp format.
* `status` - The resource attribute field that represents the resource status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Instance Member.
* `delete` - (Defaults to 1 mins) Used when delete the Instance Member.
* `update` - (Defaults to 5 mins) Used when update the Instance Member.

## Import

Cloud Firewall Instance Member can be imported using the id, e.g.

```shell
$terraform import alicloud_cloud_firewall_instance_member.example <id>
```