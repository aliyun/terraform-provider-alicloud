---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_bandwidth_limit"
sidebar_current: "docs-alicloud-resource-cen-bandwidth-limit"
description: |-
  Provides a Alicloud CEN cross-regional interconnection bandwidth configuration resource.
---

# alicloud_cen_bandwidth_limit

Provides a CEN cross-regional interconnection bandwidth resource. To connect networks in different regions, you must set cross-region interconnection bandwidth after buying a bandwidth package. The total bandwidth set for all the interconnected regions of a bandwidth package cannot exceed the bandwidth of the bandwidth package. By default, 1 Kbps bandwidth is provided for connectivity test. To run normal business, you must buy a bandwidth package and set a proper interconnection bandwidth.

For example, a CEN instance is bound to a bandwidth package of 20 Mbps and  the interconnection areas are Mainland China and North America. You can set the cross-region interconnection bandwidth between US West 1 and China East 1, China East 2, China South 1, and so on. However, the total bandwidth set for all the interconnected regions cannot exceed 20  Mbps.

For information about CEN and how to use it, see [Cross-region interconnection bandwidth](https://www.alibabacloud.com/help/doc-detail/65983.htm)

-> **NOTE:** Available since v1.18.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_bandwidth_limit&exampleId=8a3f6579-32b5-3650-e52d-00d011fa28ece1b35428&activeTab=example&spm=docs.r.cen_bandwidth_limit.0.8a3f657932&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "region1" {
  default = "eu-central-1"
}
variable "region2" {
  default = "ap-southeast-1"
}

provider "alicloud" {
  alias  = "ec"
  region = var.region1
}
provider "alicloud" {
  alias  = "as"
  region = var.region2
}

resource "alicloud_vpc" "vpc1" {
  provider   = alicloud.ec
  vpc_name   = "tf-example"
  cidr_block = "192.168.0.0/16"
}
resource "alicloud_vpc" "vpc2" {
  provider   = alicloud.as
  vpc_name   = "tf-example"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_cen_instance" "example" {
  cen_instance_name = "tf_example"
  description       = "an example for cen"
}
resource "alicloud_cen_instance_attachment" "example1" {
  instance_id              = alicloud_cen_instance.example.id
  child_instance_id        = alicloud_vpc.vpc1.id
  child_instance_type      = "VPC"
  child_instance_region_id = var.region1
}
resource "alicloud_cen_instance_attachment" "example2" {
  instance_id              = alicloud_cen_instance.example.id
  child_instance_id        = alicloud_vpc.vpc2.id
  child_instance_type      = "VPC"
  child_instance_region_id = var.region2
}
resource "alicloud_cen_bandwidth_package" "example" {
  bandwidth                  = 5
  cen_bandwidth_package_name = "tf_example"
  geographic_region_a_id     = "Europe"
  geographic_region_b_id     = "Asia-Pacific"
}

resource "alicloud_cen_bandwidth_package_attachment" "example" {
  instance_id          = alicloud_cen_instance.example.id
  bandwidth_package_id = alicloud_cen_bandwidth_package.example.id
}

resource "alicloud_cen_bandwidth_limit" "example" {
  instance_id     = alicloud_cen_bandwidth_package_attachment.example.instance_id
  region_ids      = [alicloud_cen_instance_attachment.example1.child_instance_region_id, alicloud_cen_instance_attachment.example2.child_instance_region_id]
  bandwidth_limit = 4
}
```
## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the CEN.
* `region_ids` - (Required, ForceNew) List of the two regions to interconnect. Must be two different regions.
* `bandwidth_limit` - (Required) The bandwidth configured for the interconnected regions communication.

->**NOTE:** The "alicloud_cen_bandwidthlimit" resource depends on the related "alicloud_cen_bandwidth_package_attachment" resource and "alicloud_cen_instance_attachment" resource.

## Timeouts
-> **NOTE:** Available in 1.48.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `update` - (Defaults to 10 mins) Used when activating the cen bandwidth limit when necessary during update - when changing bandwidth limit.
* `delete` - (Defaults to 10 mins) Used when terminating the cen bandwidth limit. 

## Attributes Reference

The following attributes are exported:

- `id` - ID of the resource, formatted as `<instance_id>:<region_id_1>:<region_id_2>`.

->**NOTE:** The region_id_1 and region_id_2 are sorted lexicographically.

## Import

CEN bandwidth limit can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_bandwidth_limit.example cen-abc123456:cn-beijing:eu-west-1
```

->**NOTE:** The sequence of the region_id_1 and region_id_2 makes no difference when import. But the in the id of the resource, they are sorted lexicographically.
