---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_http_incoming_request_header_modification_rule"
description: |-
  Provides a Alicloud ESA Http Incoming Request Header Modification Rule resource.
---

# alicloud_esa_http_incoming_request_header_modification_rule

Provides a ESA Http Incoming Request Header Modification Rule resource.



For information about ESA Http Incoming Request Header Modification Rule and how to use it, see [What is Http Incoming Request Header Modification Rule](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateHttpIncomingRequestHeaderModificationRule).

-> **NOTE:** Available since v1.263.0.

## Example Usage

Basic Usage

没有资源测试用例，请先通过资源测试用例后再生成示例代码。

## Argument Reference

The following arguments are supported:
* `request_header_modification` - (Required, List) The configurations of modifying request headers. You can add, delete, or modify a request header. See [`request_header_modification`](#request_header_modification) below.
* `rule` - (Optional) Rule content, using conditional expressions to match user requests. When adding global configuration, this parameter does not need to be set. There are two usage scenarios:
  - Match all incoming requests: value set to true
  - Match specified request: Set the value to a custom expression, for example: (http.host eq \"video.example.com\")
* `rule_enable` - (Optional) Rule switch. When adding global configuration, this parameter does not need to be set. Value range:
  - on: open.
  - off: close.
* `rule_name` - (Optional) Rule name. When adding global configuration, this parameter does not need to be set.
* `sequence` - (Optional, Int) Order of rule execution. The smaller the value, the higher the priority for execution.
* `site_id` - (Required, ForceNew) The site ID.
* `site_version` - (Optional, ForceNew, Int) The version number of the site configuration. For sites that have enabled configuration version management, this parameter can be used to specify the effective version of the configuration site, which defaults to version 0.

### `request_header_modification`

The request_header_modification supports the following:
* `name` - (Required) Request Header Name.
* `operation` - (Required) Mode of operation. Value range:
  - `add`: add.
  - `del`: delete
  - `modify`: change.
* `type` - (Optional) Value type. Value range:
  - `static`:static mode.
  - `dynamic`:dynamic mode.
* `value` - (Optional) Request header value

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<config_id>`.
* `config_id` - Config Id

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Http Incoming Request Header Modification Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Http Incoming Request Header Modification Rule.
* `update` - (Defaults to 5 mins) Used when update the Http Incoming Request Header Modification Rule.

## Import

ESA Http Incoming Request Header Modification Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_http_incoming_request_header_modification_rule.example <site_id>:<config_id>
```