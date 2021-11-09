---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_vbr_health_check"
sidebar_current: "docs-alicloud-resource-cen-vbr-health-check"
description: |-
  Provides a Alicloud CEN VBR HealthCheck resource.
---

# alicloud\_cen\_vbr\_health\_check

This topic describes how to configure the health check feature for a Cloud Enterprise Network (CEN) instance. 
After you attach a Virtual Border Router (VBR) to the CEN instance and configure the health check feature, you can monitor the network conditions of the on-premises data center connected to the VBR.

For information about CEN VBR HealthCheck and how to use it, see [Manage CEN VBR HealthCheck](https://www.alibabacloud.com/help/en/doc-detail/71141.htm).

-> **NOTE:** Available in 1.88.0+

## Example Usage

Basic Usage

```terraform
# Create a cen vbr HealrhCheck resource and use it.
resource "alicloud_cen_instance" "default" {
  cen_instance_name = "test_name"
}

resource "alicloud_cen_instance_attachment" "default" {
  instance_id              = alicloud_cen_instance.default.id
  child_instance_id        = "vbr-xxxxx"
  child_instance_type      = "VBR"
  child_instance_region_id = "cn-hangzhou"
}

resource "alicloud_cen_vbr_health_check" "default" {
  cen_id                 = alicloud_cen_instance.default.id
  health_check_source_ip = "192.168.1.2"
  health_check_target_ip = "10.0.0.2"
  vbr_instance_id        = "vbr-xxxxx"
  vbr_instance_region_id = "cn-hangzhou"
  health_check_interval  = 2
  healthy_threshold      = 8
  depends_on             = [alicloud_cen_instance_attachment.default]
}
```
## Argument Reference

The following arguments are supported:

* `cen_id` - (Required, ForceNew) The ID of the CEN instance.
* `health_check_interval` - (Optional, Default) Specifies the interval at which the health check sends continuous detection packets. Default value: 2. Value range: 2 to 3.
* `health_check_source_ip` - (Optional) The source IP address of health checks.
* `health_check_target_ip` - (Required) The destination IP address of health checks.
* `healthy_threshold` - (Optional, Default) Specifies the number of probe messages sent by the health check. Default value: 8. Value range: 3 to 8.
* `vbr_instance_id` - (Required, ForceNew) The ID of the VBR.
* `vbr_instance_owner_id` - (Optional) The ID of the account to which the VBR belongs.
* `vbr_instance_region_id` - (Required, ForceNew) The ID of the region to which the VBR belongs.

->**NOTE:** The `alicloud_cen_vbr_health_check` resource depends on the related `alicloud_cen_instance_attachment` resource.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the resource, formatted as `<vbr_instance_id>:<vbr_instance_region_id>`.

### Timeouts

-> **NOTE:** Available in 1.98.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 6 mins) Used when create the CEN VBR Health Check. (until it reaches the available status).
* `update` - (Defaults to 6 mins) Used when update the CEN VBR Health Check.
* `delete` - (Defaults to 6 mins) Used when delete the CEN VBR Health Check.

## Import

CEN VBR HealthCheck can be imported using the id, e.g.

```
$ terraform import alicloud_cen_vbr_health_check.example vbr-xxxxx:cn-hangzhou
```
