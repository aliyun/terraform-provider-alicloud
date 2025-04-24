---
subcategory: "Aligreen"
layout: "alicloud"
page_title: "Alicloud: alicloud_aligreen_biz_type"
description: |-
  Provides a Alicloud Aligreen Biz Type resource.
---

# alicloud_aligreen_biz_type

Provides a Aligreen Biz Type resource.



For information about Aligreen Biz Type and how to use it, see [What is Biz Type](https://next.api.alibabacloud.com/document/Green/2017-08-23/CreateBizType).

-> **NOTE:** Available since v1.228.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_aligreen_biz_type&exampleId=92992bff-4b12-f3a8-0a00-e0f4187b741f9da10dc5&activeTab=example&spm=docs.r.aligreen_biz_type.0.92992bff4b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform_example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_aligreen_biz_type" "default" {
  biz_type_name   = var.name
  description     = var.name
  cite_template   = true
  industry_info   = "社交-注册信息-昵称"
  biz_type_import = "1"
}
```

## Argument Reference

The following arguments are supported:
* `biz_type_import` - (Optional) The name of the existing business scenario that was imported from when the business scenario was created.
* `biz_type_name` - (Required, ForceNew) The name of the business scenario defined by the customer. It can contain no more than 32 characters in English, numbers, and underscores.
* `cite_template` - (Optional, ForceNew) Specifies whether to import the configuration of an industry template. Default value: false. Valid values: true: imports the configuration of an industry template. false: does not import the configuration of an industry template. If the value is true, you must specify the industryInfo parameter.
* `description` - (Optional) The description of the business scenario defined by the customer, which is a combination of Chinese and English, numbers, and underscores, and cannot exceed 32 characters.
* `industry_info` - (Optional, ForceNew) The industry classification. Valid values: Social-Registration information-Profile picture Social-Registration information-Nickname Social-Registration information-Bio Social-Instant messaging-Chat Social-Instant messaging-Group chat Social-Instant messaging-Chat room Social-Forums&Communities-Post Social-Forums&Communities-Comment Social-Forums&Communities-Tag Social-Forums&Communities-Recommendation Multimedia-Registration information-Profile picture Multimedia-Registration information-Nickname Multimedia-Registration information-Bio Multimedia-Instant messaging-Chat Multimedia-Live streaming-Heading Multimedia-Live streaming-Cover Multimedia-Live streaming-Content Multimedia-Live streaming-Comment Multimedia-Online storage-Storage content Multimedia-Online storage-Shared content Gaming-Registration information-Nickname Gaming-Registration information-Profile picture Gaming-Registration information-Signature Gaming-Instant messaging-Chat Gaming-Instant messaging-Group chat Gaming-Instant messaging-Chat room Gaming-Forums&Communities-Post Gaming-Forums&Communities-Comment Gaming-Forums&Communities-Tag Gaming-Forums&Communities-Recommendation New retail-Goods-Heading New retail-Goods-Description Reading-Books-Title Reading-Books-Heading Reading-Books-Cover Reading-Books-Content Media-News content-News content Education-Registration information-Nickname Education-Registration information-Profile picture Education-Registration information-Bio Gaming-Instant messaging-Chat Gaming-Forums&Communities-Post Education-Forums&Communities-Comment Education-Forums&Communities-Tag Education-Forums&Communities-Recommendation Education-Customer service-Voice call Others

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Biz Type.
* `delete` - (Defaults to 5 mins) Used when delete the Biz Type.
* `update` - (Defaults to 5 mins) Used when update the Biz Type.

## Import

Aligreen Biz Type can be imported using the id, e.g.

```shell
$ terraform import alicloud_aligreen_biz_type.example <id>
```