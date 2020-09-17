---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_instance_grant"
sidebar_current: "docs-alicloud-resource-cen-instance-grant"
description: |-
  Provides a Alicloud CEN child instance grant resource.
---

# alicloud\_cen_instance_grant

Provides a CEN child instance grant resource, which allow you to authorize a VPC or VBR to a CEN of a different account.

For more information about how to use it, see [Attach a network in a different account](https://www.alibabacloud.com/help/doc-detail/73645.htm). 

## Example Usage

Basic Usage

```
# Create a new instance-grant and use it to grant one child instance of account1 to a new CEN of account 2.
provider "alicloud" {
  access_key = "access123"
  secret_key = "secret123"
  alias      = "account1"
}

provider "alicloud" {
  access_key = "access456"
  secret_key = "secret456"
  alias      = "account2"
}

variable "name" {
  default = "tf-testAccCenInstanceGrantBasic"
}

resource "alicloud_cen_instance" "cen" {
  provider = alicloud.account2
  name     = var.name
}

resource "alicloud_vpc" "vpc" {
  provider   = alicloud.account1
  name       = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_cen_instance_grant" "foo" {
  provider          = alicloud.account1
  cen_id            = alicloud_cen_instance.cen.id
  child_instance_id = alicloud_vpc.vpc.id
  cen_owner_id      = "uid2"
}

resource "alicloud_cen_instance_attachment" "foo" {
  provider                 = alicloud.account2
  instance_id              = alicloud_cen_instance.cen.id
  child_instance_id        = alicloud_vpc.vpc.id
  child_instance_type      = "VPC"
  child_instance_region_id = "cn-qingdao"
  child_instance_owner_id  = "uid1"
  depends_on               = [alicloud_cen_instance_grant.foo]
}
```
## Argument Reference

The following arguments are supported:

* `cen_id` - (Required) The ID of the CEN.
* `child_instance_id` - (Required) The ID of the child instance to grant.
* `cen_owner_id` - (Required) The owner UID of the  CEN which the child instance granted to.

## Attributes Reference

The following attributes are exported:

- `id` - ID of the resource, formatted as `<cen_id>:<child_instance_id>:<cen_owner_id>`.

## Import

CEN instance can be imported using the id, e.g.

```
$ terraform import alicloud_cen_instance_grant.example cen-abc123456:vpc-abc123456:uid123456
```
