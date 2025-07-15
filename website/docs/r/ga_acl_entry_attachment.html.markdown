---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_acl_entry_attachment"
sidebar_current: "docs-alicloud-resource-ga-acl-entry-attachment"
description: |-
  Provides a Alicloud Global Accelerator (GA) Acl Entry Attachment resource.
---

# alicloud_ga_acl_entry_attachment

Provides a Global Accelerator (GA) Acl Entry Attachment resource.

For information about Global Accelerator (GA) Acl Entry Attachment and how to use it, see [What is Acl Entry Attachment](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-addentriestoacl).

-> **NOTE:** Available since v1.190.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ga_acl_entry_attachment&exampleId=100dd0e6-5005-4456-de38-7ec339c60e19ff0f9a6c&activeTab=example&spm=docs.r.ga_acl_entry_attachment.0.100dd0e650&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_ga_acl" "default" {
  address_ip_version = "IPv4"
  acl_name           = var.name
}

resource "alicloud_ga_acl_entry_attachment" "default" {
  acl_id            = alicloud_ga_acl.default.id
  entry             = "192.168.1.1/32"
  entry_description = var.name
}
```

## Argument Reference

The following arguments are supported:

* `acl_id` - (Required, ForceNew) The ID of the Acl.
* `entry` - (Required, ForceNew) The entry (IP address or CIDR block) that you want to add.
* `entry_description` - (Optional, ForceNew) The description of the entry. The description must be `1` to `256` characters in length, and can contain letters, digits, hyphens (-), forward slashes (/), periods (.), and underscores (_).

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Acl Entry Attachment. It formats as `<acl_id>:<entry>`.
* `status` - The status of the Acl Entry Attachment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Acl Entry Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Acl Entry Attachment.

## Import

Global Accelerator (GA) Acl Entry Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_acl_entry_attachment.example <acl_id>:<entry>
```
