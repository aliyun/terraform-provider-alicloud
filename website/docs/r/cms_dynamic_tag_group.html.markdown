---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_dynamic_tag_group"
sidebar_current: "docs-alicloud-resource-cms-dynamic-tag-group"
description: |-
  Provides a Alicloud Cloud Monitor Service Dynamic Tag Group resource.
---

# alicloud_cms_dynamic_tag_group

Provides a Cloud Monitor Service Dynamic Tag Group resource.

For information about Cloud Monitor Service Dynamic Tag Group and how to use it, see [What is Dynamic Tag Group](https://www.alibabacloud.com/help/en/cloudmonitor/latest/createdynamictaggroup).

-> **NOTE:** Available since v1.142.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cms_dynamic_tag_group&exampleId=a79c7aa3-0bc8-4e94-fab4-8976e8c52dead2436476&activeTab=example&spm=docs.r.cms_dynamic_tag_group.0.a79c7aa30b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_cms_alarm_contact_group" "default" {
  alarm_contact_group_name = var.name
}

resource "alicloud_cms_dynamic_tag_group" "default" {
  tag_key            = var.name
  contact_group_list = [alicloud_cms_alarm_contact_group.default.id]
  match_express {
    tag_value                = var.name
    tag_value_match_function = "all"
  }
}
```

## Argument Reference

The following arguments are supported:

* `tag_key` - (Required, ForceNew) The tag keys of the cloud resources.
* `match_express_filter_relation` - (Optional, ForceNew) The relationship between the conditional expressions for the tag values of the cloud resources. Valid values: `and`, `or`.
* `contact_group_list` - (Required, ForceNew, List) The alert contact groups. The alert notifications of the application group are sent to the alert contacts that belong to the specified alert contact groups.
* `template_id_list` - (Optional, ForceNew, List) The IDs of the alert templates.
* `match_express` - (Required, ForceNew, Set) The conditional expressions used to create an application group based on the tag. See [`match_express`](#match_express) below.

### `match_express`

The match_express supports the following: 

* `tag_value` - (Required, ForceNew) The tag values of the cloud resources.
* `tag_value_match_function` - (Required, ForceNew) The method that is used to match the tag values of the cloud resources. Valid values: `all`, `startWith`, `endWith`, `contains`, `notContains`, `equals`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Dynamic Tag Group.
* `status` - The status of the Dynamic Tag Group.

## Import

Cloud Monitor Service Dynamic Tag Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_dynamic_tag_group.example <id>
```
