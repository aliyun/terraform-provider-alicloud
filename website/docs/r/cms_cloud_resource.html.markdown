---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_cloud_resource"
sidebar_current: "docs-alicloud-resource-cms-cloud-resource"
description: |-
  Provides a resource to manage the Cloud Resource Center for Cloud Monitor.
---

# alicloud_cms_cloud_resource

Provides a Cloud Resource Center resource that enables centralized cloud resource management for the current region.

-> **NOTE:** Available since v1.247.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cms_cloud_resource" "default" {
}
```

## Argument Reference

The following arguments are supported:


This resource has no user-configurable arguments. The resource is a singleton per region and is created by calling the CreateCloudResource API. The region is determined by the provider configuration.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Cloud Resource Center. The value is the region ID.
* `region_id` - The region ID of the Cloud Resource Center.

## Import

Cloud Resource Center can be imported using the id, e.g. the region ID.

```shell
$ terraform import alicloud_cms_cloud_resource.default cn-hangzhou
```
