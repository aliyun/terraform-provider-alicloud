---
subcategory: "CDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_cdn_refresh_object_cache"
sidebar_current: "docs-alicloud-resource-cdn-refresh-object-cache"
description: |-
  Provides a Alicloud CDN Refresh Object Cache resource.
---

# alicloud_cdn_refresh_object_cache

Provides a CDN Refresh Object Cache resource.

For information about CDN Refresh Object Cache and how to use it, see [What is Refresh Object Cache](https://www.alibabacloud.com/help/en/cdn/developer-reference/api-cdn-2018-05-10-refreshobjectcaches).

-> **NOTE:** Available since v1.241.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cdn_refresh_object_cache" "example" {
  object_path = "https://www.example.com/path/file.html"
  object_type = "File"
  force       = true
}
```

Refresh a directory:

```terraform
resource "alicloud_cdn_refresh_object_cache" "example" {
  object_path = "https://www.example.com/path/"
  object_type = "Directory"
}
```

## Argument Reference

The following arguments are supported:

* `object_path` - (Required, ForceNew) The content to refresh. You can enter a URL or a directory. When `object_type` is `File`, enter a URL; when `object_type` is `Directory`, enter a directory path.
* `object_type` - (Optional, ForceNew) The type of the refresh task. Valid values: `File`, `Directory`. Default to `File`.
* `force` - (Optional, ForceNew) Specifies whether to forcibly refresh the content.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the refresh task.
* `status` - The status of the refresh task. Valid values: `Complete`, `Refreshing`, `Failed`.
* `process` - The progress of the refresh task.
* `creation_time` - The creation time of the refresh task.

-> **NOTE:** The refresh task cannot be cancelled after submission. Running `terraform destroy` will remove the resource from the Terraform state file, but the refresh task will remain on the cloud until it completes.

## Import

CDN Refresh Object Cache can be imported using the id, e.g.

```shell
$ terraform import alicloud_cdn_refresh_object_cache.example <task-id>
```
