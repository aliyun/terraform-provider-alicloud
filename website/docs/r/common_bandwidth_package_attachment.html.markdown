---
subcategory: "EIP Bandwidth Plan (CBWP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_common_bandwidth_package_attachment"
sidebar_current: "docs-alicloud-resource-common-bandwidth-package-attachment"
description: |-
  Provides an Alicloud Common  Attachment resource.
---

# alicloud\_common\_bandwidth\_package\_attachment

Provides an Alicloud Common Bandwidth Package Attachment resource for associating Common Bandwidth Package to EIP Instance.

-> **NOTE:** Terraform will auto build common bandwidth package attachment while it uses `alicloud_common_bandwidth_package_attachment` to build a common bandwidth package attachment resource.

For information about common bandwidth package and how to use it, see [What is Common Bandwidth Package](https://www.alibabacloud.com/help/product/55092.htm).

-> **NOTE:** From version 1.194.0, the resource can set the maximum bandwidth of an EIP that is associated with an EIP bandwidth plan by `bandwidth_package_bandwidth`. see [how to use it](https://www.alibabacloud.com/help/en/eip-bandwidth-plan/latest/120327).

## Example Usage

Basic Usage

```terraform
resource "alicloud_common_bandwidth_package" "foo" {
  bandwidth   = "2"
  name        = "test_common_bandwidth_package"
  description = "test_common_bandwidth_package"
}

resource "alicloud_eip_address" "foo" {
  bandwidth            = "2"
  internet_charge_type = "PayByBandwidth"
}

resource "alicloud_common_bandwidth_package_attachment" "foo" {
  bandwidth_package_id = alicloud_common_bandwidth_package.foo.id
  instance_id          = alicloud_eip_address.foo.id
}
```
## Argument Reference

The following arguments are supported:

* `bandwidth_package_id` - (Required, ForceNew) The bandwidth_package_id of the common bandwidth package attachment, the field can't be changed.
* `instance_id` - (Required, ForceNew) The instance_id of the common bandwidth package attachment, the field can't be changed.
* `bandwidth_package_bandwidth` - (Optional, Computed, Available in v1.194.0+) The maximum bandwidth for the EIP. This value cannot be larger than the maximum bandwidth of the EIP bandwidth plan. Unit: Mbit/s.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the common bandwidth package attachment id and formates as `<bandwidth_package_id>:<instance_id>`.

#### Timeouts

-> **NOTE:** Available in 1.194.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Common Bandwidth Package Attachment.
* `update` - (Defaults to 5 mins) Used when update the Common Bandwidth Package Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Common Bandwidth Package Attachment.

## Import

The common bandwidth package attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_common_bandwidth_package_attachment.foo cbwp-abc123456:eip-abc123456
```
