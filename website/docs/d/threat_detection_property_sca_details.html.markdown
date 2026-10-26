---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_property_sca_details"
sidebar_current: "docs-alicloud-datasource-threat_detection-property-sca-details"
description: |-
  Provides a list of asset fingerprint details collected by Cloud Security Center.
---

# alicloud_threat_detection_property_sca_details

This data source provides a list of asset fingerprint details owned by an Alibaba Cloud account. It calls the [DescribePropertyScaDetail](https://www.alibabacloud.com/help/en/security-center/developer-reference/api-sas-2018-12-03-describepropertyscadetail) API and returns the middleware, database and web service fingerprints collected by Cloud Security Center (SAS).

-> **NOTE:** Available since v2.2.0.

-> **NOTE:** All matching pages are collected in a single read, so `properties` is the complete result set of the query rather than one page. On accounts with a large number of fingerprints, narrow the query with `sca_name`, `biz_type`, `uuid` or `search_criteria_list` to keep the read short. The API allows 10 requests per second, and the reads back off and retry when the account hits that limit.

## Example Usage

```terraform
data "alicloud_threat_detection_property_sca_details" "example" {
  biz       = "sca"
  sca_name  = "openssl"
  page_size = 10
}

output "sca_names" {
  value = [for item in data.alicloud_threat_detection_property_sca_details.example.properties : item.name]
}
```

Query the fingerprints on a specified server:

```terraform
data "alicloud_threat_detection_property_sca_details" "by_uuid" {
  biz       = "sca"
  uuid      = "uuid-02ebabe7-1c19-4a3b-9c2d-0f1e2d3c4b5a"
  page_size = 20
}
```

Page through the results with `NextToken` instead of page numbers:

```terraform
data "alicloud_threat_detection_property_sca_details" "by_next_token" {
  biz            = "sca_database"
  use_next_token = true
  page_size      = 20
}
```

## Argument Reference

The following arguments are supported:

* `lang` - (Optional, String) The language of the request and the response. Valid values: `zh` and `en`. Defaults to `zh`.
* `biz` - (Optional, String) The type of the asset fingerprint. Valid values: `sca` (middleware), `sca_database` (database) and `sca_web` (web service). Defaults to `sca`.
* `biz_type` - (Optional, String) The type of the middleware, database or web service. Valid values: `system_service`, `software_library`, `docker_component`, `database`, `web_container`, `jar` and `web_framework`.
* `sca_name` - (Optional, String) The name of the asset fingerprint to query.
* `sca_name_pattern` - (Optional, String) The process name of the asset fingerprint to query.
* `sca_version` - (Optional, String) The version of the middleware, database or web service.
* `uuid` - (Optional, String) The UUID of the server that hosts the asset fingerprint.
* `search_item` - (Optional, String) The type of the search condition. Valid values: `name` and `type`. It must be set together with `search_info`: when `search_item` is `name`, `search_info` is the asset fingerprint name; when `search_item` is `type`, `search_info` is one of the `biz_type` values.
* `search_info` - (Optional, String) The search content. It must be set together with `search_item`, and it must not be empty: setting only one of the two, or an empty search content, is rejected at plan time because the API requires both to be set.
* `search_item_sub` - (Optional, String) The type of the sub-search condition. Valid values: `port`, `pid`, `version` and `user`. It must be set together with `search_info_sub`.
* `search_info_sub` - (Optional, String) The content of the sub-search condition, matching the type set in `search_item_sub`. It must be set together with `search_item_sub` and must not be empty.
* `search_criteria_list` - (Optional, List) The search conditions. Up to 20 conditions are supported. Each condition supports the following attributes:
  * `name` - (Required, String) The name of the search condition.
  * `value` - (Required, String) The value of the search condition.
* `remark` - (Optional, String) The server name or IP address of the asset fingerprint. Fuzzy match is supported.
* `port` - (Optional, String) The listening port of the process.
* `user` - (Optional, String) The user that runs the process.
* `pid` - (Optional, String) The ID of the process.
* `process_started_start` - (Optional, Int) The beginning of the process start time range to query, as a second-level timestamp.
* `process_started_end` - (Optional, Int) The end of the process start time range to query, as a second-level timestamp.
* `resource_directory_account_id` - (Optional, Int) The ID of the Resource Directory member account.
* `current_page` - (Optional, Int) The page to start reading from. Defaults to `1`. It cannot be set when `use_next_token` is `true`. All pages from this page to the last one are collected in one read.
* `page_size` - (Optional, Int) The number of entries to return on each page. Defaults to `10`.
* `use_next_token` - (Optional, Bool) Whether to page through the results with `NextToken` instead of `CurrentPage` and `PageSize`. Defaults to `false`. When it is `true`, `current_page` must not be set and `page_info.total_count` is not returned.
* `next_token` - (Optional, String) The `NextToken` to start reading from. It can only be set when `use_next_token` is `true`, and all remaining pages are collected in one read.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - (String) The ID of the data source.
* `properties` - (List) The asset fingerprint entries. Entries collected over multiple pages are concatenated in page order. Each element contains the following attributes:
  * `type` - (String) The type of the middleware, database or web service.
  * `biz_type` - (String) The type of the middleware, database or web service. Valid values are the same as the `biz_type` argument.
  * `biz_type_display` - (String) The display name of the type.
  * `name` - (String) The name of the middleware, database or web service.
  * `version` - (String) The version of the middleware, database or web service.
  * `path` - (String) The path of the middleware, database or web service.
  * `proof` - (String) The version verification information of the middleware, database or web service.
  * `config_path` - (String) The path of the configuration file.
  * `cmdline` - (String) The command line of the process.
  * `pid` - (String) The ID of the process.
  * `ppid` - (String) The ID of the parent process.
  * `port` - (String) The port of the middleware, database or web service.
  * `listen_ip` - (String) The IP address that the process listens on.
  * `listen_protocol` - (String) The protocol that the process listens on. Valid values: `UDP` and `TCP`.
  * `listen_status` - (String) The listening status of the process. Valid values: `NONE` and `LISTEN`.
  * `process_user` - (String) The user that runs the process.
  * `process_started` - (Int) The start time of the process, as a second-level timestamp.
  * `create_timestamp` - (Int) The last scan time of the fingerprint, as a millisecond-level timestamp.
  * `runtime_env_version` - (String) The version of the runtime environment, for example the JDK version of a Java process.
  * `instance_id` - (String) The ID of the server instance.
  * `instance_name` - (String) The name of the server instance.
  * `uuid` - (String) The UUID of the server.
  * `ip` - (String) The public IP address of the server.
  * `internet_ip` - (String) The public IP address of the server.
  * `intranet_ip` - (String) The private IP address of the server.
  * `web_path` - (String) The web directory.
  * `container_name` - (String) The name of the container.
  * `image_name` - (String) The name of the image.
  * `pod_name` - (String) The name of the pod.
* `page_info` - (Object) The pagination information of the query. It contains the following attributes:
  * `current_page` - (Int) The first page read. It is `null` when `use_next_token` is `true`.
  * `page_size` - (Int) The requested page size.
  * `total_count` - (Int) The total number of matches reported by the API. It is `null` when `use_next_token` is `true`.
  * `count` - (Int) The number of returned `properties` after all pages were collected.
  * `next_token` - (String) The trailing `NextToken` of the NextToken pagination. It is `null` when `use_next_token` is `false`.
