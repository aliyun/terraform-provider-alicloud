---
subcategory: "DNS"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_monitor_config"
sidebar_current: "docs-alicloud-resource-alidns-monitor-config"
description: |-
  Provides a Alicloud DNS Monitor Config resource.
---

# alicloud\_alidns\_monitor\_config

Provides a DNS Monitor Config resource.

For information about DNS Monitor Config and how to use it, see [What is Monitor Config](https://www.alibabacloud.com/help/en/doc-detail/198064.html).

-> **NOTE:** Available in v1.153.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testacc"
}

variable "domain_name" {
  default = "your_domain_name"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_cms_alarm_contact_group" "default" {
  alarm_contact_group_name = var.name
}

resource "alicloud_alidns_gtm_instance" "default" {
  instance_name           = var.name
  payment_type            = "Subscription"
  period                  = 1
  renewal_status          = "ManualRenewal"
  package_edition         = "ultimate"
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

resource "alicloud_alidns_address_pool" "default" {
  address_pool_name = var.name
  instance_id       = alicloud_alidns_gtm_instance.default.id
  lba_strategy      = "RATIO"
  type              = "IPV4"
  address {
    attribute_info = "{\"lineCodeRectifyType\":\"RECTIFIED\",\"lineCodes\":[\"os_namerica_us\"]}"
    remark         = "address_remark"
    address        = "1.1.1.1"
    mode           = "SMART"
    lba_weight     = 1
  }
}

resource "alicloud_alidns_monitor_config" "default" {
  addr_pool_id        = alicloud_alidns_address_pool.default.id
  evaluation_count    = "1"
  interval            = "60"
  timeout             = "5000"
  protocol_type       = "TCP"
  monitor_extend_info = "{\"failureRate\"=50,\"port\"=80}"
  isp_city_node = {
    city_code = "503"
    isp_code  = "465"
  }
}
```

## Argument Reference

The following arguments are supported:

* `addr_pool_id` - (Required, ForceNew) The ID of the address pool.
* `evaluation_count` - (Required) The number of consecutive times of failed health check attempts. Valid values: `1`, `2`, `3`.
* `interval` - (Required) The health check interval. Unit: seconds. Valid values: `60`.
* `isp_city_node` - (Required) The Monitoring node. See the following `Block isp_city_node`.
* `lang` - (Optional) The lang.
* `monitor_extend_info` - (Required) The extended information. This value follows the json format. For more details, see the [description of MonitorExtendInfo in the Request parameters table for details](https://www.alibabacloud.com/help/en/doc-detail/198064.html).
* `protocol_type` - (Required) The health check protocol. Valid values: `HTTP`, `HTTPS`, `PING`, `TCP`.
* `timeout` - (Required) The timeout period. Unit: milliseconds. Valid values: `2000`, `3000`, `5000`, `10000`.

#### Block isp_city_node

The isp_city_node supports the following: 

* `city_code` - (Required) The code of the city node to monitor.
* `isp_code` - (Required) The code of the Internet provider service (ISP) node to monitor.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Monitor Config.

## Import

DNS Monitor Config can be imported using the id, e.g.

```
$ terraform import alicloud_alidns_monitor_config.example <id>
```