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

-> **NOTE:** Available since v1.266.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_http_incoming_request_header_modification_rule&exampleId=9f98f5b0-fc72-2d6b-cb21-d1f62aefc659ad406222&activeTab=example&spm=docs.r.esa_http_incoming_request_header_modification_rule.0.9f98f5b0fc&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_esa_rate_plan_instance" "resource_HttpIncomingRequestHeaderModificationRule_example" {
  type         = "NS"
  auto_renew   = false
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  auto_pay     = true
  plan_name    = "high"
}

resource "alicloud_esa_site" "resource_Site_HttpIncomingRequestHeaderModificationRule_example" {
  site_name   = "gositecdn${random_integer.default.result}.cn"
  instance_id = alicloud_esa_rate_plan_instance.resource_HttpIncomingRequestHeaderModificationRule_example.id
  coverage    = "overseas"
  access_type = "NS"
}


resource "alicloud_esa_http_incoming_request_header_modification_rule" "default" {
  site_id      = alicloud_esa_site.resource_Site_HttpIncomingRequestHeaderModificationRule_example.id
  rule_enable  = "on"
  rule         = "(http.host eq \"video.example.com\")"
  sequence     = "1"
  site_version = "0"
  rule_name    = "example"
  request_header_modification {
    type      = "static"
    value     = "add"
    operation = "add"
    name      = "exampleadd"
  }
  request_header_modification {
    operation = "del"
    name      = "exampledel"
  }
  request_header_modification {
    type      = "dynamic"
    value     = "ip.geoip.country"
    operation = "modify"
    name      = "examplemodify"
  }
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_esa_http_incoming_request_header_modification_rule&spm=docs.r.esa_http_incoming_request_header_modification_rule.example&intl_lang=EN_US)

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