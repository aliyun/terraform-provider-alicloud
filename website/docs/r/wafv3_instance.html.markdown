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

-> **NOTE:** Available in v1.200.0+.

## Example Usage
<div class="oics-button" style="float: right;margin: 0 0 -40px 0;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_wafv3_instance&exampleId=68c7e1d5-968d-4a2d-19dd-3448f9b9e8fe6de0327f&activeTab=example&spm=docs.r.wafv3_instance.0.68c7e1d596" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>

Basic Usage

```terraform
resource "alicloud_wafv3_instance" "default" {
}
```

## Argument Reference

The following arguments are supported:


## Attributes Reference

The following attributes are exported:
* `id` - The `key` of the resource supplied above.
* `create_time` - The creation time of the resource
* `instance_id` - The first ID of the resource
* `status` - The status of the resource

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 1 mins) Used when create the Instance.
* `delete` - (Defaults to 1 mins) Used when delete the Instance.

## Import

Wafv3 Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_wafv3_instance.example <id>
```