---
subcategory: "Private Zone"
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_zone"
sidebar_current: "docs-alicloud-resource-pvtz-zone"
description: |-
  Provides a Alicloud Private Zone resource.
---

# alicloud_pvtz_zone

Provides a Private Zone resource.

-> **NOTE:** Terraform will auto Create a Private Zone while it uses `alicloud_pvtz_zone` to build a Private Zone resource.

-> **NOTE:** Available since v1.13.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_pvtz_zone&exampleId=996411e1-9008-9cba-134f-db2e142804cf4085f3e0&activeTab=example&spm=docs.r.pvtz_zone.0.996411e190&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_pvtz_zone" "foo" {
  zone_name = "foo.example.com"
}

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_pvtz_zone&spm=docs.r.pvtz_zone.example&intl_lang=EN_US)
```
## Argument Reference

The following arguments are supported:

* `name` - (Optional, ForceNew, Deprecated from v1.107.0+) The name of the Private Zone. The `name` has been deprecated from provider version 1.107.0. Please use 'zone_name' instead.
* `zone_name` - (Optional, ForceNew) The zone_name of the Private Zone. The `zone_name` is required when the value of the `name`  is Empty.
* `remark` - (Optional) The remark of the Private Zone.
* `proxy_pattern` - (Optional, Available in 1.69.0+) The recursive DNS proxy. Valid values:
    - ZONE: indicates that the recursive DNS proxy is disabled.
    - RECORD: indicates that the recursive DNS proxy is enabled.
    Default to "ZONE".
* `user_client_ip` - (Optional, Available in 1.69.0+) The IP address of the client.
* `lang` - (Optional, Available in 1.69.0+) The language. Valid values: "zh", "en", "jp".
* `resource_group_id` - (Optional, ForceNew, Available in v1.86.0+) The Id of resource group which the Private Zone belongs.
* `sync_status` - (Optional, Available in 1.146.0+) The status of the host synchronization task. Valid values:  `ON`,`OFF`. **NOTE:** You can update the `sync_status` to enable/disable the host synchronization task.
* `user_info` - (Optional, Available in 1.146.0+) The user information of the host synchronization task. See [`user_info`](#user_info) below.
* `tags` - (Optional, Map) The tags of the Private Zone.

### `user_info`

The user_info supports the following:
* `user_id` - (Optional, Available in 1.146.0+) The user ID belonging to the region is used for cross-account synchronization scenarios.
* `region_ids` - (Optional, Available in 1.146.0+) The list of the region IDs.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Private Zone.
* `record_count` - The count of the Private Zone Record.
* `is_ptr` - Whether the Private Zone is ptr.
* `creation_time` - (Removed since v1.107.0+) The create time of the Private Zone.
* `update_time` - (Removed since v1.107.0+) The update time of the Private Zone.

## Import

Private Zone can be imported using the id, e.g.

```shell
$ terraform import alicloud_pvtz_zone.example abc123456
```

