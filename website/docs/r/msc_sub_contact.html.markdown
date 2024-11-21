---
subcategory: "Message Center (MscSub)"
layout: "alicloud"
page_title: "Alicloud: alicloud_msc_sub_contact"
sidebar_current: "docs-alicloud-resource-msc-sub-contact"
description: |-
  Provides a Alicloud Message Center Contact resource.
---

# alicloud_msc_sub_contact

Provides a Msc Sub Contact resource.

-> **NOTE:** Available since v1.132.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_msc_sub_contact&exampleId=a6bdae8c-2818-8947-a1cc-9ab53f9af0086c428fef&activeTab=example&spm=docs.r.msc_sub_contact.0.a6bdae8c28&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tfexample"
}

resource "alicloud_msc_sub_contact" "default" {
  contact_name = var.name
  position     = "CEO"
  email        = "123@163.com"
  mobile       = "15388888888"
}
```

## Argument Reference

The following arguments are supported:

* `contact_name` - (Required) The User's Contact Name. **Note:** The name must be 2 to 12 characters in length.
* `email` - (Required) The User's Contact Email Address.
* `mobile` - (Required) The User's Telephone.
* `position` - (Required, ForceNew) The User's Position. Valid values: `CEO`, `Technical Director`, `Maintenance Director`, `Project Director`,`Finance Director` and `Other`.

-> **NOTE:** When the user creates a contact, the user should use `alicloud_msc_sub_contact_verification_message` to receive the verification message and confirm it.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Contact.

## Import

Msc Sub Contact can be imported using the id, e.g.

```shell
$ terraform import alicloud_msc_sub_contact.example <id>
```
