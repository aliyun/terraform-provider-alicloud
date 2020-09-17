---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_bandwidth_limit"
sidebar_current: "docs-alicloud-resource-cen-bandwidth-limit"
description: |-
  Provides a Alicloud CEN cross-regional interconnection bandwidth configuration resource.
---

# alicloud\_cen_bandwidth_limit

Provides a CEN cross-regional interconnection bandwidth resource. To connect networks in different regions, you must set cross-region interconnection bandwidth after buying a bandwidth package. The total bandwidth set for all the interconnected regions of a bandwidth package cannot exceed the bandwidth of the bandwidth package. By default, 1 Kbps bandwidth is provided for connectivity test. To run normal business, you must buy a bandwidth package and set a proper interconnection bandwidth.

For example, a CEN instance is bound to a bandwidth package of 20 Mbps and  the interconnection areas are Mainland China and North America. You can set the cross-region interconnection bandwidth between US West 1 and China East 1, China East 2, China South 1, and so on. However, the total bandwidth set for all the interconnected regions cannot exceed 20  Mbps.

For information about CEN and how to use it, see [Cross-region interconnection bandwidth](https://www.alibabacloud.com/help/doc-detail/65983.htm)

## Example Usage

Basic Usage

```
variable "name" {
  default = "tf-testAccCenBandwidthLimitConfig"
}

provider "alicloud" {
  alias  = "fra"
  region = "eu-central-1"
}

provider "alicloud" {
  alias  = "sh"
  region = "cn-shanghai"
}

resource "alicloud_vpc" "vpc1" {
  provider   = alicloud.fra
  name       = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vpc" "vpc2" {
  provider   = alicloud.sh
  name       = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_cen_instance" "cen" {
  name        = var.name
  description = "tf-testAccCenBandwidthLimitConfigDescription"
}

resource "alicloud_cen_bandwidth_package" "bwp" {
  bandwidth = 5
  geographic_region_ids = [
    "Europe",
    "China",
  ]
}

resource "alicloud_cen_bandwidth_package_attachment" "bwp_attach" {
  instance_id          = alicloud_cen_instance.cen.id
  bandwidth_package_id = alicloud_cen_bandwidth_package.bwp.id
}

resource "alicloud_cen_instance_attachment" "vpc_attach_1" {
  instance_id              = alicloud_cen_instance.cen.id
  child_instance_id        = alicloud_vpc.vpc1.id
  child_instance_type      = "VPC"
  child_instance_region_id = "eu-central-1"
}

resource "alicloud_cen_instance_attachment" "vpc_attach_2" {
  instance_id              = alicloud_cen_instance.cen.id
  child_instance_id        = alicloud_vpc.vpc2.id
  child_instance_type      = "VPC"
  child_instance_region_id = "cn-shanghai"
}

resource "alicloud_cen_bandwidth_limit" "foo" {
  instance_id = alicloud_cen_instance.cen.id
  region_ids = [
    "eu-central-1",
    "cn-shanghai",
  ]
  bandwidth_limit = 4
  depends_on = [
    alicloud_cen_bandwidth_package_attachment.bwp_attach,
    alicloud_cen_instance_attachment.vpc_attach_1,
    alicloud_cen_instance_attachment.vpc_attach_2,
  ]
}
```
## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the CEN.
* `region_ids` - (Required, ForceNew) List of the two regions to interconnect. Must be two different regions.
* `bandwidth_limit` - (Required) The bandwidth configured for the interconnected regions communication.

->**NOTE:** The "alicloud_cen_bandwidthlimit" resource depends on the related "alicloud_cen_bandwidth_package_attachment" resource and "alicloud_cen_instance_attachment" resource.

### Timeouts
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

```
terraform import alicloud_cen_bandwidth_limit.example cen-abc123456:cn-beijing:eu-west-1
```

->**NOTE:** The sequence of the region_id_1 and region_id_2 makes no difference when import. But the in the id of the resource, they are sorted lexicographically.
