---
subcategory: "Smart Access Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_sag_grant"
sidebar_current: "docs-alicloud-resource-sag-grant"
description: |-
  Provides a Alicloud Sag Grant resource.
---

# alicloud\_sag\_grant

Provides a Sag Grant resource. After setting up cross-account authorization, the other party's account can be loaded into its cloud connection network through the SAG connection gateway instance, and the cloud connection network where the other party's cloud account is located will open to your network, please proceed with caution.

For information about Sag Grant and how to use it, see [What is Sag Grant](https://www.alibabacloud.com/help/doc-detail/132028.htm).

-> **NOTE:** Available in 1.65.0+

-> **NOTE:** Only the following regions support. [`cn-shanghai`, `cn-shanghai-finance-1`, `cn-hongkong`, `ap-southeast-1`, `ap-southeast-2`, `ap-southeast-3`, `ap-southeast-5`, `ap-northeast-1`, `eu-central-1`]

## Example Usage

Basic Usage

```
provider "alicloud" {
  alias = "sag_account"
}
provider "alicloud" {
  region     = "cn-shanghai"
  access_key = "xxx"
  secret_key = "xxx"
  alias      = "ccn_account"
}
resource "alicloud_cloud_connect_network" "ccn" {
  provider   = "alicloud.ccn_account"
  name       = "tf-testAccCloudConnectNetwork-xxx"
  is_default = "true"
}
resource "alicloud_sag_grant" "default" {
  provider   = "alicloud.sag_account"
  sag_id     = "tf-testAccSagGrant-xxx"
  ccn_id     = "${alicloud_cloud_connect_network.ccn.id}"
  ccn_uid    = "xxx"
  depends_on = ["alicloud_cloud_connect_network.ccn"]
}
```
## Argument Reference

The following arguments are supported:

* `sag_id` - (Required,ForceNew) The ID of the SAG instance.
* `ccn_id` - (Required,ForceNew) The ID of the CCN instance.
* `ccn_uid` - (Required,ForceNew) The ID of the account to which the CCN instance belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the SAG grant Id and formates as `<sag_id>:<ccn_id>`.

## Import

The SAG Grant can be imported using the instance_id, e.g.

```
$ terraform import alicloud_sag_grant.example sag-abc123456:ccn-abc123456
```