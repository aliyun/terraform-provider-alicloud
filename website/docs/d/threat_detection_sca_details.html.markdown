---
subcategory: "Threat Detection"
layout: "alicloud"
page_title: "Alicloud: alicloud_threat_detection_sca_details"
description: |-
  Provides a list of Threat Detection Sca Details to the user.
---

# alicloud_threat_detection_sca_details

This data source provides the Threat Detection Sca Details of the current Alibaba Cloud user.

-> **NOTE:** Available since v2.2.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_threat_detection_sca_details" "default" {
}

output "sca_details" {
  value = data.alicloud_threat_detection_sca_details.default.details
}
```

## Argument Reference

The following arguments are supported:

* `search_item` - (Optional, String) The type of the search condition. Valid values: `name` and `type`. It must be configured together with `search_info`.
* `search_info` - (Optional, String) The content of the search condition. It must be configured together with `search_item`.
* `search_item_sub` - (Optional, String) The type of the sub-condition. Valid values: `port`, `pid`, `version` and `user`. It must be configured together with `search_info_sub`.
* `search_info_sub` - (Optional, String) The content of the sub-condition. It must be configured together with `search_item_sub`.
* `biz` - (Optional, String) The type of asset fingerprint that you want to query. Default value: `sca`. Valid values:
  - `sca`: middleware.
  - `sca_database`: database.
  - `sca_web`: web service.
* `lang` - (Optional, String) The language of the content within the request and response. Default value: `zh`. Valid values:
  - `zh`: Chinese.
  - `en`: English.
* `output_file` - (Optional, String) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported:

* `details` - A list of the asset fingerprint details. Each element contains the following attributes:
  * `type` - The type of the middleware, database or web service.
  * `biz_type_dispaly` - The display name of the type. The attribute name keeps the original spelling of the API.
  * `process_started` - The timestamp when the process started, in seconds.
  * `internet_ip` - The public IP address of the server.
  * `config_path` - The path of the configuration file.
  * `pid` - The ID of the process.
  * `port` - The listening port of the process.
  * `cmdline` - The command line of the process.
  * `biz_type` - The type of the asset fingerprint.
  * `listen_ip` - The IP address that the process listens on.
  * `version` - The version of the middleware, database or web service.
  * `instance_name` - The name of the instance.
  * `listen_status` - The listening status of the process.
  * `name` - The name of the middleware, database or web service.
  * `ip` - The IP address of the server.
  * `process_user` - The user that runs the process.
  * `instance_id` - The ID of the instance.
  * `web_path` - The path of the web service.
  * `ppid` - The ID of the parent process.
  * `intranet_ip` - The private IP address of the server.
  * `uuid` - The UUID of the server.
  * `listen_protocol` - The protocol that the process listens on.
  * `image_name` - The name of the image.
  * `path` - The path of the process.
  * `container_name` - The name of the container.
  * `proof` - The evidence of the asset fingerprint.
  * `create_timestamp` - The timestamp when the asset fingerprint was last scanned, in milliseconds.
  * `runtime_env_version` - The version of the runtime environment.
  * `pod_name` - The name of the pod.
