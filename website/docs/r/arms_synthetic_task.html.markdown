---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_synthetic_task"
description: |-
  Provides a Alicloud ARMS Synthetic Task resource.
---

# alicloud_arms_synthetic_task

Provides a ARMS Synthetic Task resource. Cloud Synthetic task resources.

For information about ARMS Synthetic Task and how to use it, see [What is Synthetic Task](https://next.api.alibabacloud.com/document/ARMS/2019-08-08/CreateTimingSyntheticTask).

-> **NOTE:** Available since v1.215.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_arms_synthetic_task&exampleId=fb00a731-b3e0-d235-c299-5b5a5dff4cf8034fd07c&activeTab=example&spm=docs.r.arms_synthetic_task.0.fb00a731b3&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_arms_synthetic_task" "default" {
  monitors {
    city_code     = "1200101"
    operator_code = "246"
    client_type   = "4"
  }
  synthetic_task_name = var.name
  custom_period {
    end_hour   = "12"
    start_hour = "11"
  }
  available_assertions {
    type     = "IcmpPackLoss"
    operator = "neq"
    expect   = "200"
    target   = "example"
  }
  available_assertions {
    type     = "IcmpPackAvgLatency"
    operator = "lte"
    expect   = "1000"
  }
  available_assertions {
    type     = "IcmpPackMaxLatency"
    operator = "lte"
    expect   = "10000"
  }
  tags = {
    Created = "TF"
    For     = "example"
  }
  status = "RUNNING"
  monitor_conf {
    net_tcp {
      tracert_timeout = "1050"
      target_url      = "www.aliyun.com"
      connect_times   = "6"
      interval        = "300"
      timeout         = "3000"
      tracert_num_max = "2"
    }
    net_dns {
      query_method       = "1"
      timeout            = "5050"
      target_url         = "www.aliyun.com"
      dns_server_ip_type = "1"
      ns_server          = "61.128.114.167"
    }
    api_http {
      timeout    = "10050"
      target_url = "https://www.aliyun.com"
      method     = "POST"
      request_headers = {
        key1 = "value1"
      }
      request_body {
        content = "example2"
        type    = "text/html"
      }
      connect_timeout = "6000"
    }
    website {
      slow_element_threshold   = "5005"
      verify_string_blacklist  = "Failed"
      element_blacklist        = "a.jpg"
      disable_compression      = "1"
      ignore_certificate_error = "0"
      monitor_timeout          = "20000"
      redirection              = "0"
      dns_hijack_whitelist     = "www.aliyun.com:203.0.3.55"
      page_tamper              = "www.aliyun.com:|/cc/bb/a.gif"
      flow_hijack_jump_times   = "10"
      custom_header            = "1"
      disable_cache            = "1"
      verify_string_whitelist  = "Senyuan"
      target_url               = "http://www.aliyun.com"
      automatic_scrolling      = "1"
      wait_completion_time     = "5005"
      flow_hijack_logo         = "senyuan1"
      custom_header_content = {
        key1 = "value1"
      }
      filter_invalid_ip = "0"
    }
    file_download {
      white_list                             = "www.aliyun.com:203.0.3.55"
      monitor_timeout                        = "1050"
      ignore_certificate_untrustworthy_error = "0"
      redirection                            = "0"
      ignore_certificate_canceled_error      = "0"
      ignore_certificate_auth_error          = "0"
      ignore_certificate_out_of_date_error   = "0"
      ignore_certificate_using_error         = "0"
      connection_timeout                     = "6090"
      ignore_invalid_host_error              = "0"
      verify_way                             = "0"
      custom_header_content = {
        key1 = "value1"
      }
      target_url                      = "https://www.aliyun.com"
      download_kernel                 = "0"
      quick_protocol                  = "2"
      ignore_certificate_status_error = "1"
      transmission_size               = "128"
      validate_keywords               = "senyuan1"
    }
    stream {
      stream_monitor_timeout = "10"
      stream_address_type    = "0"
      player_type            = "2"
      custom_header_content = {
        key1 = "value1"
      }
      white_list  = "www.aliyun.com:203.0.3.55"
      target_url  = "https://acd-assets.alicdn.com:443/2021productweek/week1_s.mp4"
      stream_type = "1"
    }
    net_icmp {
      target_url      = "www.aliyun.com"
      interval        = "200"
      package_num     = "36"
      package_size    = "512"
      timeout         = "1000"
      tracert_enable  = "true"
      tracert_num_max = "1"
      tracert_timeout = "1200"
    }
  }
  task_type        = "1"
  frequency        = "1h"
  monitor_category = "1"
  common_setting {
    xtrace_region = "cn-beijing"
    custom_host {
      hosts {
        domain = "www.a.aliyun.com"
        ips = [
          "153.3.238.102"
        ]
        ip_type = "0"
      }
      hosts {
        domain = "www.shifen.com"
        ips = [
          "153.3.238.110",
          "114.114.114.114",
          "127.0.0.1"
        ]
        ip_type = "1"
      }
      hosts {
        domain = "www.aliyun.com"
        ips = [
          "153.3.238.110",
          "180.101.50.242",
          "180.101.50.188"
        ]
        ip_type = "0"
      }
      select_type = "1"
    }
    monitor_samples   = "1"
    ip_type           = "1"
    is_open_trace     = "true"
    trace_client_type = "1"
  }
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}
```

## Argument Reference

The following arguments are supported:
* `available_assertions` - (Optional) Assertion List. See [`available_assertions`](#available_assertions) below.
* `common_setting` - (Optional) Common settings. See [`common_setting`](#common_setting) below.
* `custom_period` - (Optional) Custom Cycle. See [`custom_period`](#custom_period) below.
* `frequency` - (Required) Frequency.
* `monitor_category` - (Required, ForceNew) Classification of selected monitors.
* `monitor_conf` - (Required) Monitoring configuration. See [`monitor_conf`](#monitor_conf) below.
* `monitors` - (Required) List of selected monitors. See [`monitors`](#monitors) below.
* `resource_group_id` - (Optional, Computed) Describes which resource group the resource belongs.
* `status` - (Optional, Computed) task status.
* `synthetic_task_name` - (Required) The name of synthetic task.
* `tags` - (Optional, Map) The list of tags.
* `task_type` - (Required, ForceNew) The type of synthetic task.

### `available_assertions`

The available_assertions supports the following:
* `expect` - (Required) Expected value.
* `operator` - (Required) Condition: gt: greater than; gte: greater than or equal to; lt: less than; te: less than or equal to; eq: equal to; neq: not equal to; ctn: contains; nctn: does not contain; exist: exists; n_exist: does not exist; belong: belongs to; reg_match: regular matching.
* `target` - (Optional) Check the target. If the target is HttpResCode, HttpResBody, or httpressetime, you do not need to specify the target. If the target is HttpResHead, you need to specify the key in the header. If the target is HttpResHead, you need to use jsonPath.
* `type` - (Required) Assertion type, including: httpresead, httpresead, HttpResBody, HttpResBodyJson, httpressetime, IcmpPackLoss (packet loss rate), IcmpPackMaxLatency (maximum packet delay ms), icmppackwebscreen, fmppackavglatency (average delay rendering), TraceRouteHops (number of hops), dnsarecname, websiteOnload (full load time), see the supplement below for specific use.

### `common_setting`

The common_setting supports the following:
* `custom_host` - (Optional) Custom host. See [`custom_host`](#common_setting-custom_host) below.
* `ip_type` - (Optional) IP Type:
  - 0: Automatic
  - 1:IPv4
  - 2:IPv6.
* `is_open_trace` - (Optional) Whether to enable link tracking.
* `monitor_samples` - (Optional) Whether the monitoring samples are evenly distributed:
  - 0: No
1: Yes.
* `trace_client_type` - (Optional) Link trace client type:
  - 0:ARMS Agent
  - 1:OpenTelemetry
  - 2:Jaeger.
* `xtrace_region` - (Optional) The link data is reported to the region.

### `common_setting-custom_host`

The common_setting-custom_host supports the following:
* `hosts` - (Required) The host list. See [`hosts`](#common_setting-custom_host-hosts) below.
* `select_type` - (Required) Selection method:
  - 0: Random
  - 1: Polling.

### `common_setting-custom_host-hosts`

The common_setting-custom_host-hosts supports the following:
* `domain` - (Required) Domain Name.
* `ip_type` - (Required) IpType.
* `ips` - (Required) The IP list.

### `custom_period`

The custom_period supports the following:
* `end_hour` - (Optional) End hours, 0-24.
* `start_hour` - (Optional) Starting hours, 0-24.

### `monitor_conf`

The monitor_conf supports the following:
* `api_http` - (Optional) HTTP(S) task configuration information. See [`api_http`](#monitor_conf-api_http) below.
* `file_download` - (Optional) File download type task configuration. See [`file_download`](#monitor_conf-file_download) below.
* `net_dns` - (Optional) The configuration parameters of the DNS dial test. Required when TaskType is 3. See [`net_dns`](#monitor_conf-net_dns) below.
* `net_icmp` - (Optional) ICMP dialing configuration parameters. Required when TaskType is 1. See [`net_icmp`](#monitor_conf-net_icmp) below.
* `net_tcp` - (Optional) The configuration parameters of TCP dial test. Required when TaskType is 2. See [`net_tcp`](#monitor_conf-net_tcp) below.
* `stream` - (Optional) Streaming Media Dial Test Configuration. See [`stream`](#monitor_conf-stream) below.
* `website` - (Optional) Website speed measurement type task configuration. See [`website`](#monitor_conf-website) below.

### `monitor_conf-api_http`

The monitor_conf-api_http supports the following:
* `connect_timeout` - (Optional) Connection timeout, in ms. Default 5000. Optional range: 1000-300000ms.
* `method` - (Optional) HTTP method, GET or POST.
* `request_body` - (Optional) HTTP request body. See [`request_body`](#monitor_conf-api_http-request_body) below.
* `request_headers` - (Optional, Map) HTTP request header.
* `target_url` - (Required) Dial test target address (request path).
* `timeout` - (Optional) Timeout, unit: ms. Default 10000. Optional range: 1000-300000ms.

### `monitor_conf-file_download`

The monitor_conf-file_download supports the following:
* `connection_timeout` - (Optional) Connection timeout time, in ms. Default 5000. Optional range: 1000-120000ms.
* `custom_header_content` - (Optional, Map) Custom request header content, JSON Map.
* `download_kernel` - (Optional) Download the kernel.
  - 1:curl
  - 0:WinInet
Default 1.
* `ignore_certificate_auth_error` - (Optional) Ignore CA Certificate authorization error 0: Do not ignore, 1: ignore, default 1.
* `ignore_certificate_canceled_error` - (Optional) Ignore certificate revocation error 0: Do not ignore, 1: ignore, default 1.
* `ignore_certificate_out_of_date_error` - (Optional) Ignore certificate expiration error 0: not ignored, 1: Ignored, default 1.
* `ignore_certificate_status_error` - (Optional) The certificate status error is ignored. 0: Do not ignore, 1: IGNORE. The default value is 1.
* `ignore_certificate_untrustworthy_error` - (Optional) The certificate cannot be trusted and ignored. 0: Do not ignore, 1: IGNORE. The default value is 1.
* `ignore_certificate_using_error` - (Optional) Ignore certificate usage error 0: Do not ignore, 1: ignore, default 1.
* `ignore_invalid_host_error` - (Optional) Invalid host error ignored, 0: not ignored, 1: Ignored, default 1.
* `monitor_timeout` - (Optional) Monitoring timeout time, ms, default 60000, optional range: 1000~120000ms.
* `quick_protocol` - (Optional) Quick agreement
  - 1:http1
  - 2:http2
  - 3:http3
Default 1.
* `redirection` - (Optional) Whether to support redirection, 0: not supported, 1: Supported, default 1.
* `target_url` - (Required) File download link.
* `transmission_size` - (Optional) The transmission size, in KB. The default value is 2048KB. The transmission size of the downloaded file must be between 1 and 20480KB.
* `validate_keywords` - (Optional) Verify keywords.
* `verify_way` - (Optional) The verification method.
  - 0: Do not validate
  - 1: Validation string
  - 2:MD5 validation.
* `white_list` - (Optional) DNS hijack whitelist. Match rules support IP, IP wildcard, subnet mask, and CNAME. Multiple match rules can be filled in. Multiple match rules are separated by vertical bars (|). For example, www.aliyun.com:203.0.3.55 | 203.3.44.67 indicates that all other IP addresses under the www.aliyun.com domain except 203.0.3.55 and 203.3.44.67 are hijacked.

### `monitor_conf-net_dns`

The monitor_conf-net_dns supports the following:
* `dns_server_ip_type` - (Optional) The IP address type of the DNS server.
  - 0 (default):ipv4
  - 1:ipv6
2: Automatic.
* `ns_server` - (Optional) The IP address of the NS server. The default value is 114.114.114.114.
* `query_method` - (Optional) DNS query method.
  - 0 (default): Recursive
  - 1: Iteration.
* `target_url` - (Required) The destination address (domain name) of the DNS dial test.
* `timeout` - (Optional) DNS dial test timeout. The unit is milliseconds (ms), the minimum value is 1000, the maximum value is 45000, and the default value is 5000.

### `monitor_conf-net_icmp`

The monitor_conf-net_icmp supports the following:
* `interval` - (Optional) The interval at which ICMP(Ping) packets are sent. The unit is milliseconds (ms), the minimum value is 200, the maximum value is 2000, and the default value is 200.
* `package_num` - (Optional) Number of ICMP(Ping) packets sent. The minimum value is 1, the maximum value is 50, and the default is 4.
* `package_size` - (Optional) The size of the sent ICMP(Ping) packet. The unit is byte. The ICMP(PING) packet size is limited to 32, 64, 128, 256, 512, 1024, 1080, and 1450.
* `split_package` - (Optional) Whether to split ICMP(Ping) packets. The default is true.
* `target_url` - (Required) TargetUrl.
* `timeout` - (Optional) The timeout period of the ICMP dial test. The unit is milliseconds (ms), the minimum value is 1000, the maximum value is 300000, and the default value is 20000.
* `tracert_enable` - (Optional) Whether to enable tracert. The default is true.
* `tracert_num_max` - (Optional) The maximum number of hops for tracert. The minimum value is 1, the maximum value is 128, and the default value is 20.
* `tracert_timeout` - (Optional) The time-out of tracert. The unit is milliseconds (ms), the minimum value is 1000, the maximum value is 300000, and the default value is 60000.

### `monitor_conf-net_tcp`

The monitor_conf-net_tcp supports the following:
* `connect_times` - (Optional) The number of TCP connections established. The minimum value is 1, the maximum value is 16, and the default is 4.
* `interval` - (Optional) The interval between TCP connections. The unit is milliseconds (ms), the minimum value is 200, the maximum value is 10000, and the default value is 200.
* `target_url` - (Required) Dial test target address (host).
* `timeout` - (Optional) TCP dial test timeout. The unit is milliseconds (ms), the minimum value is 1000, the maximum value is 300000, and the default value is 20000.
* `tracert_enable` - (Optional) Whether to enable tracert. The default is true.
* `tracert_num_max` - (Optional) The maximum number of hops for tracert. The minimum value is 1, the maximum value is 128, and the default value is 20.
* `tracert_timeout` - (Optional) The time-out of tracert. The unit is milliseconds (ms), the minimum value is 1000, the maximum value is 300000, and the default value is 60000.

### `monitor_conf-stream`

The monitor_conf-stream supports the following:
* `custom_header_content` - (Optional, Map) Custom header, in JSON Map format.
* `player_type` - (Optional) Player, do not pass the default 12.
  - 12:VLC
  - 2:FlashPlayer.
* `stream_address_type` - (Optional) Resource address type:
  - 1: Resource address.
  - 0: page address, not 0 by default.
* `stream_monitor_timeout` - (Optional) Monitoring duration, in seconds, up to 60s, not 60 by default.
* `stream_type` - (Optional) Audio and video flags: 0-video, 1-audio.
* `target_url` - (Optional) The streaming media resource address.
* `white_list` - (Optional) DNS hijack whitelist. Match rules support IP, IP wildcard, subnet mask, and CNAME. Multiple match rules can be filled in. Multiple match rules are separated by vertical bars (|). For example, www.aliyun.com:203.0.3.55 | 203.3.44.67 indicates that all other IP addresses under the www.aliyun.com domain except 203.0.3.55 and 203.3.44.67 are hijacked.

### `monitor_conf-website`

The monitor_conf-website supports the following:
* `automatic_scrolling` - (Optional) Whether to support automatic scrolling screen, loading page.
  - 0 (default): No
1: Yes.
* `custom_header` - (Optional) Custom header.
  - 0 (default): Off
  - 1: Modify the first package
  - 2: Modify all packages.
* `custom_header_content` - (Optional, Map) Custom header, in JSON Map format.
* `disable_cache` - (Optional) Whether to disable caching.
  - 0: not disabled
  - 1 (default): Disabled.
* `disable_compression` - (Optional) The Accept-Encoding field is used to determine whether to Accept compressed files. 0-do not disable, 1-disable, the default is 0.
* `dns_hijack_whitelist` - (Optional) When a domain name (such as www.aliyun.com) is resolved, if the resolved IP address or CNAME is not in the DNS hijacking white list, the user will fail to access or return a target IP address that is not Aliyun. If the IP or CNAME in the resolution result is in the DNS white list, it will be determined that DNS hijacking has not occurred.  Fill in the format: Domain name: matching rules. Match rules support IP, IP wildcard, subnet mask, and CNAME. Multiple match rules can be filled in. Multiple match rules are separated by vertical bars (|). For example, www.aliyun.com:203.0.3.55 | 203.3.44.67 indicates that all other IP addresses under the www.aliyun.com domain except 203.0.3.55 and 203.3.44.67 are hijacked.
* `element_blacklist` - (Optional) If an element configured in the element blacklist appears during page loading, the element is not requested to be loaded.
* `filter_invalid_ip` - (Optional) Whether to filter invalid IP parameters. 0: filter, 1: do not filter. The default value is 0.
* `flow_hijack_jump_times` - (Optional) Identify elements: Set the total number of elements on the Browse page.
* `flow_hijack_logo` - (Optional) Hijacking ID: Set the matching key information. Enter the hijacking keyword or key element, with an asterisk (*) allowed.
* `ignore_certificate_error` - (Optional) Whether to ignore certificate errors during certificate verification in SSL Handshake and continue browsing. 0-do not ignore, 1-ignore. The default value is 1.
* `monitor_timeout` - (Optional) Monitoring timeout, in ms. Not required, 20000 by default.
* `page_tamper` - (Optional) Monitoring the page appears to be tampered with elements other than the domain settings that belong to the page. Common manifestations are pop-up advertisements, floating advertisements, jumps, etc.  Fill in the format: Domain name: Element. You can fill multiple elements separated by a vertical bar (|). For example, www.aliyun.com:|/cc/bb/a.gif |/vv/bb/cc.jpg indicates that all the other elements of the www.aliyun.com domain name except the basic document,/cc/bb/a.gif, and/vv/bb/cc.jpg are tampered.
* `redirection` - (Optional) When redirection occurs, whether to continue browsing, 0-No, 1-Yes, the default is 1.
* `slow_element_threshold` - (Optional) The slow element threshold, in ms, is 5000 by default and can be selected from 1 to 300000ms.
* `target_url` - (Required) The target URL.
* `verify_string_blacklist` - (Optional) The verification string is an arbitrary string in the source code of the monitoring page. If the source code returned by the client contains any of the blacklisted strings, 650 error is returned. Multiple strings are separated by a vertical bar (|).
* `verify_string_whitelist` - (Optional) The verification string is an arbitrary string in the source code of the monitoring page. The source code returned by the client must contain all the strings in the whitelist. Otherwise, 650 error is returned. Multiple strings are separated by a vertical bar (|).
* `wait_completion_time` - (Optional) The maximum waiting time, in ms, is 5000 by default and can be selected from 5000 ms to 300000ms.

### `monitor_conf-api_http-request_body`

The monitor_conf-api_http-request_body supports the following:
* `content` - (Optional) The request body content, in JSON string format. When the type is text/plain,application/json,application/xml,text/html, the content can be converted to a JSON string.
* `type` - (Optional) The request body type. Supported parameters include text/plain, application/json, application/x-www-form-urlencoded, multipart/form-data, application/xml, and text/html.

### `monitors`

The monitors supports the following:
* `city_code` - (Required) The city code of monitor.
* `client_type` - (Required) The type of monitor.
* `operator_code` - (Required) The operator code of monitor.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Synthetic Task.
* `delete` - (Defaults to 5 mins) Used when delete the Synthetic Task.
* `update` - (Defaults to 5 mins) Used when update the Synthetic Task.

## Import

ARMS Synthetic Task can be imported using the id, e.g.

```shell
$ terraform import alicloud_arms_synthetic_task.example <id>
```