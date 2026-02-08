---
subcategory: "Anti-DDoS Pro (DdosCoo)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddoscoo_domain_precise_access_rule"
description: |-
  Provides a Alicloud DdosCoo Domain Precise Access Rule resource.
---

# alicloud_ddoscoo_domain_precise_access_rule

Provides a DdosCoo Domain Precise Access Rule resource.

Precise access control rules for website business.

For information about DdosCoo Domain Precise Access Rule and how to use it, see [What is Domain Precise Access Rule](https://next.api.alibabacloud.com/document/ddoscoo/2020-01-01/ModifyWebPreciseAccessRule).

-> **NOTE:** Available since v1.271.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "terraform"
}

variable "domain" {
  default = "terraform-example.alibaba.com"
}

data "alicloud_ddoscoo_instances" "default" {
}

resource "alicloud_ddoscoo_domain_resource" "default" {
  domain       = var.domain
  instance_ids = [data.alicloud_ddoscoo_instances.default.ids.0]
  proxy_types {
    proxy_ports = [443]
    proxy_type  = "https"
  }
  real_servers = ["177.167.32.11"]
  rs_type      = 0
}

resource "alicloud_ddoscoo_domain_precise_access_rule" "default" {
  condition {
    match_method = "contain"
    field        = "header"
    content      = "222"
    header_name  = "15"
  }
  action  = "accept"
  expires = "0"
  domain  = alicloud_ddoscoo_domain_resource.default.id
  name    = var.name
}
```

## Argument Reference

The following arguments are supported:
* `action` - (Required) Action to take on match. Valid values:
  - `accept`: Allow.
  - `block`: Block.
  - `challenge`: Challenge
* `condition` - (Required, List) List of matching conditions. See [`condition`](#condition) below.
* `domain` - (Required, ForceNew) Domain name of the website service.
-> **NOTE:**  The domain name must already have a website service forwarding rule configured. You can call [DescribeDomains](https://help.aliyun.com/document_detail/91724.html) to query all domain names.
* `expires` - (Optional, Int) Rule validity period, in seconds. This parameter takes effect only when the rule's matching action is set to block (`action` is `block`), blocking access requests during the validity period. If this parameter is not specified, the rule remains effective permanently.
* `name` - (Required, ForceNew) Rule name.

### `condition`

The condition supports the following:
* `content` - (Required) Matching content.
* `field` - (Required) Matching field.
* `header_name` - (Optional) Custom HTTP header field name.

-> **NOTE:**  Valid only when `Field` is `header`.

* `match_method` - (Required) Matching method.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<domain>:<name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Domain Precise Access Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Domain Precise Access Rule.
* `update` - (Defaults to 5 mins) Used when update the Domain Precise Access Rule.

## Import

DdosCoo Domain Precise Access Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_ddoscoo_domain_precise_access_rule.example <domain>:<name>
```
