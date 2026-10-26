---
subcategory: "Security Center"
layout: "alicloud"
page_title: "Alicloud: alicloud_security_center_sca_details"
sidebar_current: "docs-alicloud-datasource-security-center-sca-details"
description: |-
  Provides the asset fingerprint details collected by Security Center.
---

# alicloud\_security\_center\_sca\_details

This data source provides the asset fingerprint details of the middleware, databases and web services collected by Security Center.

-> **NOTE:** Available since v2.2.0.

## Example Usage

```terraform
data "alicloud_security_center_sca_details" "default" {
  biz       = "sca_database"
  page_size = 10
}

output "sca_details" {
  value = data.alicloud_security_center_sca_details.default.propertys
}
```

## Argument Reference

The following arguments are supported:

* `lang` - (Optional, String) The language of the request and the response. Defaults to `zh`. Valid values: `zh` and `en`.
* `biz` - (Optional, String) The type of the asset fingerprint. Defaults to `sca`. Valid values: `sca` (middleware), `sca_database` (database) and `sca_web` (web service).
* `search_item` - (Optional, String) The type of the search condition. Valid values: `name` and `type`. It must be configured together with `search_info`.
* `search_info` - (Optional, String) The content of the search condition. It must be configured together with `search_item`.
* `search_criteria_list` - (Optional, List) The search conditions. A maximum of 20 conditions are supported. Each element contains the following attributes:
  * `name` - (Optional, String) The name of the search condition.
  * `value` - (Optional, String) The value of the search condition.
* `sca_name` - (Optional, String) The name of the asset fingerprint.
* `biz_type` - (Optional, String) The type of the middleware, database or web service. Valid values: `system_service`, `software_library`, `docker_component`, `database`, `web_container`, `jar` and `web_framework`.
* `search_item_sub` - (Optional, String) The type of the sub-condition. Valid values: `port`, `pid`, `version` and `user`. It must be configured together with `search_info_sub`.
* `search_info_sub` - (Optional, String) The content of the sub-condition. It must be configured together with `search_item_sub`.
* `remark` - (Optional, String) The name or IP address of the server. Fuzzy match is supported.
* `uuid` - (Optional, String) The UUID of the server on which the asset fingerprint is collected.
* `sca_version` - (Optional, String) The version of the middleware, database or web service.
* `port` - (Optional, String) The listening port of the process.
* `user` - (Optional, String) The user that runs the process.
* `pid` - (Optional, String) The ID of the process.
* `process_started_start` - (Optional, Int) The start of the process start time range, in seconds.
* `process_started_end` - (Optional, Int) The end of the process start time range, in seconds.
* `sca_name_pattern` - (Optional, String) The name of the process.
* `resource_directory_account_id` - (Optional, Int) The ID of the Alibaba Cloud account that is a member of the resource directory.
* `current_page` - (Optional, Computed, Int) The number of the page to return. Defaults to `1`.
* `page_size` - (Optional, Computed, Int) The number of entries per page. Defaults to `10`.
* `use_next_token` - (Optional, Bool) Specifies whether to use `next_token` instead of the paged query.
* `next_token` - (Optional, String) The token that starts the query. It is used together with `use_next_token`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the data source, which is derived from the query condition.
* `total_count` - (Int) The total number of the asset fingerprints that match the query condition. It is left unset when the API does not return it.
* `page_item_count` - (Int) The number of the asset fingerprints on the current page.
* `next_token_result` - (String) The token that is used to retrieve the next page of results.
* `propertys` - (List) A list of the asset fingerprint details. Each element contains the following attributes:
  * `type` - (String) The type of the middleware, database or web service.
  * `biz_type_dispaly` - (String) The display name of the type. The attribute name keeps the original spelling of the API.
  * `process_started` - (Int) The timestamp when the process started, in seconds.
  * `internet_ip` - (String) The public IP address of the server.
  * `config_path` - (String) The path of the configuration file.
  * `pid` - (String) The ID of the process.
  * `port` - (String) The listening port of the process.
  * `cmdline` - (String) The command line of the process.
  * `biz_type` - (String) The type of the asset fingerprint.
  * `listen_ip` - (String) The IP address that the process listens on.
  * `version` - (String) The version of the middleware, database or web service.
  * `instance_name` - (String) The name of the instance.
  * `listen_status` - (String) The listening status of the process. Valid values: `NONE` and `LISTEN`.
  * `name` - (String) The name of the middleware, database or web service.
  * `ip` - (String) The IP address of the server.
  * `process_user` - (String) The user that runs the process.
  * `instance_id` - (String) The ID of the instance.
  * `web_path` - (String) The path of the web service.
  * `ppid` - (String) The ID of the parent process.
  * `intranet_ip` - (String) The private IP address of the server.
  * `uuid` - (String) The UUID of the server.
  * `listen_protocol` - (String) The protocol that the process listens on. Valid values: `UDP` and `TCP`.
  * `image_name` - (String) The name of the image.
  * `path` - (String) The path of the process.
  * `container_name` - (String) The name of the container.
  * `proof` - (String) The evidence of the asset fingerprint.
  * `create_timestamp` - (Int) The timestamp when the asset fingerprint was last scanned, in milliseconds.
  * `runtime_env_version` - (String) The version of the runtime environment.
  * `pod_name` - (String) The name of the pod.
