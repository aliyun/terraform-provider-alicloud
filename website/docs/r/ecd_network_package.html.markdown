---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_network_package"
sidebar_current: "docs-alicloud-resource-ecd-network-package"
description: |-
  Provides a Alicloud ECD Network Package resource.
---

# alicloud\_ecd\_network\_package

Provides a ECD Network Package resource.

For information about ECD Network Package and how to use it, see [What is Network Package](https://help.aliyun.com/document_detail/188382.html).

-> **NOTE:** Available in v1.142.0+.

## Example Usage

Basic Usage

```terraform

resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block          = "172.16.0.0/12"
  desktop_access_type = "Internet"
  office_site_name    = "your_office_site_name"
}

resource "alicloud_ecd_network_package" "example" {
  bandwidth      = 10
  office_site_id = alicloud_ecd_simple_office_site.default.id
}

```

## Argument Reference

The following arguments are supported:

* `bandwidth` - (Required) The bandwidth of package public network bandwidth peak. Valid values: 1~200. Unit:Mbps.
* `internet_charge_type` - (Optional, ForceNew) The internet charge type  of  package.
* `office_site_id` - (Required, ForceNew) The ID of office site.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Network Package.
* `status` - The status of network package. Valid values: `Creating`, `InUse`, `Releasing`,`Released`.

## Import

ECD Network Package can be imported using the id, e.g.

```
$ terraform import alicloud_ecd_network_package.example <id>
```
