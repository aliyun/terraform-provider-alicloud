---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_fw_switch"
description: |-
  Provides a Alicloud Cloud Firewall Fw Switch resource.
---

# alicloud_cloud_firewall_fw_switch

Provides a Cloud Firewall Fw Switch resource.

The protection switch of an internet boundary firewall asset. The existence of the resource indicates that protection for the public IP asset is enabled (ProtectStatus=open), and removal of the resource indicates that protection is disabled (ProtectStatus=closed).

For information about Cloud Firewall Fw Switch and how to use it, see [What is Fw Switch](https://next.api.alibabacloud.com/document/Cloudfw/2017-12-07/PutEnableFwSwitch).

-> **NOTE:** Available since v1.295.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cloud_firewall_fw_switch" "example" {
  internet_address = "203.0.113.1"
  member_uid       = 1412345678901234
}
```

## Argument Reference

The following arguments are supported:
* `dry_run` - (Optional) Specifies whether to precheck the request only. Valid values: `true`, `false`.

  -> **NOTE:** This parameter is only evaluated during resource creation and deletion. Modifying it in isolation will not trigger any action.

* `internet_address` - (Required, ForceNew) The public IP address of the internet boundary firewall asset for which to enable or disable protection. It is also the unique ID of the resource.
* `ip_version` - (Optional, ForceNew, Computed) The IP version. Valid values: `4` (IPv4) and `6` (IPv6).
* `lang` - (Optional) The language type of the received message.
  - `zh`: Chinese
  - `en`: English

  -> **NOTE:** This parameter is immutable. Changing it after creation has no effect.

* `member_uid` - (Optional, ForceNew, Computed, Int) The UID of the Cloud Firewall member account.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `ali_uid` - AliUid.
* `bind_instance_id` - BindInstanceId.
* `bind_instance_name` - BindInstanceName.
* `intranet_address` - IntranetAddress.
* `name` - Name.
* `note` - Note.
* `protect_status` - The protection status of the asset. Valid values: `open`, `opening`, `closed`, `closing`.
* `region_id` - RegionId.
* `region_status` - RegionStatus.
* `resource_instance_id` - ResourceInstanceId.
* `resource_type` - The type of the asset public IP. Valid values: `BastionHostIP`, `EcsEIP`, `EcsPublicIP`, `EIP`, `EniEIP`, `NatEIP`, `NatPublicIP`, `SlbEIP`.
* `sg_status` - SgStatus.
* `sg_status_time` - SgStatusTime.
* `sync_status` - SyncStatus.
* `type` - Type.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Fw Switch.
* `delete` - (Defaults to 5 mins) Used when delete the Fw Switch.

## Import

Cloud Firewall Fw Switch can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_fw_switch.example <internet_address>
```
