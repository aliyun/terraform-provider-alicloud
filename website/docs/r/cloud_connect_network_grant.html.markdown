---
subcategory: "Smart Access Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_connect_network_grant"
sidebar_current: "docs-alicloud-resource-cloud-connect-network-grant"
description: |-
  Provides a Alicloud Cloud Connect Network Grant resource.
---

# alicloud\_cloud_connect_network\_grant

Provides a Cloud Connect Network Grant resource. If the CEN instance to be attached belongs to another account, authorization by the CEN instance is required.

For information about Cloud Connect Network Grant and how to use it, see [What is Cloud Connect Network Grant](https://www.alibabacloud.com/help/doc-detail/94543.htm).

-> **NOTE:** Available in 1.63.0+

-> **NOTE:** Only the following regions support create Cloud Connect Network Grant. [`cn-shanghai`, `cn-shanghai-finance-1`, `cn-hongkong`, `ap-southeast-1`, `ap-southeast-2`, `ap-southeast-3`, `ap-southeast-5`, `ap-northeast-1`, `eu-central-1`]

## Example Usage

Basic Usage

```
provider "alicloud" {
  alias = "ccn_account"
}

provider "alicloud" {
  region     = "cn-hangzhou"
  access_key = "xxxxxx"
  secret_key = "xxxxxx"
  alias      = "cen_account"
}

resource "alicloud_cen_instance" "cen" {
  provider = alicloud.cen_account
  name     = "tf-testAccCenInstance-xxx"
}

resource "alicloud_cloud_connect_network" "ccn" {
  provider   = alicloud.ccn_account
  name       = "tf-testAccCloudConnectNetwork-xxx"
  is_default = "true"
}

resource "alicloud_cloud_connect_network_grant" "default" {
  ccn_id  = alicloud_cloud_connect_network.ccn.id
  cen_id  = alicloud_cen_instance.cen.id
  cen_uid = "xxxxxx"
  depends_on = [
    alicloud_cloud_connect_network.ccn,
    alicloud_cen_instance.cen,
  ]
}
```
## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required,ForceNew) The ID of the CCN instance.
* `cen_id` - (Required,ForceNew) The ID of the CEN instance.
* `cen_uid` - (Required,ForceNew) The ID of the account to which the CEN instance belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Cloud Connect Network grant Id and formates as `<ccn_id>:<cen_id>`.

## Import

The Cloud Connect Network Grant can be imported using the instance_id, e.g.

```
$ terraform import alicloud_cloud_connect_network_grant.example ccn-abc123456:cen-abc123456
```

