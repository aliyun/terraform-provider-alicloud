---
subcategory: "DdosCoo"
layout: "alicloud"
page_title: "Alicloud: alicloud_ddoscoo_web_rule"
sidebar_current: "docs-alicloud-resource-ddoscoo-web-rule"
description: |-
  Provides a Alicloud Ddoscoo Web Rule resource.
---

# alicloud\_ddoscoo\_web\_rule

Provides a Ddoscoo Web Rule resource.
For information about Ddoscoo Web Rule and how to use it, see [What is Ddoscoo Web Rule](https://www.alibabacloud.com/help/en/doc-detail/157463.htm).

-> **NOTE:** Available in 1.84.0+

## Example Usage

Basic Usage

```
resource "alicloud_ddos_coo_web_rule" "example" {
  domain="sojson.com"
  rs_type="0"
  rules= <<EOF
          [{
            "ProxyRules": [{
                "ProxyPort": 80,
                "RealServers": ["2.2.2.2"]
            }],
            "ProxyType": "http"
        }, {
            "ProxyRules": [{
                "ProxyPort": 443,
                "RealServers": ["3.3.3.3"]
            }],
            "ProxyType": "https"
        }]
  EOF
  real_servers= <<EOF
                ["1.1.1.1", "2.2.2.2", "4.4.4.4"]
  EOF
  proxy_types= <<EOF
            [{
                "ProxyType": "http",
                "ProxyPorts": [80]
            }, {
                "ProxyType": "https",
                "ProxyPorts": [443]
            }]
  EOF
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, ForceNew) The domain name of the website.
* `rs_type` - (Required) The type of the server address. Valid values: 0: Origin IP address. 1: Origin site domain name.
* `rules` - (Required, ForceNew) The details of the forwarding rule. The structure is in JSON format. ProxyRules: Array type. Required. The protocol information. The structure is as follows: ProxyPort: Required. The port number. It is an Integer. RealServers server address. This parameter is required and of Array type. ProxyType protocol Type. Required. The value is of the String type. Valid values: HTTP,https,websocket,websockets.
* `real_servers` - (Optional) The list of server addresses.
* `proxy_types` - (Optional) The protocol of the forwarding rule. The structure is in JSON format. ProxyType protocol Type. Required. The value is of the String type. Valid values: HTTP,https,websocket,websockets. ProxyPort: Required. The port number. It is an Integer.
* `resource_group_id` - (Optional) The ID of the resource group to which the anti-DDoS pro instance belongs in resource management. By default, no value is specified, indicating that the domains in the default resource group are listed.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of web rule. The value is "`<domain>`".

## Import

Ddoscoo Web Rule can be imported using the id, e.g.

```
$ terraform import alicloud_ddos_coo_web_rule.example sojson.com
```
