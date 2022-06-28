---
subcategory: "DNS"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_gtm_instance"
sidebar_current: "docs-alicloud-resource-alidns-gtm-instance"
description: |-
  Provides a Alicloud Alidns Gtm Instance resource.
---

# alicloud\_alidns\_gtm\_instance

Provides a Alidns Gtm Instance resource.

For information about Alidns Gtm Instance and how to use it, see [What is Gtm Instance](https://www.alibabacloud.com/help/en/doc-detail/204852.html).

-> **NOTE:** Available in v1.151.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_cms_alarm_contact_group" "default" {
  alarm_contact_group_name = var.name
}

resource "alicloud_alidns_gtm_instance" "default" {
  instance_name           = var.name
  payment_type            = "Subscription"
  period                  = 1
  renewal_status          = "ManualRenewal"
  package_edition         = "standard"
  health_check_task_count = 100
  sms_notification_count  = 1000
  public_cname_mode       = "SYSTEM_ASSIGN"
  ttl                     = 60
  cname_type              = "PUBLIC"
  resource_group_id       = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  alert_group             = [alicloud_cms_alarm_contact_group.default.alarm_contact_group_name]
  public_user_domain_name = var.domain_name
  alert_config {
    sms_notice      = true
    notice_type     = "ADDR_ALERT"
    email_notice    = true
    dingtalk_notice = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `alert_config` - (Optional) The alert notification methods. See the following `Block alert_config`.
* `alert_group` - (Optional) The alert group.
* `force_update` - (Optional) The force update.
* `cname_type` - (Optional, Computed) The access type of the CNAME domain name. Valid value: `PUBLIC`.
* `instance_name` - (Required) The name of the instance.
* `lang` - (Optional) The lang.
* `payment_type` - (Required, ForceNew) The Payment Type of the resource. Valid value: `Subscription`.
* `period` - (Required) Creating a pre-paid instance, it must be set, the unit is month, please enter an integer multiple of 12 for annually paid products.
* `renew_period` - (Optional, ForceNew) Automatic renewal period, the unit is month. When setting `renewal_status` to AutoRenewal, it must be set.
* `renewal_status` - (Optional, ForceNew, Computed) Automatic renewal status. Valid values: `AutoRenewal`, `ManualRenewal`.
* `package_edition` - (Required, ForceNew) Paid package version. Valid values: `ultimate`, `standard`.
* `health_check_task_count` - (Required, ForceNew) The quota of detection tasks.
* `sms_notification_count` - (Required, ForceNew) The quota of SMS notifications.
* `strategy_mode` - (Optional, Computed) The type of the access policy. Valid values: `GEO`, `LATENCY`.
* `public_cname_mode` - (Optional, Computed) The Public Network domain name access method. Valid values: `CUSTOM`, `SYSTEM_ASSIGN`.
* `public_rr` - (Optional, Computed) The CNAME access domain name.
* `public_user_domain_name` - (Optional, Computed) The website domain name that the user uses on the Internet.
* `public_zone_name` - (Optional, Computed) The domain name that is used to access GTM over the Internet.
* `resource_group_id` - (Optional) The ID of the resource group.
* `ttl` - (Optional) The global time to live. Valid values: `60`, `120`, `300`, `600`. Unit: second.

#### Block alert_config

The alert_config supports the following: 

* `dingtalk_notice` - (Optional) Whether to configure DingTalk notifications. Valid values: `true`, `false`.
* `email_notice` - (Optional) Whether to configure mail notification. Valid values: `true`, `false`.
* `sms_notice` - (Optional) Whether to configure SMS notification. Valid values: `true`, `false`.
* `notice_type` - (Optional) The Alarm Event Type.
  - `ADDR_ALERT`: Address not available.
  - `ADDR_RESUME`: Address Recovery available.
  - `ADDR_POOL_GROUP_UNAVAILABLE`: Address pool collection not available.
  - `ADDR_POOL_GROUP_AVAILABLE`: Address pool collection recovery available.
  - `ACCESS_STRATEGY_POOL_GROUP_SWITCH`: Primary/standby address pool switch.
  - `MONITOR_NODE_IP_CHANGE`: Monitoring node IP address changes.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Gtm Instance.

## Import

Alidns Gtm Instance can be imported using the id, e.g.

```
$ terraform import alicloud_alidns_gtm_instance.example <id>
```