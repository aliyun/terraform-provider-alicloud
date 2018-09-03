---
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_bandwidthlimit"
sidebar_current: "docs-alicloud-resource-cen-bandwidthlimit"
description: |-
  Provides a Alicloud CEN cross-regional interconnection bandwidth configuration resource.
---

# alicloud\_cen_bandwidthlimit

Provides a CEN cross-regional interconnection bandwidth configuration resource.

## Example Usage

Basic Usage

```
# Create a bandwidthlimit and use it to confige the cross-regional interconnection bandwidth between cn-beijing and cn-shanghai in a new CEN
provider "alicloud" {
  alias = "bj"
  region = "cn-beijing"
}

provider "alicloud" {
  alias = "sh"
  region = "cn-shanghai"
}

resource "alicloud_vpc" "vpc1" {
  provider = "alicloud.bj"
  name = "terraform-01"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vpc" "vpc2" {
  provider = "alicloud.sh"
  name = "terraform-02"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_cen" "cen" {
     name = "terraform-yl-01"
     description = "terraform01"
}

resource "alicloud_cen_bandwidthpackage" "bwp" {
    bandwidth = 20
    geographic_region_id = [
		"China",
		"China"]
}

resource "alicloud_cen_bandwidthpackage_attachment" "bwp_attach" {
    cen_id = "${alicloud_cen.cen.id}"
    cen_bandwidthpackage_id = "${alicloud_cen_bandwidthpackage.bwp.id}"
}

resource "alicloud_cen_instance_attachment" "vpc_attach_1" {
    cen_id = "${alicloud_cen.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc1.id}"
    child_instance_type = "VPC"
    child_instance_region_id = "cn-beijing"
}

resource "alicloud_cen_instance_attachment" "vpc_attach_2" {
    cen_id = "${alicloud_cen.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc2.id}"
    child_instance_type = "VPC"
    child_instance_region_id = "cn-shanghai"
}

resource "alicloud_cen_bandwidthlimit" "foo" {
     cen_id = "${alicloud_cen.cen.id}"
     regions_id = [
                    "cn-beijing",
                    "cn-shanghai"]
     bandwidth_limit = 15
     depends_on = [
        "alicloud_cen_bandwidthpackage_attachment.bwp_attach",
        "alicloud_cen_instance_attachment.vpc_attach_1",
        "alicloud_cen_instance_attachment.vpc_attach_2"]
}
```
## Argument Reference

The following arguments are supported:

* `cen_id` - (Required) The ID of the CEN.
* `regions_id` - (Required) List of the two regions to interconnect. 
* `bandwidth_limit` - (Required) The bandwidth configured for the interconnected regions communication.

~>**NOTE:** The "alicloud_cen_bandwidthlimit" resource depends on the related "alicloud_cen_bandwidthpackage_attachment" resource and "alicloud_cen_instance_attachment" resource.

## Attributes Reference

The following attributes are exported:

- `cen_id` - (Required) The ID of the CEN.
- `regions_id` - (Required) List of the two regions to interconnect. 
- `bandwidth_limit` - (Required) The bandwidth configured for the interconnected regions communication.