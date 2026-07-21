---
subcategory: "Elastic Desktop Service (ECD)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_user"
sidebar_current: "docs-alicloud-resource-ecd-user"
description: |-
  Provides a Alicloud Elastic Desktop Service (ECD) User resource.
---

# alicloud_ecd_user

Provides a Elastic Desktop Service (ECD) User resource.

For information about Elastic Desktop Service (ECD) User and how to use it, see [What is User](https://www.alibabacloud.com/help/en/wuying-workspace/developer-reference/api-eds-user-2021-03-08-createusers-desktop).

-> **NOTE:** Available since v1.142.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecd_user&exampleId=9d28c162-c6d0-8a2f-cfd5-465713d78c448e0fa22f&activeTab=example&spm=docs.r.ecd_user.0.9d28c162c6&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-shanghai"
}

resource "alicloud_ecd_user" "default" {
  end_user_id = "terraform_example123"
  email       = "tf.example@abc.com"
  phone       = "18888888888"
  password    = "Example_123"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ecd_user&spm=docs.r.ecd_user.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `email` - (Required, ForceNew) The email of the user email.
* `end_user_id` - (Required, ForceNew) The Username. The custom setting is composed of lowercase letters, numbers and underscores, and the length is 3~24 characters.
* `password` - (Optional, ForceNew) The password of the user password.
* `phone` - (Optional, ForceNew) The phone of the mobile phone number.
* `status` - (Optional, Computed) The status of the resource. Valid values: `Unlocked`, `Locked`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of User. Its value is same as `end_user_id`.

## Import

ECD User can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecd_user.example <end_user_id>
```
