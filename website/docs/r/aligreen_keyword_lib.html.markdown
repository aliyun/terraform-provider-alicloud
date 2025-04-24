---
subcategory: "Aligreen"
layout: "alicloud"
page_title: "Alicloud: alicloud_aligreen_keyword_lib"
description: |-
  Provides a Alicloud Aligreen Keyword Lib resource.
---

# alicloud_aligreen_keyword_lib

Provides a Aligreen Keyword Lib resource.

Keyword library for text detection.

For information about Aligreen Keyword Lib and how to use it, see [What is Keyword Lib](https://next.api.alibabacloud.com/document/Green/2017-08-23/CreateKeywordLib).

-> **NOTE:** Available since v1.228.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_aligreen_keyword_lib&exampleId=a99f0cbe-b658-f0ac-406a-0a59728d2988411187b1&activeTab=example&spm=docs.r.aligreen_keyword_lib.0.a99f0cbeb6&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_aligreen_biz_type" "defaultMn8sVK" {
  biz_type_name = "${var.name}${random_integer.default.result}"
  cite_template = true
  industry_info = "社交-注册信息-昵称"
}


resource "alicloud_aligreen_keyword_lib" "default" {
  category         = "BLACK"
  resource_type    = "TEXT"
  lib_type         = "textKeyword"
  keyword_lib_name = var.name
  match_mode       = "fuzzy"
  language         = "cn"
  biz_types        = ["example_007"]
  lang             = "cn"
  enable           = true
}
```

## Argument Reference

The following arguments are supported:
* `biz_types` - (Optional, ForceNew) The business scenario. Example:["bizTypeA","bizTypeB"]
* `category` - (Optional, ForceNew) The category of the text library. Valid values: BLACK: a blacklist. WHITE: a whitelist. REVIEW: a review list
* `enable` - (Optional, ForceNew) Specifies whether to enable text library.true: Enable the text library. This is the default value.false: Disable the text library.
* `keyword_lib_name` - (Required) The name of the keyword library defined by the customer. It can contain no more than 20 characters in Chinese, English, and underscore (_).
* `lang` - (Optional) Language.
* `language` - (Optional, ForceNew) Language used by the text Library
* `lib_type` - (Optional, ForceNew) The category of the text library in each moderation scenario. Valid values: textKeyword: a text library against which terms in text are matched. similarText: a text library against which text patterns are matched. textKeyword: a text library against which terms extracted from images are matched. voiceText: a text library against which terms converted from audio are matched.
* `match_mode` - (Optional, ForceNew) The matching method. Valid values:fuzzy: fuzzy match precise: exact match
* `resource_type` - (Required, ForceNew) The moderation scenario to which the text library applies. Valid values:TEXT: text anti-spam、IMAGE: ad violation detection、VOICE: audio anti-spam

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Keyword Lib.
* `delete` - (Defaults to 5 mins) Used when delete the Keyword Lib.
* `update` - (Defaults to 5 mins) Used when update the Keyword Lib.

## Import

Aligreen Keyword Lib can be imported using the id, e.g.

```shell
$ terraform import alicloud_aligreen_keyword_lib.example <id>
```