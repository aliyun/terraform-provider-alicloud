---
subcategory: "Elastic Container Instance (ECI)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eci_image_caches"
sidebar_current: "docs-alicloud-eci-image-caches"
description: |-
  Provides a collection of ECI Image Cache to the specified filters.
---

# alicloud\_eci\_image\_caches

Provides a collection of ECI Image Cache to the specified filters.

~> **NOTE:** Available in 1.90.0+.

## Example Usage

 ```
data "alicloud_eci_image_caches" "example" {
  ids = ["imc-bp1ef0dyp7ldhb1d****"]
}

output "image_cache" {
  value = data.alicloud_eci_image_caches.example.caches.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list ids of ECI Image Cache.
* `name_regex` - (Optional) A regex string to filter results by the image cache name.
* `image` - (Optional) Find the mirror cache containing it according to the image name.
* `image_cache_name` - (Optional) The name of ECI Image Cache.
* `snapshot_id` - (Optional) The id of snapshot.
* `status` - (Optional) The status of ECI Image Cache.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

 * `ids` - A list ids of ECI Image Cache.
 * `names` - A list of ECI Image Cache names.
 * `caches` - A list of caches. Each element contains the following attributes:
   * `id` - The ID of the ECI Image Cache.
   * `container_group_id` - The id of container group. 
   * `expire_date_time` - The time of expired.
   * `image_cache_id` - The id of the ECI Image Cache.
   * `image_cache_name` - The name of the ECI Image Cache.
   * `images` - The list of cached images.
   * `progress` - The progress of ECI Image Cache.
   * `snapshot_id` - The id of snapshot.
   * `status` - The status of ECI Image Cache.
   * `events` - Image cache pulls image event information.
     * `first_timestamp` - Start time.   
     * `last_timestamp` - End time.   
     * `count` - Number of events.   
     * `name` - The name of event.   
     * `type` - The type of event.  
