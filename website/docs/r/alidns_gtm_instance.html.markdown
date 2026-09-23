---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_gtm_instance"
sidebar_current: "docs-alicloud-resource-alidns-gtm-instance"
description: |-
  Provides a Alicloud Alidns Gtm Instance resource.
---

# alicloud_alidns_gtm_instance

Provides a Alidns Gtm Instance resource.

For information about Alidns Gtm Instance and how to use it, see [What is Gtm Instance](https://www.alibabacloud.com/help/en/doc-detail/204852.html).

-> **NOTE:** Available since v1.151.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alidns_gtm_instance&exampleId=20a4adaa-6554-13bf-04da-f128fffe8ef72c331f77&activeTab=example&spm=docs.r.alidns_gtm_instance.0.20a4adaa65&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "domain_name" {
  default = "alicloud-provider.com"
}
data "alicloud_resource_manager_resource_groups" "default" {}
resource "alicloud_cms_alarm_contact_group" "default" {
  alarm_contact_group_name = "tf_example"
}

resource "alicloud_alidns_gtm_instance" "default" {
  instance_name           = "tf_example"
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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_alidns_gtm_instance&spm=docs.r.alidns_gtm_instance.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `alert_config` - (Optional) The alert notification methods. See [`alert_config`](#alert_config) below for details.
* `alert_group` - (Optional) The alert group.
* `force_update` - (Optional) The force update.
* `cname_type` - (Optional) The access type of the CNAME domain name. Valid value: `PUBLIC`.
* `instance_name` - (Required) The name of the instance.
* `lang` - (Optional) The lang.
* `payment_type` - (Required, ForceNew) The Payment Type of the resource. Valid value: `Subscription`.
* `period` - (Required) Creating a pre-paid instance, it must be set, the unit is month, please enter an integer multiple of 12 for annually paid products.
* `renew_period` - (Optional, ForceNew) Automatic renewal period, the unit is month. When setting `renewal_status` to AutoRenewal, it must be set.
* `renewal_status` - (Optional, ForceNew) Automatic renewal status. Valid values: `AutoRenewal`, `ManualRenewal`.
* `package_edition` - (Required, ForceNew) Paid package version. Valid values: `ultimate`, `standard`.
* `health_check_task_count` - (Required, ForceNew) The quota of detection tasks.
* `sms_notification_count` - (Required, ForceNew) The quota of SMS notifications.
* `strategy_mode` - (Optional) The type of the access policy. Valid values: `GEO`, `LATENCY`.
* `public_cname_mode` - (Optional) The Public Network domain name access method. Valid values: `CUSTOM`, `SYSTEM_ASSIGN`.
* `public_rr` - (Optional) The CNAME access domain name.
* `public_user_domain_name` - (Optional) The website domain name that the user uses on the Internet.
* `public_zone_name` - (Optional) The domain name that is used to access GTM over the Internet.
* `resource_group_id` - (Optional) The ID of the resource group.
* `ttl` - (Optional) The global time to live. Valid values: `60`, `120`, `300`, `600`. Unit: second.

### `alert_config`

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

```shell
$ terraform import alicloud_alidns_gtm_instance.example <id>
```