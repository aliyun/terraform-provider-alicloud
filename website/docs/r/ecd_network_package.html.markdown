---
subcategory: "Elastic Desktop Service (ECD)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_network_package"
sidebar_current: "docs-alicloud-resource-ecd-network-package"
description: |-
  Provides a Alicloud ECD Network Package resource.
---

# alicloud_ecd_network_package

Provides a ECD Network Package resource.

For information about ECD Network Package and how to use it, see [What is Network Package](https://www.alibabacloud.com/help/en/wuying-workspace/developer-reference/api-ecd-2020-09-30-createnetworkpackage).

-> **NOTE:** Available since v1.142.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block          = "172.16.0.0/12"
  enable_admin_access = true
  desktop_access_type = "Internet"
  office_site_name    = "terraform-example-${random_integer.default.result}"
}

resource "alicloud_ecd_network_package" "default" {
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

```shell
$ terraform import alicloud_ecd_network_package.example <id>
```
