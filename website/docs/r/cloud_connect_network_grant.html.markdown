---
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_connect_network_grant"
sidebar_current: "docs-alicloud-resource-cloud-connect-network-grant"
description: |-
  Provides a Alicloud Cloud Connect Network Grant resource.
---

# alicloud\_sag\_client_user

Provides a Cloud Connect Network Grant resource. If the CEN instance to be attached belongs to another account, authorization by the CEN instance is required.

For information about Cloud Connect Network Grant and how to use it, see [What is Cloud Connect Network Grant](https://www.alibabacloud.com/help/doc-detail/94543.htm).

-> **NOTE:** Available in 1.60.0+

-> **NOTE:** Only the following regions support create Cloud Connect Network Grant. [`cn-shanghai`, `cn-shanghai-finance-1`, `cn-hongkong`, `ap-southeast-1`, `ap-southeast-2`, `ap-southeast-3`, `ap-southeast-5`, `ap-northeast-1`, `eu-central-1`]

## Example Usage

Basic Usage

```
variable "name" {
  default = "tf-testAccCloudConnectNetworkGrant"
}

data "alicloud_account" "default"{
}

resource "alicloud_cen_instance" "default" {
  name = "${var.name}"
}

resource "alicloud_cloud_connect_network" "default" {
  name = "${var.name}"
  is_default = "true"
}

resource "alicloud_cloud_connect_network_grant" "default" {
  ccn_id = "${alicloud_cloud_connect_network.default.id}"
  cen_id = "${alicloud_cen_instance.default.id}"
  cen_uid = "${data.alicloud_account.default.id}"
}
```
## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required) The ID of the CCN instance.
* `cen_id` - (Required) The ID of the CEN instance.
* `cen_uid` - (Required) The ID of the account to which the CEN instance belongs.

## Attributes Reference

The following attributes are exported:

* `"instance_id` - The ID of the authorization rule instance.

## Import

The Cloud Connect Network Grant can be imported using the instance_id, e.g.

```
$ terraform import alicloud_cloud_connect_network_grant.example aaa
```

