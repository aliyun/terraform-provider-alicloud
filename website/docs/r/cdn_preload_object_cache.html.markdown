---
subcategory: "CDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_cdn_preload_object_cache"
sidebar_current: "docs-alicloud-resource-cdn-preload-object-cache"
description: |-
  Provides a Alicloud CDN Preload Object Cache resource.
---

# alicloud_cdn_preload_object_cache

Provides a CDN Preload Object Cache resource.

For information about CDN Preload Object Cache and how to use it, see [What is Preload Object Cache](https://www.alibabacloud.com/help/en/cdn/developer-reference/api-cdn-2018-05-10-pushobjectcache).

-> **NOTE:** Available since v1.241.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cdn_preload_object_cache" "example" {
  object_path = "https://www.example.com/path/file.html"
  area        = "domestic"
  l2_preload  = false
}
```

## Argument Reference

The following arguments are supported:

* `object_path` - (Required, ForceNew) The content to preload. You can enter a URL.
* `area` - (Optional, ForceNew) The region for preload. Valid values: `domestic`, `overseas`. If you do not specify this parameter, the system determines the region based on the user's location.
* `l2_preload` - (Optional, ForceNew) Specifies whether to use L2 preload. Default to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the preload task.
* `status` - The status of the preload task. Valid values: `Complete`, `Refreshing`, `Failed`.
* `process` - The progress of the preload task.
* `creation_time` - The creation time of the preload task.

-> **NOTE:** The preload task cannot be cancelled after submission. Running `terraform destroy` will remove the resource from the Terraform state file, but the preload task will remain on the cloud until it completes.

## Import

CDN Preload Object Cache can be imported using the id, e.g.

```shell
$ terraform import alicloud_cdn_preload_object_cache.example <task-id>
```
