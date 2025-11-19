---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_site_version"
description: |-
  Provides a Alicloud ESA Site Version resource.
---

# alicloud_esa_site_version

Provides a ESA Site Version resource.



For information about ESA Site Version and how to use it, see [What is Site Version](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CloneVersion).

-> **NOTE:** Available since v1.263.0.

## Example Usage

Basic Usage

没有资源测试用例，请先通过资源测试用例后再生成示例代码。

## Argument Reference

The following arguments are supported:
* `description` - (Optional) The Site version's description.
* `site_id` - (Required, ForceNew) The site ID, which can be obtained by calling the ListSites API.
* `site_version` - (Required, ForceNew, Int) The version number of the site configuration. For sites that have enabled configuration version management, this parameter can be used to specify the effective version of the configuration site, which defaults to version 0.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<site_version>`.
* `create_time` - The creation time. The date format follows ISO8601 notation and uses UTC time. The format is yyyy-MM-ddTHH:mm:ssZ.
* `status` - Site version status:：`online`.：`configuring`._faild`：`configure_faild`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Site Version.
* `delete` - (Defaults to 5 mins) Used when delete the Site Version.
* `update` - (Defaults to 5 mins) Used when update the Site Version.

## Import

ESA Site Version can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_site_version.example <site_id>:<site_version>
```