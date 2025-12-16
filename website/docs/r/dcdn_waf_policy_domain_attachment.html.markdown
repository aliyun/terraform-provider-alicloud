---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_waf_policy_domain_attachment"
sidebar_current: "docs-alicloud-resource-dcdn-waf-policy-domain-attachment"
description: |-
  Provides a Alicloud DCDN Waf Policy Domain Attachment resource.
---

# alicloud_dcdn_waf_policy_domain_attachment

Provides a DCDN Waf Policy Domain Attachment resource.

For information about DCDN Waf Policy Domain Attachment and how to use it, see [What is Waf Policy Domain Attachment](https://www.alibabacloud.com/help/en/dcdn/developer-reference/api-dcdn-2018-01-15-modifydcdnwafpolicydomains).

-> **NOTE:** Available since v1.186.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dcdn_waf_policy_domain_attachment&exampleId=2c02cf02-a544-d842-c633-f8e4f31a84c4a1eedca5&activeTab=example&spm=docs.r.dcdn_waf_policy_domain_attachment.0.2c02cf02a5&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "domain_name" {
  default = "tf-example.com"
}

variable "name" {
  default = "tf_example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_dcdn_domain" "example" {
  domain_name = "${var.domain_name}-${random_integer.default.result}"
  scope       = "overseas"
  sources {
    content  = "1.1.1.1"
    port     = "80"
    priority = "20"
    type     = "ipaddr"
    weight   = "10"
  }
}

resource "alicloud_dcdn_waf_domain" "example" {
  domain_name   = alicloud_dcdn_domain.example.domain_name
  client_ip_tag = "X-Forwarded-For"
}

resource "alicloud_dcdn_waf_policy" "example" {
  defense_scene = "waf_group"
  policy_name   = "${var.name}_${random_integer.default.result}"
  policy_type   = "custom"
  status        = "on"
}

resource "alicloud_dcdn_waf_policy_domain_attachment" "example" {
  domain_name = alicloud_dcdn_waf_domain.example.domain_name
  policy_id   = alicloud_dcdn_waf_policy.example.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_dcdn_waf_policy_domain_attachment&spm=docs.r.dcdn_waf_policy_domain_attachment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, ForceNew) Access the accelerated domain name of the specified protection policy.
* `policy_id` - (Required, ForceNew) The protection policy ID. Only one input is supported.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Waf Policy Domain Attachment. The value is formulated as `<policy_id>:<domain_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Waf Policy Domain Attachment.
* `delete` - (Defaults to 1 mins) Used when delete the Waf Policy Domain Attachment.

## Import

DCDN Waf Policy Domain Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_dcdn_waf_policy_domain_attachment.example policy_id:domain_name
```