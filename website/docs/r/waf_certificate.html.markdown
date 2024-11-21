---
subcategory: "Web Application Firewall(WAF)"
layout: "alicloud"
page_title: "Alicloud: alicloud_waf_certificate"
sidebar_current: "docs-alicloud-resource-waf-certificate"
description: |-
  Provides a Web Application Firewall Certificate resource.
---

# alicloud\_waf\_certificate

Provides a WAF Certificate resource.

For information about WAF Certificate and how to use it, see [What is Certificate](https://www.alibabacloud.com/help/doc-detail/28517.htm).

-> **NOTE:** Available in v1.135.0+.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_waf_certificate&exampleId=8fa3caa5-7920-7d0c-8dc0-50b4d30b1e9edcd44744&activeTab=example&spm=docs.r.waf_certificate.0.8fa3caa579&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_waf_certificate" "default" {
  certificate_name = "your_certificate_name"
  instance_id      = "your_instance_id"
  domain           = "your_domain_name"
  private_key      = "your_private_key"
  certificate      = "your_certificate"
}
resource "alicloud_waf_certificate" "default2" {
  instance_id    = "your_instance_id"
  domain         = "your_domain_name"
  certificate_id = "your_certificate_id"
}
```

## Argument Reference

The following arguments are supported:

* `certificate` - (Optional, ForceNew, Conflicts with `certificate_id`) Certificate file content.
* `certificate_name` - (Optional, ForceNew, Conflicts with `certificate_id`) Certificate file name.
* `instance_id` - (Required, ForceNew) The ID of the WAF instance.
* `domain` - (Required, ForceNew) The domain that you want to add to WAF.
* `private_key` - (Optional, ForceNew, Conflicts with `certificate_id`) The private key.
* `certificate_id` - (Optional, ForceNew, Conflicts with `certificate`, `certificate_name`,`private_key`) The certificate id is automatically generated when you upload your certificate content.**NOTE:** you can also use Certificate ID saved in SSL.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Certificate. The value formats as `<instance_id>:<domain>:<certificate_id>`.

## Import

WAF Certificate can be imported using the id, e.g.

```shell
$ terraform import alicloud_waf_certificate.example <instance_id>:<domain>:<certificate_id>
```
