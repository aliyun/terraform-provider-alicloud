---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_site_monitor"
sidebar_current: "docs-alicloud-resource-cms-site-monitor"
description: |-
  Provides a resource to build a site monitor rule for cloud monitor.
---

# alicloud_cms_site_monitor

This resource provides a site monitor resource and it can be used to monitor public endpoints and websites.
Details at https://www.alibabacloud.com/help/doc-detail/67907.htm

-> **NOTE:** Available since v1.72.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cms_site_monitor" "basic" {
  address   = "http://www.alibabacloud.com"
  task_name = "tf-example"
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
* `task_type` - (Required, ForceNew) The protocol of the site monitoring task. Currently, site monitoring supports the following protocols: HTTP, PING, TCP, UDP, DNS, SMTP, POP3, and FTP.
* `alert_ids` - (Optional, List) The IDs of existing alert rules to be associated with the site monitoring task.
* `interval` - (Optional, Int) The monitoring interval of the site monitoring task. Unit: minutes. Valid values: `1`, `5`, `15`, `30` and `60`. Default value: `1`. **NOTE:** From version 1.207.0, `interval` can be set to `30`, `60`.
* `options_json` - (Optional) The extended options of the protocol of the site monitoring task. The options vary according to the protocol.
* `isp_cities` - (Optional, Set) The detection points in a JSON array. For example, `[{"city":"546","isp":"465"},{"city":"572","isp":"465"},{"city":"738","isp":"465"}]` indicates the detection points in Beijing, Hangzhou, and Qingdao respectively. You can call the [DescribeSiteMonitorISPCityList](https://www.alibabacloud.com/help/en/doc-detail/115045.htm) operation to query detection point information. If this parameter is not specified, three detection points will be chosen randomly for monitoring. See [`isp_cities`](#isp_cities) below.

### `isp_cities`

The isp_cities supports the following:

* `city` - (Required) The ID of the city.
* `isp` - (Required) The ID of the carrier.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the site monitor rule.
* `task_state` - The status of the site monitoring task.
* `create_time` - The time when the site monitoring task was created.
* `update_time` - The time when the site monitoring task was updated.

## Timeouts

-> **NOTE:** Available since 1.207.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Site Monitor.
* `update` - (Defaults to 5 mins) Used when update the Site Monitor.
* `delete` - (Defaults to 3 mins) Used when delete the Site Monitor.

## Import

Cloud Monitor Service Site Monitor can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_site_monitor.example <id>
```
