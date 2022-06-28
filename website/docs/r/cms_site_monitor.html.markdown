---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_site_monitor"
sidebar_current: "docs-alicloud-resource-cms-site-monitor"
description: |-
  Provides a resource to build a site monitor rule for cloud monitor.
---

# alicloud\_cms\_site\_monitor

This resource provides a site monitor resource and it can be used to monitor public endpoints and websites.
Details at https://www.alibabacloud.com/help/doc-detail/67907.htm

Available in 1.72.0+

## Example Usage

Basic Usage

```
resource "alicloud_cms_site_monitor" "basic" {
  address   = "http://www.alibabacloud.com"
  task_name = "tf-testAccCmsSiteMonitor_basic"
  task_type = "HTTP"
  interval  = 5
  isp_cities {
    city = "546"
    isp  = "465"
  }
}
```

## Argument Reference

The following arguments are supported:

* `address` - (Required) The URL or IP address monitored by the site monitoring task.
* `task_name` - (Required) The name of the site monitoring task. The name must be 4 to 100 characters in length. The name can contain the following types of characters: letters, digits, and underscores.
* `task_type` - (Required, ForceNew) The protocol of the site monitoring task. Currently, site monitoring supports the following protocols: HTTP, Ping, TCP, UDP, DNS, SMTP, POP3, and FTP.
* `alert_ids` - The IDs of existing alert rules to be associated with the site monitoring task.
* `interval` - The monitoring interval of the site monitoring task. Unit: minutes. Valid values: 1, 5, and 15. Default value: 1.
* `isp_cities` - The detection points in a JSON array. For example, `[{"city":"546","isp":"465"},{"city":"572","isp":"465"},{"city":"738","isp":"465"}]` indicates the detection points in Beijing, Hangzhou, and Qingdao respectively. You can call the [DescribeSiteMonitorISPCityList](https://www.alibabacloud.com/help/en/doc-detail/115045.htm) operation to query detection point information. If this parameter is not specified, three detection points will be chosen randomly for monitoring.
* `options_json` - The extended options of the protocol of the site monitoring task. The options vary according to the protocol.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the site monitor rule.

## Import

Alarm rule can be imported using the id, e.g.

```
$ terraform import alicloud_cms_site_monitor.alarm abc12345
```
