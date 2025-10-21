---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_site_monitor"
description: |-
  Provides a Alicloud Cloud Monitor Service Site Monitor resource.
---

# alicloud_cms_site_monitor

Provides a Cloud Monitor Service Site Monitor resource.

Describes the SITE monitoring tasks created by the user.

For information about Cloud Monitor Service Site Monitor and how to use it, see [What is Site Monitor](https://next.api.alibabacloud.com/document/Cms/2019-01-01/CreateSiteMonitor).

-> **NOTE:** Available since v1.72.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cms_site_monitor&exampleId=bf350aa4-e5a5-2b81-3d7a-b32d7aad969600731d9d&activeTab=example&spm=docs.r.cms_site_monitor.0.bf350aa4e5&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
  options_json = <<EOT
{
    "http_method": "get",
    "waitTime_after_completion": null,
    "ipv6_task": false,
    "diagnosis_ping": false,
    "diagnosis_mtr": false,
    "assertions": [
        {
            "operator": "lessThan",
            "type": "response_time",
            "target": 1000
        }
    ],
    "time_out": 30000
}
EOT
}
```

## Argument Reference

The following arguments are supported:
* `address` - (Required) The URL or IP address monitored by the site monitoring task.
* `agent_group` - (Optional, ForceNew, Available since v1.262.0) The type of the detection point. Default value: `PC`. Valid values: `PC`, `MOBILE`.
* `custom_schedule` - (Optional, Set, Available since v1.262.0) Custom probing period. Only a certain period of time from Monday to Sunday can be selected for detection. See [`custom_schedule`](#custom_schedule) below.
* `interval` - (Optional) The monitoring interval of the site monitoring task. Unit: minutes. Valid values: `1`, `5`, `15`, `30` and `60`. Default value: `1`. **NOTE:** From version 1.207.0, `interval` can be set to `30`, `60`.
* `isp_cities` - (Optional, Set) The detection points in a JSON array. For example, `[{"city":"546","isp":"465"},{"city":"572","isp":"465"},{"city":"738","isp":"465"}]` indicates the detection points in Beijing, Hangzhou, and Qingdao respectively. You can call the [DescribeSiteMonitorISPCityList](https://www.alibabacloud.com/help/en/doc-detail/115045.htm) operation to query detection point information. If this parameter is not specified, three detection points will be chosen randomly for monitoring. See [`isp_cities`](#isp_cities) below.
* `status` - (Optional, Available since v1.262.0) The status of the site monitoring task. Valid values:
  - `1`: The task is enabled.
  - `2`: The task is disabled.
* `task_name` - (Required) The name of the site monitoring task. The name must be 4 to 100 characters in length. The name can contain the following types of characters: letters, digits, and underscores.
* `task_type` - (Required, ForceNew) The protocol of the site monitoring task. Currently, site monitoring supports the following protocols: HTTP, PING, TCP, UDP, DNS, SMTP, POP3, and FTP.
* `options_json` - (Optional) The extended options of the protocol of the site monitoring task. The options vary according to the protocol. See [extended options](https://www.alibabacloud.com/help/en/cms/developer-reference/api-cms-2019-01-01-createsitemonitor#api-detail-35).
* `alert_ids` - (Optional, List, Deprecated since v1.262.0) Field `alert_ids` has been deprecated from provider version 1.262.0.

### `custom_schedule`

The custom_schedule supports the following:
* `days` - (Optional, List, Available since v1.262.0) The days in a week.
* `end_hour` - (Optional, Int, Available since v1.262.0) The end time of the detection. Unit: hours.
* `start_hour` - (Optional, Int, Available since v1.262.0) The start time of the detection. Unit: hours.
* `time_zone` - (Optional, Available since v1.262.0) The time zone of the detection.

### `isp_cities`

The isp_cities supports the following:

* `city` - (Optional) The ID of the city.
* `isp` - (Optional) The ID of the carrier.
* `type` - (Optional, Available since v1.262.0) The network type of the detection point. Valid values: `IDC`, `LASTMILE`, and `MOBILE`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource supplied above.
* `task_state` - (Deprecated since v1.262.0) Field `task_state` has been deprecated from provider version 1.262.0. New field `status` instead
* `create_time` - (Deprecated since v1.262.0) Field `create_time` has been deprecated from provider version 1.262.0.
* `update_time` - (Deprecated since v1.262.0) Field `update_time` has been deprecated from provider version 1.262.0.

## Timeouts

-> **NOTE:** Available since 1.207.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Site Monitor.
* `delete` - (Defaults to 5 mins) Used when delete the Site Monitor.
* `update` - (Defaults to 5 mins) Used when update the Site Monitor.

## Import

Cloud Monitor Service Site Monitor can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_site_monitor.example <id>
```
