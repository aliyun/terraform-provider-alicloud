---
subcategory: "CDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_cdn_ip_info"
sidebar_current: "docs-alicloud-datasource-cdn-ip-info"
description: |-
  Verify whether an IP is a CDN node.
---

# alicloud\_cdn\_ip\_info

This data source provides the function of verifying whether an IP is a CDN node.

-> **NOTE:** Available in v1.153.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cdn_ip_info" "ip_test" {
  ip = "114.114.114.114"
}

```

## Argument Reference

The following arguments are supported:

* `ip` - (Required, ForceNew)  Specify IP address.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:
* `cdn_ip` - Whether it belongs to an Alibaba Cloud CDN node.Valid values: `True`, `False`.
* `isp` - The operator's Chinese name, e.g. `电信`.
* `isp_ename` - The operator's English name, e.g. `telecom`.
* `region` - Chinese name of the region, e.g. `中国-贵州省-贵阳市`.
* `region_ename` - English name of the region, e.g. `China-Guizhou-guiyang`.