---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_edge_container_app_version"
description: |-
  Provides a Alicloud ESA Edge Container App Version resource.
---

# alicloud_esa_edge_container_app_version

Provides a ESA Edge Container App Version resource. An edge container application version describes the container group (image, spec, probes, storage, environment variables, lifecycle hooks) to deploy for a given [edge container application](./esa_edge_container_app.html.markdown). The version is immutable once created; configuring `publish_env` publishes the version, and changing the publish parameters re-publishes it.

For information about ESA Edge Container App Version and how to use it, see [CreateEdgeContainerAppVersion](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateEdgeContainerAppVersion).

-> **NOTE:** Available since v1.249.0.

## Example Usage

Basic Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "tfexample"
}

resource "alicloud_esa_edge_container_app" "default" {
  edge_container_app_name = var.name
  service_port            = 8080
  target_port             = 80
}

resource "alicloud_esa_edge_container_app_version" "default" {
  app_id      = alicloud_esa_edge_container_app.default.id
  name        = var.name
  remarks     = "tf example version"
  publish_env = "staging"
  containers {
    name         = "tf-example-container"
    image        = "nginx:1.25"
    is_acr_image = false
    spec         = "1C2G"
    storage      = "0.5G"
    command      = "nginx -g 'daemon off;'"
    probe_type   = "tcpSocket"
    probe_content {
      initial_delay_seconds = 5
      timeout_seconds       = 3
      period_seconds        = 10
      failure_threshold     = 3
      success_threshold     = 1
      port                  = 80
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, ForceNew) The ID of the edge container application. The application must be created first via `alicloud_esa_edge_container_app`.
* `name` - (Required, ForceNew) The version name. The length must be 6 to 128 characters.
* `remarks` - (Optional, ForceNew) The remarks of the version.
* `containers` - (Required, ForceNew) The container group to deploy for this version. See [`containers`](#containers) below.
* `publish_env` - (Optional) The environment to publish the version to. Valid values: `prod`, `staging`. When set, the version is published on create and re-published whenever a publish parameter changes.
* `publish_type` - (Optional) The publish type. Valid values: `percentage`, `region`. Default value: `percentage`.
* `percentage` - (Optional) The publish percentage. Valid values: 1 to 100. Default value: 100. Only effective when `publish_type` is `percentage`.
* `regions` - (Optional) The list of regions to publish to. Only effective when `publish_type` is `region`.
* `full_release` - (Optional) Whether to perform a full release. Only effective when `publish_type` is `region`. Default value: `false`.
* `start_time` - (Optional) The time to start the publish, in ISO 8601 (UTC) format, e.g. `2023-06-05T16:00:00Z`. Defaults to the current time.

### containers

The `containers` block supports:

* `name` - (Required) The container name, unique within the container group.
* `image` - (Required) The image address.
* `is_acr_image` - (Required) Whether the image is an ACR image.
* `acr_image_info` - (Optional) The ACR image information. See [`acr_image_info`](#acr_image_info) below.
* `spec` - (Required) The container spec. Valid values: `1C2G`, `2C4G`, `2C8G`, `4C8G`, `4C16G`, `8C16G`, `8C32G`.
* `command` - (Optional) The startup command, with multiple arguments separated by spaces.
* `args` - (Optional) The startup arguments, with multiple arguments separated by spaces.
* `storage` - (Required) The storage spec. Valid values: `0.5G`, `10G`, `20G`, `30G`.
* `env_variables` - (Optional) The environment variables, in the format `key1=val1,key2=val2`.
* `post_start` - (Optional) The command executed before the container starts.
* `pre_stop` - (Optional) The command executed before the container stops.
* `probe_type` - (Required) The probe type. Valid values: `exec`, `tcpSocket`, `httpGet`.
* `probe_content` - (Required) The probe configuration. See [`probe_content`](#probe_content) below.

### acr_image_info

The `acr_image_info` block supports:

* `is_enterprise_registry` - (Optional) Whether the image is from an enterprise-grade registry.
* `domain` - (Optional) The ACR image domain.
* `instance_id` - (Optional) The ACR instance ID.
* `region_id` - (Optional) The region ID of the ACR instance.
* `repo_id` - (Optional) The repository ID.
* `repo_name` - (Optional) The repository name.
* `repo_namespace` - (Optional) The repository namespace.
* `tag` - (Optional) The image tag.
* `tag_url` - (Optional) The image tag URL.

### probe_content

The `probe_content` block supports:

* `initial_delay_seconds` - (Optional) The probe initial delay in seconds.
* `timeout_seconds` - (Optional) The probe timeout in seconds.
* `period_seconds` - (Optional) The probe interval in seconds.
* `success_threshold` - (Optional) The consecutive successes required to mark the probe successful.
* `failure_threshold` - (Optional) The consecutive failures required to mark the probe failed.
* `command` - (Optional) The probe command. Required when `probe_type` is `exec`.
* `port` - (Optional) The probe port. Required when `probe_type` is `tcpSocket` or `httpGet`.
* `host` - (Optional) The probe host. Only for `httpGet`.
* `path` - (Optional) The probe path. Only for `httpGet`.
* `scheme` - (Optional) The probe scheme, e.g. `http`. Only for `httpGet`.
* `http_headers` - (Optional) The HTTP request headers, as a JSON string. Only for `httpGet`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID, formatted as `<app_id>:<version_id>`.
* `version_id` - The version ID returned by CreateEdgeContainerAppVersion.
* `status` - The status of the version. Valid values: `created`, `failed`, `creating`.
* `create_time` - The time when the version was created.
* `update_time` - The last time the version was modified.
* `publish_time` - The time when the version was published.
* `last_publish_time` - The last time the version was published.

## Import

ESA Edge Container App Version can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_edge_container_app_version.example <app_id>:<version_id>
```
