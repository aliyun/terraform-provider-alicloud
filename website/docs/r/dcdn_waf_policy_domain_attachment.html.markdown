---
subcategory: "DCDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_dcdn_waf_policy_domain_attachment"
sidebar_current: "docs-alicloud-resource-dcdn-waf-policy-domain-attachment"
description: |-
  Provides a Alicloud DCDN Waf Policy Domain Attachment resource.
---

# alicloud\_dcdn\_waf\_policy\_domain\_attachment

Provides a DCDN Waf Policy Domain Attachment resource.

For information about DCDN Waf Policy Domain Attachment and how to use it, see [What is Waf Policy Domain Attachment](https://www.alibabacloud.com/help/en/dynamic-route-for-cdn/latest/modify-the-domain-name-bound-to-a-protection-policies).

-> **NOTE:** Available in v1.186.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_dcdn_domain" "default" {
  domain_name = "example_domain_name"
  sources {
    content  = "1.1.1.1"
    port     = "80"
    priority = "20"
    type     = "ipaddr"
  }
}
resource "alicloud_dcdn_waf_domain" "default" {
  domain_name   = alicloud_dcdn_domain.default.domain_name
  client_ip_tag = "X-Forwarded-For"
}
resource "alicloud_dcdn_waf_policy" "default" {
  policy_type   = "custom"
  policy_name   = "example_value"
  defense_scene = "waf_group"
  status        = "on"
}
resource "alicloud_dcdn_waf_policy_domain_attachment" "example" {
  domain_name = alicloud_dcdn_waf_domain.default.domain_name
  policy_id   = alicloud_dcdn_waf_policy.default.id
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, ForceNew) Access the accelerated domain name of the specified protection policy.
* `policy_id` - (Required, ForceNew) The protection policy ID. Only one input is supported.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Waf Policy Domain Attachment. The value is formulated as `<policy_id>:<domain_name>`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Waf Policy Domain Attachment.
* `delete` - (Defaults to 1 mins) Used when delete the Waf Policy Domain Attachment.

## Import

DCDN Waf Policy Domain Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_dcdn_waf_policy_domain_attachment.example policy_id:domain_name
```