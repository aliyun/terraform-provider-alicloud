---
subcategory: "Web Application Firewall(WAF)"
layout: "alicloud"
page_title: "Alicloud: alicloud_waf_protection_module"
sidebar_current: "docs-alicloud-resource-waf-protection-module"
description: |-
  Provides a Alicloud Web Application Firewall(WAF) Protection Module resource.
---

# alicloud\_waf\_protection\_module

Provides a Web Application Firewall(WAF) Protection Module resource.

For information about Web Application Firewall(WAF) Protection Module and how to use it, see [What is Protection Module](https://www.alibabacloud.com/help/en/doc-detail/160775.htm).

-> **NOTE:** Available in v1.141.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_waf_instances" "default" {
}

resource "alicloud_waf_domain" "default" {
  domain_name       = "you domain"
  instance_id       = data.alicloud_waf_instances.default.ids.0
  is_access_product = "On"
  source_ips        = ["1.1.1.1"]
  cluster_type      = "PhysicalCluster"
  http2_port        = [443]
  http_port         = [80]
  https_port        = [443]
  http_to_user_ip   = "Off"
  https_redirect    = "Off"
  load_balancing    = "IpHash"
  log_headers {
    key   = "foo"
    value = "http"
  }
}

resource "alicloud_waf_protection_module" "default" {
  instance_id  = data.alicloud_waf_instances.default.ids.0
  domain       = alicloud_waf_domain.default.domain_name
  defense_type = "ac_cc"
  mode         = 0
  status       = 0
}
```

## Argument Reference

The following arguments are supported:

* `defense_type` - (Required, ForceNew) The Protection Module. Valid values: `ac_cc`, `antifraud`, `dld`, `normalized`, `waf`.
  * `waf`: RegEx Protection Engine.
  * `dld`: Big Data Deep Learning Engine.
  * `ac_cc`: HTTP Flood Protection.
  * `antifraud`: Data Risk Control.
  * `normalized`: Positive Security Model.
* `domain` - (Required, ForceNew) The domain name that is added to WAF.
* `instance_id` - (Required, ForceNew) The ID of the WAF instance.
* `mode` - (Required) The protection mode of the specified protection module. **NOTE:** The value of the Mode parameter varies based on the value of the `defense_type` parameter.
  * The `defense_type` is `waf`. `0`: block mode. `1`: warn mode.
  * The `defense_type` is `dld`. `0`: warn mode. `1`: block mode.
  * The `defense_type` is `ac_cc`. `0`: prevention mode. `1`: protection-emergency mode.
  * The `defense_type` is `antifraud`. `0`: warn mode. `1`: block mode. `2`: strict interception mode.
  * The `defense_type` is `normalized`. `0`: warn mode. `1`: block mode.
* `status` - (Optional) The status of the resource. Valid values: `0`, `1`.
  * `0`: disables the protection module.
  * `1`: enables the protection module.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Protection Module. The value formats as `<instance_id>:<domain>:<defense_type>`.

## Import

Web Application Firewall(WAF) Protection Module can be imported using the id, e.g.

```
$ terraform import alicloud_waf_protection_module.example <instance_id>:<domain>:<defense_type>
```
