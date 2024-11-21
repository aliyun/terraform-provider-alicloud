---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_monitor_config"
sidebar_current: "docs-alicloud-resource-alidns-monitor-config"
description: |-
  Provides a Alicloud DNS Monitor Config resource.
---

# alicloud_alidns_monitor_config

Provides a DNS Monitor Config resource.

For information about DNS Monitor Config and how to use it, see [What is Monitor Config](https://www.alibabacloud.com/help/en/alibaba-cloud-dns/latest/api-alidns-2015-01-09-adddnsgtmmonitor).

-> **NOTE:** Available since v1.153.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alidns_monitor_config&exampleId=f0f8cd88-cfc8-fe2c-0f67-073663dc7f95e895eb0c&activeTab=example&spm=docs.r.alidns_monitor_config.0.f0f8cd88cf&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}
variable "domain_name" {
  default = "alicloud-provider.com"
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
  monitor_extend_info = "{\"failureRate\":50,\"port\":80}"
  isp_city_node {
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
* `isp_city_node` - (Required) The Monitoring node. See [`isp_city_node`](#isp_city_node) below for details.
* `lang` - (Optional) The lang.
* `monitor_extend_info` - (Required) The extended information. This value follows the json format. For more details, see the [description of MonitorExtendInfo in the Request parameters table for details](https://www.alibabacloud.com/help/en/alibaba-cloud-dns/latest/api-alidns-2015-01-09-adddnsgtmmonitor).
* `protocol_type` - (Required) The health check protocol. Valid values: `HTTP`, `HTTPS`, `PING`, `TCP`.
* `timeout` - (Required) The timeout period. Unit: milliseconds. Valid values: `2000`, `3000`, `5000`, `10000`.

### `isp_city_node`

The isp_city_node supports the following: 

* `city_code` - (Required) The code of the city node to monitor.
* `isp_code` - (Required) The code of the Internet provider service (ISP) node to monitor.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Monitor Config.

## Import

DNS Monitor Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_alidns_monitor_config.example <id>
```