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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecd_network_package&exampleId=7532ca24-20d0-9045-18e0-33f8b6a7ca732eda3a9a&activeTab=example&spm=docs.r.ecd_network_package.0.7532ca2420&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ecd_network_package&spm=docs.r.ecd_network_package.example&intl_lang=EN_US)

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
