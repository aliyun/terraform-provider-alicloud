---
subcategory: "Private Zone"
layout: "alicloud"
page_title: "Alicloud: alicloud_private_zone_zone_manage"
sidebar_current: "docs-alicloud-resource-private_zone_zone_manage"
description: |-
  Provides Alibaba Cloud Private Zone Zone Manage Resource
---

# alicloud\_private\_zone\_zone\_manage

Provides a Private Zone Zone Manage resource. PrivateZone is an Alibaba Cloud private domain name resolution and management service based on VPC. You can use PrivateZone to resolve private domain names to IP addresses in one or multiple specified VPCs.
For information about Private Zone Zone Manage and how to use it, see [What is Private Zone Zone Manage](https://www.alibabacloud.com/help/en/doc-detail/64611.htm).

-> **NOTE:** Available in 1.83.0+

## Example Usage

Basic Usage

```
resource "alicloud_private_zone_zone_manage" "example" {
  zone_name="demo.com"
  proxy_pattern="ZONE"
}
```

## Argument Reference

The following arguments are supported:

* `zone_name` - (Optional, ForceNew) The name of the zone.
* `proxy_pattern` - (Optional) ZONE: indicates that the recursive DNS proxy is disabled.RECORD: indicates that the recursive DNS proxy is enabled.
* `lang` - (Optional) The language.
* `remark` - (Optional) The description.
* `user_client_ip` - (Optional) The IP address of the client.
* `resource_group_id` - The ID of the resource group.

## Attributes Reference

* `id` - The ID of the zone manage.

## Import

Private Zone Zone Manage can be imported using the id, e.g.

```
$ terraform import alicloud_private_zone_zone_manage.example 2b28070**** 
```
