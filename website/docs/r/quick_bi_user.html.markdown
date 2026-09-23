---
subcategory: "Quick BI"
layout: "alicloud"
page_title: "Alicloud: alicloud_quick_bi_user"
sidebar_current: "docs-alicloud-resource-quick-bi-user"
description: |-
  Provides a Alicloud Quick BI User resource.
---

# alicloud\_quick\_bi\_user

Provides a Quick BI User resource.

For information about Quick BI User and how to use it, see [What is User](https://www.alibabacloud.com/help/doc-detail/33813.htm).

-> **NOTE:** Available in v1.136.0+.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_quick_bi_user&exampleId=ea3132ca-f214-fe0a-3d36-89e29b8f913cb39c601d&activeTab=example&spm=docs.r.quick_bi_user.0.ea3132caf2&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_quick_bi_user" "example" {
  account_name    = "example_value"
  admin_user      = false
  auth_admin_user = false
  nick_name       = "example_value"
  user_type       = "Analyst"
}

```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_quick_bi_user&spm=docs.r.quick_bi_user.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional, ForceNew) Alibaba Cloud account ID.
* `account_name` - (Required) An Alibaba Cloud account, Alibaba Cloud name.
* `admin_user` - (Required) Whether it is the administrator. Valid values: `true` and `false`.
* `auth_admin_user` - (Required) Whether this is a permissions administrator. Valid values: `false`, `true`.
* `nick_name` - (Required, ForceNew) The nickname of the user.
* `user_type` - (Required) The members of the organization of the type of role separately. Valid values: `Analyst`, `Developer` and `Visitor`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of User.

## Import

Quick BI User can be imported using the id, e.g.

```shell
$ terraform import alicloud_quick_bi_user.example <id>
```
