---
subcategory: "Web Application Firewall(WAF)"
layout: "alicloud"
page_title: "Alicloud: alicloud_wafv3_instance"
sidebar_current: "docs-alicloud-resource-wafv3-instance"
description: |-
  Provides a Alicloud Wafv3 Instance resource.
---

# alicloud_wafv3_instance

Provides a Wafv3 Instance resource.

For information about Wafv3 Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/en/web-application-firewall/latest/what-is-waf).

-> **NOTE:** Available since v1.200.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_wafv3_instance&exampleId=634b4bc3-5bec-f3c1-d5d0-e4cb4c5a71d4f0ec29f8&activeTab=example&spm=docs.r.wafv3_instance.0.634b4bc35b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  alias = "singapore"
  # when you want to create in specific region, you can set the region here.
  region = "ap-southeast-1"
}

resource "alicloud_wafv3_instance" "example" {
  provider = alicloud.singapore
}
```


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `create_time` - The creation time of the resource
* `instance_id` - The first ID of the resource
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 1 mins) Used when create the Instance.
* `delete` - (Defaults to 1 mins) Used when delete the Instance.

## Import

Wafv3 Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_wafv3_instance.example <id>
```