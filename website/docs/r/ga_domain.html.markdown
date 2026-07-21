---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_domain"
sidebar_current: "docs-alicloud-resource-ga-domain"
description: |-
  Provides a Alicloud Ga Domain resource.
---

# alicloud_ga_domain

Provides a Ga Domain resource.

For information about Ga Domain and how to use it, see [What is Domain](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-createdomain).

-> **NOTE:** Available since v1.197.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ga_domain&exampleId=46ed1550-9e68-d51d-36ef-a4fcfc9b8d6a42b2e902&activeTab=example&spm=docs.r.ga_domain.0.46ed15509e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_ga_accelerator" "default" {
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}

resource "alicloud_ga_domain" "default" {
  domain         = "changes.com.cn"
  accelerator_id = alicloud_ga_accelerator.default.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ga_domain&spm=docs.r.ga_domain.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `accelerator_id` - (Required, ForceNew) The ID of the global acceleration instance.
* `domain` - (Required, ForceNew) The accelerated domain name to be added. only top-level domain names are supported, such as 'example.com'.

## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above. The value is formulated as `<accelerator_id>:<domain>`.
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Domain.
* `delete` - (Defaults to 5 mins) Used when delete the Domain.

## Import

Ga Domain can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_domain.example <accelerator_id>:<domain>
```