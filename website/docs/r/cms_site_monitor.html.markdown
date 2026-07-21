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
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cms_site_monitor&exampleId=53d1dff0-2de6-b3f3-fb5d-e7705480a25b0d304b16&activeTab=example&spm=docs.r.cms_site_monitor.0.53d1dff02d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_cms_site_monitor" "basic" {
  address   = "https://www.alibabacloud.com"
  task_name = var.name
  task_type = "HTTP"
  interval  = 5
  isp_cities {
    isp  = "232"
    city = "641"
    type = "IDC"
  }
  option_json {
    response_content     = "example"
    expect_value         = "example"
    port                 = 81
    is_base_encode       = true
    ping_num             = 5
    match_rule           = 1
    failure_rate         = "0.3"
    request_content      = "example"
    attempts             = 4
    request_format       = "hex"
    password             = "YourPassword123!"
    diagnosis_ping       = true
    response_format      = "hex"
    cookie               = "key2=value2"
    ping_port            = 443
    user_name            = "example"
    dns_match_rule       = "DNS_IN"
    timeout              = 3000
    dns_server           = "223.6.6.6"
    diagnosis_mtr        = true
    header               = "key2:value2"
    min_tls_version      = "1.1"
    ping_type            = "udp"
    dns_type             = "NS"
    dns_hijack_whitelist = "DnsHijackWhitelist"
    http_method          = "post"
    assertions {
      operator = "lessThan"
      target   = 300
      type     = "response_time"
    }
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cms_site_monitor&spm=docs.r.cms_site_monitor.example&intl_lang=EN_US)

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
* `option_json` - (Optional, Set, Available since v1.262.0) The extended options of the protocol that is used by the site monitoring task. See [`option_json`](#option_json) below.
* `options_json` - (Optional, Deprecated since v1.262.0) Field `options_json` has been deprecated from provider version 1.262.0. New field `option_json` instead.
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

### `option_json`

The option_json supports the following:
* `assertions` - (Optional, List, Available since v1.262.0) Assertion configuration group. See [`assertions`](#option_json-assertions) below.
* `attempts` - (Optional, Int, Available since v1.262.0) Number of retries after DNS failed.
* `cookie` - (Optional, Available since v1.262.0) The Cookie that sends the HTTP request.
* `diagnosis_mtr` - (Optional, Available since v1.262.0) Whether to enable automatic MTR network diagnosis after a task failure. Value:
  - false: does not enable automatic MTR network diagnosis.
  - true to turn on automatic MTR network diagnostics.
* `diagnosis_ping` - (Optional, Available since v1.262.0) Whether to enable the automatic PING network delay detection after the task fails. Value:
  - false: does not enable automatic PING network delay detection.
  - true: Enable automatic PING network delay detection.
* `dns_hijack_whitelist` - (Optional, Available since v1.262.0) List of DNS hijacking configurations.
* `dns_match_rule` - (Optional, Available since v1.262.0) Matching Rules for DNS. Value:
  - IN_DNS: The alias or IP address that is expected to be resolved is in the DNS response.
  - DNS_IN: All DNS responses appear in the alias or IP address that is expected to be resolved.
  - EQUAL: the DNS response is exactly the same as the alias or IP address that is expected to be resolved.
  - ANY:DNS response and the alias or IP address expected to be resolved have an intersection.
* `dns_server` - (Optional, Available since v1.262.0) The IP address of the DNS server.

-> **NOTE:**  only applicable to DNS probe types.

* `dns_type` - (Optional, Available since v1.262.0) DNS resolution type. Only applicable to DNS probe types. Value:
  - A (default): specifies the IP address corresponding to the host name or domain name.
  - CNAME: maps multiple domain names to another domain name.
  - NS: specifies that the domain name is resolved by a DNS server.
  - MX: point domain name to a mail server address.
  - TXT: Description of host name or domain name. The text length is limited to 512 bytes, which is usually used as SPF(Sender Policy Framework) record, that is, anti-spam.
* `expect_value` - (Optional, Available since v1.262.0) The alias or address to be resolved.

-> **NOTE:**  This parameter applies only to DNS probe types.

* `failure_rate` - (Optional, Available since v1.262.0) Packet loss rate.

-> **NOTE:**  This parameter only applies to PING probe types.

* `header` - (Optional, Available since v1.262.0) HTTP request header.
* `http_method` - (Optional, Available since v1.262.0) HTTP request method. Value:
  - get
  - post
  - head
* `is_base_encode` - (Optional, ForceNew, Available since v1.262.0) Whether the parameter' Password' is Base64 encoded.
  - true: Yes.
  - false: No.
* `match_rule` - (Optional, Int, Available since v1.262.0) Whether alarm rules are included. Value:
  - 0: Yes.
  - 1: No.
* `min_tls_version` - (Optional, Available since v1.262.0) Minimum TLS version. By default, TLS1.2 and later versions are supported. TLS1.0 and 1.1 have been disabled. If they still need to be supported, the configuration can be changed.
* `password` - (Optional, Available since v1.262.0) The password of the SMTP, POP3, or FTP probe type.
* `ping_num` - (Optional, Int, Available since v1.262.0) The heartbeat of the PING probe type.
* `ping_port` - (Optional, Int, Available since v1.262.0) PING the port. Applies to TCP PING.
* `ping_type` - (Optional, Available since v1.262.0) The PING protocol type. Value:
  - icmp
  - tcp
  - udp
* `port` - (Optional, Int, Available since v1.262.0) Ports of TCP, UDP, SMTP, and POP3 probe types.
* `request_content` - (Optional, Available since v1.262.0) The request content of the HTTP probe type.
* `request_format` - (Optional, Available since v1.262.0) HTTP request content format. Value:
  - hex: hexadecimal format.
  - text: text format.
* `response_content` - (Optional, Available since v1.262.0) Match the response content.
* `response_format` - (Optional, Available since v1.262.0) HTTP response content format. Value:
  - hex: hexadecimal format.
  - text: text format.
* `timeout` - (Optional, Int, Available since v1.262.0) Timeout time. Unit: milliseconds.
* `user_name` - (Optional, Available since v1.262.0) The username of FTP, SMTP, or pop3.

### `option_json-assertions`

The option_json-assertions supports the following:
* `operator` - (Optional, Available since v1.262.0) Assertion comparison operator. Value:
  - contains: contains.
  - doesNotContain: does not contain.
  - matches: regular matching.
  - doesNotMatch: regular mismatch.
  - is: Numeric equals or character matches equals.
  - isNot: not equal.
  - Lesthan: less.
  - moreThan: Greater.
* `target` - (Optional, Available since v1.262.0) Assertion matches the target numeric value or character of the comparison.
* `type` - (Optional, Available since v1.262.0) Assertion type. Value:
  - response_time: determines whether the response time is as expected.
  - status_code: determines whether the HTTP response status code is as expected.
  - header: determines whether the fields in the response Header are as expected.
  - body_text: determines whether the content in the returned Body is as expected by text character matching.
  - body_json: uses JSON parsing (JSON Path) to determine whether the content in the returned Body meets expectations.
  - body_xml: uses XML parsing (XPath) to determine whether the content in the returned Body meets expectations.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource supplied above.
* `task_state` - (Deprecated since v1.262.0) Field `task_state` has been deprecated from provider version 1.262.0. New field `status` instead.
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
